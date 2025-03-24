<template>
  <!-- <ClientOnly> -->
  <el-dropdown
    v-if="menus && menus.length"
    trigger="click"
    @command="handleCommand">
    <span class="el-dropdown-link">{{ $t('publish.manage') }}</span>
    <template #dropdown>
      <el-dropdown-menu>
        <el-dropdown-item
          v-for="item in menus"
          :key="item.command"
          :command="item.command">
          {{ item.label }}
        </el-dropdown-item>
      </el-dropdown-menu>
    </template>
  </el-dropdown>
  <!-- </ClientOnly> -->
</template>

<script setup>
const i18n = useI18n()
const props = defineProps({
  modelValue: {
    type: Object,
    required: true,
  },
})

const topic = ref(props.modelValue)

const emits = defineEmits(['update:modelValue'])

const userStore = useUserStore()
const isOwner = userIsOwner(userStore.user)
const isAdmin = userIsAdmin(userStore.user)

const menus = computed(() => {
  const isTopicOwner
    = userStore.user && userStore.user.id === topic.value.user.id
  const items = []
  if (isTopicOwner && topic.value.type === 0) {
    items.push({
      command: 'edit',
      label: i18n.t('publish.action.edit'),
    })
  }
  if (isTopicOwner || isOwner || isAdmin) {
    items.push({
      command: 'delete',
      label: i18n.t('publish.action.delete'),
    })
  }
  if (isOwner || isAdmin) {
    items.push({
      command: 'recommend',
      label: topic.value.recommend ? i18n.t('publish.action.unrecommend') : i18n.t('publish.action.recommend'),
    })
  }
  if (isOwner || isAdmin) {
    items.push({
      command: 'pin',
      label: topic.value.pinned ? i18n.t('publish.action.unpin') : i18n.t('publish.action.pin'),
    })
  }
  if (isOwner || isAdmin) {
    items.push({
      command: 'forbidden7Days',
      label: i18n.t('profile.actions.mute_7days'),
    })
  }
  if (isOwner) {
    items.push({
      command: 'forbiddenForever',
      label: i18n.t('profile.actions.mute_permanent'),
    })
  }
  return items
})

async function handleCommand(command) {
  if (command === 'edit') {
    editTopic()
  }
  else if (command === 'delete') {
    deleteTopic()
  }
  else if (command === 'recommend') {
    switchRecommend()
  }
  else if (command === 'pin') {
    toggleTopicPin()
  }
  else if (command === 'forbidden7Days') {
    await forbidden(7)
  }
  else if (command === 'forbiddenForever') {
    await forbidden(-1)
  }
  else {
    console.log('click on item ' + command)
  }
}
async function forbidden(days) {
  try {
    await useHttpPostForm('/api/user/forbidden', {
      body: {
        userId: topic.value.user.id,
        days,
      },
    })
    useMsgSuccess(i18n.t('message.mute_user_success'))
  }
  catch (e) {
    useMsgError(i18n.t('message.mute_user_failure', { error: e }))
  }
}
function deleteTopic() {
  useConfirm(i18n.t('dialog.message.confirm_delete_post')).then(function () {
    useHttpPost(`/api/topic/delete/${topic.value.id}`)
      .then(() => {
        useMsg({
          message: i18n.t('message.delete_success'),
          onClose() {
            useLinkTo('/topics')
          },
        })
      })
      .catch((e) => {
        useMsgError(i18n.t('message.delete_success', { error: (e.message || e) }))
      })
  })
}
function editTopic() {
  useLinkTo(`/topic/edit/${topic.value.id}`)
}
function switchRecommend() {
  const action = topic.value.recommend ? i18n.t('publish.action.unrecommend') : i18n.t('publish.action.recommend')
  useConfirm(i18n.t('dialog.message.confirm_action_post', { action })).then(function () {
    const recommend = !topic.value.recommend
    useHttpPostForm(`/api/topic/recommend/${topic.value.id}`, {
      body: {
        recommend,
      },
    })
      .then(() => {
        topic.value.recommend = recommend
        emits('update:modelValue', topic.value)
        useMsgSuccess({
          message: `${action}成功`,
        })
      })
      .catch((e) => {
        useMsgError(`${action}失败：` + (e.message || e))
      })
  })
}
function toggleTopicPin() {
  const action = topic.value.pinned ? i18n.t('publish.action.unpin') : i18n.t('publish.action.pin')
  useConfirm(i18n.t('dialog.message.confirm_action_post', { action })).then(function () {
    const pinned = !topic.value.pinned
    useHttpPostForm(`/api/topic/${topic.value.id}`, {
      body: {
        pinned: pinned,
      },
    })
      .then(() => {
        topic.value.pinned = pinned
        emits('update:modelValue', topic.value)
        useMsgSuccess({
          message: `${action}成功`,
        })
      })
      .catch((e) => {
        useMsgError(`${action}失败：` + (e.message || e))
      })
  })
}
</script>

<style lang="scss" scoped>
.el-dropdown-link {
  cursor: pointer;
  color: var(--text-color3);
  font-size: 12px;
}

.el-dropdown-menu__item {
  font-size: 12px;
}
</style>
