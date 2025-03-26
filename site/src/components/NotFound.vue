<template>
  <section class="main">
    <div class="container">
      <div class="main-body redirect">
        <div>
          <img src="~/assets/images/logo.png" style="max-width: 100px" />
        </div>
        <div style="margin: 20px 0">
          <h2>
            {{ $t('page.not_found') }}
          </h2>
          <i18n-t keypath="message.page_not_found" tag="label">
            <a href="/">{{ $t('message.click_to_redirect') }}</a>
          </i18n-t>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
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
