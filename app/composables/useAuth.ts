export const useAuth = () => {
  const user = useState('auth:user', () => null as { name: string, avatar?: string, email: string, role?: string } | null)
  const admin = useState('auth:admin', () => null as { name: string, email: string, role: string } | null)

  // Sync with Pinia store for middleware compatibility
  const authStore = useAuthStore()

  const isLoggedIn = computed(() => !!user.value || authStore.isAuthenticated)
  const isAdminLoggedIn = computed(() => !!admin.value || authStore.userRole?.toLowerCase().includes('admin') === true)

  const login = (userData: { name: string, email: string, avatar?: string, role?: string }) => {
    user.value = userData
  }

  const adminLogin = (adminData: { name: string, email: string, role: string }) => {
    admin.value = adminData
  }

  const logout = () => {
    user.value = null
    authStore.clearAuth()
    navigateTo('/auth/login')
  }

  const adminLogout = () => {
    admin.value = null
    authStore.clearAuth()
    navigateTo('/admin/login')
  }

  return {
    user,
    admin,
    isLoggedIn,
    isAdminLoggedIn,
    login,
    adminLogin,
    logout,
    adminLogout
  }
}
