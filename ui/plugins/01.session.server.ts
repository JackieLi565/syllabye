import type { Session } from "~/types/types";

export default defineNuxtPlugin(async (nuxtApp) => {
  const session = useState<Session | null>('session', () => null);

  const event = import.meta.server ? useRequestEvent() : null;

  const headers: HeadersInit | undefined = event?.node.req.headers.cookie
    ? { cookie: event.node.req.headers.cookie }
    : undefined;

  try {
    const { data: sessionData } = await useAsyncData('session', () =>
      $fetch<Session | null>("/api/auth/session", {
        headers,
        credentials: "include",
      })
    );

    session.value = sessionData.value ?? null;

    return {
      provide: {
        session
      }
    };
  } catch (e) {
    console.error('Failed to fetch session:', e);
    session.value = null;

    return {
      provide: {
        session
      }
    };
  }
});
