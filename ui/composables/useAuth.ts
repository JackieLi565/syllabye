import type { Session, User } from "~/types/types";

export const useAuth = () => {
  const session = useState<Session | null>('session');
  const user = useState<(User & { newuser?: boolean }) | null>('user');
  return { session, user };
};