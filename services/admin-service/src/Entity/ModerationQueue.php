<?php

namespace App\Entity;

use App\Repository\ModerationQueueRepository;
use Doctrine\DBAL\Types\Types;
use Doctrine\ORM\Mapping as ORM;

#[ORM\Entity(repositoryClass: ModerationQueueRepository::class)]
#[ORM\Table(name: 'moderation_queue')]
#[ORM\HasLifecycleCallbacks]
class ModerationQueue
{
    #[ORM\Id]
    #[ORM\GeneratedValue]
    #[ORM\Column]
    private ?int $id = null;

    #[ORM\Column(length: 50)]
    private ?string $contentType = null;

    #[ORM\Column]
    private ?int $contentId = null;

    #[ORM\Column(nullable: true)]
    private ?int $reporterId = null;

    #[ORM\Column(type: Types::TEXT, nullable: true)]
    private ?string $reason = null;

    #[ORM\Column(length: 20)]
    private ?string $status = 'pending';

    #[ORM\Column(nullable: true)]
    private ?int $reviewedBy = null;

    #[ORM\Column(type: Types::DATETIME_MUTABLE, nullable: true)]
    private ?\DateTimeInterface $reviewedAt = null;

    #[ORM\Column(type: Types::DATETIME_MUTABLE)]
    private ?\DateTimeInterface $createdAt = null;

    #[ORM\Column(type: Types::DATETIME_MUTABLE)]
    private ?\DateTimeInterface $updatedAt = null;

    public function __construct()
    {
        $this->createdAt = new \DateTime();
        $this->updatedAt = new \DateTime();
    }

    #[ORM\PreUpdate]
    public function setUpdatedAtValue(): void
    {
        $this->updatedAt = new \DateTime();
    }

    public function getId(): ?int
    {
        return $this->id;
    }

    public function getContentType(): ?string
    {
        return $this->contentType;
    }

    public function setContentType(string $contentType): self
    {
        $this->contentType = $contentType;
        return $this;
    }

    public function getContentId(): ?int
    {
        return $this->contentId;
    }

    public function setContentId(int $contentId): self
    {
        $this->contentId = $contentId;
        return $this;
    }

    public function getReporterId(): ?int
    {
        return $this->reporterId;
    }

    public function setReporterId(?int $reporterId): self
    {
        $this->reporterId = $reporterId;
        return $this;
    }

    public function getReason(): ?string
    {
        return $this->reason;
    }

    public function setReason(?string $reason): self
    {
        $this->reason = $reason;
        return $this;
    }

    public function getStatus(): ?string
    {
        return $this->status;
    }

    public function setStatus(string $status): self
    {
        $this->status = $status;
        return $this;
    }

    public function getReviewedBy(): ?int
    {
        return $this->reviewedBy;
    }

    public function setReviewedBy(?int $reviewedBy): self
    {
        $this->reviewedBy = $reviewedBy;
        return $this;
    }

    public function getReviewedAt(): ?\DateTimeInterface
    {
        return $this->reviewedAt;
    }

    public function setReviewedAt(?\DateTimeInterface $reviewedAt): self
    {
        $this->reviewedAt = $reviewedAt;
        return $this;
    }

    public function getCreatedAt(): ?\DateTimeInterface
    {
        return $this->createdAt;
    }

    public function getUpdatedAt(): ?\DateTimeInterface
    {
        return $this->updatedAt;
    }

    public function toArray(): array
    {
        return [
            'id' => $this->id,
            'content_type' => $this->contentType,
            'content_id' => $this->contentId,
            'reporter_id' => $this->reporterId,
            'reason' => $this->reason,
            'status' => $this->status,
            'reviewed_by' => $this->reviewedBy,
            'reviewed_at' => $this->reviewedAt?->format('Y-m-d H:i:s'),
            'created_at' => $this->createdAt?->format('Y-m-d H:i:s'),
            'updated_at' => $this->updatedAt?->format('Y-m-d H:i:s'),
        ];
    }
}
