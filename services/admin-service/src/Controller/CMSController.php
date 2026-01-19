<?php

namespace App\Controller;

use App\Service\DatabaseService;

/**
 * CMS Controller
 * Handles blog posts, documentation, and help center articles
 */
class CMSController
{
    private DatabaseService $db;
    
    public function __construct()
    {
        $this->db = new DatabaseService();
    }
    
    /**
     * List all blog posts
     */
    public function listBlogs(): void
    {
        header('Content-Type: application/json');
        
        $page = $_GET['page'] ?? 1;
        $limit = min($_GET['limit'] ?? 20, 100);
        $offset = ($page - 1) * $limit;
        
        $blogs = $this->db->query(
            'SELECT * FROM blog_posts ORDER BY published_at DESC LIMIT ? OFFSET ?',
            [$limit, $offset]
        );
        
        echo json_encode(['blogs' => $blogs, 'page' => $page, 'limit' => $limit]);
    }
    
    /**
     * Create a new blog post
     */
    public function createBlog(): void
    {
        header('Content-Type: application/json');
        
        $data = json_decode(file_get_contents('php://input'), true);
        
        $required = ['title', 'content', 'author_id'];
        foreach ($required as $field) {
            if (empty($data[$field])) {
                http_response_code(400);
                echo json_encode(['error' => "Missing required field: {$field}"]);
                return;
            }
        }
        
        $id = $this->db->insert('blog_posts', [
            'title' => $data['title'],
            'slug' => $this->slugify($data['title']),
            'content' => $data['content'],
            'excerpt' => substr($data['content'], 0, 200),
            'author_id' => $data['author_id'],
            'category' => $data['category'] ?? 'general',
            'status' => $data['status'] ?? 'draft',
            'published_at' => $data['status'] === 'published' ? date('Y-m-d H:i:s') : null
        ]);
        
        http_response_code(201);
        echo json_encode(['id' => $id, 'message' => 'Blog post created']);
    }
    
    /**
     * Get a specific blog post
     */
    public function getBlog(int $id): void
    {
        header('Content-Type: application/json');
        
        $blog = $this->db->query('SELECT * FROM blog_posts WHERE id = ?', [$id]);
        
        if (empty($blog)) {
            http_response_code(404);
            echo json_encode(['error' => 'Blog post not found']);
            return;
        }
        
        echo json_encode($blog[0]);
    }
    
    /**
     * Update a blog post
     */
    public function updateBlog(int $id): void
    {
        header('Content-Type: application/json');
        
        $data = json_decode(file_get_contents('php://input'), true);
        
        $updates = [];
        if (isset($data['title'])) {
            $updates['title'] = $data['title'];
            $updates['slug'] = $this->slugify($data['title']);
        }
        if (isset($data['content'])) $updates['content'] = $data['content'];
        if (isset($data['category'])) $updates['category'] = $data['category'];
        if (isset($data['status'])) {
            $updates['status'] = $data['status'];
            if ($data['status'] === 'published') {
                $updates['published_at'] = date('Y-m-d H:i:s');
            }
        }
        
        if (empty($updates)) {
            http_response_code(400);
            echo json_encode(['error' => 'No fields to update']);
            return;
        }
        
        $updates['updated_at'] = date('Y-m-d H:i:s');
        
        $this->db->update('blog_posts', $updates, ['id' => $id]);
        
        echo json_encode(['message' => 'Blog post updated']);
    }
    
    /**
     * Delete a blog post
     */
    public function deleteBlog(int $id): void
    {
        header('Content-Type: application/json');
        
        $this->db->delete('blog_posts', ['id' => $id]);
        
        echo json_encode(['message' => 'Blog post deleted']);
    }
    
    /**
     * List documentation articles
     */
    public function listDocs(): void
    {
        header('Content-Type: application/json');
        
        $docs = $this->db->query(
            'SELECT * FROM documentation ORDER BY category, sort_order'
        );
        
        echo json_encode(['docs' => $docs]);
    }
    
    /**
     * Create documentation
     */
    public function createDoc(): void
    {
        header('Content-Type: application/json');
        
        $data = json_decode(file_get_contents('php://input'), true);
        
        $id = $this->db->insert('documentation', [
            'title' => $data['title'],
            'slug' => $this->slugify($data['title']),
            'content' => $data['content'],
            'category' => $data['category'] ?? 'general',
            'sort_order' => $data['sort_order'] ?? 0
        ]);
        
        http_response_code(201);
        echo json_encode(['id' => $id, 'message' => 'Documentation created']);
    }
    
    /**
     * List help articles
     */
    public function listHelpArticles(): void
    {
        header('Content-Type: application/json');
        
        $articles = $this->db->query(
            'SELECT * FROM help_articles ORDER BY category, title'
        );
        
        echo json_encode(['articles' => $articles]);
    }
    
    /**
     * Create help article
     */
    public function createHelpArticle(): void
    {
        header('Content-Type: application/json');
        
        $data = json_decode(file_get_contents('php://input'), true);
        
        $id = $this->db->insert('help_articles', [
            'title' => $data['title'],
            'slug' => $this->slugify($data['title']),
            'content' => $data['content'],
            'category' => $data['category'] ?? 'general'
        ]);
        
        http_response_code(201);
        echo json_encode(['id' => $id, 'message' => 'Help article created']);
    }
    
    /**
     * Create URL-friendly slug
     */
    private function slugify(string $text): string
    {
        $text = preg_replace('~[^\pL\d]+~u', '-', $text);
        $text = iconv('utf-8', 'us-ascii//TRANSLIT', $text);
        $text = preg_replace('~[^-\w]+~', '', $text);
        $text = trim($text, '-');
        $text = preg_replace('~-+~', '-', $text);
        $text = strtolower($text);
        
        return empty($text) ? 'n-a' : $text;
    }
}
