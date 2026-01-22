import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import VideoCard from '../VideoCard'

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
    // Mock fetch
    global.fetch = vi.fn()
  })

  it('renders video card with correct information', () => {
    render(<VideoCard video={mockVideo} />)
    
    expect(screen.getByText('Test Video Title')).toBeInTheDocument()
    expect(screen.getByText('Test Channel')).toBeInTheDocument()
    expect(screen.getByText('10:30')).toBeInTheDocument()
    expect(screen.getByAltText('Test Video Title')).toHaveAttribute('src', mockVideo.thumbnail)
    expect(screen.getByAltText('Test Channel')).toHaveAttribute('src', mockVideo.channel_avatar)
  })

  it('formats view count correctly for millions', () => {
    render(<VideoCard video={mockVideo} />)
    expect(screen.getByText(/1.5M views/)).toBeInTheDocument()
  })

  it('formats view count correctly for thousands', () => {
    const videoWithThousands = { ...mockVideo, views: 5400 }
    render(<VideoCard video={videoWithThousands} />)
    expect(screen.getByText(/5.4K views/)).toBeInTheDocument()
  })

  it('formats view count correctly for small numbers', () => {
    const videoWithSmallViews = { ...mockVideo, views: 42 }
    render(<VideoCard video={videoWithSmallViews} />)
    expect(screen.getByText(/42 views/)).toBeInTheDocument()
  })

  it('formats time ago correctly', () => {
    render(<VideoCard video={mockVideo} />)
    expect(screen.getByText(/2 days ago/)).toBeInTheDocument()
  })

  it('displays placeholder thumbnail when thumbnail is missing', () => {
    const videoWithoutThumbnail = { ...mockVideo, thumbnail: '' }
    render(<VideoCard video={videoWithoutThumbnail} />)
    
    const img = screen.getByAltText('Test Video Title')
    expect(img).toHaveAttribute('src')
    expect(img.getAttribute('src')).toContain('placeholder')
  })

  it('displays placeholder channel avatar when avatar is missing', () => {
    const videoWithoutAvatar = { ...mockVideo, channel_avatar: '' }
    render(<VideoCard video={videoWithoutAvatar} />)
    
    const img = screen.getByAltText('Test Channel')
    expect(img).toHaveAttribute('src')
    expect(img.getAttribute('src')).toContain('placeholder')
  })

  it('calls API to increment views and add to history when clicked', async () => {
    global.fetch = vi.fn(() => Promise.resolve({ ok: true }))
    
    render(<VideoCard video={mockVideo} />)
    const card = screen.getByText('Test Video Title').closest('div').parentElement
    
    fireEvent.click(card)
    
    await waitFor(() => {
      expect(global.fetch).toHaveBeenCalledTimes(2)
    })
    
    // Check view increment call
    expect(global.fetch).toHaveBeenCalledWith(
      expect.stringContaining(`/api/videos/${mockVideo.id}/views`),
      expect.objectContaining({ method: 'POST' })
    )
    
    // Check history call
    expect(global.fetch).toHaveBeenCalledWith(
      expect.stringContaining('/api/users/1/history'),
      expect.objectContaining({
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ video_id: mockVideo.id })
      })
    )
  })

  it('handles API errors gracefully when tracking interactions', async () => {
    const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation(() => {})
    global.fetch = vi.fn(() => Promise.reject(new Error('API Error')))
    
    render(<VideoCard video={mockVideo} />)
    const card = screen.getByText('Test Video Title').closest('div').parentElement
    
    fireEvent.click(card)
    
    await waitFor(() => {
      expect(consoleErrorSpy).toHaveBeenCalled()
    })
    
    consoleErrorSpy.mockRestore()
  })

  it('does not display duration if not provided', () => {
    const videoWithoutDuration = { ...mockVideo, duration: '' }
    render(<VideoCard video={videoWithoutDuration} />)
    
    expect(screen.queryByText('10:30')).not.toBeInTheDocument()
  })
})
