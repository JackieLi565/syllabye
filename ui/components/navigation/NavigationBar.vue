<script setup lang="ts">
import SyllabyeLogo from '@/components/logo/SyllabyeLogo.vue'
</script>

<template>
  <nav class="fixed top-0 bg-secondary/40 backdrop-blur-md px-4 w-full flex items-center justify-between h-14 border-b">
    <SyllabyeLogo size="small" :slogan="false"/>
    <div class="flex items-center space-x-2">
      <LoginButton />
      <ThemeButton />
    </div>
  </nav>
</template>