<?php

namespace App\Controller;

use App\Entity\EmailTemplate;
use App\Entity\EmailLog;
use App\Repository\EmailTemplateRepository;
use Doctrine\ORM\EntityManagerInterface;
use Symfony\Bundle\FrameworkBundle\Controller\AbstractController;
use Symfony\Component\HttpFoundation\JsonResponse;
use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\HttpFoundation\Response;
use Symfony\Component\Routing\Annotation\Route;

#[Route('/email', name: 'email_')]
class EmailController extends AbstractController
{
    public function __construct(
        private EntityManagerInterface $em,
        private EmailTemplateRepository $templateRepo
    ) {
    }

    #[Route('/templates', name: 'templates_list', methods: ['GET'])]
    public function listTemplates(): JsonResponse
    {
        $templates = $this->templateRepo->findAll();
        
        return $this->json([
            'templates' => array_map(fn($t) => $t->toArray(), $templates)
        ]);
    }

    #[Route('/templates', name: 'templates_create', methods: ['POST'])]
    public function createTemplate(Request $request): JsonResponse
    {
        $data = json_decode($request->getContent(), true);
        
        if (empty($data['name']) || empty($data['subject']) || empty($data['html_content'])) {
            return $this->json(['error' => 'Missing required fields'], Response::HTTP_BAD_REQUEST);
        }
        
        $template = new EmailTemplate();
        $template->setName($data['name']);
        $template->setSubject($data['subject']);
        $template->setHtmlContent($data['html_content']);
        $template->setTextContent($data['text_content'] ?? strip_tags($data['html_content']));
        $template->setCategory($data['category'] ?? 'general');
        
        $this->em->persist($template);
        $this->em->flush();
        
        return $this->json([
            'id' => $template->getId(),
            'message' => 'Email template created'
        ], Response::HTTP_CREATED);
    }

    #[Route('/templates/{id}', name: 'templates_get', methods: ['GET'])]
    public function getTemplate(int $id): JsonResponse
    {
        $template = $this->templateRepo->find($id);
        
        if (!$template) {
            return $this->json(['error' => 'Template not found'], Response::HTTP_NOT_FOUND);
        }
        
        return $this->json($template->toArray());
    }

    #[Route('/templates/{id}', name: 'templates_update', methods: ['PUT'])]
    public function updateTemplate(int $id, Request $request): JsonResponse
    {
        $template = $this->templateRepo->find($id);
        
        if (!$template) {
            return $this->json(['error' => 'Template not found'], Response::HTTP_NOT_FOUND);
        }
        
        $data = json_decode($request->getContent(), true);
        
        if (isset($data['name'])) $template->setName($data['name']);
        if (isset($data['subject'])) $template->setSubject($data['subject']);
        if (isset($data['html_content'])) $template->setHtmlContent($data['html_content']);
        if (isset($data['text_content'])) $template->setTextContent($data['text_content']);
        if (isset($data['category'])) $template->setCategory($data['category']);
        
        $this->em->flush();
        
        return $this->json(['message' => 'Template updated']);
    }

    #[Route('/templates/{id}', name: 'templates_delete', methods: ['DELETE'])]
    public function deleteTemplate(int $id): JsonResponse
    {
        $template = $this->templateRepo->find($id);
        
        if (!$template) {
            return $this->json(['error' => 'Template not found'], Response::HTTP_NOT_FOUND);
        }
        
        $this->em->remove($template);
        $this->em->flush();
        
        return $this->json(['message' => 'Template deleted']);
    }

    #[Route('/send', name: 'send', methods: ['POST'])]
    public function sendEmail(Request $request): JsonResponse
    {
        $data = json_decode($request->getContent(), true);
        
        if (empty($data['to']) || empty($data['template_id'])) {
            return $this->json(['error' => 'Missing required fields'], Response::HTTP_BAD_REQUEST);
        }
        
        $template = $this->templateRepo->find($data['template_id']);
        
        if (!$template) {
            return $this->json(['error' => 'Template not found'], Response::HTTP_NOT_FOUND);
        }
        
        $variables = $data['variables'] ?? [];
        $subject = $this->replaceVariables($template->getSubject(), $variables);
        $htmlContent = $this->replaceVariables($template->getHtmlContent(), $variables);
        
        // Log email
        $log = new EmailLog();
        $log->setToEmail($data['to']);
        $log->setTemplateId($template->getId());
        $log->setSubject($subject);
        $log->setStatus('sent');
        $log->setSentAt(new \DateTime());
        
        $this->em->persist($log);
        $this->em->flush();
        
        return $this->json([
            'message' => 'Email sent successfully',
            'to' => $data['to'],
            'subject' => $subject
        ]);
    }

    private function replaceVariables(string $content, array $variables): string
    {
        foreach ($variables as $key => $value) {
            $content = str_replace("{{" . $key . "}}", $value, $content);
        }
        return $content;
    }
}
