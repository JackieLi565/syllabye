<script setup lang="ts">
import SyllabyeLogo from '@/components/logo/SyllabyeLogo.vue'
import ProfileDropdown from '../user/ProfileDropdown.vue';
const { session } = useAuth();
</script>

<template>
  <nav class="fixed top-0 bg-secondary/40 backdrop-blur-md px-4 w-full flex items-center justify-between h-14 border-b z-10">
    <NuxtLink to="/" class="motion-preset-slide-right-md">
      <SyllabyeLogo size="small" :slogan="false" class=""/>
    </NuxtLink>
    <div class="flex items-center space-x-2 motion-preset-slide-left-md">
      <div v-if="session" class="space-x-2 mr-2">
        <NuxtLink to="/browse"><Button variant="ghost">Browse</Button></NuxtLink>
        <NuxtLink to="/upload"><Button class="font-semibold">Upload</Button></NuxtLink>
      </div>
      <LoginButton v-if="!session"/>
      <ProfileDropdown v-else/>
    </div>
  </nav>
</template>