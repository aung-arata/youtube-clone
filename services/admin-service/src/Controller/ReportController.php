<?php

namespace App\Controller;

use App\Entity\Report;
use App\Service\GoServiceClient;
use Doctrine\ORM\EntityManagerInterface;
use Symfony\Bundle\FrameworkBundle\Controller\AbstractController;
use Symfony\Component\HttpFoundation\JsonResponse;
use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\Routing\Annotation\Route;

#[Route('/reports', name: 'reports_')]
class ReportController extends AbstractController
{
    public function __construct(
        private GoServiceClient $goClient,
        private EntityManagerInterface $em
    ) {
    }

    #[Route('/analytics', name: 'analytics', methods: ['GET'])]
    public function generateAnalytics(Request $request): JsonResponse
    {
        $period = $request->query->get('period', 'week');
        
        $videos = $this->goClient->get('/videos', 'video');
        
        $totalViews = 0;
        $totalLikes = 0;
        $totalVideos = is_array($videos) ? count($videos) : 0;
        
        if (is_array($videos)) {
            foreach ($videos as $video) {
                $totalViews += $video['views'] ?? 0;
                $totalLikes += $video['likes'] ?? 0;
            }
        }
        
        $reportData = [
            'period' => $period,
            'total_videos' => $totalVideos,
            'total_views' => $totalViews,
            'total_likes' => $totalLikes,
            'avg_views_per_video' => $totalVideos > 0 ? round($totalViews / $totalVideos, 2) : 0,
            'generated_at' => (new \DateTime())->format('Y-m-d H:i:s')
        ];
        
        // Save report
        $report = new Report();
        $report->setType('analytics');
        $report->setPeriod($period);
        $report->setData($reportData);
        
        $this->em->persist($report);
        $this->em->flush();
        
        return $this->json($reportData);
    }

    #[Route('/users', name: 'users', methods: ['GET'])]
    public function generateUserReport(): JsonResponse
    {
        // User statistics would come from database and Go services
        $reportData = [
            'total_users' => 0,
            'active_users' => 0,
            'suspended_users' => 0,
            'generated_at' => (new \DateTime())->format('Y-m-d H:i:s')
        ];
        
        return $this->json($reportData);
    }

    #[Route('/videos', name: 'videos', methods: ['GET'])]
    public function generateVideoReport(): JsonResponse
    {
        $videos = $this->goClient->get('/videos', 'video');
        
        $byCategory = [];
        if (is_array($videos)) {
            foreach ($videos as $video) {
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
        }
        
        $reportData = [
            'by_category' => $byCategory,
            'total_categories' => count($byCategory),
            'generated_at' => (new \DateTime())->format('Y-m-d H:i:s')
        ];
        
        return $this->json($reportData);
    }
}
