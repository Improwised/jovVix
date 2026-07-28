import { useUsersStore } from "~~/store/users";

export const useEmailVerified = () => useState("email-verified", () => null);

export const refreshEmailVerified = async () => {
  const emailVerified = useEmailVerified();
  const { apiUrl } = useRuntimeConfig().public;
  const { getUserData } = useUsersStore();

  if (getUserData()?.role !== "admin-user") {
    emailVerified.value = null;
    return;
  }

  try {
    const response = await $fetch(`${apiUrl}/kratos/whoami`, {
      method: "GET",
      credentials: "include",
    });
    emailVerified.value =
      !!response?.data?.identity?.verifiable_addresses?.some(
        (address) => address.verified
      );
  } catch (error) {
    console.log("unable to read email verification status", error.message);
    emailVerified.value = null;
  }
};
