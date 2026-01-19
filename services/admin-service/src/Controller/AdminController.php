<?php

namespace App\Controller;

use App\Service\DatabaseService;
use App\Service\GoServiceClient;

/**
 * Admin Dashboard Controller
 * Handles admin operations: user management, content moderation, system config
 */
class AdminController
{
    private DatabaseService $db;
    private GoServiceClient $goClient;
    
    public function __construct()
    {
        $this->db = new DatabaseService();
        $this->goClient = new GoServiceClient();
    }
    
    /**
     * Admin Dashboard Home
     */
    public function dashboard(): void
    {
        header('Content-Type: application/json');
        
        // Get stats from Go services
        $videoStats = $this->goClient->get('/videos', $_ENV['VIDEO_SERVICE_URL']);
        $userStats = $this->goClient->get('/users', $_ENV['USER_SERVICE_URL']);
        
        $stats = [
            'total_videos' => count($videoStats ?? []),
            'total_users' => count($userStats ?? []),
            'pending_moderation' => $this->db->query('SELECT COUNT(*) as count FROM moderation_queue WHERE status = ?', ['pending'])[0]['count'] ?? 0,
            'service' => 'admin-dashboard'
        ];
        
        echo json_encode($stats);
    }
    
    /**
     * List all users with pagination
     */
    public function listUsers(): void
    {
        header('Content-Type: application/json');
        
        $page = $_GET['page'] ?? 1;
        $limit = min($_GET['limit'] ?? 20, 100);
        $offset = ($page - 1) * $limit;
        
        // Get users from User Service
        $users = $this->goClient->get("/users?page={$page}&limit={$limit}", $_ENV['USER_SERVICE_URL']);
        
        // Enrich with admin data
        $adminData = $this->db->query(
            'SELECT user_id, role, status, last_login FROM admin_users'
        );
        
        $userMap = [];
        foreach ($adminData as $row) {
            $userMap[$row['user_id']] = $row;
        }
        
        // Combine data
        if (is_array($users)) {
            foreach ($users as &$user) {
                if (isset($userMap[$user['id']])) {
                    $user = array_merge($user, $userMap[$user['id']]);
                }
            }
        }
        
        echo json_encode([
            'users' => $users,
            'page' => $page,
            'limit' => $limit
        ]);
    }
    
    /**
     * Content Moderation Queue
     */
    public function moderationQueue(): void
    {
        header('Content-Type: application/json');
        
        $status = $_GET['status'] ?? 'pending';
        
        $queue = $this->db->query(
            'SELECT * FROM moderation_queue WHERE status = ? ORDER BY created_at DESC LIMIT 50',
            [$status]
        );
        
        // Fetch associated content from Go services
        foreach ($queue as &$item) {
            if ($item['content_type'] === 'video') {
                $item['content'] = $this->goClient->get("/videos/{$item['content_id']}", $_ENV['VIDEO_SERVICE_URL']);
            } elseif ($item['content_type'] === 'comment') {
                $item['content'] = $this->goClient->get("/comments/{$item['content_id']}", $_ENV['COMMENT_SERVICE_URL']);
            }
        }
        
        echo json_encode(['queue' => $queue]);
    }
}
