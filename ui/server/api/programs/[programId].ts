import { Program } from "~/types/types";

export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig();
  const programId = getRouterParam(event, 'programId')
  const apiUrl = `${config.public.apiUrl}/programs/${programId}`;

  const cookie = getRequestHeader(event, "cookie");

  try {
    const programs = await $fetch<Program[] | null>(apiUrl, {
      headers: {
        cookie: cookie || "",
      },
    });
    return programs ?? null
  } catch (error) {
      console.warn(`Program ${programId} fetch failed in API route:`, error);
      return null;
  }
});
