import type { User } from "~/types/types";

export async function useUser() {
  return useState<User | null>('user');
}

