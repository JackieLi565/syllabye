import { ref } from 'vue'
import { debounce } from 'lodash'
import { z } from 'zod'

interface CheckUsernameResponse {
  Exists: boolean
}

export const usernameSchema = z
  .string()
  .min(3, 'Username must be at least 3 characters long')
  .max(30, 'Username cannot be longer than 30 characters')
  .regex(/^[a-z0-9._-]{3,30}$/, 'Username can only contain lowercase letters, numbers, periods, underscores, and dashes')

export function useCheckUsername() {
  const { user } = useAuth()
  const isValid = ref<{
    valid: boolean, // true if no duplicate username in database
    zod: boolean, // true if zod validation is being used
    message: string
  } | null>(null)

  const _checkUsername = async (username: string) => {
    const validation = usernameSchema.safeParse(username)

    if (!validation.success) {
      isValid.value = {
        valid: false,
        zod: true,
        message: 'Username is invalid'
      }
      return isValid.value
    }

    if (username === user.value?.nickname) {
      isValid.value = {
        valid: true,
        zod: false,
        message: 'This is your current username'
      }
      return isValid.value
    }

    try {
      const response = await useFetch<CheckUsernameResponse>(`/api/user/nickname?search=${username}`)

      if (response.data.value?.Exists) {
        isValid.value = {
          valid: false,
          zod: false,
          message: 'Username already exists'
        }
      } else {
        isValid.value = {
          valid: true,
          zod: false,
          message: 'Username is available!'
        }
      }
    } catch (error) {
      isValid.value = {
        valid: false,
        zod: false,
        message: 'Something went wrong, please try again later'
      }
    }

    return isValid.value
  }

  const checkUsername = debounce(_checkUsername, 700)

  return {
    isValid,
    checkUsername
  }
}
