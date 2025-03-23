<template>
  <div class="widget no-margin">
    <div class="widget-header">
      {{ $t('page.points_history') }}
    </div>
    <div class="widget-content">
      <load-more-async v-slot="{ items }" url="/api/user/score_logs">
        <ul class="score-logs">
          <li
            v-for="scoreLog in items"
            :key="scoreLog.id"
            :class="{ plus: scoreLog.type === 0 }">
            <span class="score-type">
              {{
                scoreLog.type === 0 ? $t('profile.info.points_earned') : $t('profile.info.points_deducted')
              }}
            </span>
            <span class="score-score">
              <icon name="Trophy" />
              <span>{{ scoreLog.score }}</span>
            </span>
            <span class="score-description">{{ scoreLog.description }}</span>
            <span class="score-time">@{{ useFormatDate(scoreLog.createTime) }}</span>
          </li>
        </ul>
      </load-more-async>
    </div>
  </div>
</template>

<script setup>
const i18n = useI18n()

useHead({
  title: useSiteTitle(i18n.t('page.points_history')),
})

definePageMeta({
  middleware: ['auth'],
  layout: 'ucenter',
})
</script>

<style lang="scss" scoped>
.score-logs {
  // margin-top: 10px;
  font-size: 1rem;

  li {
    line-height: 2;
    margin-bottom: 10px;

    .score-type {
      color: green;
    }

    .score-score {
      margin: 0 3px;

      span {
        font-weight: bold;
      }
    }

    .score-time {
      color: var(--text-color3);
    }

    .score-description {
      color: var(--text-color3);
    }

    &.plus {
      .score-type {
        color: red;
      }
    }
  }
}
</style>
