import React, { useState, useRef, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'

const CATEGORIES = [
  'Gaming', 'Music', 'Sports', 'Education', 'Entertainment',
  'Science & Technology', 'News & Politics', 'Travel & Events',
  'Howto & Style', 'Comedy', 'Film & Animation', 'Autos & Vehicles',
  'People & Blogs', 'Pets & Animals', 'Other',
]

export default function UploadPage() {
  const { user, token } = useAuth()
  const navigate = useNavigate()

  const [videoFile, setVideoFile] = useState(null)
  const [thumbnailFile, setThumbnailFile] = useState(null)
  const [thumbnailPreview, setThumbnailPreview] = useState(null)
  const [videoDragging, setVideoDragging] = useState(false)
  const [form, setForm] = useState({
    title: '',
    description: '',
    category: '',
  })
  const [uploading, setUploading] = useState(false)
  const [progress, setProgress] = useState(0)
  const [error, setError] = useState('')

  const videoInputRef = useRef()
  const thumbnailInputRef = useRef()

  useEffect(() => {
    if (!user) navigate('/login')
  }, [user, navigate])

  const handleVideoDrop = (e) => {
    e.preventDefault()
    setVideoDragging(false)
    const file = e.dataTransfer.files[0]
    if (file && file.type.startsWith('video/')) {
      setVideoFile(file)
      if (!form.title) setForm(f => ({ ...f, title: file.name.replace(/\.[^.]+$/, '') }))
    }
  }

  const handleVideoSelect = (e) => {
    const file = e.target.files[0]
    if (file) {
      setVideoFile(file)
      if (!form.title) setForm(f => ({ ...f, title: file.name.replace(/\.[^.]+$/, '') }))
    }
  }

  const handleThumbnailSelect = (e) => {
    const file = e.target.files[0]
    if (file) {
      setThumbnailFile(file)
      setThumbnailPreview(URL.createObjectURL(file))
    }
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    if (!videoFile) { setError('Please select a video file.'); return }
    if (!form.title.trim()) { setError('Title is required.'); return }

    setError('')
    setUploading(true)
    setProgress(0)

    const data = new FormData()
    data.append('video', videoFile)
    data.append('title', form.title.trim())
    data.append('description', form.description.trim())
    data.append('category', form.category)
    data.append('channel_name', user.username)
    data.append('channel_avatar', user.avatar || '')
    if (thumbnailFile) data.append('thumbnail', thumbnailFile)

    try {
      const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'

      await new Promise((resolve, reject) => {
        const xhr = new XMLHttpRequest()
        xhr.open('POST', `${apiUrl}/api/upload/video`)
        xhr.setRequestHeader('Authorization', `Bearer ${token}`)

        xhr.upload.onprogress = (e) => {
          if (e.lengthComputable) setProgress(Math.round((e.loaded / e.total) * 100))
        }

        xhr.onload = () => {
          if (xhr.status === 201) {
            const result = JSON.parse(xhr.responseText)
            resolve(result)
            navigate(`/video/${result.id}`)
          } else {
            reject(new Error(xhr.responseText || 'Upload failed'))
          }
        }
        xhr.onerror = () => reject(new Error('Network error'))
        xhr.send(data)
      })
    } catch (err) {
      setError(err.message || 'Upload failed. Please try again.')
      setUploading(false)
      setProgress(0)
    }
  }

  if (!user) return null

  return (
    <div className="max-w-3xl mx-auto py-8 px-4">
      <h1 className="text-2xl font-bold mb-8 dark:text-white">Upload Video</h1>

      <form onSubmit={handleSubmit} className="space-y-6">
        {/* Video drop zone */}
        <div
          onDragOver={(e) => { e.preventDefault(); setVideoDragging(true) }}
          onDragLeave={() => setVideoDragging(false)}
          onDrop={handleVideoDrop}
          onClick={() => !videoFile && videoInputRef.current.click()}
          className={`border-2 border-dashed rounded-xl p-10 text-center transition-colors cursor-pointer
            ${videoDragging ? 'border-blue-500 bg-blue-50 dark:bg-blue-900/20' : 'border-gray-300 dark:border-gray-600 hover:border-blue-400 dark:hover:border-blue-500'}
            ${videoFile ? 'cursor-default' : ''}`}
        >
          {videoFile ? (
            <div className="space-y-2">
              <svg className="w-12 h-12 text-green-500 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <p className="font-medium dark:text-white">{videoFile.name}</p>
              <p className="text-sm text-gray-500 dark:text-gray-400">{(videoFile.size / 1024 / 1024).toFixed(1)} MB</p>
              <button
                type="button"
                onClick={(e) => { e.stopPropagation(); setVideoFile(null) }}
                className="text-sm text-red-500 hover:underline"
              >
                Remove
              </button>
            </div>
          ) : (
            <div className="space-y-3">
              <svg className="w-16 h-16 text-gray-400 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
              </svg>
              <p className="text-lg font-medium dark:text-gray-200">Drag and drop your video here</p>
              <p className="text-sm text-gray-500 dark:text-gray-400">or click to select a file</p>
              <p className="text-xs text-gray-400 dark:text-gray-500">MP4, MOV, AVI, MKV · up to 500 MB</p>
            </div>
          )}
          <input ref={videoInputRef} type="file" accept="video/*" className="hidden" onChange={handleVideoSelect} />
        </div>

        {/* Title */}
        <div>
          <label className="block text-sm font-medium mb-1.5 dark:text-gray-200">
            Title <span className="text-red-500">*</span>
          </label>
          <input
            type="text"
            value={form.title}
            onChange={(e) => setForm(f => ({ ...f, title: e.target.value }))}
            maxLength={100}
            placeholder="Add a title that describes your video"
            className="w-full px-4 py-2.5 border border-gray-300 dark:border-gray-600 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-gray-700 dark:text-white dark:placeholder-gray-400"
          />
          <p className="text-xs text-gray-400 text-right mt-1">{form.title.length}/100</p>
        </div>

        {/* Description */}
        <div>
          <label className="block text-sm font-medium mb-1.5 dark:text-gray-200">Description</label>
          <textarea
            value={form.description}
            onChange={(e) => setForm(f => ({ ...f, description: e.target.value }))}
            maxLength={5000}
            rows={4}
            placeholder="Tell viewers about your video"
            className="w-full px-4 py-2.5 border border-gray-300 dark:border-gray-600 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-gray-700 dark:text-white dark:placeholder-gray-400 resize-y"
          />
          <p className="text-xs text-gray-400 text-right mt-1">{form.description.length}/5000</p>
        </div>

        {/* Category */}
        <div>
          <label className="block text-sm font-medium mb-1.5 dark:text-gray-200">Category</label>
          <select
            value={form.category}
            onChange={(e) => setForm(f => ({ ...f, category: e.target.value }))}
            className="w-full px-4 py-2.5 border border-gray-300 dark:border-gray-600 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-gray-700 dark:text-white"
          >
            <option value="">Select a category</option>
            {CATEGORIES.map(c => <option key={c} value={c}>{c}</option>)}
          </select>
        </div>

        {/* Thumbnail */}
        <div>
          <label className="block text-sm font-medium mb-1.5 dark:text-gray-200">Thumbnail</label>
          <div className="flex items-start gap-4">
            <div
              onClick={() => thumbnailInputRef.current.click()}
              className="w-40 h-24 border-2 border-dashed border-gray-300 dark:border-gray-600 rounded-lg flex items-center justify-center cursor-pointer hover:border-blue-400 dark:hover:border-blue-500 overflow-hidden flex-shrink-0"
            >
              {thumbnailPreview ? (
                <img src={thumbnailPreview} alt="Thumbnail preview" className="w-full h-full object-cover" />
              ) : (
                <div className="text-center p-2">
                  <svg className="w-8 h-8 text-gray-400 mx-auto mb-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                  </svg>
                  <p className="text-xs text-gray-400">Upload</p>
                </div>
              )}
            </div>
            <div className="text-sm text-gray-500 dark:text-gray-400 pt-2">
              <p>Upload a thumbnail for your video.</p>
              <p className="mt-1">Recommended: 1280×720 (16:9), JPG or PNG.</p>
              {thumbnailFile && (
                <button
                  type="button"
                  onClick={() => { setThumbnailFile(null); setThumbnailPreview(null) }}
                  className="mt-2 text-red-500 hover:underline text-xs"
                >
                  Remove
                </button>
              )}
            </div>
          </div>
          <input ref={thumbnailInputRef} type="file" accept="image/*" className="hidden" onChange={handleThumbnailSelect} />
        </div>

        {/* Channel info (read-only) */}
        <div className="flex items-center gap-3 p-3 bg-gray-50 dark:bg-gray-800 rounded-lg">
          <div className="w-8 h-8 bg-blue-500 rounded-full flex items-center justify-center text-white font-semibold text-sm flex-shrink-0">
            {user.username.charAt(0).toUpperCase()}
          </div>
          <div>
            <p className="text-sm font-medium dark:text-white">{user.username}</p>
            <p className="text-xs text-gray-500 dark:text-gray-400">Publishing as this channel</p>
          </div>
        </div>

        {/* Error */}
        {error && (
          <div className="p-3 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg text-red-600 dark:text-red-400 text-sm">
            {error}
          </div>
        )}

        {/* Progress */}
        {uploading && (
          <div className="space-y-2">
            <div className="flex justify-between text-sm dark:text-gray-300">
              <span>Uploading…</span>
              <span>{progress}%</span>
            </div>
            <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2">
              <div
                className="bg-blue-500 h-2 rounded-full transition-all duration-300"
                style={{ width: `${progress}%` }}
              />
            </div>
          </div>
        )}

        {/* Actions */}
        <div className="flex justify-end gap-3 pt-2">
          <button
            type="button"
            onClick={() => navigate('/')}
            disabled={uploading}
            className="px-6 py-2.5 border border-gray-300 dark:border-gray-600 rounded-lg text-sm font-medium dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-700 disabled:opacity-50"
          >
            Cancel
          </button>
          <button
            type="submit"
            disabled={uploading || !videoFile}
            className="px-6 py-2.5 bg-blue-600 hover:bg-blue-700 disabled:bg-blue-400 text-white rounded-lg text-sm font-medium transition-colors disabled:cursor-not-allowed"
          >
            {uploading ? 'Uploading…' : 'Upload'}
          </button>
        </div>
      </form>
    </div>
  )
}
