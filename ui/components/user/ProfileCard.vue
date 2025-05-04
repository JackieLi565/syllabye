<script setup lang="ts">
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge';
import { Icon } from '@iconify/vue/dist/iconify.js';
import ProfileEditForm from './ProfileEditForm.vue';
import SkeletonBar from '../loading/SkeletonBar.vue';
import { schoolYearFormatter } from '~/lib/utils';
import Button from '../ui/button/Button.vue';
import SkeletonPage from '../loading/SkeletonPage.vue';

const { user } = useAuth()
const { program, programStatus } = usePrograms({ programId: user.value?.programId })
const { syllabi, syllabiStatus } = useSyllabi({ userId: user.value?.id })
</script>

<template>
  <main class="space-y-6 motion-preset-slide-up-sm">
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
                <template v-if="user?.instagram">
                  <NuxtLink :to="user?.instagram" target="blank">
                    <Badge variant="secondary" class="rounded-xl py-1 space-x-1">
                      <Icon icon="mdi:instagram" class="text-base" />
                      <p>@{{ user?.instagram?.split('/').slice(-1).pop() }}</p>
                    </Badge>
                  </NuxtLink>
                </template>
              </div>
            </div>
          </div>
  
          <div class="mt-4">
            <ProfileEditForm />
          </div>
        </div>
        <div v-if="programStatus !== 'success'" class="mx-8 my-2">
          <SkeletonBar />
        </div>
        <div class="mx-8 space-y-4" v-else>
          <div class="mt-2 text-muted-foreground text-sm">
            <div class="flex items-center space-x-1.5"><Icon icon="lucide:book-text" class="text-lg" /><span>{{ program?.name }}</span></div>
            <h1 class="flex items-center space-x-1.5"><Icon icon="mdi:education-outline" class="text-xl" /><span>{{ schoolYearFormatter(user?.currentYear) }}</span></h1>
          </div>
          <div>
            <p class="font-light text-sm">{{ user?.bio }}</p>
          </div>
        </div>
      </div>  
    </div>

    <div class="px-8 py-6 border border-border w-full md:w-3/5 xl:w-2/5 m-auto rounded-xl overflow-hidden h-fit">
      <h1 class="text-xl font-semibold">Your Syllabi</h1>
      <p class="text-muted-foreground text-sm">View syllabi you have uploaded</p>
      <div v-if="syllabi?.length === 0" class="flex flex-col justify-center items-center space-y-3 h-48">
        <div class="flex flex-col items-center">
          <div class="rounded-full bg-primary-alternative p-3 mb-3">
            <Icon icon="lucide:file-text" width="36" height="36" class="text-primary"/>
          </div>
          <h1 class="flex space-x-1 items-center font-medium">Wow, such empty</h1>
          <h1 class="flex space-x-1 items-center text-sm text-muted-foreground">No syllabi to see here...yet</h1>
        </div>
        <NuxtLink to="/upload">
          <Button>
            <Icon icon="lucide:plus"/>
            <span>Upload a syllabus</span>
          </Button>
        </NuxtLink>
      </div>
      <template v-else-if="syllabiStatus === 'idle'">
        <SkeletonPage :rows="3" class="py-4"></SkeletonPage>
      </template>
      <template v-else-if="syllabiStatus === 'success' && syllabi?.length !== 0">
        <div>Syllabi content goes here</div>
      </template>
    </div>
  </main>
</template>