import type { Course } from "~/types/types";

export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig();
  const query = getQuery(event);
  const queryString = new URLSearchParams(query as Record<string, string>).toString();
  const apiUrl = `${config.public.apiUrl}/courses${queryString ? `?${queryString}` : ''}`;

  const cookie = getRequestHeader(event, "cookie");

  try {
    const courses = await $fetch<Course[] | null>(apiUrl, {
      headers: {
        cookie: cookie || "",
      },
    });
    return courses ?? null
  } catch (error) {
      console.warn(`Courses fetch failed in API route:`, error);
      return null;
  }
});
