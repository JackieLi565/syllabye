<script setup lang="ts">
import SkeletonPage from '~/components/loading/SkeletonPage.vue';
import type { Session, User } from '~/types/types';
</script>

<template>
  <div class="flex flex-col min-h-screen">
    <NavigationBar/>
    <NuxtLoadingIndicator :color="true" :height="5" />
    <div class="md:mt-14 flex-grow">
      <template>
        <!-- <div v-if="!isAuthResolved">Loading...</div>
        <div v-else-if="session">Welcome back! {{ user?.fullname }}</div>
        <div v-else>Please log in</div> -->
      </template>
      <slot />
    </div>
    <Footer/>
  </div>
</template>