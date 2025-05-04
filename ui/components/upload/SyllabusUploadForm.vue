<script setup lang="ts">
import * as z from 'zod'
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { ref } from 'vue' // Ensure ref is imported

import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'

import CourseSearchBar from '../search/CourseSearchBar.vue'
import UploadArea from './UploadArea.vue'
import { Separator } from '@/components/ui/separator'
import type { SyllabusUploadData } from '~/types/types'

const courseSearchRef = ref<InstanceType<typeof CourseSearchBar> | null>(null)
const uploadedFile = ref<File | null>(null)

const schema = toTypedSchema(z.object({
  courseId: z.string({ message: 'Please pick a course from the search results' }).uuid({ message: 'Invalid course' }),
  semester: z.string({ message: 'Semester is required'}),
  year: z.coerce.number({ message: 'Year is required' }).min(1993, { message: 'Earliest year is 1993' }).max(2025, { message: 'Latest year is 2025' }),
  syllabus: z.instanceof(File, { message: 'Please upload a syllabus file.' }),
}))

const form = useForm({
  validationSchema: schema,
  initialValues: {
    syllabus: undefined as File | undefined,
    courseId: '',
    semester: '',
    year: undefined,
  },
})

const { handleSubmit, resetForm, setFieldValue } = form

watch(uploadedFile, (newFile) => {
  setFieldValue('syllabus', newFile as File | undefined)
})

const onSubmit = handleSubmit(async (values) => {
  const { uploadSyllabus } = useSyllabusUpload()

  const payload = {
    courseId: values.courseId,
    semester: values.semester,
    year: values.year,
  }

  const file = form.values.syllabus as File
  const result = await uploadSyllabus(payload as SyllabusUploadData, file)
  console.log(result)
  // if (result.success) {
  //   console.log('Uploaded successfully')
  //   resetForm()
  //   courseSearchRef.value?.reset()
  //   uploadedFile.value = null
  // } else {
  //   console.error('Upload failed:', result.errorText)
  // }
})
</script>


<template>
  <form @submit="onSubmit">
    <div class="flex w-full justify-center gap-6">
      <div class="space-y-4 w-full md:w-3/4 xl:w-1/2">
        <div>
          <h1 class="text-xl font-semibold">Upload Syllabus</h1>
          <p class="text-sm text-muted-foreground">Share a course syllabus with us TMU students ❤️</p>
        </div>

        <div class="border rounded-lg w-full p-6">
          <FormField v-slot="{ componentField }" name="syllabus">
            <FormItem>
              <FormControl>
                <UploadArea v-model="uploadedFile" />
              </FormControl>
              <FormMessage />
            </FormItem>
          </FormField>

          <Separator class="my-5" />

          <section class="flex flex-col gap-y-2 md:flex-row md:justify-between md:gap-x-2">
            <FormField v-slot="{ componentField }" name="courseId">
              <FormItem class="w-full md:w-2/5">
                <FormLabel>Course Code</FormLabel>
                <FormControl>
                  <CourseSearchBar
                    ref="courseSearchRef"
                    @update:selected="(value) => componentField.onChange(value)"
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            </FormField>

            <FormField v-slot="{ componentField }" name="semester">
              <FormItem class="w-full md:w-1/3">
                <FormLabel>Semester</FormLabel>
                <Select v-bind="componentField">
                  <FormControl>
                    <SelectTrigger>
                      <SelectValue placeholder="Select a semester" />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    <SelectGroup>
                      <SelectItem value="Fall">Fall</SelectItem>
                      <SelectItem value="Winter">Winter</SelectItem>
                      <SelectItem value="Spring/Summer">Spring/Summer</SelectItem>
                    </SelectGroup>
                  </SelectContent>
                </Select>
                <FormMessage />
              </FormItem>
            </FormField>

            <FormField v-slot="{ componentField }" name="year">
              <FormItem class="w-full md:w-1/3">
                <FormLabel>Year</FormLabel>
                <FormControl>
                  <Input type="number" placeholder="2025" v-bind="componentField" />
                </FormControl>
                <FormMessage />
              </FormItem>
            </FormField>
          </section>

          <div class="flex items-center justify-end mt-4">
            <Button size="sm" type="submit">Submit</Button>
          </div>
        </div>
      </div>
    </div>
  </form>
</template>