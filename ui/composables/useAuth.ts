import { useState, useFetch, useRuntimeConfig } from '#app'

interface Session {
  userId: string
}

interface User {
  id: string
  currentYear?: number | null
  email: string
  fullname: string
  gender?: 'male' | 'female' | 'other'
  nickname?: string
  picture: string
  programId?: string
}

export function useAuth() {
  const config = useRuntimeConfig()

  const session = useState<Session | null>('session', () => null)
  const user = useState<User | null>('user', () => null)

  async function fetchSession() {
    console.log('Fetching session...') // Debug log
    const { data, error } = await useFetch<Session>(`${config.public.apiUrl}/me`, {
      credentials: 'include',
    })
  
    console.log('Session response:', data.value) // Debugging
  
    if (!error.value && data.value) {
      session.value = data.value
      await fetchUser(data.value.userId)
    } else {
      session.value = null
      user.value = null
    }
  }

  async function fetchUser(userId: string) {
    const { data, error } = await useFetch<User>(`${config.public.apiUrl}/users/${userId}`, {
      credentials: 'include',
    })

    if (!error.value && data.value) {
      user.value = data.value
    } else {
      user.value = null
    }
  }

  return {
    session,
    user,
    fetchSession,
  }
}
