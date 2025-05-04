<script setup lang="ts">
import { ref, watch } from 'vue'
import { useDebounceFn } from '@vueuse/core'
import type { Course } from '~/types/types'
import { Input } from '@/components/ui/input'
import { Search } from 'lucide-vue-next'

const search = ref('')
const debouncedSearch = ref('')
const courses = ref<Course[]>([])
const selectedCourseId = ref<string>('')
const selected = ref<string | null>(null)
const props = defineProps<{ selected?: string }>()
defineExpose({ reset })

watch(() => props.selected, (val) => {
  selectedCourseId.value = val ?? ''
}, { immediate: true })

const emit = defineEmits<{
  (e: 'update:selected', value: string): void
}>()

const fetchCourses = async () => {
  const { courses: result } = useCourses({ search: debouncedSearch.value, size: 8 })
  watch(result, (newCourses) => {
    courses.value = newCourses || []
  }, { immediate: true })
}

const debouncedFetch = useDebounceFn(() => {
  debouncedSearch.value = search.value.trim()
  fetchCourses()
}, 500)

watch(search, () => {
  if (!search.value || !selectedCourseId) {
    courses.value = []
  } else {
    debouncedFetch()
  }
})

const isFocused = ref(false)

function handleFocus() {
  isFocused.value = true
}

function handleBlur() {
  setTimeout(() => {
    isFocused.value = false
  }, 100)
}

function handleCourseSelect(course: Course) {
  console.log(course.title, course.id)
  selectedCourseId.value = course.id ?? ''
  emit('update:selected', selectedCourseId.value)
  search.value = `${course.course} ${course.title}`
  isFocused.value = false
}

function reset() {
  selected.value = null
  search.value = ''
  debouncedSearch.value = ''
}

</script>

<template>
  <div>
    <div class="relative w-full items-center">
      <Input 
        id="search" 
        type="text" 
        placeholder="Search for a course..." 
        class="pl-10"
        v-model="search"
        @focus="handleFocus"
        @blur="handleBlur"
      />
      <span class="absolute start-0 inset-y-0 flex items-center justify-center px-2">
        <Search class="size-6 text-muted-foreground" />
      </span>
    </div>

    <ul
      v-if="courses?.length && isFocused"
      class="absolute border border-border bg-background rounded-xl p-1 text-sm w-[385px] mt-1 motion-preset-fade"
    >
      <li
        v-for="course in courses"
        :key="course.id"
        class="rounded-lg p-2 space-x-2 hover:bg-secondary ease duration-150 cursor-pointer"
        @click="handleCourseSelect(course)"
      >
        <span class="font-semibold">{{ course.course }}</span>
        <span class="font-light">{{ course.title }}</span>
      </li>
    </ul>
  </div>
</template>