import type { Category } from "~/types/types";

interface Options {
  server: boolean;
  query?: {
    search?: string;
  };
}

export const useCategories = (options: Options) => {
  const cookieKey = "syllabye.session";
  const config = useRuntimeConfig();
  const cookie = useCookie(cookieKey);

  const query = new URLSearchParams(options.query);

  const {
    data: categories,
    status,
    error,
  } = useAsyncData(
    "course-categories",
    () => {
      return $fetch<Category[]>(
        `${config.public.apiUrl}/courses/categories${
          query && query.toString()
        }`,
        {
          headers: {
            Cookie:
              options.server && cookie.value
                ? `${cookieKey}=${cookie.value}`
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
    categories,
    status,
    error,
  };
};
