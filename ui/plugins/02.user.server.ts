import type { Session, User } from "~/types/types";

export default defineNuxtPlugin(nuxtApp => {
  const session = useState<Session | null>('session');
  const user = useState<(User & { newuser?: boolean }) | null>('user', () => null);
  const hasFetched = ref(false);

  const fetchUser = async () => {
    if (session.value?.userId) {
      try {
        const { data: userData } = await useAsyncData('user', async () => {
          if (!session.value?.userId) return null;

          const event = import.meta.server ? useRequestEvent() : null;

          const headers: HeadersInit | undefined = event?.node.req.headers.cookie
            ? { cookie: event.node.req.headers.cookie }
            : undefined;

          const user = await $fetch<User | null>(`/api/user/${session.value.userId}`, {
            headers,
            credentials: "include"
          });

          return user ? { ...user, newuser: !user.nickname } : null;
        });

        user.value = userData.value;

      } catch (e) {
        console.error('Failed to fetch user:', e);
        user.value = null;
      }
    }
  };

  watchEffect(async () => {
    if (session.value && !hasFetched.value) {
      hasFetched.value = true;
      await fetchUser();
    }
  });
});
