import React, { useState, useRef } from 'react'
import Header from './components/Header'
import Sidebar from './components/Sidebar'
import VideoGrid from './components/VideoGrid'

function App() {
  const [sidebarOpen, setSidebarOpen] = useState(true)
  const videoGridRef = useRef()

  const handleSearch = (query) => {
    if (videoGridRef.current && videoGridRef.current.fetchVideos) {
      videoGridRef.current.fetchVideos(query)
    }
  }

  return (
    <div className="flex flex-col h-screen bg-gray-100">
      <Header onMenuClick={() => setSidebarOpen(!sidebarOpen)} onSearch={handleSearch} />
      <div className="flex flex-1 overflow-hidden">
        <Sidebar isOpen={sidebarOpen} />
        <main className="flex-1 overflow-y-auto p-6">
          <VideoGrid ref={videoGridRef} />
        </main>
      </div>
    </div>
  )
}

export default App
