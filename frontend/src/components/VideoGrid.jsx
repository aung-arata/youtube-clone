import React, { useState, useEffect, forwardRef, useImperativeHandle } from 'react'
import VideoCard from './VideoCard'

const VideoGrid = forwardRef((props, ref) => {
  const [videos, setVideos] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [categories, setCategories] = useState([])
  const [selectedCategory, setSelectedCategory] = useState('')

  useEffect(() => {
    fetchVideos()
    fetchCategories()
  }, [])

  const fetchCategories = async () => {
    try {
      const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
      const response = await fetch(`${apiUrl}/api/videos/categories`)
      if (response.ok) {
        const data = await response.json()
        setCategories(data)
      }
    } catch (err) {
      console.error('Error fetching categories:', err)
    }
  }

  const fetchVideos = async (searchQuery = '', category = '') => {
    try {
      setLoading(true)
      setError(null)
      const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
      
      let url = `${apiUrl}/api/videos?`
      const params = []
      if (searchQuery) params.push(`q=${encodeURIComponent(searchQuery)}`)
      if (category) params.push(`category=${encodeURIComponent(category)}`)
      url += params.join('&')
      
      const response = await fetch(url)
      if (!response.ok) {
        throw new Error('Failed to fetch videos')
      }
      const data = await response.json()
      setVideos(data)
    } catch (err) {
      setError(err.message)
      console.error('Error fetching videos:', err)
    } finally {
      setLoading(false)
    }
  }

  const handleCategoryChange = (category) => {
    setSelectedCategory(category)
    fetchVideos('', category)
  }

  // Expose fetchVideos to parent via ref
  useImperativeHandle(ref, () => ({
    fetchVideos
  }))

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-lg text-gray-600 dark:text-gray-400">Loading videos...</div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-lg text-red-600 dark:text-red-400">Error: {error}</div>
      </div>
    )
  }

  if (videos.length === 0) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-lg text-gray-600 dark:text-gray-400">No videos found</div>
      </div>
    )
  }

  return (
    <>
      {categories.length > 0 && (
        <div className="mb-6 flex gap-2 overflow-x-auto pb-2">
          <button
            onClick={() => handleCategoryChange('')}
            className={`px-4 py-2 rounded-full text-sm font-medium whitespace-nowrap ${
              selectedCategory === ''
                ? 'bg-gray-900 dark:bg-gray-100 text-white dark:text-gray-900'
                : 'bg-gray-200 dark:bg-gray-700 text-gray-900 dark:text-gray-100 hover:bg-gray-300 dark:hover:bg-gray-600'
            }`}
          >
            All
          </button>
          {categories.map((category) => (
            <button
              key={category}
              onClick={() => handleCategoryChange(category)}
              className={`px-4 py-2 rounded-full text-sm font-medium whitespace-nowrap ${
                selectedCategory === category
                  ? 'bg-gray-900 dark:bg-gray-100 text-white dark:text-gray-900'
                  : 'bg-gray-200 dark:bg-gray-700 text-gray-900 dark:text-gray-100 hover:bg-gray-300 dark:hover:bg-gray-600'
              }`}
            >
              {category}
            </button>
          ))}
        </div>
      )}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
        {videos.map((video) => (
          <VideoCard key={video.id} video={video} />
        ))}
      </div>
    </>
  )
})

VideoGrid.displayName = 'VideoGrid'

export default VideoGrid
