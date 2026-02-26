<template>
  <Frame page-title="Change Password" page-message="Please enter your new password information.">
    <!-- Form -->
    <form @submit.prevent="handleChangePassword">
      <div class="mb-3">
        <label for="newPassword" class="form-label">New Password</label>
        <input 
          id="newPassword" 
          v-model="password.new" 
          class="form-control mb-2" 
          type="password" 
          placeholder="Enter new password"
          required 
        />
      </div>

      <div class="mb-3">
        <label for="confirmPassword" class="form-label">Confirm Password</label>
        <input 
          id="confirmPassword" 
          v-model="password.confirm" 
          class="form-control mb-2" 
          type="password" 
          placeholder="Confirm new password"
          required 
        />
        <div v-if="passwordRequestError">
          • {{ passwordRequestError }}
        </div>
        <div v-if="passwordSubmitted && passwordErrors.length">
          <div v-for="(err, i) in passwordErrors" :key="i">
            • {{ err }}
          </div>
        </div>
      </div>
      <div>
        <div class="d-grid gap-2">
          <button 
            type="submit" 
            class="btn btn-primary text-white"
          >
            Change Password
          </button>
          <button 
            type="button" 
            class="btn btn-secondary text-white"
            @click="handleCancel"
          >
            Cancel
          </button>
        </div>
      </div>
    </form>
  </Frame>
</template>

<script setup>
import { useToast } from "vue-toastification";
import { toRef } from "vue";
import { useUserPasswordRules } from "@/composables/user_password_rules";

const toast = useToast();
const url = useRuntimeConfig().public;

const password = reactive({
  new: "",
  confirm: ""
});

const newPasswordRef = toRef(password, "new");

const { passwordErrors } = useUserPasswordRules(newPasswordRef);

const passwordRequestError = ref("");
const passwordSubmitted = ref(false);
const flow = ref("");
const csrfToken = ref("");

const fetchFlowIdAndCsrfToken = async () => {
  try {
    const response = await fetch(
      `${url.kratosUrl}/self-service/settings/browser?aal=&refresh=&return_to=`,
      {
        method: "GET",
        headers: {
          Accept: "application/json",
        },
        credentials: "include",
      }
    );

    if (!response.ok) {
      throw new Error(`Error: ${response.statusText}`);
    }

    const data = await response.json();
    flow.value = data.id;
    csrfToken.value = data.ui.nodes.find(
      (node) => node.attributes.name === "csrf_token"
    ).attributes.value;
  } catch (error) {
    console.error("Error fetching flow ID and CSRF token:", error);
  }
};

const handleChangePassword = async () => {
  passwordSubmitted.value = true;
  
  try {
    if (passwordErrors.value.length > 0) {
      return;
    }

    if (password.new !== password.confirm) {
      passwordRequestError.value = "Passwords do not match.";
      return;
    }

    await fetchFlowIdAndCsrfToken();

    const response = await fetch(
      `${url.kratosUrl}/self-service/settings?flow=${flow.value}`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Accept: "application/json",
        },
        credentials: "include",
        body: JSON.stringify({
          password: password.new,
          csrf_token: csrfToken.value,
          method: "password",
        }),
      }
    );

    if (!response.ok) {
      const errorData = await response.json();
      if (errorData.error.id === "session_refresh_required") {
        window.location.href = errorData.redirect_browser_to;
      } else {
        throw new Error(errorData.error.message);
      }
    }

    passwordRequestError.value = "";
    
    password.new = "";
    password.confirm = "";
    
    toast.success("Password updated successfully.");
    
    navigateTo("/admin");

  } catch (error) {
    passwordRequestError.value = error.message;
  }
};

const handleCancel = () => {
  navigateTo("/admin");
};
</script>
