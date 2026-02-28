import { describe, it, expect, vi } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import Sidebar from '../Sidebar'

const mockNavigate = vi.fn()
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom')
  return { ...actual, useNavigate: () => mockNavigate }
})

describe('Sidebar', () => {
  beforeEach(() => {
    mockNavigate.mockClear()
  })

  it('renders sidebar with navigation items', () => {
    render(<MemoryRouter><Sidebar isOpen={true} /></MemoryRouter>)
    
    expect(screen.getByText('Home')).toBeInTheDocument()
    expect(screen.getByText('Trending')).toBeInTheDocument()
    expect(screen.getByText('Subscriptions')).toBeInTheDocument()
  })

  it('shows sidebar when isOpen is true', () => {
    const { container } = render(<MemoryRouter><Sidebar isOpen={true} /></MemoryRouter>)
    const sidebar = container.querySelector('aside')
    
    expect(sidebar).toBeInTheDocument()
    expect(sidebar).not.toHaveClass('hidden')
  })

  it('renders Library section', () => {
    render(<MemoryRouter><Sidebar isOpen={true} /></MemoryRouter>)
    
    expect(screen.getByText('Library')).toBeInTheDocument()
    expect(screen.getByText('History')).toBeInTheDocument()
  })

  it('navigates when menu items are clicked', () => {
    render(<MemoryRouter><Sidebar isOpen={true} /></MemoryRouter>)
    
    fireEvent.click(screen.getByText('Trending'))
    expect(mockNavigate).toHaveBeenCalledWith('/trending')

    fireEvent.click(screen.getByText('Home'))
    expect(mockNavigate).toHaveBeenCalledWith('/')

    fireEvent.click(screen.getByText('History'))
    expect(mockNavigate).toHaveBeenCalledWith('/history')
  })

  it('applies correct styling classes', () => {
    const { container } = render(<MemoryRouter><Sidebar isOpen={true} /></MemoryRouter>)
    const sidebar = container.querySelector('aside')
    
    expect(sidebar).toHaveClass('bg-white', 'dark:bg-gray-800')
  })
})
