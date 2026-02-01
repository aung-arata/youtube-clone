package transcoding

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

// QualityPreset defines video quality presets
type QualityPreset struct {
	Name       string
	Width      int
	Height     int
	Bitrate    int    // in kbps
	AudioRate  int    // in kbps
	CRF        int    // Constant Rate Factor (0-51, lower is better quality)
	Preset     string // FFmpeg preset (ultrafast, superfast, veryfast, faster, fast, medium, slow, slower, veryslow)
}

// Predefined quality presets
var QualityPresets = map[string]QualityPreset{
	"240p": {
		Name:      "240p",
		Width:     426,
		Height:    240,
		Bitrate:   400,
		AudioRate: 64,
		CRF:       28,
		Preset:    "medium",
	},
	"360p": {
		Name:      "360p",
		Width:     640,
		Height:    360,
		Bitrate:   800,
		AudioRate: 96,
		CRF:       27,
		Preset:    "medium",
	},
	"480p": {
		Name:      "480p",
		Width:     854,
		Height:    480,
		Bitrate:   1500,
		AudioRate: 128,
		CRF:       26,
		Preset:    "medium",
	},
	"720p": {
		Name:      "720p",
		Width:     1280,
		Height:    720,
		Bitrate:   3000,
		AudioRate: 128,
		CRF:       24,
		Preset:    "medium",
	},
	"1080p": {
		Name:      "1080p",
		Width:     1920,
		Height:    1080,
		Bitrate:   6000,
		AudioRate: 192,
		CRF:       22,
		Preset:    "medium",
	},
	"1440p": {
		Name:      "1440p",
		Width:     2560,
		Height:    1440,
		Bitrate:   10000,
		AudioRate: 256,
		CRF:       20,
		Preset:    "medium",
	},
	"4K": {
		Name:      "4K",
		Width:     3840,
		Height:    2160,
		Bitrate:   20000,
		AudioRate: 320,
		CRF:       18,
		Preset:    "medium",
	},
}

// TranscodingJob represents a video transcoding job
type TranscodingJob struct {
	ID            int
	VideoID       int
	TargetQuality string
	Status        string // pending, processing, completed, failed
	Progress      int    // 0-100
	ErrorMessage  string
	SourcePath    string
	OutputPath    string
	StartedAt     *time.Time
	CompletedAt   *time.Time
	CreatedAt     time.Time
}

// VideoQuality represents a transcoded video quality variant
type VideoQuality struct {
	ID        int
	VideoID   int
	Quality   string
	URL       string
	Bitrate   int
	Width     int
	Height    int
	Format    string
	FileSize  int64
	Status    string
	CreatedAt time.Time
}

// TranscodingService handles video transcoding operations
type TranscodingService struct {
	db            *sql.DB
	outputDir     string
	ffmpegPath    string
	maxConcurrent int
	jobQueue      chan *TranscodingJob
	wg            sync.WaitGroup
	ctx           context.Context
	cancel        context.CancelFunc
}

// NewTranscodingService creates a new transcoding service
func NewTranscodingService(db *sql.DB, outputDir string, maxConcurrent int) *TranscodingService {
	ctx, cancel := context.WithCancel(context.Background())

	// Find FFmpeg path
	ffmpegPath := "ffmpeg"
	if path := os.Getenv("FFMPEG_PATH"); path != "" {
		ffmpegPath = path
	}

	service := &TranscodingService{
		db:            db,
		outputDir:     outputDir,
		ffmpegPath:    ffmpegPath,
		maxConcurrent: maxConcurrent,
		jobQueue:      make(chan *TranscodingJob, 100),
		ctx:           ctx,
		cancel:        cancel,
	}

	// Start worker goroutines
	for i := 0; i < maxConcurrent; i++ {
		go service.worker(i)
	}

	return service
}

// worker processes transcoding jobs from the queue
func (s *TranscodingService) worker(id int) {
	for {
		select {
		case job := <-s.jobQueue:
			log.Printf("Worker %d: Processing job %d (video %d, quality %s)", id, job.ID, job.VideoID, job.TargetQuality)
			s.processJob(job)
		case <-s.ctx.Done():
			log.Printf("Worker %d: Shutting down", id)
			return
		}
	}
}

// QueueTranscoding adds a video to the transcoding queue for specified qualities
func (s *TranscodingService) QueueTranscoding(videoID int, sourcePath string, qualities []string) error {
	for _, quality := range qualities {
		preset, ok := QualityPresets[quality]
		if !ok {
			log.Printf("Unknown quality preset: %s", quality)
			continue
		}

		// Create job record in database
		outputPath := filepath.Join(s.outputDir, fmt.Sprintf("video_%d_%s.mp4", videoID, quality))

		query := `
			INSERT INTO transcoding_jobs (video_id, target_quality, status)
			VALUES ($1, $2, 'pending')
			ON CONFLICT (video_id, target_quality) DO UPDATE SET status = 'pending'
			RETURNING id
		`
		var jobID int
		err := s.db.QueryRow(query, videoID, quality).Scan(&jobID)
		if err != nil {
			return fmt.Errorf("failed to create transcoding job: %w", err)
		}

		// Create video quality record
		qualityQuery := `
			INSERT INTO video_qualities (video_id, quality, url, bitrate, width, height, format, status)
			VALUES ($1, $2, $3, $4, $5, $6, 'mp4', 'pending')
			ON CONFLICT (video_id, quality) DO UPDATE SET status = 'pending'
		`
		_, err = s.db.Exec(qualityQuery, videoID, quality, outputPath, preset.Bitrate, preset.Width, preset.Height)
		if err != nil {
			return fmt.Errorf("failed to create video quality record: %w", err)
		}

		job := &TranscodingJob{
			ID:            jobID,
			VideoID:       videoID,
			TargetQuality: quality,
			Status:        "pending",
			SourcePath:    sourcePath,
			OutputPath:    outputPath,
		}

		// Add to queue
		s.jobQueue <- job
	}

	return nil
}

// processJob handles the actual transcoding
func (s *TranscodingService) processJob(job *TranscodingJob) {
	// Update job status to processing
	now := time.Now()
	job.StartedAt = &now
	s.updateJobStatus(job.ID, "processing", 0, "")

	preset, ok := QualityPresets[job.TargetQuality]
	if !ok {
		s.updateJobStatus(job.ID, "failed", 0, "Unknown quality preset")
		return
	}

	// Build FFmpeg command
	args := []string{
		"-i", job.SourcePath,
		"-vf", fmt.Sprintf("scale=%d:%d", preset.Width, preset.Height),
		"-c:v", "libx264",
		"-preset", preset.Preset,
		"-crf", fmt.Sprintf("%d", preset.CRF),
		"-b:v", fmt.Sprintf("%dk", preset.Bitrate),
		"-c:a", "aac",
		"-b:a", fmt.Sprintf("%dk", preset.AudioRate),
		"-movflags", "+faststart",
		"-y", // Overwrite output file
		job.OutputPath,
	}

	cmd := exec.CommandContext(s.ctx, s.ffmpegPath, args...)

	// Run the command
	output, err := cmd.CombinedOutput()
	if err != nil {
		errMsg := fmt.Sprintf("FFmpeg error: %v\nOutput: %s", err, string(output))
		s.updateJobStatus(job.ID, "failed", 0, errMsg)
		s.updateQualityStatus(job.VideoID, job.TargetQuality, "failed")
		log.Printf("Transcoding failed for video %d quality %s: %s", job.VideoID, job.TargetQuality, errMsg)
		return
	}

	// Get file size
	fileInfo, err := os.Stat(job.OutputPath)
	var fileSize int64
	if err == nil {
		fileSize = fileInfo.Size()
	}

	// Update quality record with file size
	s.updateQualityFileSize(job.VideoID, job.TargetQuality, fileSize)
	s.updateQualityStatus(job.VideoID, job.TargetQuality, "ready")

	// Mark job as completed
	completedAt := time.Now()
	job.CompletedAt = &completedAt
	s.updateJobStatus(job.ID, "completed", 100, "")

	log.Printf("Transcoding completed for video %d quality %s", job.VideoID, job.TargetQuality)
}

// updateJobStatus updates a job's status in the database
func (s *TranscodingService) updateJobStatus(jobID int, status string, progress int, errorMsg string) {
	var query string
	var args []interface{}

	if status == "processing" {
		query = `UPDATE transcoding_jobs SET status = $1, progress = $2, started_at = NOW() WHERE id = $3`
		args = []interface{}{status, progress, jobID}
	} else if status == "completed" {
		query = `UPDATE transcoding_jobs SET status = $1, progress = $2, completed_at = NOW() WHERE id = $3`
		args = []interface{}{status, progress, jobID}
	} else if status == "failed" {
		query = `UPDATE transcoding_jobs SET status = $1, progress = $2, error_message = $3 WHERE id = $4`
		args = []interface{}{status, progress, errorMsg, jobID}
	} else {
		query = `UPDATE transcoding_jobs SET status = $1, progress = $2 WHERE id = $3`
		args = []interface{}{status, progress, jobID}
	}

	_, err := s.db.Exec(query, args...)
	if err != nil {
		log.Printf("Failed to update job status: %v", err)
	}
}

// updateQualityStatus updates a video quality's status
func (s *TranscodingService) updateQualityStatus(videoID int, quality, status string) {
	query := `UPDATE video_qualities SET status = $1 WHERE video_id = $2 AND quality = $3`
	_, err := s.db.Exec(query, status, videoID, quality)
	if err != nil {
		log.Printf("Failed to update quality status: %v", err)
	}
}

// updateQualityFileSize updates the file size for a video quality
func (s *TranscodingService) updateQualityFileSize(videoID int, quality string, fileSize int64) {
	query := `UPDATE video_qualities SET file_size = $1 WHERE video_id = $2 AND quality = $3`
	_, err := s.db.Exec(query, fileSize, videoID, quality)
	if err != nil {
		log.Printf("Failed to update quality file size: %v", err)
	}
}

// GetVideoQualities returns all available qualities for a video
func (s *TranscodingService) GetVideoQualities(videoID int) ([]VideoQuality, error) {
	query := `
		SELECT id, video_id, quality, url, bitrate, width, height, format, COALESCE(file_size, 0), status, created_at
		FROM video_qualities
		WHERE video_id = $1 AND status = 'ready'
		ORDER BY height DESC
	`

	rows, err := s.db.Query(query, videoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var qualities []VideoQuality
	for rows.Next() {
		var q VideoQuality
		err := rows.Scan(&q.ID, &q.VideoID, &q.Quality, &q.URL, &q.Bitrate, &q.Width, &q.Height, &q.Format, &q.FileSize, &q.Status, &q.CreatedAt)
		if err != nil {
			return nil, err
		}
		qualities = append(qualities, q)
	}

	return qualities, rows.Err()
}

// GetTranscodingStatus returns the transcoding status for a video
func (s *TranscodingService) GetTranscodingStatus(videoID int) ([]TranscodingJob, error) {
	query := `
		SELECT id, video_id, target_quality, status, progress, COALESCE(error_message, ''), started_at, completed_at, created_at
		FROM transcoding_jobs
		WHERE video_id = $1
		ORDER BY created_at DESC
	`

	rows, err := s.db.Query(query, videoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []TranscodingJob
	for rows.Next() {
		var j TranscodingJob
		var startedAt, completedAt sql.NullTime
		err := rows.Scan(&j.ID, &j.VideoID, &j.TargetQuality, &j.Status, &j.Progress, &j.ErrorMessage, &startedAt, &completedAt, &j.CreatedAt)
		if err != nil {
			return nil, err
		}
		if startedAt.Valid {
			j.StartedAt = &startedAt.Time
		}
		if completedAt.Valid {
			j.CompletedAt = &completedAt.Time
		}
		jobs = append(jobs, j)
	}

	return jobs, rows.Err()
}

// Shutdown gracefully shuts down the transcoding service
func (s *TranscodingService) Shutdown() {
	s.cancel()
	s.wg.Wait()
}
