<script setup>
import { useListUserstore } from "~/store/userlist";
import { storeToRefs } from "pinia";

const listUserStore = useListUserstore();
const { listUsers } = storeToRefs(listUserStore);
import { getAvatarUrlByName } from "~~/composables/avatar";
</script>

<template>
  <div class="container" style="max-width: 800px">
    <div class="row justify-content-center mb-2">
      <div v-if="listUsers.length == 0" class="col-7 col-md-4 mt-5 mb-5">
        <div
          class="d-flex border border-1 justify-content-center align-items-center px-3 py-2 py-md-4 gap-3 border-radius"
        >
          <font-awesome-icon icon="fa-solid fa-users" size="xl" />
          <h5 class="text-center mb-0">Waiting for Participants..</h5>
        </div>
      </div>

      <div v-else class="col-6 col-md-4 mt-5">
        <div
          class="d-flex border border-1 justify-content-center align-items-center px-3 py-2 py-md-4 gap-3 border-radius"
        >
          <font-awesome-icon icon="fa-solid fa-users" size="xl" />
          <h5 class="text-center text-sm fs-5 mb-0">
            {{ listUsers.length }} Participants
          </h5>
        </div>
      </div>
    </div>

    <v-card
      v-if="listUsers.length"
      :flat="true"
      class="mb-2 d-flex flex-wrap justify-content-center"
    >
      <div v-for="user in listUsers" :key="user.UserId" class="chip m-2">
        <img
          :src="getAvatarUrlByName(user?.Avatar)"
          alt="Person"
          width="96"
          height="96"
        />
        {{ user.UserName }}
      </div>
    </v-card>
  </div>
</template>

<style scoped>
.border-radius {
  border-radius: 2rem !important;
}

.chip {
  display: inline-block;
  padding: 0 25px;
  height: 50px;
  font-size: 16px;
  line-height: 50px;
  border-radius: 25px;
  max-width: 600px;
  background-color: #f1f1f1;
}

.chip img {
  float: left;
  margin: 0 10px 0 -25px;
  height: 50px;
  width: 50px;
  border-radius: 50%;
}
</style>
