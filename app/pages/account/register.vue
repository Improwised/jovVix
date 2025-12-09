<script setup>
import { useToast } from "vue-toastification";
import { useUsersStore } from "~~/store/users";
const userData = useUsersStore();
const { getUserData } = userData;
const route = useRoute();
const router = useRouter();
const toast = useToast();
const firstname = ref();
const lastname = ref();
const email = ref();
const password = ref();
const csrfToken = ref();
const component = ref("waiting");
const { kratosUrl } = useRuntimeConfig().public;
// const errors = ref({
//   email: "",
//   password: "",
//   firstname: "",
//   lastname: "",
// });
const { errors, validate } = useRegisterValidation();

const registerURLWithFlowQuery = ref("");
console.log(); // this console.log is required because without this, nuxt will give 5xx error as async function is called afterwards

function onSubmit(event) {
  const isValid = validate({
    firstname: firstname.value,
    lastname: lastname.value,
    email: email.value,
    password: password.value,
  });

  console.log(isValid);
  if (!isValid) {
    // Stop form submission
    return;
  }

  // If valid → submit the form manually
  event.target.submit();
}

(async () => {
  if (process.client) {
    const user = getUserData();
    if (user && user?.role == "admin-user") {
      navigateTo("/");
      return;
    }
    if (route.query.flow) {
      try {
        await $fetch(kratosUrl + "/self-service/registration/flows", {
          method: "GET",
          credentials: "include",
          headers: {
            Accept: "application/json",
          },
          query: {
            id: route.query.flow,
          },
          onResponseError({ response }) {
            if (response._data?.error?.code === 410) {
              navigateTo("/account/register");
            }

            if (response.status === 400) {
              toast.warning(
                "Please fill out the form correctly, password or email was incorrect"
              );
            }
            component.value = "form";
          },
          onResponse({ response }) {
            response?._data?.ui?.messages?.forEach((message) => {
              if (message.type === "error") {
                if (message.id === 4000007) {
                  toast.error("An account with the same email exists already!");
                } else {
                  toast.error(message.text);
                }
              }
            });
            response?._data?.ui?.nodes?.forEach((node) => {
              if (node.attributes.name === "traits.email") {
                errors.value.email = node.messages
                  .map((message) => message.text)
                  .join(", ");
              }
              if (node.attributes.name === "password") {
                errors.value.password = node.messages
                  .map((message) => message.text)
                  .join(", ");
              }
              if (node.attributes.name === "traits.name.first") {
                errors.value.firstname = node.messages
                  .map((message) => message.text)
                  .join(", ");
              }
              if (node.attributes.name === "traits.name.last") {
                errors.value.lastname = node.messages
                  .map((message) => message.text)
                  .join(", ");
              }
            });
          },
        });
        await setFlowIDAndCSRFToken();
        component.value = "form";
      } catch (error) {
        console.error(error);
        await setFlowIDAndCSRFToken();
        component.value = "form";
      }
    } else {
      await setFlowIDAndCSRFToken();
      component.value = "form";
    }
  } else {
    await setFlowIDAndCSRFToken();
    component.value = "form";
  }
})();

async function setFlowIDAndCSRFToken() {
  try {
    const kratosResponse = await $fetch(
      kratosUrl + "/self-service/registration/browser",
      {
        method: "GET",
        headers: {
          Accept: "application/json",
        },
        credentials: "include",
        onResponseError({ response }) {
          console.error(
            "error while getting the flow id from the server",
            response
          );
        },
      }
    );

    router.push("?flow=" + kratosResponse?.id);
    csrfToken.value = kratosResponse?.ui?.nodes[0]?.attributes?.value;
    registerURLWithFlowQuery.value = kratosResponse?.ui?.action;
  } catch (error) {
    console.error(error);
  }
}
</script>

<template>
  <QuizLoadingSpace v-if="component === 'waiting'"></QuizLoadingSpace>
  <Frame
    v-else
    page-title="Register"
    page-message="We Would Be Happy To Have You"
  >
    <form
      method="POST"
      @submit.prevent="onSubmit"
      :action="registerURLWithFlowQuery"
      enctype="application/json"
    >
      <div class="mb-3">
        <label for="firstname" class="form-label">First Name</label>
        <input
          id="firstname"
          v-model="firstname"
          type="text"
          name="traits.name.first"
          class="form-control"
        />
        <div v-if="errors.firstname" class="text-danger">
          {{ errors.firstname }}
        </div>
      </div>
      <div class="mb-3">
        <label for="lastname" class="form-label">Last Name</label>
        <input
          id="lastname"
          v-model="lastname"
          type="text"
          name="traits.name.last"
          class="form-control"
        />
        <div v-if="errors.lastname" class="text-danger">
          {{ errors.lastname }}
        </div>
      </div>
      <div class="mb-3">
        <label for="traits.email" class="form-label">Email</label>
        <input
          id="email"
          v-model="email"
          type="text"
          name="traits.email"
          class="form-control"
        />
        <div v-if="errors.email" class="text-danger">{{ errors.email }}</div>
      </div>
      <div class="mb-3">
        <label for="password" class="form-label">Password</label>
        <input
          id="password"
          v-model="password"
          type="password"
          name="password"
          class="form-control"
        />
        <div v-if="errors.password" class="text-danger">
          {{ errors.password }}
        </div>
      </div>
      <div>
        <input type="hidden" name="csrf_token" :value="csrfToken" />
      </div>
      <div>
        <input type="hidden" name="method" value="password" />
      </div>
      <button type="submit" class="btn text-white btn-primary">Submit</button>
      <button type="reset" class="btn text-white btn-primary ms-2">
        Clear
      </button>
    </form>
  </Frame>
</template>
