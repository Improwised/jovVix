<script setup>
const props = defineProps({
  user: {
    type: Object,
    required: true,
    default: () => {
      return {};
    },
  },
});

const permission = ref(props.user.permission);
</script>

<template>
  <v-list-item-title>
    <div class="d-flex align-center py-3">
      <div class="mr-3">
        <v-badge bordered bottom color="success" dot offset-x="0" offset-y="0">
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

      <!-- User Permission -->
      <div class="d-flex">
        <select
          id="permission"
          v-model="permission"
          class="form-select"
          required
          disabled
        >
          <option value="" disabled>Select permission level</option>
          <option value="read">Read</option>
          <option value="write">Write</option>
          <option value="share">Share</option>
        </select>
      </div>
    </div>
  </v-list-item-title>
</template>
