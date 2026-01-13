import React, { useState, useEffect } from 'react'

function UserProfile({ userId }) {
  const [user, setUser] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [editing, setEditing] = useState(false)
  const [formData, setFormData] = useState({
    username: '',
    email: '',
    avatar: ''
  })

  useEffect(() => {
    if (userId) {
      fetchUser()
    }
  }, [userId])

  const fetchUser = async () => {
    try {
      setLoading(true)
      const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
      const response = await fetch(`${apiUrl}/api/users/${userId}`)
      if (!response.ok) {
        throw new Error('Failed to fetch user')
      }
      const data = await response.json()
      setUser(data)
      setFormData({
        username: data.username,
        email: data.email,
        avatar: data.avatar || ''
      })
    } catch (err) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  const handleUpdate = async (e) => {
    e.preventDefault()
    try {
      const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
      const response = await fetch(`${apiUrl}/api/users/${userId}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(formData)
      })
      if (!response.ok) {
        throw new Error('Failed to update user')
      }
      const data = await response.json()
      setUser(data)
      setEditing(false)
    } catch (err) {
      setError(err.message)
    }
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-lg text-gray-600 dark:text-gray-400">Loading profile...</div>
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

  if (!user) {
    return null
  }

  return (
    <div className="max-w-2xl mx-auto p-6">
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
        <div className="flex items-center gap-6 mb-6">
          <img
            src={user.avatar || `https://via.placeholder.com/100/4299E1/FFFFFF?text=${user.username.charAt(0)}`}
            alt={user.username}
            className="w-24 h-24 rounded-full"
          />
          <div>
            <h1 className="text-2xl font-bold dark:text-white">{user.username}</h1>
            <p className="text-gray-600 dark:text-gray-400">{user.email}</p>
          </div>
        </div>

        {editing ? (
          <form onSubmit={handleUpdate} className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Username
              </label>
              <input
                type="text"
                value={formData.username}
                onChange={(e) => setFormData({ ...formData, username: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md dark:bg-gray-700 dark:text-white"
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Email
              </label>
              <input
                type="email"
                value={formData.email}
                onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md dark:bg-gray-700 dark:text-white"
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Avatar URL
              </label>
              <input
                type="text"
                value={formData.avatar}
                onChange={(e) => setFormData({ ...formData, avatar: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md dark:bg-gray-700 dark:text-white"
              />
            </div>
            <div className="flex gap-2">
              <button
                type="submit"
                className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
              >
                Save
              </button>
              <button
                type="button"
                onClick={() => setEditing(false)}
                className="px-4 py-2 bg-gray-300 dark:bg-gray-600 text-gray-900 dark:text-white rounded-md hover:bg-gray-400 dark:hover:bg-gray-500"
              >
                Cancel
              </button>
            </div>
          </form>
        ) : (
          <button
            onClick={() => setEditing(true)}
            className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
          >
            Edit Profile
          </button>
        )}
      </div>
    </div>
  )
}

export default UserProfile
