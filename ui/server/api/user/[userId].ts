import type { User } from "~/types/types";

export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig();
  const userId = getRouterParam(event, 'userId')
  const cookie = getRequestHeader(event, "cookie");

  try {
    const user = await $fetch<User | null>(`${config.public.apiUrl}/users/${userId}`, {
      headers: {
        cookie: cookie || "",
      },
    });

    return user ?? null;
  } catch (error) {
    console.error("User fetch failed:", error);
    return null;
  }
});
