import { Session } from "~/types/types";

export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig();
  const apiUrl = `${config.public.apiUrl}/me`;

  const cookie = getRequestHeader(event, "cookie");

  try {
    const session = await $fetch<Session | null>(apiUrl, {
      headers: {
        cookie: cookie || "", // 🔁 Forward the user's session cookie
      },
    });

    return session ?? null;
  } catch (error) {
    console.warn("Session fetch failed in API route:", error);
    return null;
  }
});
