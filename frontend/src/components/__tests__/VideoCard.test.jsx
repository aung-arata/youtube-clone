import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import VideoCard from '../VideoCard'

const mockNavigate = vi.fn()
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom')
  return { ...actual, useNavigate: () => mockNavigate }
})

describe('VideoCard', () => {
  const mockVideo = {
    id: 1,
    title: 'Test Video Title',
    description: 'Test video description',
    thumbnail: 'https://example.com/thumbnail.jpg',
    channel_name: 'Test Channel',
    channel_avatar: 'https://example.com/avatar.jpg',
    views: 1500000,
    duration: '10:30',
    uploaded_at: new Date(Date.now() - 86400000 * 2).toISOString(), // 2 days ago
  }

  beforeEach(() => {
    mockNavigate.mockClear()
  })

  it('renders video card with correct information', () => {
    render(<MemoryRouter><VideoCard video={mockVideo} /></MemoryRouter>)
    
    expect(screen.getByText('Test Video Title')).toBeInTheDocument()
    expect(screen.getByText('Test Channel')).toBeInTheDocument()
    expect(screen.getByText('10:30')).toBeInTheDocument()
    expect(screen.getByAltText('Test Video Title')).toHaveAttribute('src', mockVideo.thumbnail)
    expect(screen.getByAltText('Test Channel')).toHaveAttribute('src', mockVideo.channel_avatar)
  })

  it('formats view count correctly for millions', () => {
    render(<MemoryRouter><VideoCard video={mockVideo} /></MemoryRouter>)
    expect(screen.getByText(/1.5M views/)).toBeInTheDocument()
  })

  it('formats view count correctly for thousands', () => {
    const videoWithThousands = { ...mockVideo, views: 5400 }
    render(<MemoryRouter><VideoCard video={videoWithThousands} /></MemoryRouter>)
    expect(screen.getByText(/5.4K views/)).toBeInTheDocument()
  })

  it('formats view count correctly for small numbers', () => {
    const videoWithSmallViews = { ...mockVideo, views: 42 }
    render(<MemoryRouter><VideoCard video={videoWithSmallViews} /></MemoryRouter>)
    expect(screen.getByText(/42 views/)).toBeInTheDocument()
  })

  it('formats time ago correctly', () => {
    render(<MemoryRouter><VideoCard video={mockVideo} /></MemoryRouter>)
    expect(screen.getByText(/2 days ago/)).toBeInTheDocument()
  })

  it('displays placeholder thumbnail when thumbnail is missing', () => {
    const videoWithoutThumbnail = { ...mockVideo, thumbnail: '' }
    render(<MemoryRouter><VideoCard video={videoWithoutThumbnail} /></MemoryRouter>)
    
    const img = screen.getByAltText('Test Video Title')
    expect(img).toHaveAttribute('src')
    expect(img.getAttribute('src')).toContain('placeholder')
  })

  it('displays placeholder channel avatar when avatar is missing', () => {
    const videoWithoutAvatar = { ...mockVideo, channel_avatar: '' }
    render(<MemoryRouter><VideoCard video={videoWithoutAvatar} /></MemoryRouter>)
    
    const img = screen.getByAltText('Test Channel')
    expect(img).toHaveAttribute('src')
    expect(img.getAttribute('src')).toContain('placeholder')
  })

  it('navigates to video page when clicked', () => {
    render(<MemoryRouter><VideoCard video={mockVideo} /></MemoryRouter>)
    const card = screen.getByText('Test Video Title').closest('div').parentElement
    
    fireEvent.click(card)
    
    expect(mockNavigate).toHaveBeenCalledWith('/video/1')
  })

  it('does not display duration if not provided', () => {
    const videoWithoutDuration = { ...mockVideo, duration: '' }
    render(<MemoryRouter><VideoCard video={videoWithoutDuration} /></MemoryRouter>)
    
    expect(screen.queryByText('10:30')).not.toBeInTheDocument()
  })

  it('renders compact mode correctly', () => {
    render(<MemoryRouter><VideoCard video={mockVideo} compact /></MemoryRouter>)
    
    expect(screen.getByText('Test Video Title')).toBeInTheDocument()
    expect(screen.getByText('Test Channel')).toBeInTheDocument()
  })
})
