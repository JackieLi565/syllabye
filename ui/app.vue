<script setup lang="ts">
useHead({
  title: 'Syllabye',
  meta: [
    { name: "Syllabye", content: 'Say goodbye to searching for syllabi!' }
  ],
})
</script>

<template>
  <NuxtLayout>
    <NuxtPage />
  </NuxtLayout>
</template>

<style>
/* Page transitions; remove if no longer want blur transition when navigating pages */
.page-enter-active,
.page-leave-active {
  transition: all 150ms;
}

.page-enter-from,
.page-leave-to {
  opacity: 0;
  filter: blur(0.4rem);
}
</style>