<script setup lang="ts">
import { useCourse } from "~/composables/useCourse";

const router = useRouter();

const route = useRoute();
const { course } = useCourse(String(route.params.courseId));

const currentYear = new Date().getFullYear();
const uriPrefix = `https://www.torontomu.ca/calendar/${currentYear}-${
  currentYear + 1
}/courses`;
</script>

<template>
  <div v-if="!course">
    <div
      class="flex flex-col items-center justify-center py-12 px-4 sm:px-6 lg:px-8 min-h-[400px] bg-gray-50 rounded-lg border border-gray-200"
    >
      <div class="text-center">
        <h2 class="text-2xl font-bold text-gray-900 mb-2">
          Course Not Available
        </h2>
        <p class="text-gray-600 mb-6 max-w-md mx-auto">
          Looks like we couldn't find this course.
        </p>

        <div class="flex flex-col sm:flex-row gap-3 justify-center">
          <button
            @click="router.back"
            class="inline-flex items-center justify-center px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary"
          >
            <ArrowLeft class="h-4 w-4 mr-2" />
            Go Back
          </button>
          <button
            class="inline-flex items-center justify-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-primary hover:bg-primary/90 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary"
          >
            Request Course
          </button>
        </div>
      </div>
    </div>
  </div>

  <div v-else class="py-6 max-w-6xl m-auto">
    <!-- Course Header -->
    <header class="mb-2">
      <div class="container mx-auto px-4 py-8">
        <div class="flex flex-col md:flex-row md:items-center justify-between">
          <div>
            <div class="inline-flex items-center mb-2">
              <span
                class="text-sm font-medium bg-primary text-white px-2 py-1 rounded"
                >{{ course.course }}</span
              >
            </div>
            <h1 class="text-2xl font-bold">
              {{ course.title }}
            </h1>
          </div>
        </div>
      </div>
    </header>

    <div class="space-y-6">
      <section class="mx-auto px-4 py-8">
        <div class="max-w-4xl">
          <h2 class="text-xl font-semibold text-gray-800 mb-4">
            Course Overview
          </h2>
          <p class="text-gray-700 leading-relaxed">
            {{ course.description }}
          </p>
        </div>
      </section>

      <section class="mx-auto px-4 py-8">
        <div class="pt-6">
          <h2 class="text-xl font-semibold text-gray-800 mb-6">
            Available Syllabi
          </h2>

          <div class="bg-gray-50 rounded-lg p-6 text-center">
            <p class="text-gray-600">
              No syllabi are currently available for this course.
            </p>
            <p class="text-sm text-gray-500 mt-2">
              Check back later or upload one yourself!
            </p>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>
