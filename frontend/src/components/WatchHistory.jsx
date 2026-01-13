import React, { useState, useEffect } from 'react'
import VideoCard from './VideoCard'

function WatchHistory({ userId }) {
  const [history, setHistory] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  useEffect(() => {
    if (userId) {
      fetchHistory()
    }
  }, [userId])

  const fetchHistory = async () => {
    try {
      setLoading(true)
      setError(null)
      const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
      const response = await fetch(`${apiUrl}/api/users/${userId}/history`)
      if (!response.ok) {
        throw new Error('Failed to fetch watch history')
      }
      const data = await response.json()
      setHistory(data)
    } catch (err) {
      setError(err.message)
      console.error('Error fetching history:', err)
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-lg text-gray-600 dark:text-gray-400">Loading watch history...</div>
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

  if (history.length === 0) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-lg text-gray-600 dark:text-gray-400">No watch history yet</div>
      </div>
    )
  }

  return (
    <div>
      <h2 className="text-2xl font-bold mb-6 dark:text-white">Watch History</h2>
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
        {history.map((item) => (
          <div key={item.id}>
            <VideoCard video={item} />
            <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
              Watched {new Date(item.watched_at).toLocaleDateString()}
            </p>
          </div>
        ))}
      </div>
    </div>
  )
}

export default WatchHistory
