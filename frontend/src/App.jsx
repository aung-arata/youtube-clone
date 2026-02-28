import React, { useState, useRef, useEffect, useCallback } from 'react'
import { Routes, Route, useNavigate, useLocation } from 'react-router-dom'
import Header from './components/Header'
import Sidebar from './components/Sidebar'
import VideoGrid from './components/VideoGrid'
import WatchHistory from './components/WatchHistory'
import UserProfile from './components/UserProfile'
import VideoPage from './pages/VideoPage'
import TrendingPage from './pages/TrendingPage'
import LoginPage from './pages/LoginPage'
import SignupPage from './pages/SignupPage'
import { useAuth } from './contexts/AuthContext'

function AppLayout() {
  const [sidebarOpen, setSidebarOpen] = useState(true)
  const [darkMode, setDarkMode] = useState(() => {
    const saved = localStorage.getItem('darkMode')
    return saved === 'true'
  })
  const [pendingSearch, setPendingSearch] = useState(null)
  const videoGridRef = useRef()
  const navigate = useNavigate()
  const location = useLocation()
  const { user } = useAuth()

  useEffect(() => {
    localStorage.setItem('darkMode', darkMode)
    if (darkMode) {
      document.documentElement.classList.add('dark')
    } else {
      document.documentElement.classList.remove('dark')
    }
  }, [darkMode])

  // Execute pending search once we're on the home page and VideoGrid is mounted
  useEffect(() => {
    if (pendingSearch !== null && location.pathname === '/' && videoGridRef.current?.fetchVideos) {
      videoGridRef.current.fetchVideos(pendingSearch)
      setPendingSearch(null)
    }
  }, [pendingSearch, location.pathname])

  const handleSearch = useCallback((query) => {
    if (location.pathname === '/' && videoGridRef.current?.fetchVideos) {
      videoGridRef.current.fetchVideos(query)
    } else {
      setPendingSearch(query)
      navigate('/')
    }
  }, [location.pathname, navigate])

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
          <Routes>
            <Route path="/" element={<VideoGrid ref={videoGridRef} />} />
            <Route path="/trending" element={<TrendingPage />} />
            <Route path="/video/:id" element={<VideoPage />} />
            <Route path="/history" element={<WatchHistory userId={user?.id} />} />
            <Route path="/profile/:userId" element={<UserProfile userId={user?.id} />} />
            <Route path="/subscriptions" element={<SubscriptionsPlaceholder />} />
          </Routes>
        </main>
      </div>
    </div>
  )
}

function SubscriptionsPlaceholder() {
  const { user } = useAuth()
  if (!user) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-lg text-gray-600 dark:text-gray-400">
          Sign in to view your subscriptions
        </div>
      </div>
    )
  }
  return (
    <div className="flex items-center justify-center h-64">
      <div className="text-lg text-gray-600 dark:text-gray-400">
        No subscriptions yet. Browse videos and subscribe to channels!
      </div>
    </div>
  )
}

function App() {
  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route path="/signup" element={<SignupPage />} />
      <Route path="/*" element={<AppLayout />} />
    </Routes>
  )
}

export default App
