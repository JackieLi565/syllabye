import type { Course } from "~/types/types";

export const useCourse = (id: string) => {
  const config = useRuntimeConfig();
  const cookie = useCookie(config.public.sessionKey);

  const {
    data: course,
    status,
    error,
    refresh,
  } = useAsyncData(
    `course-${id}`,
    () => {
      return $fetch<Course>(`${config.public.apiUrl}/courses/${id}`, {
        headers: {
          Cookie: cookie.value
            ? `${config.public.sessionKey}=${cookie.value}`
            : "",
        },
        credentials: "include",
      });
    },
    {
      server: true,
    }
  );

  return {
    course,
    status,
    error,
    refresh,
  };
};
