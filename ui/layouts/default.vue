<template>
  <div class="flex flex-col min-h-screen">
    <NavigationBar/>
    <NuxtLoadingIndicator :color="true" :height="5" />
    <div class="md:mt-14 flex-grow">
      <slot />
    </div>
    <Footer/>
  </div>
</template>