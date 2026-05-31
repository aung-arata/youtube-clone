import React, { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'

function NotificationsPage() {
  const { user, token } = useAuth()
  const navigate = useNavigate()
  const [notifications, setNotifications] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  useEffect(() => {
    if (!user) return
    fetchNotifications()
  }, [user])

  const fetchNotifications = async () => {
    try {
      setLoading(true)
      const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
      const response = await fetch(`${apiUrl}/api/users/${user.id}/notifications`, {
        headers: { Authorization: `Bearer ${token}` },
      })
      if (!response.ok) throw new Error('Failed to fetch notifications')
      const data = await response.json()
      setNotifications(data || [])
    } catch (err) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  const markAllRead = async () => {
    try {
      const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
      await fetch(`${apiUrl}/api/users/${user.id}/notifications/read`, {
        method: 'POST',
        headers: { Authorization: `Bearer ${token}` },
      })
      setNotifications((prev) => prev.map((n) => ({ ...n, is_read: true })))
    } catch {
      // Silently fail
    }
  }

  if (!user) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-lg text-gray-600 dark:text-gray-400">
          Sign in to view your notifications
        </div>
      </div>
    )
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-lg text-gray-600 dark:text-gray-400">Loading notifications...</div>
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

  const unreadCount = notifications.filter((n) => !n.is_read).length

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold dark:text-white">
          Notifications {unreadCount > 0 && <span className="text-base font-normal text-red-600">({unreadCount} unread)</span>}
        </h1>
        {unreadCount > 0 && (
          <button
            onClick={markAllRead}
            className="text-sm text-blue-500 hover:underline"
          >
            Mark all as read
          </button>
        )}
      </div>
      {notifications.length === 0 ? (
        <div className="flex items-center justify-center h-64">
          <div className="text-lg text-gray-600 dark:text-gray-400">No notifications yet.</div>
        </div>
      ) : (
        <ul className="space-y-2">
          {notifications.map((n) => (
            <li
              key={n.id}
              className={`bg-white dark:bg-gray-800 rounded-lg p-4 border border-gray-200 dark:border-gray-700 flex items-start gap-3 ${
                !n.is_read ? 'border-l-4 border-l-blue-500' : ''
              }`}
            >
              <div className="flex-1 min-w-0">
                <p className={`text-sm dark:text-gray-200 ${!n.is_read ? 'font-semibold' : ''}`}>
                  {n.message}
                </p>
                {n.created_at && (
                  <p className="text-xs text-gray-400 mt-1">
                    {new Date(n.created_at).toLocaleString()}
                  </p>
                )}
              </div>
              {!n.is_read && (
                <span className="w-2 h-2 rounded-full bg-blue-500 mt-2 shrink-0" />
              )}
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}

export default NotificationsPage
