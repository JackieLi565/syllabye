<template>
  <Dialog :open="props.open">
    <DialogContent 
      class="[&>button:last-child]:hidden outline-none"
      @interact-outside="(event) => event.preventDefault()" 
      @escape-key-down="(event) => event.preventDefault()"
    >
      <DialogHeader class="-space-y-2">
        <h1 class="text-lg">Hey there, first time user!</h1>
        <p class="text-muted-foreground font-light text-sm">Enter your details to continue</p>
      </DialogHeader>
      <Separator class="h-[0.5px]"/>
      
      <SkeletonPage :rows='7' v-if="loading"/>
      <template v-else>
        <form class="space-y-4" @submit="onSubmit">
          <!-- Username Field -->
          <FormField v-slot="{ componentField }" name="nickname">
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
  
          <!-- Gender Field -->
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
  
          <!-- Program Field -->
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
    
                  <ComboboxList class="overflow-y-auto max-h-64 w-[235px] ml-9">
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
                      </ComboboxItem>
                    </ComboboxGroup>
                  </ComboboxList>
                </Combobox>
    
                <FormMessage />
              </FormItem>
            </FormField>
  
            <!-- Year of Study Field -->
            <div class="w-2/5">
              <FormField v-slot="{ componentField }" name="currentYear">
                <FormItem>
                  <FormLabel>Year of Study</FormLabel>
                  <Select v-bind="componentField">
                    <FormControl>
                      <SelectTrigger>
                        <SelectValue placeholder="Select year..." />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      <SelectGroup>
                        <SelectItem v-for="year in years"
                          :key="year.value"
                          :value="year.value"
                          @select="() => {
                            setFieldValue('currentYear', year.value)
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
            <Button type="submit" :disabled="loading">
              Submit
            </Button>
          </div>
        </form>
      </template>
    </DialogContent>
  </Dialog>
</template>

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
import { Combobox, ComboboxAnchor, ComboboxEmpty, ComboboxGroup, ComboboxInput, ComboboxItem, ComboboxList, ComboboxTrigger } from '@/components/ui/combobox'
import { ChevronsUpDown } from 'lucide-vue-next'
import { useToast } from '@/components/ui/toast/use-toast'
import SkeletonPage from '../loading/SkeletonPage.vue'

const { toast } = useToast()
const loading = ref(false)
const user = await useUser()
const event = useRequestEvent(); 

const props = defineProps<{
  open: boolean
}>()

const formSchema = toTypedSchema(z.object({
  nickname: z
    .string()
    .min(3, 'Username must be at least 3 characters long')
    .max(30, 'Username cannot be longer than 30 characters')
    .regex(/^[a-z0-9._-]{3,30}$/, 'Username can only contain lowercase letters, numbers, periods, underscores, and dashes'),
  gender: z.string(),
  program: z.string(),
  currentYear: z.number().min(1).max(6)
}))

const { data: programsData } = await useFetch('/api/programs/programs')

const programs = programsData.value?.map(program => ({
  label: program.name,
  value: program.id,
}))
const years = [
  { label: '1st Year', value: 1 },
  { label: '2nd Year', value: 2 },
  { label: '3rd Year', value: 3 },
  { label: '4th Year', value: 4 },
  { label: '5th Year', value: 5 },
  { label: 'Other', value: 6 }
]

const form = useForm({
  validationSchema: formSchema,
})

const { setFieldValue, handleSubmit } = form

const onSubmit = handleSubmit(async (formValues) => {
  
  try {
    const { data, error, status } = await useFetch(`/api/user/update/${user.value?.id}`, {
      headers: {
        'Content-Type': 'application/json'
      },
      credentials: "include",
      method: 'PATCH',
      body: formValues,
    });
    if (status.value === 'pending') {
      loading.value = true
    }

    if (data.value) {
      toast({
        title: 'Welcome to Syllabye!',
        description: 'Your profile has been created.'
      })
    }
    if (error) {
      toast({
        title: 'Error creating account :(',
        description: 'Please try again'
      })
    }
  } catch (err) {
    console.error('Error sending PATCH request:', err);
    toast({
      title: 'Something went wrong in our server',
      description: 'Please try again later'
    })
  } finally {
    loading.value = false;
  }
})
</script>