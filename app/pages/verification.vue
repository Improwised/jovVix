<template>
  <div class="row mt-5 justify-content-center g-0 min-vh-100">
    <div class="col-12 col-md-8 col-lg-6 col-xxl-4 py-8 py-xl-0">
      <div class="card smooth-shadow-md">
        <div class="card-body p-6">
          <form
            method="post"
            :action="verificationUrl"
            enctype="application/json"
          >
            <h3>Please enter the verification code sent to your email</h3>
            <div class="mb-3">
              <label for="code" class="form-label">Verification Code</label>
              <input
                id="code"
                v-model="code"
                name="code"
                type="text"
                class="form-control"
                maxlength="6"
                required
              />
            </div>
            <input type="hidden" name="method" value="code" />
            <input type="hidden" name="csrf_token" :value="csrfToken" />
            <button type="submit" class="btn btn-primary text-light">
              Verify
            </button>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useToast } from "vue-toastification";

const { kratos_url } = useRuntimeConfig().public;
const route = useRoute("");
const toast = useToast();

const code = ref("");
const csrfToken = ref("");

onMounted(async () => {
  await fetchFlowIdAndCsrfToken();
});

// Method to fetch flow ID and CSRF token
const fetchFlowIdAndCsrfToken = async () => {
  try {
    // Example endpoint to fetch CSRF token (replace with your actual endpoint)
    const response = await fetch(
      `${kratos_url}/self-service/verification/flows?id=${route.query.flow}`,
      {
        method: "GET",
        headers: {
          Accept: "application/json",
        },
        credentials: "include",
      }
    );

    if (!response.ok) {
      throw new Error(`Failed to fetch CSRF token: ${response.statusText}`);
    }

    // Parse the response data
    const data = await response.json();
    csrfToken.value = data.ui.nodes.find(
      (node) => node.attributes.name === "csrf_token"
    ).attributes.value;

    data?.ui?.messages?.forEach((element) => {
      if (element.type === "error") {
        toast.error(element.text);
      } else {
        toast.success(element.text);
      }
    });
    if (data.state === "passed_challenge") {
      navigateTo(data.ui.action, { external: true });
    }
  } catch (error) {
    console.error("Error fetching flow ID and CSRF token:", error.message);
  }
};

const verificationUrl = computed(
  () =>
    `${kratos_url}/self-service/verification?token=${code.value}&flow=${route.query.flow}`
);
</script>
