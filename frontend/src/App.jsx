import React, { useState, useRef, useEffect } from 'react'
import Header from './components/Header'
import Sidebar from './components/Sidebar'
import VideoGrid from './components/VideoGrid'

function App() {
  const [sidebarOpen, setSidebarOpen] = useState(true)
  const [darkMode, setDarkMode] = useState(() => {
    const saved = localStorage.getItem('darkMode')
    return saved === 'true'
  })
  const videoGridRef = useRef()

  useEffect(() => {
    localStorage.setItem('darkMode', darkMode)
    if (darkMode) {
      document.documentElement.classList.add('dark')
    } else {
      document.documentElement.classList.remove('dark')
    }
  }, [darkMode])

  const handleSearch = (query) => {
    if (videoGridRef.current && videoGridRef.current.fetchVideos) {
      videoGridRef.current.fetchVideos(query)
    }
  }

  const toggleDarkMode = () => {
    setDarkMode(!darkMode)
  }

  return (
    <div className="flex flex-col h-screen bg-gray-100 dark:bg-gray-900">
      <Header 
        onMenuClick={() => setSidebarOpen(!sidebarOpen)} 
        onSearch={handleSearch}
        darkMode={darkMode}
        onToggleDarkMode={toggleDarkMode}
      />
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
