import type { Session, User } from "~/types/types";

export default defineNuxtPlugin(nuxtApp => {
  const session = useState<Session | null>('session');
  const user = useState<(User & { newuser?: boolean }) | null>('user', () => null);

  const fetchUser = async () => {
    if (session.value?.userId) {
      const event = useRequestEvent();
      try {
        const userData = await $fetch<User | null>(`/api/auth/user/${session.value.userId}`, {
          headers: event?.node.req.headers.cookie
            ? { cookie: event.node.req.headers.cookie }
            : undefined,
          credentials: "include",
        });

        // If userData exists and has no nickname, mark as new user
        const isNewUser = !!userData && !userData.nickname;
        user.value = userData
          ? { ...userData, newuser: isNewUser }
          : null;

      } catch (e) {
        console.error('Failed to fetch user due to no session:', e);
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
