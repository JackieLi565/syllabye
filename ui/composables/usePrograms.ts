import type { Program } from '~/types/types'

export function usePrograms(params?: {
  search?: string
  faculty?: string
  programId?: string
}) {
  const query = new URLSearchParams()

  if (params?.search) query.set('search', params.search)
  if (params?.faculty) query.set('category', params.faculty)

  const queryString = query.toString()

  const key = params?.programId
    ? `program-${params.programId}`
    : `programs-${queryString}`

  const { data: program, status: programStatus, error: programError } = useAsyncData<Program | null>(
    key,
    async () => {
      if (params?.programId) {
        return await $fetch<Program>(`/api/programs/${params.programId}`)
      }
      return null
    },
    {
      server: false,
    }
  )

  const { data: programs, status: programsStatus, error: programsError } = useAsyncData<Program[]>(
    `programs-${queryString}`,
    async () => {
      if (!params?.programId) {
        return await $fetch<Program[]>(`/api/programs/programs?${queryString}`)
      }
      return []
    },
    {
      server: false,
    }
  )

  return {
    program,
    programStatus,
    programError,
    programs,
    programsStatus,
    programsError,
  }
}
