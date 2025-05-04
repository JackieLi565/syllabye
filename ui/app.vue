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