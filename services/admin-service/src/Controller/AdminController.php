<?php

namespace App\Controller;

use App\Repository\AdminUserRepository;
use App\Repository\ModerationQueueRepository;
use App\Service\GoServiceClient;
use Symfony\Bundle\FrameworkBundle\Controller\AbstractController;
use Symfony\Component\HttpFoundation\JsonResponse;
use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\Routing\Annotation\Route;

#[Route('/admin', name: 'admin_')]
class AdminController extends AbstractController
{
    public function __construct(
        private GoServiceClient $goClient,
        private AdminUserRepository $adminUserRepo,
        private ModerationQueueRepository $moderationRepo
    ) {
    }

    #[Route('', name: 'dashboard', methods: ['GET'])]
    public function dashboard(): JsonResponse
    {
        // Get stats from Go services
        $videoStats = $this->goClient->get('/videos', 'video');
        $userStats = $this->goClient->get('/users', 'user');
        
        $stats = [
            'total_videos' => is_array($videoStats) ? count($videoStats) : 0,
            'total_users' => is_array($userStats) ? count($userStats) : 0,
            'pending_moderation' => $this->moderationRepo->countByStatus('pending'),
            'service' => 'admin-dashboard'
        ];
        
        return $this->json($stats);
    }

    #[Route('/users', name: 'users', methods: ['GET'])]
    public function listUsers(Request $request): JsonResponse
    {
        $page = max(1, (int) $request->query->get('page', 1));
        $limit = min((int) $request->query->get('limit', 20), 100);
        
        // Get users from User Service
        $users = $this->goClient->get("/users?page={$page}&limit={$limit}", 'user');
        
        // Enrich with admin data
        $adminData = $this->adminUserRepo->findAll();
        $userMap = [];
        foreach ($adminData as $adminUser) {
            $userMap[$adminUser->getUserId()] = [
                'role' => $adminUser->getRole(),
                'status' => $adminUser->getStatus(),
                'last_login' => $adminUser->getLastLogin()?->format('Y-m-d H:i:s')
            ];
        }
        
        // Combine data
        if (is_array($users)) {
            foreach ($users as &$user) {
                if (isset($userMap[$user['id'] ?? null])) {
                    $user = array_merge($user, $userMap[$user['id']]);
                }
            }
        }
        
        return $this->json([
            'users' => $users,
            'page' => $page,
            'limit' => $limit
        ]);
    }

    #[Route('/moderation', name: 'moderation', methods: ['GET'])]
    public function moderationQueue(Request $request): JsonResponse
    {
        $status = $request->query->get('status', 'pending');
        
        $queue = $this->moderationRepo->findByStatus($status, 50);
        
        $result = [];
        foreach ($queue as $item) {
            $itemData = [
                'id' => $item->getId(),
                'content_type' => $item->getContentType(),
                'content_id' => $item->getContentId(),
                'reporter_id' => $item->getReporterId(),
                'reason' => $item->getReason(),
                'status' => $item->getStatus(),
                'created_at' => $item->getCreatedAt()->format('Y-m-d H:i:s')
            ];
            
            // Fetch associated content from Go services
            if ($item->getContentType() === 'video') {
                $itemData['content'] = $this->goClient->get("/videos/{$item->getContentId()}", 'video');
            } elseif ($item->getContentType() === 'comment') {
                $itemData['content'] = $this->goClient->get("/comments/{$item->getContentId()}", 'comment');
            }
            
            $result[] = $itemData;
        }
        
        return $this->json(['queue' => $result]);
    }
}
