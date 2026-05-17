import React, { useState, useEffect, useCallback } from 'react'
import { useParams, Link } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'
import VideoCard from '../components/VideoCard'

// ── Reply thread for a single comment ────────────────────────────────────────
function CommentReplies({ videoId, comment, user, token, apiUrl, formatDate }) {
  const [replies, setReplies] = useState([])
  const [expanded, setExpanded] = useState(false)
  const [replyText, setReplyText] = useState('')
  const [showReplyBox, setShowReplyBox] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const [loadingReplies, setLoadingReplies] = useState(false)

  const fetchReplies = useCallback(async () => {
    setLoadingReplies(true)
    try {
      const res = await fetch(`${apiUrl}/api/videos/${videoId}/comments/${comment.id}/replies`)
      if (res.ok) setReplies(await res.json())
    } catch { /* silent */ } finally {
      setLoadingReplies(false)
    }
  }, [apiUrl, videoId, comment.id])

  const handleExpand = () => {
    if (!expanded) fetchReplies()
    setExpanded(e => !e)
  }

  const handleReplySubmit = async (e) => {
    e.preventDefault()
    if (!replyText.trim() || !user) return
    setSubmitting(true)
    try {
      const res = await fetch(`${apiUrl}/api/videos/${videoId}/comments/${comment.id}/replies`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
        body: JSON.stringify({ user_id: user.id, username: user.username, content: replyText.trim() }),
      })
      if (res.ok) {
        setReplyText('')
        setShowReplyBox(false)
        setExpanded(true)
        fetchReplies()
      }
    } catch { /* silent */ } finally {
      setSubmitting(false)
    }
  }

  return (
    <div>
      {/* Reply / expand controls */}
      <div className="flex items-center gap-4 mt-2 ml-13">
        {user && (
          <button
            onClick={() => setShowReplyBox(s => !s)}
            className="text-xs font-semibold text-gray-600 dark:text-gray-400 hover:text-blue-500 dark:hover:text-blue-400"
          >
            Reply
          </button>
        )}
        {comment.reply_count > 0 && (
          <button
            onClick={handleExpand}
            className="flex items-center gap-1 text-xs font-semibold text-blue-600 dark:text-blue-400 hover:underline"
          >
            <svg className={`w-3 h-3 transition-transform ${expanded ? 'rotate-180' : ''}`} fill="currentColor" viewBox="0 0 20 20">
              <path fillRule="evenodd" d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" clipRule="evenodd" />
            </svg>
            {expanded ? 'Hide' : `${comment.reply_count} ${comment.reply_count === 1 ? 'reply' : 'replies'}`}
          </button>
        )}
      </div>

      {/* Inline reply input */}
      {showReplyBox && user && (
        <form onSubmit={handleReplySubmit} className="flex gap-2 mt-2 ml-13">
          <div className="w-7 h-7 bg-blue-500 rounded-full flex items-center justify-center text-white text-xs font-semibold shrink-0 mt-1">
            {user.username.charAt(0).toUpperCase()}
          </div>
          <div className="flex-1">
            <input
              autoFocus
              type="text"
              value={replyText}
              onChange={e => setReplyText(e.target.value)}
              placeholder={`@${comment.username} `}
              className="w-full px-0 py-1 border-b border-gray-300 dark:border-gray-600 bg-transparent text-sm dark:text-white focus:outline-none focus:border-blue-500"
            />
            <div className="flex justify-end gap-2 mt-1">
              <button type="button" onClick={() => { setShowReplyBox(false); setReplyText('') }}
                className="px-3 py-1 text-xs rounded-full hover:bg-gray-100 dark:hover:bg-gray-700 dark:text-gray-300">
                Cancel
              </button>
              <button type="submit" disabled={!replyText.trim() || submitting}
                className="px-3 py-1 text-xs bg-blue-600 text-white rounded-full hover:bg-blue-700 disabled:opacity-50">
                Reply
              </button>
            </div>
          </div>
        </form>
      )}

      {/* Replies list */}
      {expanded && (
        <div className="ml-13 mt-2 space-y-3">
          {loadingReplies ? (
            <p className="text-xs text-gray-500 dark:text-gray-400">Loading…</p>
          ) : replies.map(reply => (
            <div key={reply.id} className="flex gap-2">
              <div className="w-7 h-7 bg-gray-500 rounded-full flex items-center justify-center text-white text-xs font-semibold shrink-0 mt-0.5">
                {(reply.username || 'U').charAt(0).toUpperCase()}
              </div>
              <div>
                <div className="flex items-center gap-2">
                  <span className="font-semibold text-xs dark:text-white">{reply.username || 'User'}</span>
                  <span className="text-xs text-gray-500 dark:text-gray-400">{formatDate(reply.created_at)}</span>
                </div>
                <p className="text-sm dark:text-gray-300 mt-0.5">{reply.content}</p>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}

function VideoPage() {
  const { id } = useParams()
  const { user, token } = useAuth()
  const [video, setVideo] = useState(null)
  const [comments, setComments] = useState([])
  const [recommendations, setRecommendations] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [newComment, setNewComment] = useState('')
  const [submitting, setSubmitting] = useState(false)

  const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'

  useEffect(() => {
    fetchVideo()
    fetchComments()
    fetchRecommendations()
    trackView()
  }, [id])

  const trackView = async () => {
    const viewKey = `view_counted_${id}`
    if (sessionStorage.getItem(viewKey)) return
    try {
      await fetch(`${apiUrl}/api/videos/${id}/views`, { method: 'POST' })
      sessionStorage.setItem(viewKey, '1')
      if (user) {
        await fetch(`${apiUrl}/api/users/${user.id}/history`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ video_id: parseInt(id) }),
        })
      }
    } catch {
      // Silently fail for tracking
    }
  }

  const fetchVideo = async () => {
    try {
      setLoading(true)
      const response = await fetch(`${apiUrl}/api/videos/${id}`)
      if (!response.ok) throw new Error('Video not found')
      const data = await response.json()
      setVideo(data)
    } catch (err) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  const fetchComments = async () => {
    try {
      const response = await fetch(`${apiUrl}/api/videos/${id}/comments`)
      if (response.ok) {
        const data = await response.json()
        setComments(data || [])
      }
    } catch {
      // Silently fail
    }
  }

  const fetchRecommendations = async () => {
    try {
      const response = await fetch(`${apiUrl}/api/videos/${id}/recommendations`)
      if (response.ok) {
        const data = await response.json()
        setRecommendations(data || [])
      }
    } catch {
      // Silently fail
    }
  }

  const handleLike = async () => {
    try {
      await fetch(`${apiUrl}/api/videos/${id}/like`, { method: 'POST' })
      fetchVideo()
    } catch {
      // Silently fail
    }
  }

  const handleDislike = async () => {
    try {
      await fetch(`${apiUrl}/api/videos/${id}/dislike`, { method: 'POST' })
      fetchVideo()
    } catch {
      // Silently fail
    }
  }

  const handleCommentSubmit = async (e) => {
    e.preventDefault()
    if (!newComment.trim() || !user) return
    setSubmitting(true)
    try {
      const response = await fetch(`${apiUrl}/api/videos/${id}/comments`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          content: newComment,
          user_id: user.id,
          username: user.username,
        }),
      })
      if (response.ok) {
        setNewComment('')
        fetchComments()
      }
    } catch {
      // Silently fail
    } finally {
      setSubmitting(false)
    }
  }

  const formatViews = (count) => {
    if (count >= 1000000) return `${(count / 1000000).toFixed(1)}M views`
    if (count >= 1000) return `${(count / 1000).toFixed(1)}K views`
    return `${count} views`
  }

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    })
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-lg text-gray-600 dark:text-gray-400">Loading video...</div>
      </div>
    )
  }

  if (error || !video) {
    return (
      <div className="flex flex-col items-center justify-center h-64 gap-4">
        <div className="text-lg text-red-600 dark:text-red-400">Error: {error || 'Video not found'}</div>
        <Link to="/" className="text-blue-600 hover:underline">Back to Home</Link>
      </div>
    )
  }

  return (
    <div className="flex flex-col lg:flex-row gap-6">
      <div className="flex-1">
        {/* Video Player */}
        <div className="relative bg-black rounded-lg overflow-hidden aspect-video mb-4">
          {video.url ? (
            <video
              src={video.url.startsWith('/uploads') ? `${apiUrl}/api${video.url}` : video.url}
              controls
              autoPlay
              className="w-full h-full"
              poster={video.thumbnail ? `${apiUrl}/api${video.thumbnail}` : undefined}
            />
          ) : (
            <img
              src={video.thumbnail ? `${apiUrl}/api${video.thumbnail}` : `https://via.placeholder.com/854x480/333/FFFFFF?text=Video+Player`}
              alt={video.title}
              className="w-full h-full object-cover"
            />
          )}
        </div>

        {/* Video Info */}
        <h1 className="text-xl font-bold dark:text-white mb-2">{video.title}</h1>

        <div className="flex flex-wrap items-center justify-between gap-4 mb-4">
          <div className="flex items-center gap-3">
            {video.channel_avatar && !video.channel_avatar.startsWith('/uploads') ? (
              <img
                src={video.channel_avatar}
                alt={video.channel_name}
                className="w-10 h-10 rounded-full object-cover"
              />
            ) : video.channel_avatar ? (
              <img
                src={`${apiUrl}/api${video.channel_avatar}`}
                alt={video.channel_name}
                className="w-10 h-10 rounded-full object-cover"
              />
            ) : (
              <div className="w-10 h-10 bg-blue-500 rounded-full flex items-center justify-center text-white font-semibold shrink-0">
                {(video.channel_name || 'U').charAt(0).toUpperCase()}
              </div>
            )}
            <div>
              <p className="font-semibold dark:text-white">{video.channel_name}</p>
              <p className="text-sm text-gray-600 dark:text-gray-400">
                {formatViews(video.views)} • {formatDate(video.uploaded_at)}
              </p>
            </div>
          </div>

          <div className="flex items-center gap-2">
            <button
              onClick={handleLike}
              className="flex items-center gap-1 px-4 py-2 bg-gray-100 dark:bg-gray-700 rounded-full hover:bg-gray-200 dark:hover:bg-gray-600"
            >
              <svg className="w-5 h-5 dark:text-gray-200" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M14 10h4.764a2 2 0 011.789 2.894l-3.5 7A2 2 0 0115.263 21h-4.017c-.163 0-.326-.02-.485-.06L7 20m7-10V5a2 2 0 00-2-2h-.095c-.5 0-.905.405-.905.905 0 .714-.211 1.412-.608 2.006L7 11v9m7-10h-2M7 20H5a2 2 0 01-2-2v-6a2 2 0 012-2h2.5" />
              </svg>
              <span className="dark:text-gray-200">{video.likes || 0}</span>
            </button>
            <button
              onClick={handleDislike}
              className="flex items-center gap-1 px-4 py-2 bg-gray-100 dark:bg-gray-700 rounded-full hover:bg-gray-200 dark:hover:bg-gray-600"
            >
              <svg className="w-5 h-5 dark:text-gray-200 rotate-180" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M14 10h4.764a2 2 0 011.789 2.894l-3.5 7A2 2 0 0115.263 21h-4.017c-.163 0-.326-.02-.485-.06L7 20m7-10V5a2 2 0 00-2-2h-.095c-.5 0-.905.405-.905.905 0 .714-.211 1.412-.608 2.006L7 11v9m7-10h-2M7 20H5a2 2 0 01-2-2v-6a2 2 0 012-2h2.5" />
              </svg>
              <span className="dark:text-gray-200">{video.dislikes || 0}</span>
            </button>
          </div>
        </div>

        {/* Description */}
        {video.description && (
          <div className="bg-gray-100 dark:bg-gray-800 rounded-lg p-4 mb-6">
            <p className="text-sm dark:text-gray-300 whitespace-pre-wrap">{video.description}</p>
          </div>
        )}

        {/* Comments */}
        <div className="mb-6">
          <h2 className="text-lg font-bold dark:text-white mb-4">
            {comments.length} Comment{comments.length !== 1 ? 's' : ''}
          </h2>

          {user ? (
            <form onSubmit={handleCommentSubmit} className="mb-6">
              <div className="flex gap-3">
                <div className="w-10 h-10 bg-blue-500 rounded-full flex items-center justify-center text-white font-semibold shrink-0">
                  {user.username.charAt(0).toUpperCase()}
                </div>
                <div className="flex-1">
                  <input
                    type="text"
                    value={newComment}
                    onChange={(e) => setNewComment(e.target.value)}
                    placeholder="Add a comment..."
                    className="w-full px-0 py-2 border-b border-gray-300 dark:border-gray-600 bg-transparent dark:text-white focus:outline-none focus:border-blue-500"
                  />
                  <div className="flex justify-end gap-2 mt-2">
                    <button
                      type="button"
                      onClick={() => setNewComment('')}
                      className="px-4 py-2 text-sm rounded-full hover:bg-gray-100 dark:hover:bg-gray-700 dark:text-gray-300"
                    >
                      Cancel
                    </button>
                    <button
                      type="submit"
                      disabled={!newComment.trim() || submitting}
                      className="px-4 py-2 text-sm bg-blue-600 text-white rounded-full hover:bg-blue-700 disabled:opacity-50"
                    >
                      Comment
                    </button>
                  </div>
                </div>
              </div>
            </form>
          ) : (
            <p className="text-sm text-gray-500 dark:text-gray-400 mb-4">
              <Link to="/login" className="text-blue-600 hover:underline">Sign in</Link> to comment.
            </p>
          )}

          <div className="space-y-4">
            {comments.map((comment) => (
              <div key={comment.id} className="flex gap-3">
                <div className="w-10 h-10 bg-gray-400 rounded-full flex items-center justify-center text-white font-semibold shrink-0">
                  {(comment.username || 'U').charAt(0).toUpperCase()}
                </div>
                <div className="flex-1">
                  <div className="flex items-center gap-2">
                    <span className="font-semibold text-sm dark:text-white">{comment.username || 'User'}</span>
                    <span className="text-xs text-gray-500 dark:text-gray-400">
                      {comment.created_at ? formatDate(comment.created_at) : ''}
                    </span>
                  </div>
                  <p className="text-sm dark:text-gray-300 mt-1">{comment.content}</p>
                  <CommentReplies
                    videoId={id}
                    comment={comment}
                    user={user}
                    token={token}
                    apiUrl={apiUrl}
                    formatDate={formatDate}
                  />
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* Recommendations Sidebar */}
      <div className="lg:w-96">
        <h3 className="text-lg font-bold dark:text-white mb-4">Recommended</h3>
        <div className="space-y-3">
          {recommendations.map((rec) => (
            <VideoCard key={rec.id} video={rec} compact />
          ))}
          {recommendations.length === 0 && (
            <p className="text-sm text-gray-500 dark:text-gray-400">No recommendations available</p>
          )}
        </div>
      </div>
    </div>
  )
}

export default VideoPage
