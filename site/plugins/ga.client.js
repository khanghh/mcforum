// plugins/ga.client.ts
export default defineNuxtPlugin(() => {
  // Inject GA script
  useHead({
    script: [
      { src: 'https://www.googletagmanager.com/gtag/js?id=G-XCTKYPF766', async: true },
      {
        innerHTML: `
        window.dataLayer = window.dataLayer || [];
        function gtag(){dataLayer.push(arguments);}
        gtag('js', new Date());
        gtag('config', 'G-XCTKYPF766', { send_page_view: false });
      ` }
    ]
  })

  // Track SPA route changes
  const router = useRouter()
  router.afterEach((to) => {
    if (window.gtag) {
      window.gtag('event', 'page_view', { page_path: to.fullPath })
    }
  })
})
