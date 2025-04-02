<script setup lang="ts">
const { $session, $user } = useNuxtApp()
// console.log($session)
// console.log($user)

</script>

<template>
  <div class="flex flex-col justify-start items-center px-4 md:px-0">
    <div v-if="$session">
      <p>Welcome, {{ $user.fullname }}!</p>
      <p>Your email: {{ $user.email }}</p>
    </div>
    <div v-else>
      <p>You are not logged in.</p>
    </div>
    <div class="h-[250px] md:h-[220px] pb-8 flex flex-col justify-end items-center w-full text-center">
      <h1 class="text-2xl md:text-5xl font-semibold font-geist bg-gradient-to-r from-foreground to-violet-400 text-transparent bg-clip-text pb-1">All your TMU course syllabi in one place</h1>
      <p class="text-muted-foreground text-sm md:text-xl">Because searching for a course syllabus on Reddit shouldn't be a thing</p>
    </div>
    <SearchBar />
    <div class="w-full md:w-1/2 mt-20 space-y-4">
      <h1 class="font-semibold md:text-xl">Popular syllabi</h1>
      <div class="w-full flex flex-col md:flex-row gap-2">
        <CourseCard/>
        <CourseCard/>
        <CourseCard/>
      </div>
    </div>
    <!-- <LoginCard /> -->
  </div>
</template>