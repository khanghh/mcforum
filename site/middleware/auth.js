export default defineNuxtRouteMiddleware(async () => {
  const userStore = useUserStore()
  if (!userStore.user) {
    throw createError({ statusCode: 401, statusMessage: 'Unauthorized', fatal: true })
  }
})
