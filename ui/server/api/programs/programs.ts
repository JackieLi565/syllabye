import { Program } from "~/types/types";

export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig();
  const query = getQuery(event);
  const queryString = new URLSearchParams(query as Record<string, string>).toString();
  const apiUrl = `${config.public.apiUrl}/programs${queryString ? `?${queryString}` : ''}`;
  const cookie = getRequestHeader(event, "cookie");

  try {
    const programs = await $fetch<Program[] | null>(apiUrl, {
      headers: {
        cookie: cookie || "",
      },
    });
    return programs ?? null
  } catch (error) {
      console.warn("Programs fetch failed in API route:", error);
      return null;
  }
});
