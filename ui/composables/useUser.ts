import type { User } from "~/types/types";

export function useUser() {
  return useState<User | null>('user');
}

