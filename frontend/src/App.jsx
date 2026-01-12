import React, { useState } from 'react'
import Header from './components/Header'
import Sidebar from './components/Sidebar'
import VideoGrid from './components/VideoGrid'

function App() {
  const [sidebarOpen, setSidebarOpen] = useState(true)

  return (
    <div className="flex flex-col h-screen bg-gray-100">
      <Header onMenuClick={() => setSidebarOpen(!sidebarOpen)} />
      <div className="flex flex-1 overflow-hidden">
        <Sidebar isOpen={sidebarOpen} />
        <main className="flex-1 overflow-y-auto p-6">
          <VideoGrid />
        </main>
      </div>
    </div>
  )
}

export default App
