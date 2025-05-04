<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted } from 'vue'
import { Input } from '@/components/ui/input'
import { Search } from 'lucide-vue-next'
import Typewriter from '@btjspr/vue-typewriter'

const placeholderTexts = [
  'CHY599',
  'CPS109',
  'RTA201',
  'PHL214',
  'CPS420',
  'MTH201',
]

const inputValue = ref('')
const inputFocused = ref(false)
const inputRef = ref(null)
const placeholder = ref(placeholderTexts[Math.floor(Math.random() * placeholderTexts.length)])
const placeholderKey = ref(0)
const placeholderVisible = ref(true)

const handleFocus = () => {
  inputFocused.value = true
}

const handleBlur = () => {
  inputFocused.value = false
}

watch([inputValue, inputFocused], ([value, focused]) => {
  placeholderVisible.value = !value && !focused
})

let interval: ReturnType<typeof setInterval> | null = null

onMounted(() => {
  interval = setInterval(() => {
    placeholder.value = placeholderTexts[Math.floor(Math.random() * placeholderTexts.length)]
    placeholderKey.value++
  }, 4000)
})

onUnmounted(() => {
  if (interval) clearInterval(interval)
})
</script>

<template>
  <div class="relative w-full max-w-md">
    <Input 
      id="search" 
      ref="inputRef"
      v-model="inputValue"
      type="text" 
      placeholder=""
      @focus="handleFocus"
      @blur="handleBlur"
      class="pl-10 py-5 relative z-10 bg-background" 
    />
    
    <span class="absolute start-0 inset-y-0 flex items-center justify-center pl-3 z-20">
      <Search class="size-4 text-muted-foreground" />
    </span>
    
    <div 
      v-show="placeholderVisible" 
      class="absolute inset-y-0 left-0 flex items-center pl-10 pointer-events-none z-10"
    >
      <Typewriter 
        :key="placeholderKey"
        :speed="300" 
        :delay="0"
        :loop="true"
        class="text-muted-foreground"
        :cursorStyles="{ width: '1px', height: '0.8em' }"
      >
        <span>{{ placeholder }}</span>
      </Typewriter>
    </div>
  </div>
</template>
