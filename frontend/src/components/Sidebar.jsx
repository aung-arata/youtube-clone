import React from 'react'

function Sidebar({ isOpen }) {
  const menuItems = [
    { icon: 'home', label: 'Home' },
    { icon: 'trending', label: 'Trending' },
    { icon: 'subscriptions', label: 'Subscriptions' },
  ]

  const libraryItems = [
    { icon: 'library', label: 'Library' },
    { icon: 'history', label: 'History' },
    { icon: 'watch-later', label: 'Watch Later' },
    { icon: 'liked', label: 'Liked Videos' },
  ]

  if (!isOpen) return null

  return (
    <aside className="w-64 bg-white dark:bg-gray-800 border-r border-gray-200 dark:border-gray-700 overflow-y-auto">
      <nav className="py-2">
        <div className="px-3 py-2">
          {menuItems.map((item) => (
            <button
              key={item.label}
              className="w-full flex items-center gap-6 px-3 py-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg text-sm dark:text-gray-200"
            >
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                {item.icon === 'home' && (
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
                )}
                {item.icon === 'trending' && (
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
                )}
                {item.icon === 'subscriptions' && (
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 4v16M17 4v16M3 8h4m10 0h4M3 12h18M3 16h4m10 0h4M4 20h16a1 1 0 001-1V5a1 1 0 00-1-1H4a1 1 0 00-1 1v14a1 1 0 001 1z" />
                )}
              </svg>
              <span>{item.label}</span>
            </button>
          ))}
        </div>

        <hr className="my-2 border-gray-200 dark:border-gray-700" />

        <div className="px-3 py-2">
          {libraryItems.map((item) => (
            <button
              key={item.label}
              className="w-full flex items-center gap-6 px-3 py-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg text-sm dark:text-gray-200"
            >
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                {item.icon === 'library' && (
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
                )}
                {item.icon === 'history' && (
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                )}
                {item.icon === 'watch-later' && (
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                )}
                {item.icon === 'liked' && (
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M14 10h4.764a2 2 0 011.789 2.894l-3.5 7A2 2 0 0115.263 21h-4.017c-.163 0-.326-.02-.485-.06L7 20m7-10V5a2 2 0 00-2-2h-.095c-.5 0-.905.405-.905.905 0 .714-.211 1.412-.608 2.006L7 11v9m7-10h-2M7 20H5a2 2 0 01-2-2v-6a2 2 0 012-2h2.5" />
                )}
              </svg>
              <span>{item.label}</span>
            </button>
          ))}
        </div>
      </nav>
    </aside>
  )
}

export default Sidebar
