import type { Session } from '~/types/types';

export function useSession() {
  return useState<Session | null>('session');
}