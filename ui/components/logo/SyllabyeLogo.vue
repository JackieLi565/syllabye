<script setup lang="ts">
import clsx from 'clsx/lite';
defineProps<{
  size: "small" | "medium" | "large",
  slogan: boolean,
}>()
</script>

<template>
  <div class="-space-y-1 text-left select-none">   
    <h1 :class="clsx('tracking-wider font-display bg-gradient-to-r from-foreground to-primary text-transparent bg-clip-text', size==='small' && 'text-md', size==='medium' && 'text-lg', size==='large' && 'text-2xl')">syllabye</h1>
    <p :class="clsx('font-light text-muted-foreground', slogan ? 'block' : 'hidden', size==='small' && 'text-[10px]', size==='medium' && 'text-xs', size==='large' && 'text-sm')">say goodbye to searching for syllabi!</p>
  </div>
</template>