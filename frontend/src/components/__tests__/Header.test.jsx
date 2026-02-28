import { describe, it, expect, vi } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import Header from '../Header'

const mockNavigate = vi.fn()
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom')
  return { ...actual, useNavigate: () => mockNavigate }
})

vi.mock('../../contexts/AuthContext', () => ({
  useAuth: () => ({ user: null, logout: vi.fn() }),
}))

describe('Header', () => {
  it('renders header with YouTube Clone branding', () => {
    render(<MemoryRouter><Header /></MemoryRouter>)
    expect(screen.getByText('YouTube Clone')).toBeInTheDocument()
  })

  it('renders search input with correct placeholder', () => {
    render(<MemoryRouter><Header /></MemoryRouter>)
    const searchInput = screen.getByPlaceholderText('Search')
    expect(searchInput).toBeInTheDocument()
  })

  it('calls onMenuClick when menu button is clicked', () => {
    const handleMenuClick = vi.fn()
    render(<MemoryRouter><Header onMenuClick={handleMenuClick} /></MemoryRouter>)
    
    const menuButton = screen.getAllByRole('button')[0]
    fireEvent.click(menuButton)
    
    expect(handleMenuClick).toHaveBeenCalledTimes(1)
  })

  it('updates search query when typing', () => {
    render(<MemoryRouter><Header /></MemoryRouter>)
    const searchInput = screen.getByPlaceholderText('Search')
    
    fireEvent.change(searchInput, { target: { value: 'test query' } })
    
    expect(searchInput.value).toBe('test query')
  })

  it('calls onSearch with query when form is submitted', () => {
    const handleSearch = vi.fn()
    render(<MemoryRouter><Header onSearch={handleSearch} /></MemoryRouter>)
    
    const searchInput = screen.getByPlaceholderText('Search')
    const form = searchInput.closest('form')
    
    fireEvent.change(searchInput, { target: { value: 'react tutorial' } })
    fireEvent.submit(form)
    
    expect(handleSearch).toHaveBeenCalledWith('react tutorial')
  })

  it('calls onSearch when search button is clicked', () => {
    const handleSearch = vi.fn()
    render(<MemoryRouter><Header onSearch={handleSearch} /></MemoryRouter>)
    
    const searchInput = screen.getByPlaceholderText('Search')
    fireEvent.change(searchInput, { target: { value: 'golang' } })
    
    const searchButton = screen.getAllByRole('button')[1]
    fireEvent.click(searchButton)
    
    expect(handleSearch).toHaveBeenCalledWith('golang')
  })

  it('shows Sign In button when user is not logged in', () => {
    render(<MemoryRouter><Header /></MemoryRouter>)
    expect(screen.getByText('Sign In')).toBeInTheDocument()
  })
})
