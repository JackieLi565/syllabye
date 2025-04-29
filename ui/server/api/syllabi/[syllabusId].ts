import { Syllabus } from "~/types/types";

export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig();
  const syllabusId = getRouterParam(event, 'syllabusId')
  const apiUrl = `${config.public.apiUrl}/syllabi/${syllabusId}`;

  const cookie = getRequestHeader(event, "cookie");

  try {
    const syllabus = await $fetch<Syllabus | null>(apiUrl, {
      headers: {
        cookie: cookie || "",
      },
    });
    return syllabus ?? null
  } catch (error) {
      console.warn(`Syllabus ${syllabusId} fetch failed in API route:`, error);
      return null;
  }
});
