import type { Course } from "~/types/types"

export function useCourses(params?: {
  search?: string
  category?: string
  page?: number
  size?: number
  courseId?: string
}) {
  const query = new URLSearchParams()

  if (params?.search) query.set('search', params.search)
  if (params?.category) query.set('category', params.category)
  if (params?.page) query.set('page', params.page.toString())
  if (params?.size) query.set('size', params.size.toString())
  if (params?.courseId) query.set('courseId', params.courseId)

  const queryString = query.toString()

  const key = params?.courseId
    ? `courses-${params.courseId}`
    : `courses-${queryString}`

  const { data: course, status: courseStatus, error: courseError } = useAsyncData(
    key,
    async () => {
      if (params?.courseId) {
        return await $fetch<Course>(`/api/courses/${params.courseId}`)
      }
      return null
    },
    {
      server: true,
    }
  )

  const { data: courses, status: coursesStatus, error: coursesError } = useAsyncData(
    key,
    async () => {
      if (!params?.courseId) {
        return await $fetch<Course[]>(`/api/courses/courses${queryString ? `?${queryString}` : ''}`)
      }
      return []
    },
    {
      server: true,
    }
  )

  return {
    courses,
    coursesStatus,
    coursesError,
    course,
    courseStatus,
    courseError,
  }
}
