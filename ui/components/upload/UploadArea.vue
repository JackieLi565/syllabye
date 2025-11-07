<script lang="ts" setup>
import { ref, watch } from 'vue'
import { FileText, Upload, X, Check, AlertCircle } from 'lucide-vue-next'
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger
} from '@/components/ui/tooltip'
import { formatFileSize } from '~/lib/utils'
import { Icon } from '@iconify/vue/dist/iconify.js'

const props = defineProps({
  modelValue: {
    type: [File, null],
    default: null
  }
})

const emit = defineEmits(['update:modelValue'])

const fileInputRef = ref<HTMLInputElement | null>(null);
const isDragging = ref(false)
const file = ref(props.modelValue)
const fileError = ref<{ type: string; message: string } | null>(null);

watch(() => props.modelValue, (newVal) => {
  file.value = newVal
})

watch(file, (newVal) => {
  emit('update:modelValue', newVal)
})

const handleFileChange = (e: any) => {
  const selectedFile = e.target.files?.[0] || null
  validateAndSetFile(selectedFile)
}

const validateAndSetFile = (selectedFile: any) => {
  fileError.value = null
  
  if (!selectedFile) {
    file.value = null
    return
  }
  
  const validTypes = ['application/pdf']
  if (!validTypes.includes(selectedFile.type)) {
    fileError.value = {
      type: 'invalid file',
      message: 'We currently only accept PDF files.'
    }
    return
  }
  
  if (selectedFile.size > 50 * 1024 * 1024) {
    fileError.value = {
      type: 'invalid size',
      message: 'File size must be less than 50MB.'
    }
    return
  }
  
  file.value = selectedFile
}

const handleDragOver = (e: any) => {
  e.preventDefault()
  e.stopPropagation()
  isDragging.value = true
}

const handleDragLeave = (e: any) => {
  e.preventDefault()
  e.stopPropagation()
  isDragging.value = false
}

const handleDrop = (e: any) => {
  e.preventDefault()
  e.stopPropagation()
  isDragging.value = false
  
  const droppedFile = e.dataTransfer.files[0]
  validateAndSetFile(droppedFile)
}

const clearFile = () => {
  file.value = null
  fileError.value = null
  if (fileInputRef.value) {
    fileInputRef.value.value = ''
  }
}

const triggerFileInput = () => {
  fileInputRef.value?.click()
}

const getFileIcon = () => {
  if (!file.value) return null
  
  const fileType = file.value.type
  if (fileType === 'application/pdf') {
    return 'pdf'
  } 
  return 'other'
}
</script>

<template>
  <div class="w-full">
    <label class="block text-sm font-medium  mb-2">File Upload</label>
  
    <input
      ref="fileInputRef"
      type="file"
      class="hidden"
      @change="handleFileChange"
    />
    
    <div v-if="fileError" class="flex items-center gap-2 text-destructive text-sm mb-3">
      <AlertCircle class="h-4 w-4" />
      <template v-if="fileError.type === 'invalid file'">
        {{ fileError.message }}
        <NuxtLink class="underline hover:text-destructive/90 ease duration-150" href="https://www.freepdfconvert.com/" target="blank">Try converting your file before uploading.</NuxtLink>
      </template>
    </div>
    
    <div
      v-if="!file"
      :class="[
        'border-2 border-dashed rounded-lg transition-all',
        isDragging ? 'border-primary bg-secondary/60' : 'border-muted hover:border-primary',
      ]"
      @dragover="handleDragOver"
      @dragleave="handleDragLeave"
      @drop="handleDrop"
      @click="triggerFileInput"
    >
      <div class="flex flex-col items-center justify-center py-8 px-4 cursor-pointer">
        <div class="rounded-full bg-primary-alternative p-3 mb-3">
          <Upload class="h-6 w-6 text-primary" />
        </div>
        <p class="text-lg font-medium">Click to browse files</p>
        <p class="text-sm text-muted-foreground mb-5">or drag and drop your syllabus here</p>
        <p class="text-xs text-muted-foreground">Accepts PDF (max 50MB)</p>
      </div>
    </div>
    
    <div
      v-else
      class="border border-border rounded-lg bg-secondary/60 p-4 flex items-center motion-opacity-in-0 motion-translate-y-in-25 motion-blur-in-sm"
    >
      <div class="h-12 w-12 rounded-lg bg-primary/40 flex items-center justify-center mr-4">
        <Icon icon="prime:file-pdf" width="36" height="36" v-if="getFileIcon() === 'pdf'" class="text-foreground"/>
        <FileText v-else class="h-6 w-6 text-primary" />
      </div>
      
      <div class="flex-1 min-w-0">
        <div class="flex space-x-1 items-center">
          <p class="text-sm font-medium truncate">{{ file.name }}</p>
          <TooltipProvider>
            <Tooltip>
              <TooltipTrigger as-child><Check class="h-4 w-4 text-success mr-2 hover:cursor-pointer outline-none"/></TooltipTrigger>
              <TooltipContent>
                <p>File is valid</p>
              </TooltipContent>
            </Tooltip>
          </TooltipProvider>
        </div>
        <p class="text-xs text-muted-foreground">{{ formatFileSize(file.size) }}</p>
      </div>
      
      <TooltipProvider>
        <Tooltip>
          <TooltipTrigger as-child>
            <Button
              size="icon"
              icon
              variant="ghost"
              @click.stop="clearFile"
            >
              <X />
            </Button>
          </TooltipTrigger>
          <TooltipContent>
            <p>Remove file</p>
          </TooltipContent>
        </Tooltip>
      </TooltipProvider>
    </div>
  </div>
</template>