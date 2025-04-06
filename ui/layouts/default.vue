<script setup lang="ts">
import SignUpForm from '~/components/signup/SignUpForm.vue';
const user = await useUser()
const newUser = ref(user.value?.newuser)
</script>

<template>
  <div class="flex flex-col min-h-screen">
    <NavigationBar/>
    <NuxtLoadingIndicator :color="true" :height="5" />
    <div class="md:mt-14 flex-grow">
      <SignUpForm :open="newUser ? true : false"/>
      <slot />
    </div>
    <Footer/>
  </div>
</template>