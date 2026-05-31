package transcoding

import (
	"testing"
)

func TestQualityPresets_AllDefined(t *testing.T) {
	expected := []string{"240p", "360p", "480p", "720p", "1080p", "1440p", "4K"}
	for _, name := range expected {
		if _, ok := QualityPresets[name]; !ok {
			t.Errorf("expected quality preset %q not found in QualityPresets", name)
		}
	}
}

func TestQualityPresets_ValidDimensions(t *testing.T) {
	for name, preset := range QualityPresets {
		if preset.Width <= 0 {
			t.Errorf("preset %q has invalid Width %d", name, preset.Width)
		}
		if preset.Height <= 0 {
			t.Errorf("preset %q has invalid Height %d", name, preset.Height)
		}
		if preset.Bitrate <= 0 {
			t.Errorf("preset %q has invalid Bitrate %d", name, preset.Bitrate)
		}
		if preset.AudioRate <= 0 {
			t.Errorf("preset %q has invalid AudioRate %d", name, preset.AudioRate)
		}
		if preset.Name == "" {
			t.Errorf("preset %q has empty Name field", name)
		}
	}
}

func TestQualityPresets_Resolution720p(t *testing.T) {
	p, ok := QualityPresets["720p"]
	if !ok {
		t.Fatal("720p preset not found")
	}
	if p.Width != 1280 {
		t.Errorf("720p width: got %d, want 1280", p.Width)
	}
	if p.Height != 720 {
		t.Errorf("720p height: got %d, want 720", p.Height)
	}
}

func TestQualityPresets_Resolution1080p(t *testing.T) {
	p, ok := QualityPresets["1080p"]
	if !ok {
		t.Fatal("1080p preset not found")
	}
	if p.Width != 1920 {
		t.Errorf("1080p width: got %d, want 1920", p.Width)
	}
	if p.Height != 1080 {
		t.Errorf("1080p height: got %d, want 1080", p.Height)
	}
}

func TestQualityPresets_CRFRange(t *testing.T) {
	for name, preset := range QualityPresets {
		if preset.CRF < 0 || preset.CRF > 51 {
			t.Errorf("preset %q CRF %d is outside valid range [0, 51]", name, preset.CRF)
		}
	}
}

func TestQualityPresets_HigherResolutionHigherBitrate(t *testing.T) {
	p360, ok360 := QualityPresets["360p"]
	p720, ok720 := QualityPresets["720p"]
	p1080, ok1080 := QualityPresets["1080p"]

	if !ok360 || !ok720 || !ok1080 {
		t.Fatal("missing expected presets")
	}
	if p360.Bitrate >= p720.Bitrate {
		t.Errorf("360p bitrate (%d) should be less than 720p bitrate (%d)", p360.Bitrate, p720.Bitrate)
	}
	if p720.Bitrate >= p1080.Bitrate {
		t.Errorf("720p bitrate (%d) should be less than 1080p bitrate (%d)", p720.Bitrate, p1080.Bitrate)
	}
}
