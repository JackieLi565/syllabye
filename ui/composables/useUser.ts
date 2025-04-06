


// // CLIENT SIDE SESSION FETCHING SOLUTION

import type { Session, User } from "~/types/types";

// interface Session {
//   userId: string;
// }

// interface User {
//   id: string;
//   currentYear?: number | null;
//   email: string;
//   fullname: string;
//   gender?: "male" | "female" | "other";
//   nickname?: string;
//   picture: string;
//   programId?: string;
// }

// export function useAuth() {
//   const config = useRuntimeConfig();

//   const session = useState<Session | null>("session", () => null);
//   const user = useState<User | null>("user", () => null);
//   const showLoginModal = useState<boolean>("showLoginModal", () => false);
//   const isFetching = ref(false); // Prevents duplicate fetch calls

//   function openLoginModal() {
//     showLoginModal.value = true;
//   }

//   function closeLoginModal() {
//     showLoginModal.value = false;
//   }

//   async function fetchSession() {
//     if (isFetching.value) return;
//     isFetching.value = true;

//     try {
//       const data = await $fetch<Session | null>(`${config.public.apiUrl}/me`, {
//         method: "GET",
//         credentials: "include",
//       });

//       if (!data) {
//         console.warn("Session fetch failed. Logging out...");
//         session.value = null;
//         user.value = null;
//         return;
//       }

//       session.value = data;

//       if (session.value) {
//         await fetchUser(session.value.userId);
//       }
//     } catch (error: any) {
//       if (error.response?.status === 401) {
//         console.warn("Logged out, not authorized to access session");
//       } else {
//         console.error("Error fetching session:", error);
//       }
  
//       session.value = null;
//       user.value = null;
//     } finally {
//       isFetching.value = false;
//     }
//   }

//   async function fetchUser(userId: string) {
//     console.log('Fetching user: ${userId}');

//     try {
//       const data = await $fetch<User | null>(`${config.public.apiUrl}/users/${userId}`, {
//         method: "GET",
//         credentials: "include",
//       });

//       if (!data) {
//         console.warn("User fetch failed.");
//         user.value = null;
//         return;
//       }

//       console.log("User response:", data);
//       user.value = data;
//     } catch (error: any) {
//       if (error.response?.status === 401) {
//         console.warn("Not authorized to access data due to expired session");
//       } else {
//         console.error("Error fetching user:", error);
//       }
//     }
//   }

//   if (!session.value) {
//     fetchSession();
//   }

//   return {
//     session,
//     user,
//     showLoginModal,
//     openLoginModal,
//     closeLoginModal,
//     fetchSession,
//   };
// }

// composables/useUser.ts
export async function useUser() {
  return useState<User | null>('user');
  // const session = await useSession()
  // const user = useState<User | null>('user', () => null)
  // const event = useRequestEvent()

  // // If user already exists in state, return early
  // if (user.value) return user

  // if (!session.value) {
  //   console.warn("No session available, skipping user fetch")
  //   return null
  // }

  // try {
  //   const userData = await $fetch<User>(`/api/auth/user/${session.value.userId}`, {
  //     headers: event?.node.req.headers.cookie
  //       ? { cookie: event.node.req.headers.cookie }
  //       : undefined,
  //     credentials: "include",
  //   })

  //   user.value = userData
  // } catch (error: any) {
  //   if (error.response?.status === 401) {
  //     console.warn("Not authorized to access user data")
  //   } else {
  //     console.error("Error fetching user:", error)
  //   }
  //   user.value = null
  //   return null
  // }

  // return user
}

