export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig();
  const { search } = getQuery(event)
  const apiUrl = `${config.public.apiUrl}/users/exists?search=${search}`;
  
  const cookie = getRequestHeader(event, "cookie");

  try {
    const response = await $fetch(`${apiUrl}`, {
      headers: {
        cookie: cookie || "", // 🔁 Forward the user's session cookie
      },
    });
    return response ?? null
  } catch (error) {
      console.warn("Nickname search failed in API route:", error);
      return { error: 'Failed to search for nickname' };
  }
});
