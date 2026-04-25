export default defineNuxtRouteMiddleware((to) => {
  if (!process.client) return

  const token = localStorage.getItem('token')
  const publicPaths = [
    '/', '/login', '/signup',
    '/recipes', '/chefs', '/news', '/shows', '/sweepstakes'
  ]
  const protectedPrefixes = [
    '/profile', '/favorites', '/recipes/create', '/recipes/my', '/recipes/edit', '/admin'
  ]

  // Do not clear auth globally; only guard protected routes
  const isPublic = publicPaths.includes(to.path)
  const isProtected = protectedPrefixes.some(p => to.path.startsWith(p))

  if (isProtected && !token) {
    return navigateTo('/login')
  }

  // Additional role-based guard for admin pages
  if (to.path.startsWith('/admin')) {
    const role = localStorage.getItem('role') || ''
    if (role !== 'admin') {
      return navigateTo('/')
    }
  }
})
