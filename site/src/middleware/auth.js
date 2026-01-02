export default defineNuxtRouteMiddleware(async (to) => {
  const userStore = useUserStore();
  if (!userStore.user) {
    // Call composables inside the middleware (top-level usage causes setup errors)
    const dialog = useConfirmDialog()
    const t = useI18n().t

    const confirmed = await dialog.show({
      title: t('dialog.title.login_required'),
      message: t('dialog.message.authentication_required'),
      confirmText: t('dialog.button.confirm'),
      cancelText: t('dialog.button.cancel'),
    })

    if (confirmed) {
      return navigateTo('/signin?redirect=' + encodeURIComponent(to.fullPath))
    }

    return abortNavigation()
  }
})
