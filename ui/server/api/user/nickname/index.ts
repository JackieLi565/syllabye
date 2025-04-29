export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig();
  const { search } = getQuery(event)
  const apiUrl = `${config.public.apiUrl}/users/exists?search=${search}`;
  
  const cookie = getRequestHeader(event, "cookie");

  try {
    const response = await $fetch(`${apiUrl}`, {
      headers: {
        cookie: cookie || "",
      },
    });
    console.log(`nickname being searched: ${search}`, response)
    return response ?? null
  } catch (error) {
      console.warn("Nickname search failed in API route:", error);
      return { error: 'Failed to search for nickname' };
  }
});
