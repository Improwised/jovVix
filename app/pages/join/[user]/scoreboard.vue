<script setup>
const headers = useRequestHeaders(["cookie"]);
const url = "/final_score/user";
const configs = useRuntimeConfig();

const userMeta = ref({});
const whoEndpoint = "/user/who";

async function getUserNameData() {
  const response = await $fetch(configs.public.api_url + whoEndpoint, {
    method: "GET",
    headers: headers,
    credentials: "include",
    mode: "cors",
    onResponseError(response, request) {
      console.log(response);
      console.log(request);
    },
  });
  userMeta.value = response.data;
}

setTimeout(async () => {
  await getUserNameData();
}, 500);
</script>

<template>
  <div class="w-100">
    <FinalScoreBoard
      :is-admin="false"
      :user-name="userMeta.username"
      :user-u-r-l="url"
    />
  </div>
</template>
