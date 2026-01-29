import { useUsersStore } from "~~/store/users";


export const handleLogout = async () => {
  const userData = useUsersStore();
  const { setUserData } = userData;
  const { kratosUrl } = useRuntimeConfig().public;
  try {
    // Step 1: Fetch logout URL and token from the first API endpoint
    const response = await fetch(`${kratosUrl}/self-service/logout/browser`, {
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

export const setUserDataStore = async () => {
  const userData = useUsersStore();
  const { setUserData } = userData;
  const { apiUrl } = useRuntimeConfig().public;
  const headers = useRequestHeaders(["cookie"]);
  try {
    const response = await fetch(apiUrl + "/user/who", {
      method: "GET",
      credentials: "include",
      headers: headers,
      mode: "cors",
    });
    if (response.status != 200) {
      throw new Error(response.status);
    } else if (response.status == 200) {
      const data = await response.json();
      setUserData({ role: data?.data?.role, avatar: data?.data?.avatar });
    }
  } catch (error) {
    if (error.message == 401) {
      console.log(error.message);
      setUserData(null);
      return;
    }
    console.log(error.message);
  }
};
