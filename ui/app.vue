<script setup lang="ts">
import Toaster from '@/components/ui/toast/Toaster.vue'

useHead({
  title: 'Syllabye',
  meta: [
    { name: "Syllabye", content: 'Say goodbye to searching for syllabi!' }
  ],
})
</script>

<template>
  <Toaster />
  <NuxtLoadingIndicator color="#6f20d5"/>
  <NuxtLayout>
    <NuxtPage />
  </NuxtLayout>
</template>

<style>
.page-enter-active,
.page-leave-active {
  transition: all 0.4s;
}
.page-enter-from,
.page-leave-to {
  opacity: 0;
  filter: blur(1rem);
}
</style>