import type { Course } from "~/types/types";

interface Options {
  server: boolean;
  query?: {
    search?: string;
    category?: string;
    page?: string;
    size?: string;
  };
}

export const useCourses = (options: Options) => {
  const config = useRuntimeConfig();
  const cookie = useCookie(config.public.sessionKey);

  const query = new URLSearchParams();
  if (options.query)
    Object.entries(options.query)
      .filter(([_, value]) => value)
      .forEach(([key, value]) => query.set(key, value));

  const {
    data: courses,
    status,
    error,
    refresh,
  } = useAsyncData(
    "courses",
    () => {
      return $fetch<Course[]>(
        `${config.public.apiUrl}/courses?${query && query.toString()}`,
        {
          headers: {
            Cookie:
              options.server && cookie.value
                ? `${config.public.sessionKey}=${cookie.value}`
                : "",
          },
          credentials: "include",
        }
      );
    },
    {
      server: options.server,
      default: () => [],
    }
  );

  return {
    courses,
    status,
    error,
    refresh,
  };
};
