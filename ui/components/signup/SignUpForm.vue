<script setup lang="ts">
import { useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import * as z from 'zod'

import { Button } from '@/components/ui/button'
import {
  FormControl,
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
const { user } = useAuth()
const newUser = ref(user.value?.newuser)
const easterEgg = ref(false)

const { isValid, checkUsername } = useCheckUsername()

const formSchema = toTypedSchema(z.object({
  nickname: usernameSchema,
  gender: z.string(),
  programId: z.string(),
  currentYear: z.number().min(1).max(8),
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
  { label: '6th Year', value: 6 },
  { label: '7th Year', value: 7 },
  { label: '8th Year', value: 8 },
]

const form = useForm({
  validationSchema: formSchema,
})

const { setFieldValue, handleSubmit } = form

watchEffect(() => {
  if (form.values.currentYear === 8) {
    easterEgg.value = true
  } else {
    easterEgg.value = false
  }
})

const onSubmit = handleSubmit(async (formValues) => {
  loading.value = true
  try {
    const { status } = await useFetch(`/api/user/update/${user.value?.id}`, {
      headers: {
        'Content-Type': 'application/json'
      },
      credentials: "include",
      method: 'PATCH',
      body: formValues,
    });
    
    if (status.value === 'success') {
      toast({
        title: 'Welcome to Syllabye!',
        description: 'Your profile has been created.',
        duration: 1500
      })
      setInterval(reloadNuxtApp, 1500)
    } else if (status.value === 'pending') {
      loading.value = true
    } else {
      toast({
        title: 'Error creating account',
        description: 'Please try again'
      })
      loading.value = false
    }
  } catch (err) {
    console.error('Error sending PATCH request:', err);
    toast({
      title: 'Something went wrong in our server',
      description: 'Please try again later'
    })
  }
})
</script>

<template>
  <Dialog :open="newUser ? true : false">
    <DialogContent 
      class="[&>button:last-child]:hidden outline-none"
      @interact-outside="(event) => event.preventDefault()" 
      @escape-key-down="(event) => event.preventDefault()"
    >
      <DialogHeader class="-space-y-2">
        <template v-if="loading">
          <h1 class="text-lg">Loading...</h1>
        </template>
        <template v-else>
          <h1 class="text-lg">Hey there, first time user!</h1>
          <p class="text-muted-foreground font-light text-sm">Enter your details to continue</p>
        </template>
      </DialogHeader>
      <Separator class="h-[0.5px]"/>
      
      <SkeletonPage :rows='7' v-if="loading"/>
      <template v-else>
        <form class="space-y-4" @submit="onSubmit">
          <!-- Username Field -->
          <FormField v-slot="{ componentField }" name="nickname">
            <FormItem>
              <FormLabel class="flex justify-between">
                <p>Username</p>
                <p class="text-xs text-destructive motion-preset-fade-sm" v-if="isValid?.valid === false && !isValid.zod">{{ `❌ ${isValid.message}` }}</p>
                <p class="text-xs motion-preset-fade-sm" v-if="isValid?.valid === true && !isValid.zod">{{ `✅ ${isValid.message}` }}</p>
              </FormLabel>
              <FormControl>
                <Input
                  type="text"
                  placeholder="johndoe1"
                  v-bind="componentField"
                  @input="checkUsername($event.target.value)"
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          </FormField>
          <!-- Gender Field -->
          <div class="flex items-center gap-x-2">
            <FormField v-slot="{ componentField }" name="gender">
              <FormItem class="w-1/2">
                <FormLabel>Gender</FormLabel>
                <Select v-bind="componentField">
                  <FormControl>
                    <SelectTrigger>
                      <SelectValue placeholder="Gender" />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    <SelectGroup>
                      <SelectItem value="Male">
                        Male
                      </SelectItem>
                      <SelectItem value="Female">
                        Female
                      </SelectItem>
                      <SelectItem value="Other">
                        Other
                      </SelectItem>
                    </SelectGroup>
                  </SelectContent>
                </Select>
              </FormItem>
            </FormField>
            
            <FormField v-slot="{ componentField }" name="currentYear">
              <FormItem class="w-1/2">
                <FormLabel>
                  Year of Study
                  <span v-if="easterEgg">👀</span>
                </FormLabel>
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
  
          <!-- Program Field -->
          <div class="relative w-full">
            <FormField name="program">
              <FormItem>
                <FormLabel>Program of Study</FormLabel>
                <Combobox by="label">
                  <FormControl>
                    <ComboboxAnchor class="w-full">
                      <div class="relative w-full">
                        <ComboboxInput
                          class="w-full"
                          :display-value="(val) => val?.label ?? ''"
                          placeholder="Select program..."
                        />
                        <ComboboxTrigger class="absolute end-0 inset-y-0 flex items-center justify-center px-3">
                          <ChevronsUpDown class="size-4 text-muted-foreground" />
                        </ComboboxTrigger>
                      </div>
                    </ComboboxAnchor>
                  </FormControl>
                  
                  <ComboboxList class="overflow-y-auto max-h-64 w-[464px]">
                    <ComboboxEmpty>
                      Nothing found.
                    </ComboboxEmpty>

                    <ComboboxGroup>
                      <ComboboxItem
                        v-for="program in programs"
                        :key="program.value"
                        :value="program"
                        @select="() => {
                          setFieldValue('programId', program.value)
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
          </div>
          <div class="flex justify-end pt-3">
            <Button type="submit">
              Submit
            </Button>
          </div>
        </form>
      </template>
    </DialogContent>
  </Dialog>
</template>