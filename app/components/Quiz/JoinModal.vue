<template>
  <div class="modal-dialog modal-fullscreen modal-dialog-centered">
    <div class="modal-content">
      <div class="modal-header">
        <h1 id="exampleModalLabel" class="modal-title fs-3">Join Quiz</h1>
        <button
          type="button"
          class="btn-close"
          data-bs-dismiss="modal"
          aria-label="Close"
        ></button>
      </div>
      <div class="modal-body">
        <div class="row mt-5">
          <div
            class="col-md-6 p-5 d-flex align-items-center justify-content-center"
          >
            <form class="scale-modal">
              <div class="mb-3 pe-3">
                <div class="divider my-5 text-dark">Invitation Code</div>
                <div
                  class="d-flex align-items-center justify-content-center gap-2"
                >
                  <h2 class="display-4 code">{{ props.code }}</h2>
                  <font-awesome-icon
                    id="modal-OTP-input-container"
                    icon="fa-solid fa-copy"
                    size="xl"
                    style="color: #0c6efd"
                    class="copy-icon"
                    role="button"
                  />
                </div>
                <div class="divider my-5 text-dark">using link</div>
                <div
                  class="d-flex align-items-center justify-content-center gap-2"
                >
                  <div class="fs-1 text-dark text-decoration-underline">
                    quiz.i8d.in/join
                  </div>
                  <font-awesome-icon
                    id="modal-URL-input-container"
                    icon="fa-solid fa-copy"
                    size="xl"
                    style="color: #0c6efd"
                    class="copy-icon"
                    role="button"
                  />
                </div>
              </div>
            </form>
          </div>
          <div class="col-md-6">
            <div
              class="d-flex align-items-center justify-content-center qr-scale-down"
            >
              <QrCode
                :scan-u-r-l="props.joinURL"
                :quiz-code="code"
                :size="600"
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import usecopyToClipboard from "~~/composables/copy_to_clipboard";

const props = defineProps({
  code: {
    type: Number,
    required: true,
    default: 0,
  },
  joinURL: {
    type: String,
    required: true,
    default: "joinmyquiz.com",
  },
});

onMounted(() => {
  const copyBtn = document.getElementById("modal-OTP-input-container");
  const urlCopyBtn = document.getElementById("modal-URL-input-container");
  if (process.client && copyBtn && urlCopyBtn) {
    copyBtn.addEventListener("click", () => {
      usecopyToClipboard(props.code);
    });
    urlCopyBtn.addEventListener("click", () => {
      usecopyToClipboard(`${props.joinURL}?code=${props.code}`);
    });
  }
});
</script>

<style scoped>
.scale-modal {
  transform: scale(2);
}

.code {
  letter-spacing: 0.5rem;
}

@media (max-width: 768px) {
  .scale-modal {
    transform: scale(1);
  }

  .qr-scale-down {
    transform: scale(0.5);
  }
}
</style>
