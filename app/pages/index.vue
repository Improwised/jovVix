<script setup>
useSystemEnv1();
const envs = useState('envs')
console.log(envs.value.api_url)
console.log(envs.value.base_url)
  
const systemEnv = useSystemEnv();
console.log(systemEnv.value.api_url)
console.log(systemEnv.value.base_url)
</script>

<template>
  <div class="text-center homepage">Welcome to nuxt3</div>
</template>
<style scoped>
.homepage {
  margin-top: 400px;
  font-size: 30px;
}
</style>
