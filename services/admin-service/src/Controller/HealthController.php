<?php

namespace App\Controller;

use Symfony\Bundle\FrameworkBundle\Controller\AbstractController;
use Symfony\Component\HttpFoundation\JsonResponse;
use Symfony\Component\Routing\Annotation\Route;

class HealthController extends AbstractController
{
    #[Route('/health', name: 'health_check', methods: ['GET'])]
    public function check(): JsonResponse
    {
        return $this->json([
            'status' => 'ok',
            'service' => 'admin-service',
            'framework' => 'Symfony 6.4'
        ]);
    }
}
