<template>
  <nav class="dock-nav">
    <ul>
      <!-- <li class="dock-nav-divider"></li> -->
      <li
        v-for="forum in forums"
        :key="forum.slug"
        :class="{ active: forum.slug === slug }">
        <nuxt-link :to="forumUrl(forum.slug)">
          <img
            class="forum-logo"
            :src="forumLogo(forum)">
          <span>{{ forum.name }}</span>
        </nuxt-link>
      </li>
    </ul>
  </nav>
</template>

<script setup>

const route = useRoute()
const slug = route.params.slug || 'whats-new'

const { data: forums } = await useAsyncData('forums', () =>
  useHttpGet(`/api/forums/menu`),
)

function forumLogo(forum) {
  return forum.logo || '/images/node.png'
}

function forumUrl(slug) {
  return `/forums/${slug}`
}
</script>

<style lang="scss" scoped>
.dock-nav {
  display: block;
  position: -webkit-sticky;
  position: sticky;
  top: calc(52px + 1rem);

  width: 150px;
  border-radius: 2px;
  background-color: var(--bg-color);
  transition: all 0.2s linear;

  ul {
    height: 100%;
    display: flex;
    flex-direction: column;
    padding: 16px 12px;

    li:not(.dock-nav-divider) {
      position: relative;
      cursor: pointer;
      height: 30px;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 14px;
      color: var(--text-color);
      //padding: 0 12px;
      border-radius: 3px;
      transition: background-color 0.2s, color 0.2s;
      font-weight: 500;

      &:not(:first-child) {
        margin-top: 10px;
      }

      &.active {
        background-color: #ea6f5a;
        color: var(--text-color5);

        a {
          color: var(--text-color5);
        }
      }

      &:not(.active):hover {
        background-color: hsla(0, 0%, 94.9%, 0.6);
      }

      a {
        text-decoration: none;
        cursor: pointer;
        color: var(--text-color3);
        width: 100%;
        height: 100%;
        text-align: center;
        line-height: 30px;
        padding-left: 10px;

        display: flex;
        align-items: center;

        //justify-content: center;
        .forum-logo {
          width: 24px;
          height: 24px;
          border-radius: 4px;
          margin-right: 10px;
        }
      }
    }

    li.dock-nav-divider {
      height: 15px;
      border-bottom: 1px solid var(--border-color);
    }
  }
}
</style>
