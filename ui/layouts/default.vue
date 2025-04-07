<script setup lang="ts">
import SignUpForm from '~/components/signup/SignUpForm.vue';
import LoginModal from '~/components/login/LoginModal.vue';
</script>

<template>
  <div class="flex flex-col min-h-screen">
    <NavigationBar/>
    <NuxtLoadingIndicator :color="true" :height="5" />
    <div class="md:mt-14 flex-grow">
      <SignUpForm />
      <LoginModal />
      <slot />
    </div>
    <Footer/>
  </div>
</template>