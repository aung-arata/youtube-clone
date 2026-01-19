<?php

namespace App\Repository;

use App\Entity\ModerationQueue;
use Doctrine\Bundle\DoctrineBundle\Repository\ServiceEntityRepository;
use Doctrine\Persistence\ManagerRegistry;

/**
 * @extends ServiceEntityRepository<ModerationQueue>
 */
class ModerationQueueRepository extends ServiceEntityRepository
{
    public function __construct(ManagerRegistry $registry)
    {
        parent::__construct($registry, ModerationQueue::class);
    }

    /**
     * Count moderation items by status
     */
    public function countByStatus(string $status): int
    {
        return $this->count(['status' => $status]);
    }

    /**
     * Find moderation items by status with limit
     */
    public function findByStatus(string $status, int $limit = 50): array
    {
        return $this->createQueryBuilder('m')
            ->where('m.status = :status')
            ->setParameter('status', $status)
            ->orderBy('m.createdAt', 'DESC')
            ->setMaxResults($limit)
            ->getQuery()
            ->getResult();
    }
}
