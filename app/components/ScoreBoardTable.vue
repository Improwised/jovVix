<template>
  <v-card elevation="0" class="table-responsive container pb-5">
    <v-card-item>
      <div class="d-flex align-center justify-center">
        <div>
          <h5 class="text-h5 mb-1 font-weight-semibold">Rankings</h5>
        </div>
      </div>
      <table class="table align-middle table-bordered">
        <thead>
          <tr>
            <th>Rank</th>
            <th>User</th>
            <th>Score</th>
          </tr>
        </thead>
        <tbody class="table-group-divider">
          <tr
            v-for="(user, index) in scoreboardData"
            :key="index"
            :class="{
              'table-primary':
                user.username === props.userName && !props.isAdmin,
            }"
          >
            <td>
              {{ user.rank }}
            </td>
            <td v-if="props.isAdmin">
              {{ user.firstname }} <span>({{ user.username }})</span>
            </td>
            <td v-else>
              {{ user.firstname }}
              <span v-if="props?.userName === user.username">
                &nbsp;({{ user.username }})
              </span>
            </td>
            <td>{{ user.score }}</td>
          </tr>
        </tbody>
      </table>
    </v-card-item>
  </v-card>
</template>

<script setup>
const props = defineProps({
  scoreboardData: {
    type: Array,
    required: true,
    default: () => {
      return [];
    },
  },
  isAdmin: {
    type: Boolean,
    required: false,
    default: false,
  },
  userName: {
    type: String,
    required: true,
    default: "",
  },
});
</script>
