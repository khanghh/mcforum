<script setup></script>

<template>
  <footer class="footer">
    <div class="content container">
      <div>
        <nuxt-link to="/about">{{ $t('footer.about') }}&nbsp;</nuxt-link>
        <nuxt-link to="/links">{{ $t('footer.links') }}</nuxt-link>
      </div>
      <div>
        © 2022 Powered by
        <a href="https://bbs-go.com" target="_blank" class="light">BBS-GO</a>
      </div>
    </div>
  </footer>
</template>

<style lang="scss" scoped>
.footer {
  font-size: 14px;
  color: var(--text-color3);
  background: none;
  text-align: left;
  margin: 0 10px;

  a {
    color: var(--text-color3);
    text-decoration: none;
  }

  .light {
    color: #eb5424; // TODO
    font-weight: bold;
  }
}
</style>
