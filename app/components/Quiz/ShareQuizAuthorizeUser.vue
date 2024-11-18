<script setup>
// define props and emits
const props = defineProps({
  user: {
    type: Object,
    required: true,
    default: () => {
      return {};
    },
  },
});
const emits = defineEmits(["showEditForm", "deleteUserPermission"]);

// to show edit form for change user permission for perticular quiz
const showEditForm = () => {
  emits(
    "showEditForm",
    props.user.id,
    props.user.shared_to,
    props.user.permission
  );
};
</script>

<template>
  <v-list-item-title>
    <div class="d-flex align-center py-3 justify-content-between">
      <div class="d-flex">
        <div class="mr-3">
          <v-badge
            bordered
            bottom
            color="success"
            dot
            offset-x="0"
            offset-y="0"
          >
            <v-avatar size="40">
              <img
                src="https://api.dicebear.com/9.x/bottts/svg?seed=Jade"
                alt="props.user.title"
                width="40"
              />
            </v-avatar>
          </v-badge>
        </div>

        <!-- User Details -->
        <div class="mx-3">
          <h4 v-if="props.user.first_name.Valid">
            {{ props.user.first_name.String }}
            {{ props.user.last_name.String }}
          </h4>
          <h4 v-else class="mt-n1 mb-1">Unknown</h4>
          <div class="truncate-text text-subtitle-2 textSecondary">
            {{ props.user.shared_to }}
          </div>
        </div>
      </div>

      <!-- User Permission -->
      <div class="d-flex align-items-center">
        {{ props.user.permission }}
        <!-- Button for edit permission -->
        <button
          type="button"
          title="Edit Permission"
          class="ml-2 badge rounded-pill bg-light-info text-dark m-1 px-2 fs-5"
          @click="showEditForm"
        >
          <font-awesome-icon :icon="['fas', 'pencil']" />
        </button>
        <button
          type="button"
          title="Edit Permission"
          class="ml-2 badge rounded-pill bg-light-danger text-dark m-1 px-2 fs-5"
          @click="emits('deleteUserPermission', props.user.id)"
        >
          <font-awesome-icon :icon="['fas', 'trash-can']" />
        </button>
      </div>
    </div>
  </v-list-item-title>
</template>
