<script setup lang="ts">
import { onMounted } from 'vue'
import init, { greet } from '~/wasm/pkg'

onMounted(async () => {
  await init()
  greet('please fill in email and password')
})
</script>
<template>
  <v-app>
    <NuxtPage />
  </v-app>
</template>