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
  // try {
  //   const user = await $fetch<User | null>(`${config.public.apiUrl}/users/${userId}`, {
  //     credentials: "include",
  //   });
  //   console.log(user)

  //   return user;
  // } catch (error: any) {
  //   if (error.response?.status === 401) {
  //     console.warn("Not authorized to access user, check session");
  //   } else {
  //     console.error("Error fetching user:", error);
  //   }
  //   return null;
  // }
});
