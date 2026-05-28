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
      baseURL: config.public.userApiBase as string,
      headers: getHeaders(),
    }),
    core: $fetch.create({
      baseURL: config.public.coreApiBase as string,
      headers: getHeaders(),
    }),
    booking: $fetch.create({
      baseURL: config.public.bookingApiBase as string,
      headers: getHeaders(),
    }),
  };
};