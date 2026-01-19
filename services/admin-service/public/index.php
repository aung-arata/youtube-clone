<?php
/**
 * Admin Service - Main Entry Point
 * Lightweight admin dashboard and CMS for YouTube Clone
 */

require_once __DIR__ . '/../vendor/autoload.php';

use App\Controller\AdminController;
use App\Controller\CMSController;
use App\Controller\EmailController;
use App\Controller\ReportController;

// Load environment variables
$dotenv = parse_ini_file(__DIR__ . '/../.env');
foreach ($dotenv as $key => $value) {
    $_ENV[$key] = $value;
}

// Simple router
$requestUri = parse_url($_SERVER['REQUEST_URI'], PHP_URL_PATH);
$requestMethod = $_SERVER['REQUEST_METHOD'];

// CORS headers
header('Access-Control-Allow-Origin: *');
header('Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS');
header('Access-Control-Allow-Headers: Content-Type, Authorization');

if ($requestMethod === 'OPTIONS') {
    http_response_code(200);
    exit;
}

// Route handling
try {
    switch (true) {
        // Health check
        case $requestUri === '/health':
            header('Content-Type: application/json');
            echo json_encode(['status' => 'ok', 'service' => 'admin-service']);
            break;
            
        // Admin Dashboard Routes
        case preg_match('#^/admin/?$#', $requestUri):
            $controller = new AdminController();
            $controller->dashboard();
            break;
            
        case preg_match('#^/admin/users/?$#', $requestUri):
            $controller = new AdminController();
            if ($requestMethod === 'GET') {
                $controller->listUsers();
            }
            break;
            
        case preg_match('#^/admin/moderation/?$#', $requestUri):
            $controller = new AdminController();
            $controller->moderationQueue();
            break;
            
        // CMS Routes
        case preg_match('#^/cms/blog/?$#', $requestUri):
            $controller = new CMSController();
            if ($requestMethod === 'GET') {
                $controller->listBlogs();
            } elseif ($requestMethod === 'POST') {
                $controller->createBlog();
            }
            break;
            
        case preg_match('#^/cms/blog/(\d+)/?$#', $requestUri, $matches):
            $controller = new CMSController();
            $id = $matches[1];
            if ($requestMethod === 'GET') {
                $controller->getBlog($id);
            } elseif ($requestMethod === 'PUT') {
                $controller->updateBlog($id);
            } elseif ($requestMethod === 'DELETE') {
                $controller->deleteBlog($id);
            }
            break;
            
        case preg_match('#^/cms/docs/?$#', $requestUri):
            $controller = new CMSController();
            if ($requestMethod === 'GET') {
                $controller->listDocs();
            } elseif ($requestMethod === 'POST') {
                $controller->createDoc();
            }
            break;
            
        case preg_match('#^/cms/help/?$#', $requestUri):
            $controller = new CMSController();
            if ($requestMethod === 'GET') {
                $controller->listHelpArticles();
            } elseif ($requestMethod === 'POST') {
                $controller->createHelpArticle();
            }
            break;
            
        // Email Template Routes
        case preg_match('#^/email/templates/?$#', $requestUri):
            $controller = new EmailController();
            if ($requestMethod === 'GET') {
                $controller->listTemplates();
            } elseif ($requestMethod === 'POST') {
                $controller->createTemplate();
            }
            break;
            
        case preg_match('#^/email/templates/(\d+)/?$#', $requestUri, $matches):
            $controller = new EmailController();
            $id = $matches[1];
            if ($requestMethod === 'GET') {
                $controller->getTemplate($id);
            } elseif ($requestMethod === 'PUT') {
                $controller->updateTemplate($id);
            } elseif ($requestMethod === 'DELETE') {
                $controller->deleteTemplate($id);
            }
            break;
            
        case preg_match('#^/email/send/?$#', $requestUri):
            $controller = new EmailController();
            if ($requestMethod === 'POST') {
                $controller->sendEmail();
            }
            break;
            
        // Report Routes
        case preg_match('#^/reports/analytics/?$#', $requestUri):
            $controller = new ReportController();
            $controller->generateAnalytics();
            break;
            
        case preg_match('#^/reports/users/?$#', $requestUri):
            $controller = new ReportController();
            $controller->generateUserReport();
            break;
            
        case preg_match('#^/reports/videos/?$#', $requestUri):
            $controller = new ReportController();
            $controller->generateVideoReport();
            break;
            
        default:
            http_response_code(404);
            header('Content-Type: application/json');
            echo json_encode(['error' => 'Not Found']);
            break;
    }
} catch (Exception $e) {
    http_response_code(500);
    header('Content-Type: application/json');
    echo json_encode(['error' => $e->getMessage()]);
}
