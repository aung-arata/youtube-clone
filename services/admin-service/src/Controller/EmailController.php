<?php

namespace App\Controller;

use App\Service\DatabaseService;

/**
 * Email Template Controller
 * Manages email templates and sending
 */
class EmailController
{
    private DatabaseService $db;
    
    public function __construct()
    {
        $this->db = new DatabaseService();
    }
    
    /**
     * List all email templates
     */
    public function listTemplates(): void
    {
        header('Content-Type: application/json');
        
        $templates = $this->db->query(
            'SELECT id, name, subject, category, created_at, updated_at FROM email_templates ORDER BY category, name'
        );
        
        echo json_encode(['templates' => $templates]);
    }
    
    /**
     * Create email template
     */
    public function createTemplate(): void
    {
        header('Content-Type: application/json');
        
        $data = json_decode(file_get_contents('php://input'), true);
        
        $required = ['name', 'subject', 'html_content'];
        foreach ($required as $field) {
            if (empty($data[$field])) {
                http_response_code(400);
                echo json_encode(['error' => "Missing required field: {$field}"]);
                return;
            }
        }
        
        $id = $this->db->insert('email_templates', [
            'name' => $data['name'],
            'subject' => $data['subject'],
            'html_content' => $data['html_content'],
            'text_content' => $data['text_content'] ?? strip_tags($data['html_content']),
            'category' => $data['category'] ?? 'general'
        ]);
        
        http_response_code(201);
        echo json_encode(['id' => $id, 'message' => 'Email template created']);
    }
    
    /**
     * Get email template
     */
    public function getTemplate(int $id): void
    {
        header('Content-Type: application/json');
        
        $template = $this->db->query('SELECT * FROM email_templates WHERE id = ?', [$id]);
        
        if (empty($template)) {
            http_response_code(404);
            echo json_encode(['error' => 'Template not found']);
            return;
        }
        
        echo json_encode($template[0]);
    }
    
    /**
     * Update email template
     */
    public function updateTemplate(int $id): void
    {
        header('Content-Type: application/json');
        
        $data = json_decode(file_get_contents('php://input'), true);
        
        $updates = [];
        if (isset($data['name'])) $updates['name'] = $data['name'];
        if (isset($data['subject'])) $updates['subject'] = $data['subject'];
        if (isset($data['html_content'])) $updates['html_content'] = $data['html_content'];
        if (isset($data['text_content'])) $updates['text_content'] = $data['text_content'];
        if (isset($data['category'])) $updates['category'] = $data['category'];
        
        if (empty($updates)) {
            http_response_code(400);
            echo json_encode(['error' => 'No fields to update']);
            return;
        }
        
        $updates['updated_at'] = date('Y-m-d H:i:s');
        
        $this->db->update('email_templates', $updates, ['id' => $id]);
        
        echo json_encode(['message' => 'Template updated']);
    }
    
    /**
     * Delete email template
     */
    public function deleteTemplate(int $id): void
    {
        header('Content-Type: application/json');
        
        $this->db->delete('email_templates', ['id' => $id]);
        
        echo json_encode(['message' => 'Template deleted']);
    }
    
    /**
     * Send email
     */
    public function sendEmail(): void
    {
        header('Content-Type: application/json');
        
        $data = json_decode(file_get_contents('php://input'), true);
        
        $required = ['to', 'template_id'];
        foreach ($required as $field) {
            if (empty($data[$field])) {
                http_response_code(400);
                echo json_encode(['error' => "Missing required field: {$field}"]);
                return;
            }
        }
        
        // Get template
        $template = $this->db->query('SELECT * FROM email_templates WHERE id = ?', [$data['template_id']]);
        
        if (empty($template)) {
            http_response_code(404);
            echo json_encode(['error' => 'Template not found']);
            return;
        }
        
        $template = $template[0];
        
        // Replace variables in template
        $variables = $data['variables'] ?? [];
        $subject = $this->replaceVariables($template['subject'], $variables);
        $htmlContent = $this->replaceVariables($template['html_content'], $variables);
        $textContent = $this->replaceVariables($template['text_content'], $variables);
        
        // Log email (in production, this would actually send)
        $this->db->insert('email_log', [
            'to_email' => $data['to'],
            'template_id' => $data['template_id'],
            'subject' => $subject,
            'status' => 'sent',
            'sent_at' => date('Y-m-d H:i:s')
        ]);
        
        echo json_encode([
            'message' => 'Email sent successfully',
            'to' => $data['to'],
            'subject' => $subject
        ]);
    }
    
    /**
     * Replace variables in template
     */
    private function replaceVariables(string $content, array $variables): string
    {
        foreach ($variables as $key => $value) {
            $content = str_replace("{{" . $key . "}}", $value, $content);
        }
        return $content;
    }
}
