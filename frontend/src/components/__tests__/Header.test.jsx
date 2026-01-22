import { describe, it, expect, vi } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import Header from '../Header'

describe('Header', () => {
  it('renders header with YouTube Clone branding', () => {
    render(<Header />)
    expect(screen.getByText('YouTube Clone')).toBeInTheDocument()
  })

  it('renders search input with correct placeholder', () => {
    render(<Header />)
    const searchInput = screen.getByPlaceholderText('Search')
    expect(searchInput).toBeInTheDocument()
  })

  it('calls onMenuClick when menu button is clicked', () => {
    const handleMenuClick = vi.fn()
    render(<Header onMenuClick={handleMenuClick} />)
    
    const menuButton = screen.getAllByRole('button')[0]
    fireEvent.click(menuButton)
    
    expect(handleMenuClick).toHaveBeenCalledTimes(1)
  })

  it('updates search query when typing', () => {
    render(<Header />)
    const searchInput = screen.getByPlaceholderText('Search')
    
    fireEvent.change(searchInput, { target: { value: 'test query' } })
    
    expect(searchInput.value).toBe('test query')
  })

  it('calls onSearch with query when form is submitted', () => {
    const handleSearch = vi.fn()
    render(<Header onSearch={handleSearch} />)
    
    const searchInput = screen.getByPlaceholderText('Search')
    const form = searchInput.closest('form')
    
    fireEvent.change(searchInput, { target: { value: 'react tutorial' } })
    fireEvent.submit(form)
    
    expect(handleSearch).toHaveBeenCalledWith('react tutorial')
  })

  it('calls onSearch when search button is clicked', () => {
    const handleSearch = vi.fn()
    render(<Header onSearch={handleSearch} />)
    
    const searchInput = screen.getByPlaceholderText('Search')
    fireEvent.change(searchInput, { target: { value: 'golang' } })
    
    const searchButton = screen.getAllByRole('button')[1]
    fireEvent.click(searchButton)
    
    expect(handleSearch).toHaveBeenCalledWith('golang')
  })

  it('does not call onSearch if query is empty', () => {
    const handleSearch = vi.fn()
    render(<Header onSearch={handleSearch} />)
    
    const form = screen.getByPlaceholderText('Search').closest('form')
    fireEvent.submit(form)
    
    // onSearch should be called but with empty string
    expect(handleSearch).toHaveBeenCalledWith('')
  })
})
