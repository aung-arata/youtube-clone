import React from 'react'
import VideoCard from './VideoCard'

function VideoGrid() {
  // Sample video data
  const videos = [
    {
      id: 1,
      title: 'Building a YouTube Clone with React and Tailwind CSS',
      thumbnail: 'https://via.placeholder.com/320x180/FF0000/FFFFFF?text=Video+1',
      channelName: 'Code Master',
      channelAvatar: 'https://via.placeholder.com/36/4299E1/FFFFFF?text=CM',
      views: '1.2M views',
      uploadedAt: '2 days ago',
      duration: '12:34',
    },
    {
      id: 2,
      title: 'Golang Backend Tutorial - Building RESTful APIs',
      thumbnail: 'https://via.placeholder.com/320x180/00ADD8/FFFFFF?text=Video+2',
      channelName: 'Go Developer',
      channelAvatar: 'https://via.placeholder.com/36/48BB78/FFFFFF?text=GD',
      views: '854K views',
      uploadedAt: '1 week ago',
      duration: '24:15',
    },
    {
      id: 3,
      title: 'PostgreSQL Database Design Best Practices',
      thumbnail: 'https://via.placeholder.com/320x180/336791/FFFFFF?text=Video+3',
      channelName: 'Database Pro',
      channelAvatar: 'https://via.placeholder.com/36/9F7AEA/FFFFFF?text=DP',
      views: '623K views',
      uploadedAt: '3 days ago',
      duration: '18:45',
    },
    {
      id: 4,
      title: 'React Hooks Deep Dive - useState and useEffect',
      thumbnail: 'https://via.placeholder.com/320x180/61DAFB/000000?text=Video+4',
      channelName: 'React Expert',
      channelAvatar: 'https://via.placeholder.com/36/ED8936/FFFFFF?text=RE',
      views: '2.1M views',
      uploadedAt: '5 days ago',
      duration: '15:23',
    },
    {
      id: 5,
      title: 'Tailwind CSS Complete Guide for Beginners',
      thumbnail: 'https://via.placeholder.com/320x180/38B2AC/FFFFFF?text=Video+5',
      channelName: 'CSS Wizard',
      channelAvatar: 'https://via.placeholder.com/36/F56565/FFFFFF?text=CW',
      views: '1.8M views',
      uploadedAt: '1 week ago',
      duration: '32:10',
    },
    {
      id: 6,
      title: 'Full Stack Development Roadmap 2024',
      thumbnail: 'https://via.placeholder.com/320x180/805AD5/FFFFFF?text=Video+6',
      channelName: 'Tech Career',
      channelAvatar: 'https://via.placeholder.com/36/38B2AC/FFFFFF?text=TC',
      views: '3.5M views',
      uploadedAt: '2 weeks ago',
      duration: '45:00',
    },
  ]

  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
      {videos.map((video) => (
        <VideoCard key={video.id} video={video} />
      ))}
    </div>
  )
}

export default VideoGrid
