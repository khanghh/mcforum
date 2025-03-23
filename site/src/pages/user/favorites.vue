<template>
  <div class="widget no-margin">
    <div class="widget-header">
      <span>
        <icon name="Star" />
        <span>&nbsp;{{ $t('widget.title.favorites') }}</span>
      </span>
    </div>

    <div class="widget-content">
      <load-more-async v-slot="{ items }" url="/api/user/favorites">
        <ul v-if="items && items.length" class="favorite-list">
          <li
            v-for="item in items"
            :key="item.favoriteId"
            class="favorite-item">
            <div v-if="item.deleted" class="favorite-item">
              <div class="favorite-summary">
                {{ $t('message.content_deleted') }}
              </div>
            </div>
            <div v-else>
              <div class="favorite-title">
                <a :href="item.url">{{ item.title }}</a>
              </div>
              <div class="favorite-summary">
                {{ item.content }}
              </div>
              <div class="favorite-meta">
                <span class="favorite-meta-item">
                  <nuxt-link :to="'/user/' + item.user.id">
                    {{
                      item.user.nickname
                    }}
                  </nuxt-link>
                </span>
                <span class="favorite-meta-item"><time>{{ usePrettyDate(item.createTime) }}</time></span>
              </div>
            </div>
          </li>
        </ul>
      </load-more-async>
    </div>
  </div>
</template>

<script setup>
const i18n = useI18n()
definePageMeta({
  layout: 'ucenter',
  middleware: ['auth'],
})

useHead({
  title: useSiteTitle(i18n.t('page.favorites')),
})
</script>

<style lang="scss" scoped>
.favorite-list {
  margin: 0 !important;

  .favorite-item {
    overflow: hidden;
    zoom: 1;
    line-height: 24px;

    padding: 8px 0;
    zoom: 1;
    position: relative;
    overflow: hidden;

    &:not(:last-child) {
      border-bottom: 1px solid var(--border-color);
    }

    &.more {
      text-align: center;

      a {
        font-size: 15px;
        font-weight: bold;
      }
    }

    .favorite-title {
      a {
        color: var(--text-color3);
        font-weight: normal;
        overflow: hidden;
        text-overflow: ellipsis;
        font-size: 18px;
        line-height: 30px;
      }
    }

    .favorite-summary {
      color: var(--text-color);
      overflow: hidden;
      display: -webkit-box;
      -webkit-box-orient: vertical;
      -webkit-line-clamp: 3;
      text-align: justify;
      padding-top: 6px;
      word-break: break-all;
      text-overflow: ellipsis;
      font-size: 14px;
    }

    .favorite-meta {
      display: inline-block;
      font-size: 13px;
      padding-top: 6px;

      .favorite-meta-item {
        padding: 0 6px 0 0;
      }

      a {
        color: var(--text-link-color);
      }

      span {
        color: var(--text-color3);
      }
    }
  }
}
</style>
