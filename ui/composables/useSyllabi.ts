import type { Syllabus } from '~/types/types'

export function useSyllabi(params?: {
  userId?: string
  courseId?: string
  year?: string
  semester?: string
  page?: number
  size?: number
  syllabusId?: string
}) {
  const query = new URLSearchParams()

  if (params?.userId) query.set('userId', params.userId)
  if (params?.courseId) query.set('courseId', params.courseId)
  if (params?.year) query.set('year', params.year)
  if (params?.semester) query.set('semester', params.semester)
  if (params?.page) query.set('page', params.page.toString())
  if (params?.size) query.set('size', params.size.toString())

  const queryString = query.toString()

  const key = params?.syllabusId
    ? `syllabi-${params.syllabusId}`
    : `syllabi-${queryString}`

  const { data: syllabus, status: syllabusStatus, error: syllabusError } = useAsyncData<Syllabus | null>(
    key,
    async () => {
      if (params?.syllabusId) {
        return await $fetch<Syllabus>(`/api/syllabi/syllabi/${params.syllabusId}`)
      }
      return null
    },
    {
      server: false,
    }
  )

  const { data: syllabi, status: syllabiStatus, error: syllabiError } = useAsyncData<Syllabus[]>(
    `syllabi-${queryString}`,
    async () => {
      if (!params?.syllabusId) {
        const data = await $fetch<Syllabus[]>(`/api/syllabi/syllabi?${queryString}`)
        return data
      }
      return []
    },
    {
      server: false,
    }
  )

  return {
    syllabus,
    syllabusStatus,
    syllabusError,
    syllabi,
    syllabiStatus,
    syllabiError
  }
}
