export const useAuth = () => {
  // Use Pinia store as the single source of truth
  const authStore = useAuthStore()

  const user = computed(() => {
    if (!authStore.user) return null
    return {
      name: `${authStore.user.firstName} ${authStore.user.lastName}`,
      email: authStore.user.email,
      avatar: (authStore.user as any).avatar || undefined,
      role: authStore.userRole || undefined
    }
  })

  const isLoggedIn = computed(() => authStore.isAuthenticated)
  const isAdminLoggedIn = computed(() => authStore.userRole?.toLowerCase().includes('admin') === true)

  const logout = async () => {
    await authStore.logout()
    return navigateTo('/auth/login')
  }

  const adminLogout = async () => {
    await authStore.logout()
    return navigateTo('/admin/login')
  }

  return {
    user,
    isLoggedIn,
    isAdminLoggedIn,
    logout,
    adminLogout
  }
}
