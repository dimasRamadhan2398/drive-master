import { useAuthStore } from "~/stores/auth";

export const useApiClients = () => {
  const config = useRuntimeConfig();
  const authStore = useAuthStore();

  const getHeaders = () => ({
    ...(authStore.accessToken
      ? { Authorization: `Bearer ${authStore.accessToken}` }
      : {}),
  });

  return {
    user: $fetch.create({
      baseURL: config.public.userApiBase,
      headers: getHeaders(),
    }),
    core: $fetch.create({
      baseURL: config.public.coreApiBase,
      headers: getHeaders(),
    }),
    booking: $fetch.create({
      baseURL: config.public.bookingApiBase,
      headers: getHeaders(),
    }),
  };
};