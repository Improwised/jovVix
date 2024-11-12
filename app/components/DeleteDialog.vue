<template>
  <div
    :id="props.id"
    class="modal fade"
    tabindex="-1"
    role="dialog"
    aria-labelledby="exampleModalCenterTitle"
    aria-hidden="true"
  >
    <div class="modal-dialog modal-dialog-centered" role="document">
      <div class="modal-content">
        <div class="modal-header">
          <h5 id="exampleModalCenterTitle" class="modal-title">Delete Item</h5>
          <button
            type="button"
            class="btn-close"
            data-bs-dismiss="modal"
            aria-label="Close"
          >
            <span aria-hidden="true"></span>
          </button>
        </div>
        <div class="modal-body">
          Once deleted, this item cannot be recovered. Do you want to proceed?
        </div>
        <div class="modal-footer">
          <button
            type="button"
            class="btn btn-light-secondary"
            data-bs-dismiss="modal"
          >
            Close
          </button>
          <button
            type="button"
            class="btn btn-danger text-white"
            data-bs-dismiss="modal"
            @click="$emit('confirmDelete')"
          >
            Delete
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  id: {
    type: String,
    default: "",
    required: true,
  },
});

defineEmits(["confirmDelete"]);
</script>
