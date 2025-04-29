<template>
  <div class="w-[800px] border">
    <blockquote class="reddit-embed-bq w-full" style="height:316px" data-embed-height="316">
      <a href="https://www.reddit.com/r/TorontoMetU/comments/1ed7wta/looking_for_course_syllabuscourse_outlines/">
        Looking for course syllabus/course outlines
      </a>
      <br>by <a href=""></a> in
      <a href="https://www.reddit.com/r/TorontoMetU/">TorontoMetU</a>
    </blockquote>

    <blockquote class="reddit-embed-bq" style="height:316px" data-embed-height="316">
      <a href="https://www.reddit.com/r/TorontoMetU/comments/1ex1daj/where_to_find_course_syllabus_information/">
        Where to find course syllabus information
      </a>
      <br>by <a href="https://www.reddit.com/user/roboticbone/">u/roboticbone</a> in
      <a href="https://www.reddit.com/r/TorontoMetU/">TorontoMetU</a>
    </blockquote>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'

onMounted(() => {
  const script = document.createElement('script')
  script.src = 'https://embed.reddit.com/widgets.js'
  script.async = true
  document.body.appendChild(script)
})
</script>
