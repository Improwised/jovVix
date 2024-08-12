<script setup lang="js">
// core dependencies
const { $bootstrap } = useNuxtApp();

// define props and emits
const props = defineProps({
    'modalTitle': {
        type: String,
        default: "This is a modal Title",
        required: true
    }, 'modalMessage': {
        type: String,
        default: "Hello this is a modal message",
        required: true
    }, "modelPositiveMessage": {
        type: String,
        default: "Save",
        required: false
    }, "modelNegativeMessage": {
        type: String,
        default: "Cancel",
        required: false
    }
})
const emits = defineEmits(['confirmMessage'])

// refs
const modal = ref()
const isSent = ref(false)

// handlers
const handleClose = (confirm) => {
    if (!isSent.value) {
        isSent.value = true;
        emits('confirmMessage', confirm);
        modal.value.hide();
    }
}


// core logic
onMounted(() => {
    try {
        modal.value = new $bootstrap.Modal('#confirmModal');
        modal.value.show();
        modal.value._element.addEventListener('hidden.bs.modal', () => {
            handleClose(false)
        })
    } catch (e) {
        console.error('Bootstrap error: ', e);
    }
});
</script>

<template>
  <div>
    <div
      id="confirmModal"
      class="modal fade"
      tabindex="-1"
      aria-labelledby="confirmModalLabel"
      aria-hidden="true"
    >
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h1 id="confirmModalLabel" class="modal-title fs-5">
              {{ props["modalTitle"] }}
            </h1>
            <button
              type="button"
              class="btn-close text-white"
              data-bs-dismiss="modal"
              aria-label="Close"
              @click="() => handleClose(false)"
            ></button>
          </div>
          <div class="modal-body">
            {{ props["modalMessage"] }}
          </div>
          <div class="modal-footer">
            <button
              type="button"
              class="btn btn-secondary text-white"
              data-bs-dismiss="modal"
              @click="() => handleClose(false)"
            >
              {{ props.modelNegativeMessage }}
            </button>
            <button
              type="button"
              class="btn btn-primary text-white"
              @click="() => handleClose(true)"
            >
              {{ props.modelPositiveMessage }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
