import type { Session } from "~/types/types";

export default defineNuxtPlugin({
  name: 'session',
  setup (nuxtApp) {
    const session = useState<Session | null>('session', () => null);
    useAsyncData('session', async () => {
      const event = useRequestEvent();

      try {
        const sessionData = await $fetch<Session | null>("/api/auth/session", {
          headers: event?.node.req.headers.cookie
            ? { cookie: event.node.req.headers.cookie }
            : undefined,
          credentials: "include",
        });

        session.value = sessionData ?? null;
        return sessionData;
      } catch (e) {
        console.error('Failed to fetch session:', e);
        session.value = null;
        return null;
      }
    });
  }
});
