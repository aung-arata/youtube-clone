import React from 'react'

function VideoCard({ video }) {
  // Format views count
  const formatViews = (count) => {
    if (count >= 1000000) {
      return `${(count / 1000000).toFixed(1)}M views`
    } else if (count >= 1000) {
      return `${(count / 1000).toFixed(1)}K views`
    }
    return `${count} views`
  }

  // Format time ago
  const formatTimeAgo = (dateString) => {
    const date = new Date(dateString)
    const now = new Date()
    const seconds = Math.floor((now - date) / 1000)

    if (seconds < 60) return 'just now'
    if (seconds < 3600) return `${Math.floor(seconds / 60)} minutes ago`
    if (seconds < 86400) return `${Math.floor(seconds / 3600)} hours ago`
    if (seconds < 604800) return `${Math.floor(seconds / 86400)} days ago`
    if (seconds < 2592000) return `${Math.floor(seconds / 604800)} weeks ago`
    if (seconds < 31536000) return `${Math.floor(seconds / 2592000)} months ago`
    return `${Math.floor(seconds / 31536000)} years ago`
  }

  const handleClick = async () => {
    // Increment view count when video is clicked
    try {
      const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
      await fetch(`${apiUrl}/api/videos/${video.id}/views`, {
        method: 'POST',
      })
    } catch (err) {
      console.error('Error incrementing views:', err)
    }
  }

  // Use placeholder if no thumbnail
  const thumbnail = video.thumbnail || 'https://via.placeholder.com/320x180/666/FFFFFF?text=No+Thumbnail'
  const channelAvatar = video.channel_avatar || `https://via.placeholder.com/36/4299E1/FFFFFF?text=${video.channel_name?.charAt(0) || 'U'}`

  return (
    <div className="cursor-pointer" onClick={handleClick}>
      <div className="relative mb-3">
        <img
          src={thumbnail}
          alt={video.title}
          className="w-full aspect-video object-cover rounded-lg"
        />
        {video.duration && (
          <span className="absolute bottom-2 right-2 bg-black bg-opacity-80 text-white text-xs px-1 py-0.5 rounded">
            {video.duration}
          </span>
        )}
      </div>
      <div className="flex gap-3">
        <img
          src={channelAvatar}
          alt={video.channel_name}
          className="w-9 h-9 rounded-full"
        />
        <div className="flex-1">
          <h3 className="font-semibold text-sm line-clamp-2 mb-1">
            {video.title}
          </h3>
          <p className="text-sm text-gray-600">{video.channel_name}</p>
          <p className="text-sm text-gray-600">
            {formatViews(video.views)} • {formatTimeAgo(video.uploaded_at)}
          </p>
        </div>
      </div>
    </div>
  )
}

export default VideoCard
