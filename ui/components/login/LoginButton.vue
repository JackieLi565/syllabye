<script setup lang="ts">
import {
  CardDescription,
} from '@/components/ui/card'
import { Icon } from '@iconify/vue'
import { Button } from '@/components/ui/button'
import SyllabyeLogo from '@/components/logo/SyllabyeLogo.vue'
import { Separator } from '@/components/ui/separator'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTrigger,
} from '@/components/ui/dialog/index'

const config = useRuntimeConfig();
console.log(config.public.googleAuth)
</script>

<template>
  <Dialog>
    <DialogTrigger as-child>
      <Button variant="default" class="text-xs md:text-sm">
        Sign In
      </Button>
    </DialogTrigger>
    <DialogContent class="sm:max-w-[500px] px-4">
        <DialogHeader class="pb-2">
          <SyllabyeLogo size="large" :slogan="true"/>
        </DialogHeader>
        <Separator />
        <div class="space-y-2 py-5">
          <CardDescription>Sign in with your <b>torontomu.ca</b> email</CardDescription>
          <NuxtLink :to="config.public.googleAuth">
            <Button variant="secondary">
              <Icon icon="logos:google-icon" class="h-[1.6rem] w-[1.6rem]"/>
              <span>Sign in with Google</span>
            </Button>
          </NuxtLink>
        </div>
        <Separator />
        <div class="py-2 text-muted-foreground">
          <p class="text-xs">By signing in, you agree with our <NuxtLink to="https://www.google.com" class="text-primary hover:underline">terms of service</NuxtLink>.</p>
        </div>
    </DialogContent>
  </Dialog>
  
</template>