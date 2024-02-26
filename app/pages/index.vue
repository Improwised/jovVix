<script setup>
useSystemEnv();
const systemEnv = useSystemEnv();
console.log(systemEnv.value.api_url)
console.log(systemEnv.value.base_url)
</script>

<template>
  <LayoutsBase>
    <div class="text-center homepage">Welcome to nuxt3</div>
  </LayoutsBase>
</template>
<style scoped>
.homepage {
  margin-top: 400px;
  font-size: 30px;
}
</style>
