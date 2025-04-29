<template>
  <footer class="mt-20 px-4 w-full flex items-center justify-between h-14 border-t">
    <p class="text-xs text-muted-foreground">made with ❤️ by 2 tmu students</p>
    <div class="flex items-center gap-2">
      <NuxtLink to="/tos" class="text-xs text-muted-foreground hover:text-primary ease duration-150">terms of service</NuxtLink>
      <NuxtLink to="/about" class="text-xs text-muted-foreground hover:text-primary ease duration-150">about</NuxtLink>
      <ThemeButton />
    </div>
  </footer>
</template>