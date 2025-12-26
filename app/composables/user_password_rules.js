import { computed } from "vue";

export function useUserPasswordRules(passwordRef, firstnameRef, lastnameRef) {
  const rules = {
    minLength: 8,
    alphanumeric: /^(?=.*[a-zA-Z])(?=.*\d)/,
    specialChar: /[^a-zA-Z0-9]/,
  };

  const passwordErrors = computed(() => {
    const errors = [];

    if (!passwordRef.value) {
      return errors;
    }

    if (passwordRef.value.length < rules.minLength) {
      errors.push("Password must be at least 8 characters long");
    }

    if (!rules.alphanumeric.test(passwordRef.value)) {
      errors.push("Password must contain letters and numbers");
    }

    if (!rules.specialChar.test(passwordRef.value)) {
      errors.push("Password must contain at least one special character");
    }

    const first = firstnameRef?.value?.toLowerCase();
    const last = lastnameRef?.value?.toLowerCase();

    if (
      (first && passwordRef.value.toLowerCase().includes(first)) ||
      (last && passwordRef.value.toLowerCase().includes(last))
    ) {
      errors.push("Password must not contain your name");
    }

    return errors;
  });

  return {
    passwordErrors,
  };
}