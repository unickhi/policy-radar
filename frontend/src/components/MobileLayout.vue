<template>
  <div class="mobile-app-container">
    <router-view />
  </div>
</template>

<script setup lang="ts">
</script>

<style scoped>
.mobile-app-container {
  max-width: 450px;
  margin: 0 auto;
  background: #fff;
  min-height: 100vh;
  position: relative;
}
</style>