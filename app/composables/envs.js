export const useSystemEnv = () => {
  const data = {
    base_url: process.env.BASE_URL,
    api_url: process.env.API_URL,
    socket_url: process.env.API_SOCKET_URL,
  };
  return useState("urls", () => data);
};

export const useSystemEnv1 = () => {
  const data = {
    base_url: process.env.BASE_URL,
    api_url: process.env.API_URL,
  };
  return useState("envs", () => data);
};
