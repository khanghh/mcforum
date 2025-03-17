<template>
  <div class="widget">
    <div class="widget-header">{{ $t("profile.publish.statistics") }}</div>
    <div class="widget-content extra-info">
      <ul>
        <li>
          <span>{{ $t("profile.publish.points") }}</span><br />
          <nuxt-link to="/user/scores">
            <b>{{ user.score }}</b>
          </nuxt-link>
        </li>
        <li>
          <span>{{ $t("profile.publish.topics") }}</span><br />
          <b>{{ user.topicCount }}</b>
        </li>
        <li>
          <span>{{ $t("profile.publish.comments") }}</span><br />
          <b>{{ user.commentCount }}</b>
        </li>
      </ul>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    user: {
      type: Object,
      required: true,
    },
  },
};
</script>

<style lang="scss" scoped>
.extra-info {
  ul {
    display: flex;

    li {
      width: 100%;
      text-align: center;

      span {
        font-size: 13px;
        font-weight: 400;
        color: var(--text-color3);
      }

      a,
      b {
        color: var(--text-color);
      }
    }
  }
}
</style>
