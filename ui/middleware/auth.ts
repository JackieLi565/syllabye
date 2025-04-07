// export default defineNuxtRouteMiddleware((to) => {
//   const session = useNuxtApp().$session.value

//   if (!session) {
//     return navigateTo(`/?redirect=${encodeURIComponent(to.fullPath)}`, { replace: true })
//   }
// })