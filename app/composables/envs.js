export const useSystemEnv = () => {
  const data = {
    base_url: process.env.BASE_URL,
    api_url: process.env.API_URL,
  };
  return useState(() => data);
};
