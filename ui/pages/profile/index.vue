<script setup lang="ts">
import ProfileCard from '~/components/user/ProfileCard.vue';

// definePageMeta({
//   middleware: 'auth'
// })
// const user = useUser()
// useHead({
//   title: `Profile | ${user.value?.nickname}`,
//   meta: [
//     { name: "Syllabye", content: 'Say goodbye to searching for syllabi!' }
//   ],
// })
</script>

<template>
  <div>Profile Page</div>
  <ProfileCard />
</template>