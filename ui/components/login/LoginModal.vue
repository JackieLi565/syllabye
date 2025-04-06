<script setup lang="ts">
const showLoginModal = true;

const config = useRuntimeConfig()
const authUrl = config.public.googleAuth + "?redirect=" + config.public.googleRedirectUrl
</script>

<template>
  <Dialog :open="showLoginModal">
    <DialogContent 
      class="[&>button:last-child]:hidden"
      @interact-outside="(event) => event.preventDefault()" 
      @escape-key-down="(event) => event.preventDefault()"
    >
      <div class="sm:max-w-[500px] px-4">
        <DialogHeader class="pb-2">
          <SyllabyeLogo size="large" :slogan="true"/>
        </DialogHeader>
        <Separator />
        <div class="space-y-2 py-5">
          <CardDescription>Sign in with your <b>torontomu.ca</b> email</CardDescription>
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
      </div>
    </DialogContent>
  </Dialog>
</template>