export function useAuth() {
  const isLoggedIn = () => !!localStorage.getItem('token')
  const requireAuth = () => {
    if (!isLoggedIn()) window.location.href = '/login'
  }
  return { isLoggedIn, requireAuth }
}
