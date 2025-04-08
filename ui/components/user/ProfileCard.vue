<script setup lang="ts">
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge';
import { Icon } from '@iconify/vue/dist/iconify.js';
import ProfileEditForm from './ProfileEditForm.vue';

const { user } = useAuth()
const { program } = usePrograms({ programId: user.value?.programId })
console.log(user.value)
</script>

<template>
  <main class="space-y-6">
    <div class="border-b md:border border-border w-full md:w-3/5 xl:w-2/5 m-auto md:rounded-xl overflow-hidden">
      <div class="h-36 bg-gradient-to-r from-primary to-indigo-800"></div>
      <div class="pb-8">
        <div class="flex flex-row justify-between mx-8">
          <div>
            <div class="-mt-16">
              <Avatar size="lg" class="border border-border">
                <AvatarImage :src="user?.picture || ''" alt="Profile" />
                <AvatarFallback>{{ user?.fullname.charAt(0) }}</AvatarFallback>
              </Avatar>
            </div>
            <div class="mt-1">
              <h1 class="text-xl font-semibold">{{ user?.fullname }}</h1>
              <div class="text-muted-foreground flex items-start space-x-2">
                <span>{{ user?.nickname }}</span>
                <NuxtLink to="https://instagram.com/" target="blank">
                  <Badge variant="secondary" class="rounded-xl py-1 space-x-1">
                    <Icon icon="mdi:instagram" class="text-base" />
                    <p>johndoe</p>
                  </Badge>
                </NuxtLink>
              </div>
            </div>
          </div>
  
          <div class="mt-4">
            <ProfileEditForm />
          </div>
        </div>
        <div class="mx-8 space-y-4">
          <div class="mt-2 text-muted-foreground text-sm">
            <h1 class="flex items-center space-x-1.5"><Icon icon="lucide:book-text" class="text-lg" /><span>{{ program?.name }}</span></h1>
            <h1 class="flex items-center space-x-1.5"><Icon icon="mdi:education-outline" class="text-xl" /><span>{{ user?.currentYear || '4th Year' }}</span></h1>
          </div>
          <div>
            <p class="font-light text-sm">Some bio thingy</p>
          </div>
        </div>
      </div>
    </div>
  
    <div class="px-8 py-6 border border-border w-full md:w-3/5 xl:w-2/5 m-auto rounded-xl overflow-hidden h-96">
      <h1 class="text-xl font-semibold">Syllabi Uploaded</h1>
    </div>
  </main>
</template>