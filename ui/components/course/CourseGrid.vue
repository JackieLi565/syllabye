<script setup lang="ts">
import { useCourses } from "~/composables/useCourses";
import CourseCard from "./CourseCard.vue";

const { courses, status } = useCourses({
  server: false,
});
</script>

<template>
  <div class="p-6 max-w-6xl m-auto">
    <h2 class="text-2xl font-bold mb-6">Undergraduate Courses 🎓</h2>
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
      <CourseCard v-for="course in courses" :key="course.id" :course="course" />
    </div>

    <div v-if="status === 'pending'" class="text-center py-10">
      <p class="text-gray-500">Loading courses...</p>
    </div>
  </div>
</template>
