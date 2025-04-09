<script lang="ts" setup>
import { Skeleton } from '@/components/ui/skeleton'
import { ref, onMounted } from 'vue'

const props = defineProps<{
  rows: number
}>()

</script>

<template>
  <div class="flex flex-col space-y-5">
    <Skeleton class="h-10 w-full rounded-xl" />
    <div class="space-y-3.5">
      <Skeleton
        v-for="row in rows"
        :key="row"
        :class="`h-6 rounded-xl w-full`"
      />
    </div>
  </div>
</template>
