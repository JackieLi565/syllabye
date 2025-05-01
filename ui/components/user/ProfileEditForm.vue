<script setup lang="ts">
import { useForm } from 'vee-validate'
import { Icon } from '@iconify/vue/dist/iconify.js'
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
import { Textarea } from '#components'

const { toast } = useToast()
const loading = ref(false)
const isDialogOpen = ref(false)

const { user } = useAuth()
const { programs: programsData } = usePrograms()
const { isValid, checkUsername } = useCheckUsername()

const formSchema = toTypedSchema(z.object({
  nickname: usernameSchema,
  gender: z.string(),
  programId: z.string(),
  currentYear: z.number().min(1).max(8),
  instagram: z.string()
    .regex(/^(https?:\/\/)?(www\.)?instagram\.com\/[a-zA-Z0-9._]+\/?$/, {
      message: 'Please enter a valid Instagram link like instagram.com/username',
    })
    .or(z.literal(''))
    .optional(),
  bio: z.string().max(200, {
      message: 'Your bio can only contain at most 200 characters'
    }).optional(),
}))

const values = {
  nickname: user.value?.nickname ?? '',
  gender: user.value?.gender ?? '',
  programId: user.value?.programId ?? '',
  instagram: user.value?.instagram ?? '',
  bio: user.value?.bio ?? '',
  currentYear: user.value?.currentYear ?? undefined,
}

const programs = computed(() =>
  programsData.value?.map(program => ({
    label: program.name,
    value: program.id
  })) ?? []
)

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

const originalValues = ref(values)

const form = useForm({
  initialValues: originalValues.value,
  validationSchema: formSchema
})

const { setFieldValue, handleSubmit, resetForm } = form

watch(isDialogOpen, (open) => {
  if (open) {
    const freshValues = values
    originalValues.value = { ...freshValues }
    resetForm({ values: freshValues })
    isValid.value = null // reset custom validation state if needed
  }
})

const getChangedValues = () => {
  return Object.fromEntries(
    Object.entries(form.values).filter(
      ([key, value]) => value !== originalValues.value[key as keyof typeof originalValues.value]
    )
  )
}

const onSubmit = handleSubmit(async () => {
  console.log(form.values)
  loading.value = true
  try {
    const changedValues = getChangedValues()
    
    if (Object.keys(changedValues).length === 0) {
      toast({
        title: 'No changes made',
        description: 'Nothing was updated.',
      })
      loading.value = false
      return
    }

    const { status } = await useFetch(`/api/user/update/${user.value?.id}`, {
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
      method: 'PATCH',
      body: changedValues,
    })

    if (status.value === 'success') {
      await refreshNuxtData('user')

      toast({
        title: 'Profile updated!',
        description: 'Your changes were saved.',
        duration: 1000
      })
      setInterval(reloadNuxtApp, 1000)
    } else {
      toast({
        title: 'Error updating profile',
        description: 'Please try again',
      })
    }
  } catch (err) {
    console.error('Error sending PATCH request:', err)
    toast({
      title: 'Something went wrong',
      description: 'Please try again later',
    })
  }
})
</script>

<template>
  <Dialog v-model:open="isDialogOpen">
    <DialogTrigger>
      <Button size="default" variant="secondary">
        <Icon icon="lucide:user-pen" />
        <span>Edit</span>
      </Button>
    </DialogTrigger>

    <DialogContent class="outline-none">
      <template v-if="isDialogOpen">
        <DialogHeader class="-space-y-2">
          <h1 class="text-lg">{{ loading ? 'Loading...' : 'Edit your profile' }}</h1>
        </DialogHeader>
        <Separator class="h-[0.5px]" />

        <SkeletonPage :rows="7" v-if="loading" />
        <template v-else>
          <form class="space-y-4" @submit="onSubmit">
            <FormField v-slot="{ componentField }" name="nickname">
              <FormItem>
                <FormLabel class="flex justify-between">
                  <p>Username</p>
                  <p class="text-xs text-destructive motion-preset-fade-sm" v-if="isValid && isValid.valid === false && !isValid.zod">{{ `❌ ${isValid.message}` }}</p>
                  <p class="text-xs motion-preset-fade-sm" v-if="isValid && isValid.valid === true && !isValid.zod">{{ `✅ ${isValid.message}` }}</p>
                </FormLabel>
                <FormControl>
                  <Input
                    type="text"
                    :placeholder="user?.nickname"
                    v-bind="componentField"
                    @input="checkUsername($event.target.value)"
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            </FormField>
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
                        <SelectItem value="Male">Male</SelectItem>
                        <SelectItem value="Female">Female</SelectItem>
                        <SelectItem value="Other">Other</SelectItem>
                      </SelectGroup>
                    </SelectContent>
                  </Select>
                </FormItem>
              </FormField>

              <FormField v-slot="{ componentField }" name="currentYear">
                <FormItem class="w-1/2">
                  <FormLabel>Year of Study</FormLabel>
                  <Select v-bind="componentField">
                    <FormControl>
                      <SelectTrigger>
                        <SelectValue placeholder="Select year..." />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      <SelectGroup>
                        <SelectItem
                          v-for="year in years"
                          :key="year.value"
                          :value="year.value"
                          @select="() => setFieldValue('currentYear', year.value)"
                        >
                          {{ year.label }}
                        </SelectItem>
                      </SelectGroup>
                    </SelectContent>
                  </Select>
                </FormItem>
              </FormField>
            </div>

            <div class="relative w-full">
              <FormField name="programId">
                <FormItem>
                  <FormLabel>Program of Study</FormLabel>
                  <Combobox v-model="form.values.programId">
                    <FormControl>
                      <ComboboxAnchor class="w-full">
                        <div class="relative w-full">
                          <ComboboxInput
                            class="w-full"
                            :display-value="(programId) => {
                              const program = programs.find(p => p.value === programId)
                              return program ? program.label : ''
                            }"
                            placeholder="Select program..."
                          />
                          <ComboboxTrigger class="absolute end-0 inset-y-0 flex items-center justify-center px-3">
                            <ChevronsUpDown class="size-4 text-muted-foreground" />
                          </ComboboxTrigger>
                        </div>
                      </ComboboxAnchor>
                    </FormControl>

                    <ComboboxList class="overflow-y-auto max-h-64 w-[464px]">
                      <ComboboxEmpty>Nothing found.</ComboboxEmpty>

                      <ComboboxGroup>
                        <ComboboxItem
                          v-for="program in programs"
                          :key="program.value"
                          :value="program.value" 
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

            <div class="relative w-full">
              <FormField v-slot="{ componentField }" name="instagram">
                <FormItem>
                  <FormLabel class="flex justify-between">
                    <p>Instagram Link</p>
                  </FormLabel>
                  <FormControl>
                    <Input
                      type="text"
                      placeholder="instagram.com/username"
                      v-bind="componentField"
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              </FormField>
            </div>

            <div class="relative w-full">
              <FormField v-slot="{ componentField }" name="bio">
                <FormItem>
                  <FormLabel class="flex justify-between">
                    <p>Bio</p>
                  </FormLabel>
                  <FormControl>
                    <Textarea
                      type="text"
                      :placeholder="user?.bio"
                      v-bind="componentField"
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              </FormField>
            </div>

            <div class="flex justify-end pt-3">
              <Button type="submit" :disabled="Object.keys(getChangedValues()).length === 0 || (isValid && isValid.valid === false && !isValid.zod)">Submit</Button>
            </div>
          </form>
        </template>
      </template>
    </DialogContent>
  </Dialog>
</template>
