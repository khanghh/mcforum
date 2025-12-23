<template>
  <div class="widget">
    <div class="widget-header">
      <span>{{ $t("widget.title.personal_info") }}</span>
      <div class="slot">
        <nuxt-link to="/users/me/profile">{{ $t('profile.actions.edit') }}</nuxt-link>
      </div>
    </div>
    <div class="widget-content stable">
      <div class="str">
        <div class="slabel">
          {{ $t('form.label.nickname') }}
        </div>
        <div class="svalue">
          {{ user.nickname }}
        </div>
      </div>
      <div class="str">
        <div class="slabel">
          {{ $t('form.label.bio') }}
        </div>
        <div class="svalue">
          {{ user.bio }}
        </div>
      </div>
      <div
        v-if="user.homePage"
        class="str">
        <div class="slabel">
          {{ $t("form.label.website") }}
        </div>
        <div class="svalue">
          <a :href="user.homePage"
            target="_blank"
            rel="nofollow">
            {{ user.homePage }}
          </a>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'MyProfile',
  props: {
    user: {
      type: Object,
      required: true,
    },
  },
}
</script>

<style scoped></style>
