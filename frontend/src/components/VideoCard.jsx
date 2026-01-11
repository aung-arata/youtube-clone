import React from 'react'

function VideoCard({ video }) {
  return (
    <div className="cursor-pointer">
      <div className="relative mb-3">
        <img
          src={video.thumbnail}
          alt={video.title}
          className="w-full aspect-video object-cover rounded-lg"
        />
        <span className="absolute bottom-2 right-2 bg-black bg-opacity-80 text-white text-xs px-1 py-0.5 rounded">
          {video.duration}
        </span>
      </div>
      <div className="flex gap-3">
        <img
          src={video.channelAvatar}
          alt={video.channelName}
          className="w-9 h-9 rounded-full"
        />
        <div className="flex-1">
          <h3 className="font-semibold text-sm line-clamp-2 mb-1">
            {video.title}
          </h3>
          <p className="text-sm text-gray-600">{video.channelName}</p>
          <p className="text-sm text-gray-600">
            {video.views} • {video.uploadedAt}
          </p>
        </div>
      </div>
    </div>
  )
}

export default VideoCard
