import { useQuery } from "@tanstack/vue-query"

export function useCourses(size?: number, category?: string) {
  const config = useRuntimeConfig()

  const coursesQuery = useQuery({
    queryKey: ['courses', size],
    queryFn: () => $fetch(`${config.public.apiUrl}/courses?size=${size}?category=${category}`)
  })

  return { coursesQuery }
}
