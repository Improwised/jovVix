<script setup>
import { useToast } from "vue-toastification";

const route = useRoute();
const router = useRouter();
const toast = useToast();
const firstname = ref();
const lastname = ref();
const email = ref();
const password = ref();
const csrfToken = ref();
const component = ref("waiting");
const urls = useRuntimeConfig().public;
const errors = ref({
  email: "",
  password: "",
  firstname: "",
  lastname: "",
});
const registerURLWithFlowQuery = ref("");
const kratos_url = urls.kratos_url;
console.log(urls); // remove this after checking in production
console.log(); // this console.log is required because without this, nuxt will give 5xx error as async function is called afterwards

(async () => {
  if (process.client) {
    try {
      await $fetch(kratos_url + "/sessions/whoami", {
        method: "GET",
        credentials: "include",
        headers: {
          Accept: "application/json",
        },
        onResponse({ response }) {
          console.log(response);
          if (response.status >= 200 && response.status < 300) {
            toast.success("You are already logged in");
            navigateTo("/");
          }
        },
      });
    } catch (error) {
      console.log(); // this error will be produced when user is not logged in so not logging anything here or not doing anything here
    }

    if (route.query.flow) {
      try {
        await $fetch(kratos_url + "/self-service/registration/flows", {
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
      kratos_url + "/self-service/registration/browser",
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

// async function register_user(e) {
// e.preventDefault();

// if (password.value !== confirmPassword.value) {
//   toast.error("password and confirmPassword are not equal, please try again")
//   return
// }

// if (email.value.trim() == "" || password.value.trim() == "" || confirmPassword.value.trim() == "" || firstname.value.trim() == "") {
//   toast.error(app.$IncorrectCredentials);
//   return;
// }

// console.log(route.query.flow);
// const formdata = new FormData();
// formdata.append("csrf_token", csrfToken.value)
// formdata.append("traits.email", email.value)
// formdata.append("traits.name.first", firstname.value)
// formdata.append("traits.name.last", lastname.value)
// formdata.append("password", password.value)
// formdata.append("method", "password")

// const { error: error } = await useFetch(register_url, {
//   method: "POST",
//   credentials: "include",
//   query:{
//     flow: route.query.flow
//   },
//   headers: {
//     // "Access-Control-Allow-Origin": "127.0.0.1:3000",
//     "Content-Type": "application/json",
//     Accept: "application/json",
//   },

//   body: {
//     csrf_token: csrfToken.value,
//     traits: {
//       email: email.value,
//       name: {
//         first: firstname.value,
//         last: lastname.value
//       }
//     },
//     password: password.value,
//     method: "password",
//   },
//   // body: formdata,

//   // mode: "cors",
//   onResponseError: function ({ request, response }) {
//     status = response.status;
//     console.log(response);
//     console.log(request.headers);
//     response?._data?.ui?.nodes?.forEach(node => {
//       if (node.attributes.name === "traits.email") {
//         errors.value.email = node.messages.map(message => message.text).join(", ");
//       }
//       if (node.attributes.name === "password") {
//         errors.value.password = node.messages.map(message => message.text).join(", ");
//       }
//       if (node.attributes.name === "traits.name.first") {
//         errors.value.firstname = node.messages.map(message => message.text).join(", ");
//       }
//       if (node.attributes.name === "traits.name.last") {
//         errors.value.lastname = node.messages.map(message => message.text).join(", ");
//       }
//     });

//   },
//   onResponse: function (response, request){
//     console.log(response);
//     console.log(request);
//   }
// });

// }
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
          required
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
          required
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
          type="email"
          name="traits.email"
          class="form-control"
          required
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
          required
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
