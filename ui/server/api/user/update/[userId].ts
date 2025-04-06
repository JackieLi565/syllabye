import type { User } from "~/types/types";

export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig();
  const userId = getRouterParam(event, 'userId')
  const body = await readBody(event);
  const cookie = getRequestHeader(event, "cookie");
  console.log(JSON.stringify(body))
  try {
    const userPatch = await $fetch<User | null>(`${config.public.apiUrl}/users/${userId}`, {
      method: 'PATCH',
      headers: {
        cookie: cookie || "",
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(body)
    });

    return userPatch ?? null;
  } catch (error) {
    console.error("User patch failed:", error);
    return null;
  }
});
