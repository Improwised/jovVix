<script setup>
useSystemEnv1();
let envs = useState('envs')
console.log(1, envs.value.api_url)
console.log(2, envs.value.base_url)
  
let systemEnv = useSystemEnv();
console.log(3, systemEnv.value.api_url)
console.log(4, systemEnv.value.base_url)

onMounted(() => {
  if (process.client) {
    envs = useState('envs')
    console.log(1, envs.value.api_url)
    console.log(2, envs.value.base_url)

    systemEnv = useSystemEnv();
    console.log(3, systemEnv.value.api_url)
    console.log(4, systemEnv.value.base_url)
})
  
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
