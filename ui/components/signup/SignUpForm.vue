<script setup lang="ts">
import { useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import * as z from 'zod'

import { Button } from '@/components/ui/button'
import {
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Separator } from '@/components/ui/separator'
import DialogContent from '../ui/dialog/DialogContent.vue'
import { cn } from '@/lib/utils'
import { Combobox, ComboboxAnchor, ComboboxEmpty, ComboboxGroup, ComboboxInput, ComboboxItem, ComboboxItemIndicator, ComboboxList, ComboboxTrigger } from '@/components/ui/combobox'
import { Check, ChevronsUpDown, Search } from 'lucide-vue-next'
import { ref } from 'vue'

const formSchema = toTypedSchema(z.object({
  username: z.string().min(2).max(50),
  gender: z.string(),
  program: z.string(),
  year: z.string()
}))

const form = useForm({
  validationSchema: formSchema,
})

const onSubmit = form.handleSubmit((values) => {
  console.log('Form submitted!', values)
})

const programs = [
  { label: 'English', value: 'en' },
  { label: 'French', value: 'fr' },
  { label: 'German', value: 'de' },
  { label: 'Spanish', value: 'es' },
  { label: 'Portuguese', value: 'pt' },
  { label: 'Russian', value: 'ru' },
  { label: 'Japanese', value: 'ja' },
  { label: 'Korean', value: 'ko' },
  { label: 'Chinese', value: 'zh' },
] as const

const years = [
  { label: '1st Year', value: '1' },
  { label: '2nd Year', value: '2' },
  { label: '3rd Year', value: '3' },
  { label: '4th Year', value: '4' },
  { label: '5th Year', value: '5' },
  { label: 'Other', value: 'other'}
]

const { handleSubmit, setFieldValue } = useForm({
  validationSchema: formSchema,
  initialValues: {
    program: '',
  },
})

const value = ref<typeof programs[0]>()
</script>

<template>
  <Dialog>
    <DialogTrigger>Sign Up Form</DialogTrigger>
    <DialogContent 
      class="[&>button:last-child]:hidden"
      @interact-outside="(event) => event.preventDefault()" 
      @escape-key-down="(event) => event.preventDefault()"
    >
      <DialogHeader class="-space-y-2">
        <h1 class="text-lg">Hey there, first time user!</h1>
        <p class="text-muted-foreground font-light text-sm">Enter your details to continue</p>
      </DialogHeader>
      <Separator class="h-[0.5px]"/>
      
      <form @submit="onSubmit" class="space-y-4">
        <FormField v-slot="{ componentField }" name="username">
          <FormItem>
            <FormLabel>Username</FormLabel>
            <FormControl>
              <Input type="text" placeholder="johndoe1" v-bind="componentField" />
            </FormControl>
            <FormDescription class="text-xs">
              This is your public display name.
            </FormDescription>
            <FormMessage />
          </FormItem>
        </FormField>

        <FormField v-slot="{ componentField }" name="gender">
          <FormItem>
            <FormLabel>Gender</FormLabel>

            <Select v-bind="componentField">
              <FormControl>
                <SelectTrigger>
                  <SelectValue placeholder="Gender" />
                </SelectTrigger>
              </FormControl>
              <SelectContent>
                <SelectGroup>
                  <SelectItem value="male">
                    Male
                  </SelectItem>
                  <SelectItem value="female">
                    Female
                  </SelectItem>
                  <SelectItem value="other">
                    Other
                  </SelectItem>
                </SelectGroup>
              </SelectContent>
            </Select>
          </FormItem>
        </FormField>

        <div class="flex items-center justify-between">
          <FormField name="program">
            <FormItem>
              <FormLabel>Program of Study</FormLabel>
  
              <Combobox by="label" class="relative w-full">
                <FormControl>
                  <ComboboxAnchor>
                    <div class="relative w-max">
                      <ComboboxInput class="" :display-value="(val) => val?.label ?? ''" placeholder="Select program..." />
                      <ComboboxTrigger class="absolute end-0 inset-y-0 flex items-center justify-center px-3">
                        <ChevronsUpDown class="size-4 text-muted-foreground" />
                      </ComboboxTrigger>
                    </div>
                  </ComboboxAnchor>
                </FormControl>
  
                <ComboboxList>
                  <ComboboxEmpty>
                    Nothing found.
                  </ComboboxEmpty>
  
                  <ComboboxGroup>
                    <ComboboxItem
                      v-for="program in programs"
                      :key="program.value"
                      :value="program"
                      @select="() => {
                        setFieldValue('program', program.value)
                      }"
                    >
                      {{ program.label }}
  
                      <ComboboxItemIndicator>
                        <Check :class="cn('ml-auto h-4 w-4')" />
                      </ComboboxItemIndicator>
                    </ComboboxItem>
                  </ComboboxGroup>
                </ComboboxList>
              </Combobox>
  
              <FormMessage />
            </FormItem>
          </FormField>

          <div class="w-2/5">
            <FormField v-slot="{ componentField }" name="year">
              <FormItem>
                <FormLabel>Year of Study</FormLabel>
  
                <Select v-bind="componentField">
                  <FormControl>
                    <SelectTrigger>
                      <SelectValue placeholder="Select a verified email to display" />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    <SelectGroup>
                      <SelectItem v-for="year in years"
                        :key="year.value"
                        :value="year.value"
                        @select="() => {
                          setFieldValue('year', year.value)
                        }"
                      >
                        {{ year.label }}
                      </SelectItem>
                    </SelectGroup>
                  </SelectContent>
                </Select>
              </FormItem>
            </FormField>
          </div>
        </div>
        <div class="flex justify-end pt-3">
          <Button type="submit">
            Submit
          </Button>
        </div>
      </form>
    </DialogContent>
  </Dialog>
</template>