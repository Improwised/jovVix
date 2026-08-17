export const readApiError = (error, fallback) =>
  error?.data?.data || error?.data?.message || error?.message || fallback;
