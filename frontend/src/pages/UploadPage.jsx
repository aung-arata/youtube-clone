import React, { useState, useRef, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'

const CATEGORIES = [
  'Gaming', 'Music', 'Sports', 'Education', 'Entertainment',
  'Science & Technology', 'News & Politics', 'Travel & Events',
  'Howto & Style', 'Comedy', 'Film & Animation', 'Autos & Vehicles',
  'People & Blogs', 'Pets & Animals', 'Other',
]

const QUALITY_ORDER = ['360p', '480p', '720p', '1080p']

export default function UploadPage() {
  const { user, token } = useAuth()
  const navigate = useNavigate()

  const [phase, setPhase] = useState('upload')
  const [videoId, setVideoId] = useState(null)
  const [transcodingStatus, setTranscodingStatus] = useState(null)
  const [visibility, setVisibility] = useState('public')

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
  const [savingMetadata, setSavingMetadata] = useState(false)
  const [saveMessage, setSaveMessage] = useState('')
  const [error, setError] = useState('')

  const videoInputRef = useRef()
  const thumbnailInputRef = useRef()

  useEffect(() => {
    if (!user) navigate('/login')
  }, [user, navigate])

  useEffect(() => {
    if (!thumbnailPreview) return
    return () => URL.revokeObjectURL(thumbnailPreview)
  }, [thumbnailPreview])

  useEffect(() => {
    if (phase !== 'processing' || !videoId) return

    const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
    let intervalId
    const poll = async () => {
      try {
        const res = await fetch(`${apiUrl}/api/upload/video/${videoId}/status`, {
          headers: { Authorization: 'Bearer ' + token },
        })
        if (!res.ok) return
        const data = await res.json()
        setTranscodingStatus(data)
        if ((data.processing_status === 'ready' || data.processing_status === 'failed') && intervalId) clearInterval(intervalId)
      } catch {
        // silent
      }
    }

    poll()
    intervalId = setInterval(poll, 3000)
    return () => clearInterval(intervalId)
  }, [phase, videoId, token])

  useEffect(() => {
    if (phase !== 'processing' || transcodingStatus?.processing_status === 'ready' || transcodingStatus?.processing_status === 'failed') return

    const onBeforeUnload = (event) => {
      event.preventDefault()
      event.returnValue = "Your video is still processing. Navigating away won't cancel it."
      return event.returnValue
    }

    window.addEventListener('beforeunload', onBeforeUnload)
    return () => window.removeEventListener('beforeunload', onBeforeUnload)
  }, [phase, transcodingStatus?.processing_status])

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

  const handleUpload = async (e) => {
    e.preventDefault()
    if (!videoFile) { setError('Please select a video file.'); return }
    if (!form.title.trim()) { setError('Title is required.'); return }

    setError('')
    setSaveMessage('')
    setUploading(true)
    setProgress(0)

    const data = new FormData()
    data.append('video', videoFile)
    data.append('title', form.title.trim())
    data.append('description', form.description.trim())
    data.append('category', form.category)
    data.append('channel_name', user.username)
    data.append('channel_avatar', user.avatar || '')
    data.append('visibility', visibility)
    if (thumbnailFile) data.append('thumbnail', thumbnailFile)

    try {
      const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'

      const result = await new Promise((resolve, reject) => {
        const xhr = new XMLHttpRequest()
        xhr.open('POST', `${apiUrl}/api/upload/video`)
        xhr.setRequestHeader('Authorization', 'Bearer ' + token)

        xhr.upload.onprogress = (event) => {
          if (event.lengthComputable) setProgress(Math.round((event.loaded / event.total) * 100))
        }

        xhr.onload = () => {
          if (xhr.status === 201) {
            resolve(JSON.parse(xhr.responseText))
          } else {
            reject(new Error(xhr.responseText || 'Upload failed'))
          }
        }
        xhr.onerror = () => reject(new Error('Network error'))
        xhr.send(data)
      })

      setVideoId(result.id)
      setVisibility(result.visibility || 'public')
      setTranscodingStatus({
        video_id: result.id,
        processing_status: result.processing_status || 'pending',
        jobs: [],
        overall_progress: 0,
      })
      setUploading(false)
      setProgress(100)
      setPhase('processing')
    } catch (err) {
      setError(err.message || 'Upload failed. Please try again.')
      setUploading(false)
      setProgress(0)
    }
  }

  const handleSaveMetadata = async () => {
    if (!videoId) return
    if (!form.title.trim()) {
      setError('Title is required.')
      return
    }

    setSavingMetadata(true)
    setError('')
    setSaveMessage('')

    try {
      const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
      const response = await fetch(`${apiUrl}/api/upload/video/${videoId}/metadata`, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
          Authorization: 'Bearer ' + token,
        },
        body: JSON.stringify({
          title: form.title.trim(),
          description: form.description.trim(),
          category: form.category,
          visibility,
        }),
      })

      if (!response.ok) throw new Error('Failed to save metadata')
      const updated = await response.json()
      setVisibility(updated.visibility || visibility)
      setSaveMessage('Metadata saved.')
    } catch (err) {
      setError(err.message || 'Failed to save metadata.')
    } finally {
      setSavingMetadata(false)
    }
  }

  const handleCancel = () => {
    if (phase === 'processing') {
      const shouldLeave = window.confirm("Your video is still processing. Navigating away won't cancel it.")
      if (!shouldLeave) return
    }
    navigate('/')
  }

  const isReady = transcodingStatus?.processing_status === 'ready'
  const isFailed = transcodingStatus?.processing_status === 'failed'

  const jobsByQuality = new Map((transcodingStatus?.jobs || []).map(job => [job.quality, job]))

  const getJobLabel = (job) => {
    if (!job) return '…'
    if (job.status === 'completed') return '✓'
    if (job.status === 'processing') return `⚙ ${job.progress}%`
    if (job.status === 'failed') return '✗ failed'
    return '…'
  }

  if (!user) return null

  return (
    <div className="max-w-3xl mx-auto py-8 px-4">
      <h1 className="text-2xl font-bold mb-8 dark:text-white">Upload Video</h1>

      {phase === 'processing' && (
        <div className="mb-6 p-4 rounded-lg border border-blue-200 bg-blue-50 text-blue-700 dark:bg-blue-900/20 dark:border-blue-700 dark:text-blue-300">
          Upload complete! Processing your video...
        </div>
      )}

      <form
        onSubmit={(e) => {
          if (phase === 'upload') {
            handleUpload(e)
            return
          }
          e.preventDefault()
          handleSaveMetadata()
        }}
        className="space-y-6"
      >
        <div
          onDragOver={(e) => { e.preventDefault(); setVideoDragging(true) }}
          onDragLeave={() => setVideoDragging(false)}
          onDrop={handleVideoDrop}
          onClick={() => phase === 'upload' && !videoFile && videoInputRef.current.click()}
          className={`border-2 border-dashed rounded-xl p-10 text-center transition-colors ${phase === 'upload' ? 'cursor-pointer' : 'cursor-default'}
            ${videoDragging ? 'border-blue-500 bg-blue-50 dark:bg-blue-900/20' : 'border-gray-300 dark:border-gray-600 hover:border-blue-400 dark:hover:border-blue-500'}`}
        >
          {videoFile ? (
            <div className="space-y-2">
              <svg className="w-12 h-12 text-green-500 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <p className="font-medium dark:text-white">{videoFile.name}</p>
              <p className="text-sm text-gray-500 dark:text-gray-400">{(videoFile.size / 1024 / 1024).toFixed(1)} MB</p>
              {phase === 'upload' && (
                <button
                  type="button"
                  onClick={(e) => { e.stopPropagation(); setVideoFile(null) }}
                  className="text-sm text-red-500 hover:underline"
                >
                  Remove
                </button>
              )}
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
          <input ref={videoInputRef} type="file" accept="video/*" className="hidden" onChange={handleVideoSelect} disabled={phase !== 'upload'} />
        </div>

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

        <div>
          <label className="block text-sm font-medium mb-1.5 dark:text-gray-200">Visibility</label>
          <select
            value={visibility}
            onChange={(e) => setVisibility(e.target.value)}
            className="w-full px-4 py-2.5 border border-gray-300 dark:border-gray-600 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-gray-700 dark:text-white"
          >
            <option value="public">Public</option>
            <option value="unlisted">Unlisted</option>
            <option value="private">Private</option>
          </select>
        </div>

        <div>
          <label className="block text-sm font-medium mb-1.5 dark:text-gray-200">Thumbnail</label>
          <div className="flex items-start gap-4">
            <div
              onClick={() => phase === 'upload' && thumbnailInputRef.current.click()}
              className="w-40 h-24 border-2 border-dashed border-gray-300 dark:border-gray-600 rounded-lg flex items-center justify-center hover:border-blue-400 dark:hover:border-blue-500 overflow-hidden flex-shrink-0"
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
              {thumbnailFile && phase === 'upload' && (
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
          <input ref={thumbnailInputRef} type="file" accept="image/*" className="hidden" onChange={handleThumbnailSelect} disabled={phase !== 'upload'} />
        </div>

        <div className="flex items-center gap-3 p-3 bg-gray-50 dark:bg-gray-800 rounded-lg">
          <div className="w-8 h-8 bg-blue-500 rounded-full flex items-center justify-center text-white font-semibold text-sm flex-shrink-0">
            {user.username.charAt(0).toUpperCase()}
          </div>
          <div>
            <p className="text-sm font-medium dark:text-white">{user.username}</p>
            <p className="text-xs text-gray-500 dark:text-gray-400">Publishing as this channel</p>
          </div>
        </div>

        {error && (
          <div className="p-3 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg text-red-600 dark:text-red-400 text-sm">
            {error}
          </div>
        )}

        {saveMessage && (
          <div className="p-3 bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-lg text-green-700 dark:text-green-300 text-sm">
            {saveMessage}
          </div>
        )}

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

        {phase === 'processing' && (
          <div className="space-y-4 p-4 border border-gray-200 dark:border-gray-700 rounded-lg">
            {!isReady && !isFailed ? (
              <>
                <div className="flex justify-between text-sm dark:text-gray-300">
                  <span>Transcoding progress</span>
                  <span>{transcodingStatus?.overall_progress || 0}%</span>
                </div>
                <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2">
                  <div
                    className="bg-blue-500 h-2 rounded-full transition-all duration-300"
                    style={{ width: `${transcodingStatus?.overall_progress || 0}%` }}
                  />
                </div>
                <div className="space-y-2 text-sm dark:text-gray-300">
                  {QUALITY_ORDER.map((quality) => {
                    const job = jobsByQuality.get(quality)
                    return (
                      <div key={quality} className="flex items-center justify-between">
                        <span>{quality}</span>
                        <span>{getJobLabel(job)}</span>
                      </div>
                    )
                  })}
                </div>
              </>
            ) : isReady ? (
              <div className="space-y-3">
                <p className="text-green-700 dark:text-green-400 font-medium">Your video is ready!</p>
                <button
                  type="button"
                  onClick={() => navigate(`/video/${videoId}`)}
                  className="px-5 py-2 bg-green-600 hover:bg-green-700 text-white rounded-lg text-sm font-medium"
                >
                  Go to Video
                </button>
              </div>
            ) : (
              <p className="text-red-700 dark:text-red-400 font-medium">Video processing failed. You can still save metadata and retry upload later.</p>
            )}
          </div>
        )}

        <div className="flex justify-end gap-3 pt-2">
          <button
            type="button"
            onClick={handleCancel}
            disabled={uploading || savingMetadata}
            className="px-6 py-2.5 border border-gray-300 dark:border-gray-600 rounded-lg text-sm font-medium dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-700 disabled:opacity-50"
          >
            Cancel
          </button>

          {phase === 'upload' ? (
            <button
              type="submit"
              disabled={uploading || !videoFile}
              className="px-6 py-2.5 bg-blue-600 hover:bg-blue-700 disabled:bg-blue-400 text-white rounded-lg text-sm font-medium transition-colors disabled:cursor-not-allowed"
            >
              {uploading ? 'Uploading…' : 'Upload'}
            </button>
          ) : (
            <button
              type="button"
              disabled={savingMetadata}
              onClick={handleSaveMetadata}
              className="px-6 py-2.5 bg-blue-600 hover:bg-blue-700 disabled:bg-blue-400 text-white rounded-lg text-sm font-medium transition-colors disabled:cursor-not-allowed"
            >
              {savingMetadata ? 'Saving…' : 'Save'}
            </button>
          )}
        </div>
      </form>
    </div>
  )
}
