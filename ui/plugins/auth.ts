import { defineNuxtPlugin } from '#app'
import { useAuth } from '~/composables/useAuth'

export default defineNuxtPlugin(async () => {
  const { session, user, fetchSession } = useAuth()

  await fetchSession()

  return {
    provide: {
      session,
      user,
      refreshAuth: fetchSession,
    },
  }
})
