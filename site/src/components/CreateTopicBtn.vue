<template>
  <!-- <client-only> -->
  <el-dropdown placement="bottom" trigger="click" @command="handlePostCommand">
    <el-button type="primary" :icon="Plus">
      {{ $t('navbar.publish') }}
    </el-button>
    <template #dropdown>
      <el-dropdown-menu>
        <el-dropdown-item v-for="(item, i) in modules" :key="i" :command="item.command" :class="item.icon">
          <icon :name="item.icon" />&nbsp;
          {{ item.name }}
        </el-dropdown-item>
      </el-dropdown-menu>
    </template>
  </el-dropdown>
  <!-- </client-only> -->
</template>

<script setup>
import { Plus } from '@element-plus/icons-vue'

const i18n = useI18n()

const modules = ref([
  {
    command: 'tweet',
    name: i18n.t('navbar.post_status'),
    icon: 'SquarePen',
  },
  {
    command: 'topic',
    name: i18n.t('navbar.create_topic'),
    icon: 'MessageSquareText',
  },
  {
    command: 'article',
    name: i18n.t('navbar.create_article'),
    icon: 'FileText',
  },
])

function handlePostCommand(cmd) {
  const router = useRouter()
  if (cmd === 'topic') {
    router.push('/t/create')
  }
  else if (cmd === 'tweet') {
    router.push('/t/create?type=1')
  }
  else if (cmd === 'article') {
    router.push('/article/create')
  }
}
</script>

<style lang="scss" scoped></style>
