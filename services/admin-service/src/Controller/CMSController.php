<?php

namespace App\Controller;

use App\Entity\BlogPost;
use App\Entity\Documentation;
use App\Entity\HelpArticle;
use App\Repository\BlogPostRepository;
use App\Repository\DocumentationRepository;
use App\Repository\HelpArticleRepository;
use Doctrine\ORM\EntityManagerInterface;
use Symfony\Bundle\FrameworkBundle\Controller\AbstractController;
use Symfony\Component\HttpFoundation\JsonResponse;
use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\HttpFoundation\Response;
use Symfony\Component\Routing\Annotation\Route;
use Symfony\Component\String\Slugger\SluggerInterface;

#[Route('/cms', name: 'cms_')]
class CMSController extends AbstractController
{
    public function __construct(
        private EntityManagerInterface $em,
        private SluggerInterface $slugger,
        private BlogPostRepository $blogRepo,
        private DocumentationRepository $docRepo,
        private HelpArticleRepository $helpRepo
    ) {
    }

    #[Route('/blog', name: 'blog_list', methods: ['GET'])]
    public function listBlogs(Request $request): JsonResponse
    {
        $page = max(1, (int) $request->query->get('page', 1));
        $limit = min((int) $request->query->get('limit', 20), 100);
        
        $blogs = $this->blogRepo->findBy([], ['publishedAt' => 'DESC'], $limit, ($page - 1) * $limit);
        
        return $this->json([
            'blogs' => array_map(fn($blog) => $blog->toArray(), $blogs),
            'page' => $page,
            'limit' => $limit
        ]);
    }

    #[Route('/blog', name: 'blog_create', methods: ['POST'])]
    public function createBlog(Request $request): JsonResponse
    {
        $data = json_decode($request->getContent(), true);
        
        if (empty($data['title']) || empty($data['content']) || empty($data['author_id'])) {
            return $this->json(['error' => 'Missing required fields'], Response::HTTP_BAD_REQUEST);
        }
        
        $blog = new BlogPost();
        $blog->setTitle($data['title']);
        $blog->setSlug($this->slugger->slug($data['title'])->lower()->toString());
        $blog->setContent($data['content']);
        $blog->setExcerpt(substr($data['content'], 0, 200));
        $blog->setAuthorId($data['author_id']);
        $blog->setCategory($data['category'] ?? 'general');
        $blog->setStatus($data['status'] ?? 'draft');
        
        if (($data['status'] ?? 'draft') === 'published') {
            $blog->setPublishedAt(new \DateTime());
        }
        
        $this->em->persist($blog);
        $this->em->flush();
        
        return $this->json([
            'id' => $blog->getId(),
            'message' => 'Blog post created'
        ], Response::HTTP_CREATED);
    }

    #[Route('/blog/{id}', name: 'blog_get', methods: ['GET'])]
    public function getBlog(int $id): JsonResponse
    {
        $blog = $this->blogRepo->find($id);
        
        if (!$blog) {
            return $this->json(['error' => 'Blog post not found'], Response::HTTP_NOT_FOUND);
        }
        
        return $this->json($blog->toArray());
    }

    #[Route('/blog/{id}', name: 'blog_update', methods: ['PUT'])]
    public function updateBlog(int $id, Request $request): JsonResponse
    {
        $blog = $this->blogRepo->find($id);
        
        if (!$blog) {
            return $this->json(['error' => 'Blog post not found'], Response::HTTP_NOT_FOUND);
        }
        
        $data = json_decode($request->getContent(), true);
        
        if (isset($data['title'])) {
            $blog->setTitle($data['title']);
            $blog->setSlug($this->slugger->slug($data['title'])->lower()->toString());
        }
        if (isset($data['content'])) {
            $blog->setContent($data['content']);
        }
        if (isset($data['category'])) {
            $blog->setCategory($data['category']);
        }
        if (isset($data['status'])) {
            $blog->setStatus($data['status']);
            if ($data['status'] === 'published' && !$blog->getPublishedAt()) {
                $blog->setPublishedAt(new \DateTime());
            }
        }
        
        $this->em->flush();
        
        return $this->json(['message' => 'Blog post updated']);
    }

    #[Route('/blog/{id}', name: 'blog_delete', methods: ['DELETE'])]
    public function deleteBlog(int $id): JsonResponse
    {
        $blog = $this->blogRepo->find($id);
        
        if (!$blog) {
            return $this->json(['error' => 'Blog post not found'], Response::HTTP_NOT_FOUND);
        }
        
        $this->em->remove($blog);
        $this->em->flush();
        
        return $this->json(['message' => 'Blog post deleted']);
    }

    #[Route('/docs', name: 'docs_list', methods: ['GET'])]
    public function listDocs(): JsonResponse
    {
        $docs = $this->docRepo->findBy([], ['category' => 'ASC', 'sortOrder' => 'ASC']);
        
        return $this->json([
            'docs' => array_map(fn($doc) => $doc->toArray(), $docs)
        ]);
    }

    #[Route('/docs', name: 'docs_create', methods: ['POST'])]
    public function createDoc(Request $request): JsonResponse
    {
        $data = json_decode($request->getContent(), true);
        
        if (empty($data['title']) || empty($data['content'])) {
            return $this->json(['error' => 'Missing required fields'], Response::HTTP_BAD_REQUEST);
        }
        
        $doc = new Documentation();
        $doc->setTitle($data['title']);
        $doc->setSlug($this->slugger->slug($data['title'])->lower()->toString());
        $doc->setContent($data['content']);
        $doc->setCategory($data['category'] ?? 'general');
        $doc->setSortOrder($data['sort_order'] ?? 0);
        
        $this->em->persist($doc);
        $this->em->flush();
        
        return $this->json([
            'id' => $doc->getId(),
            'message' => 'Documentation created'
        ], Response::HTTP_CREATED);
    }

    #[Route('/help', name: 'help_list', methods: ['GET'])]
    public function listHelpArticles(): JsonResponse
    {
        $articles = $this->helpRepo->findBy([], ['category' => 'ASC', 'title' => 'ASC']);
        
        return $this->json([
            'articles' => array_map(fn($article) => $article->toArray(), $articles)
        ]);
    }

    #[Route('/help', name: 'help_create', methods: ['POST'])]
    public function createHelpArticle(Request $request): JsonResponse
    {
        $data = json_decode($request->getContent(), true);
        
        if (empty($data['title']) || empty($data['content'])) {
            return $this->json(['error' => 'Missing required fields'], Response::HTTP_BAD_REQUEST);
        }
        
        $article = new HelpArticle();
        $article->setTitle($data['title']);
        $article->setSlug($this->slugger->slug($data['title'])->lower()->toString());
        $article->setContent($data['content']);
        $article->setCategory($data['category'] ?? 'general');
        
        $this->em->persist($article);
        $this->em->flush();
        
        return $this->json([
            'id' => $article->getId(),
            'message' => 'Help article created'
        ], Response::HTTP_CREATED);
    }
}
