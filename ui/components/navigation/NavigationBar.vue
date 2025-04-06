<script setup lang="ts">
import SyllabyeLogo from '@/components/logo/SyllabyeLogo.vue'
const user = await useUser();
const session = await useSession();
</script>

<template>
  <nav class="fixed top-0 bg-secondary/40 backdrop-blur-md px-4 w-full flex items-center justify-between h-14 border-b">
    <SyllabyeLogo size="small" :slogan="false"/>
    <div class="flex items-center space-x-2">
      <LoginButton v-if="!session"/>
      <h1 v-else>{{ user?.fullname }}</h1>
      <ThemeButton />
    </div>
  </nav>
</template>