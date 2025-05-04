<script setup lang="ts">
import MainSearchBar from '~/components/search/MainSearchBar.vue'
const { session } = useAuth()
const route = useRoute()
const { openLoginModal } = useLoginModal()

onMounted(() => {
  if (route.query.redirect) {
    openLoginModal()
  }
})

function handleProtectedClick() {
  if (!session.value) {
    openLoginModal()
    return
  }
}

</script>

<template>
  <div class="flex flex-col justify-start items-center px-4 md:px-0 motion-preset-slide-up-sm">
    <!-- background div -->
    <div class="-z-10 -mt-20 absolute h-full w-full bg-[radial-gradient(hsl(var(--border))_1px,transparent_1px)] dark:bg-[radial-gradient(hsl(var(--muted))_1px,transparent_1px)] [background-size:16px_16px] [mask-image:radial-gradient(ellipse_50%_50%_at_50%_50%,#fff_70%,transparent_100%)]"></div>
    <div class="h-[250px] md:h-[220px] pb-8 flex flex-col justify-end items-center w-full text-center">
      <h1 class="text-2xl md:text-5xl font-semibold font-inter bg-gradient-to-r from-foreground to-violet-400 text-transparent bg-clip-text pb-1">All your TMU course syllabi in one place</h1>
      <p class="text-muted-foreground text-sm md:text-xl">Because searching for a course syllabus on Reddit shouldn't be a thing</p>
    </div>
    <MainSearchBar @click="handleProtectedClick"/>
    <div class="w-full md:w-1/2 mt-20 space-y-4">
      <h1 class="font-semibold md:text-xl">Popular syllabi</h1>
      <div class="w-full flex flex-col md:flex-row gap-2">
        <CourseCard @click="handleProtectedClick"/>
        <CourseCard @click="handleProtectedClick"/>
        <CourseCard @click="handleProtectedClick"/>
      </div>
    </div>
    <RedditEmbeds/>
  </div>
</template>


