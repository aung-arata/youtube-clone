<?php

namespace App\Controller;

use App\Service\DatabaseService;
use App\Service\GoServiceClient;

/**
 * Report Controller
 * Handles batch reporting and analytics
 */
class ReportController
{
    private DatabaseService $db;
    private GoServiceClient $goClient;
    
    public function __construct()
    {
        $this->db = new DatabaseService();
        $this->goClient = new GoServiceClient();
    }
    
    /**
     * Generate analytics report
     */
    public function generateAnalytics(): void
    {
        header('Content-Type: application/json');
        
        $period = $_GET['period'] ?? 'week'; // day, week, month
        
        // Fetch data from Go services
        $videos = $this->goClient->get('/videos', $_ENV['VIDEO_SERVICE_URL']);
        
        // Calculate statistics
        $totalViews = array_sum(array_column($videos ?? [], 'views'));
        $totalLikes = array_sum(array_column($videos ?? [], 'likes'));
        $totalVideos = count($videos ?? []);
        
        $report = [
            'period' => $period,
            'total_videos' => $totalVideos,
            'total_views' => $totalViews,
            'total_likes' => $totalLikes,
            'avg_views_per_video' => $totalVideos > 0 ? round($totalViews / $totalVideos, 2) : 0,
            'generated_at' => date('Y-m-d H:i:s')
        ];
        
        // Save report
        $this->db->insert('reports', [
            'type' => 'analytics',
            'period' => $period,
            'data' => json_encode($report),
            'generated_at' => date('Y-m-d H:i:s')
        ]);
        
        echo json_encode($report);
    }
    
    /**
     * Generate user report
     */
    public function generateUserReport(): void
    {
        header('Content-Type: application/json');
        
        // Get user statistics
        $adminUsers = $this->db->query(
            'SELECT COUNT(*) as total, 
                    SUM(CASE WHEN status = ? THEN 1 ELSE 0 END) as active,
                    SUM(CASE WHEN status = ? THEN 1 ELSE 0 END) as suspended
             FROM admin_users',
            ['active', 'suspended']
        );
        
        $report = [
            'total_users' => $adminUsers[0]['total'] ?? 0,
            'active_users' => $adminUsers[0]['active'] ?? 0,
            'suspended_users' => $adminUsers[0]['suspended'] ?? 0,
            'generated_at' => date('Y-m-d H:i:s')
        ];
        
        echo json_encode($report);
    }
    
    /**
     * Generate video report
     */
    public function generateVideoReport(): void
    {
        header('Content-Type: application/json');
        
        $videos = $this->goClient->get('/videos', $_ENV['VIDEO_SERVICE_URL']);
        
        // Group by category
        $byCategory = [];
        foreach ($videos ?? [] as $video) {
            $category = $video['category'] ?? 'Uncategorized';
            if (!isset($byCategory[$category])) {
                $byCategory[$category] = [
                    'count' => 0,
                    'total_views' => 0,
                    'total_likes' => 0
                ];
            }
            $byCategory[$category]['count']++;
            $byCategory[$category]['total_views'] += $video['views'] ?? 0;
            $byCategory[$category]['total_likes'] += $video['likes'] ?? 0;
        }
        
        $report = [
            'by_category' => $byCategory,
            'total_categories' => count($byCategory),
            'generated_at' => date('Y-m-d H:i:s')
        ];
        
        echo json_encode($report);
    }
}
