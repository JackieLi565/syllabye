export function useCourses(params?: {
  search?: string
  category?: string
  page?: number
  size?: number
}) {
  const config = useRuntimeConfig()
  const query = new URLSearchParams()

  if (params?.search) query.set('search', params.search)
  if (params?.category) query.set('category', params.category)
  if (params?.page) query.set('page', String(params.page))
  if (params?.size) query.set('size', String(params.size))

  const queryString = query.toString()
  // const endpoint = `${config.public.apiUrl}/courses${queryString ? ``}`

  const { data: courses, status: coursesStatus, error: coursesError } = useAsyncData(
    'courses',
    () => $fetch(`${config.public.apiUrl}/courses`),
    {
      server: true,
      default: () => [],
    }
  )

  const { data: categories, status: categoriesStatus, error: categoriesError } = useAsyncData(
    'course-categories',
    () => $fetch(`${config.public.apiUrl}/courses/categories`),
    {
      server: true,
      default: () => [],
    }
  )

  return {
    courses,
    coursesStatus,
    coursesError,
    categories,
    categoriesStatus,
    categoriesError,
  }
}
