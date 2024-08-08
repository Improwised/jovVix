import { useUsersStore } from "~~/store/users";

export const handleLogout = async () => {
  const { kratos_url } = useRuntimeConfig().public;
  const userData = useUsersStore();
  const { setUserData } = userData;

  try {
    // Step 1: Fetch logout URL and token from the first API endpoint
    const response = await fetch(`${kratos_url}/self-service/logout/browser`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      credentials: "include",
    });

    if (!response.ok) {
      throw new Error("Failed to fetch logout URL and token");
    }

    const { logout_url } = await response.json();

    // Step 2: Use the fetched URL and token to make another API call
    const secondApiResponse = await fetch(`${logout_url}`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      credentials: "include",
    });

    if (!secondApiResponse.ok) {
      throw new Error(
        "Failed to perform logout with the provided URL and token"
      );
    }
    console.log("Logged out successfully");
    setUserData(null);
    navigateTo("/");
  } catch (error) {
    console.error("Error during logout:", error);
  }
};
