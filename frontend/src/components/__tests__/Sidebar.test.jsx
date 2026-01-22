import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import Sidebar from '../Sidebar'

describe('Sidebar', () => {
  it('renders sidebar with navigation items', () => {
    render(<Sidebar isOpen={true} />)
    
    expect(screen.getByText('Home')).toBeInTheDocument()
    expect(screen.getByText('Trending')).toBeInTheDocument()
    expect(screen.getByText('Subscriptions')).toBeInTheDocument()
  })

  it('shows sidebar when isOpen is true', () => {
    const { container } = render(<Sidebar isOpen={true} />)
    const sidebar = container.querySelector('aside')
    
    expect(sidebar).toBeInTheDocument()
    expect(sidebar).not.toHaveClass('hidden')
  })

  it('renders Library section', () => {
    render(<Sidebar isOpen={true} />)
    
    expect(screen.getByText('Library')).toBeInTheDocument()
    expect(screen.getByText('History')).toBeInTheDocument()
  })

  it('renders navigation with icons', () => {
    render(<Sidebar isOpen={true} />)
    
    const homeLink = screen.getByText('Home')
    expect(homeLink).toBeInTheDocument()
    
    const trendingLink = screen.getByText('Trending')
    expect(trendingLink).toBeInTheDocument()
  })

  it('applies correct styling classes', () => {
    const { container } = render(<Sidebar isOpen={true} />)
    const sidebar = container.querySelector('aside')
    
    expect(sidebar).toHaveClass('bg-white', 'dark:bg-gray-800')
  })
})
