export default defineNuxtRouteMiddleware((from, to) => {
  if (import.meta.client) return
  const config = useRuntimeConfig()
  const { session } = useAuth()


  if (!session.value) {
    return navigateTo(`/?redirect=${encodeURIComponent(config.public.siteUrl + to.fullPath)}`, { replace: true })
  }
})