<script setup lang="ts">
import { useLoginModal } from '@/composables/useLoginModal'
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
} from '@/components/ui/dialog/index'

const { showLoginModal } = useLoginModal()

const config = useRuntimeConfig()
const authUrl = config.public.googleAuth + "?redirect=" + config.public.googleRedirectUrl
</script>

<template>
  <Dialog :open="showLoginModal">
    <DialogContent class="sm:max-w-[500px] px-4 [&>button:last-child]:hidden outline-none">
      <DialogHeader class="pb-2">
        <SyllabyeLogo size="large" :slogan="true"/>
      </DialogHeader>
      <Separator />
      <div class="space-y-2 py-5">
        <CardDescription>Please sign in with your <b>torontomu.ca</b> email to continue</CardDescription>
        <NuxtLink :to="authUrl">
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
