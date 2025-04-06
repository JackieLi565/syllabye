import { Program } from "~/types/types";

export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig();
  const apiUrl = `${config.public.apiUrl}/programs`;

  const cookie = getRequestHeader(event, "cookie");

  try {
    const programs = await $fetch<Program[] | null>(apiUrl, {
      headers: {
        cookie: cookie || "", // 🔁 Forward the user's session cookie
      },
    });
    return programs ?? null
  } catch (error) {
      console.warn("Programs fetch failed in API route:", error);
      return null;
  }
});
