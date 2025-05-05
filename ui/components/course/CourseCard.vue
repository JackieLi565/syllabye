<script setup lang="ts">
import type { Course } from "~/types/types";

const props = defineProps<{
  course: Course;
}>();
</script>

<template>
  <div
    class="relative border border-secondary rounded-lg shadow-md overflow-hidden hover:shadow-lg transition-shadow duration-300"
  >
    <div class="p-5">
      <h1 class="text-lg font-semibold mb-2 text-foreground truncate">
        {{ props.course.title }}
      </h1>

      <div class="mb-3">
        <p class="text-sm text-muted-foreground line-clamp-3">
          {{ course.description || "No description available" }}
        </p>
      </div>
      <div class="flex items-center mt-4">
        <NuxtLink
          :to="'/courses/' + course.id"
          class="text-sm font-medium text-primary hover:text-primary/90 flex items-center transition-color"
        >
          View {{ props.course.course }}
          <svg
            xmlns="http://www.w3.org/2000/svg"
            class="h-4 w-4 ml-1"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M14 5l7 7m0 0l-7 7m7-7H3"
            />
          </svg>
        </NuxtLink>
      </div>
    </div>
  </div>
</template>
