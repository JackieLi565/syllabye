import { Syllabus } from "~/types/types";

export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig();
  const query = getQuery(event);
  const queryString = new URLSearchParams(
    query as Record<string, string>
  ).toString();
  const apiUrl = `${config.public.apiUrl}/syllabi${
    queryString ? `?${queryString}` : ""
  }`;
  const cookie = getRequestHeader(event, "cookie");

  try {
    const syllabi = await $fetch<Syllabus[] | null>(apiUrl, {
      headers: {
        cookie: cookie || "",
      },
    });
    return syllabi ?? null;
  } catch (error) {
    console.warn(`Syllabi fetch failed in API route:`, error);
    return null;
  }
});
