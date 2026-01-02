<template>
  <section class="main">
    <div class="container">
      <div class="main-body redirect">
        <div>
          <img src="/images/logo.png" style="max-width: 100px" />
        </div>
        <div style="margin: 20px 0">
          <h2>
            {{ $t('message.redirecting_message') }}
          </h2>
          <i18n-t keypath="message.external_redirect" tag="label">
            <a :href="url">{{ $t('message.click_to_redirect') }}</a>
          </i18n-t>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
const i18n = useI18n()
const route = useRoute()
const url = route.query.url || ''
const temp = url.toLowerCase()
const autoRedirectTime = 5000

if (!temp.startsWith('http://') && !temp.startsWith('https://')) {
  throw createError({
    statusCode: 500,
    message: i18n.t('message.internal_server_error'),
    fatal: true,
  })
}

onMounted(() => {
  setTimeout(() => {
    window.location.href = url
  }, autoRedirectTime)
})
</script>

<style lang="scss" scoped>
.redirect {
  text-align: center;
  vertical-align: center;
  padding: 100px 0;
}

h2 {
  font-size: 24px;
  margin-bottom: 20px;
}
</style>
