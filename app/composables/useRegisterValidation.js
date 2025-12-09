import { reactive } from "vue";

export function useRegisterValidation() {
  const errors = reactive({
    firstname: "",
    lastname: "",
    email: "",
    password: "",
  });

  function resetErrors() {
    errors.firstname = "";
    errors.lastname = "";
    errors.email = "";
    errors.password = "";
  }

  function validate(form) {
    resetErrors();

    let valid = true;

    // First Name
    if (!form.firstname?.trim()) {
      errors.firstname = "First name is required.";
      valid = false;
    }

    // Last Name
    if (!form.lastname?.trim()) {
      errors.lastname = "Last name is required.";
      valid = false;
    }

    // Email
    if (!form.email?.trim()) {
      errors.email = "Email is required.";
      valid = false;
    } else if (!/^\S+@\S+\.\S+$/.test(form.email)) {
      errors.email = "Enter a valid email.";
      valid = false;
    }

    // Password (single combined rule)
    const passwordRegex =
      /^(?=.*[A-Z])(?=.*[a-z])(?=.*\d)(?=.*[^A-Za-z0-9]).{8,}$/;

    if (!form.password?.trim()) {
      errors.password = "Password is required.";
      valid = false;
    } else if (!passwordRegex.test(form.password)) {
      errors.password =
        "Password must contain 8 characters, 1 uppercase, 1 lowercase, 1 number, and 1 special character.";
      valid = false;
    }

    return valid;
  }

  return {
    errors,
    validate,
    resetErrors,
  };
}
