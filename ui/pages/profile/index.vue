<script setup lang="ts">
import ProfileCard from '~/components/user/ProfileCard.vue'

definePageMeta({
  middleware: 'auth'
})
const { user } = useAuth()
useHead({
  title: `${user.value?.nickname} | Profile`,
  meta: [
    { name: "Syllabye", content: 'Say goodbye to searching for syllabi!' }
  ],
})
</script>

<template>
  <div class="md:py-8">
    <ProfileCard />
  </div>
</template>