<script setup lang="ts">
import { Button } from '@/components/ui/button'
import { Icon } from '@iconify/vue'
</script>

<template>
  <footer class="mt-20 px-4 w-full flex items-center justify-between h-14 border-t">
    <p class="text-xs text-muted-foreground">made with love by 2 tmu students ❤️</p>
    <div class="flex items-center gap-2">
      <NuxtLink to="https://google.com" class="text-xs text-muted-foreground hover:text-primary ease duration-150">terms of service</NuxtLink>
      <NuxtLink to="https://github.com/JackieLi565/syllabye" target="_blank">
        <Button variant="ghost" size="icon">
          <Icon icon="octicon:mark-github-24" class="h-[1.2rem] w-[1.2rem] rotate-0 scale-100 transition-all dark:-rotate-90 dark:scale-0" />
        </Button>
      </NuxtLink>
    </div>
  </footer>
</template>