import type { Session, User } from "~/types/types";

export default defineNuxtPlugin(nuxtApp => {
  const session = useState<Session | null>('session');
  const user = useState<(User & { newuser?: boolean }) | null>('user', () => null);

  const fetchUser = async () => {
    if (session.value?.userId) {
      try {
        let headers: Record<string, string> | undefined;
  
        if (import.meta.server) {
          const event = useRequestEvent();
          if (event?.node.req.headers.cookie) {
            headers = { cookie: event.node.req.headers.cookie };
          }
        }
  
        const userData = await $fetch<User | null>(`/api/user/${session.value.userId}`, {
          headers,
          credentials: "include",
        });
  
        const isNewUser = !!userData && !userData.nickname;
        user.value = userData ? { ...userData, newuser: isNewUser } : null;
  
      } catch (e) {
        console.error('Failed to fetch user:', e);
        user.value = null;
      }
    }
  };  

  watchEffect(() => {
    if (session.value) {
      fetchUser();
    }
  });
});
