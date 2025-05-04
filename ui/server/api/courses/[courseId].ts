import type { Course } from "~/types/types";

export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig();
  const courseId = getRouterParam(event, 'courseId')
  const apiUrl = `${config.public.apiUrl}/courses/${courseId}`;

  const cookie = getRequestHeader(event, "cookie");

  try {
    const course = await $fetch<Course | null>(apiUrl, {
      headers: {
        cookie: cookie || "",
      },
    });
    return course ?? null
  } catch (error) {
      console.warn(`Course ${courseId} fetch failed in API route:`, error);
      return null;
  }
});
