<template>
  <div v-if="menus && menus.length" ref="root" class="relative inline-block text-left">
    <button type="button"
      :class="['gap-2 px-2 py-1 rounded text-sm border border-purple-500/20 hover:bg-purple-600/20', open ? 'bg-purple-600/20' : 'text-purple-300']"
      aria-haspopup="true"
      :aria-expanded="String(open)"
      :aria-label="$t('publish.manage')"
      @click="open = !open">
      <Icon name="TablerSettings" class="w-4 h-4 mr-2" />
      <span class="whitespace-nowrap">{{ $t('publish.manage') }}</span>
    </button>

    <transition name="fade">
      <div v-if="open"
        class="origin-top-right absolute right-0 mt-2 w-40 rounded-md shadow-lg bg-white ring-1 ring-black ring-opacity-5 z-50">
        <div class="py-1">
          <button v-for="item in menus"
            :key="item.command"
            class="w-full text-left text-sm text-gray-700 px-3 py-2 hover:bg-gray-100"
            @click="select(item.command)">
            {{ item.label }}
          </button>
        </div>
      </div>
    </transition>
  </div>
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

const emit = defineEmits(['update:modelValue', 'onSwitchEditMode'])

// dropdown state
const open = ref(false)
const root = ref(null)

const userStore = useUserStore()
const isOwner = userIsOwner(userStore.user)
const isAdmin = userIsAdmin(userStore.user)

const menus = computed(() => {
  const isTopicOwner = userStore.user && userStore.user.id === topic.value.user.id
  const items = []
  if (isTopicOwner) {
    items.push({
      command: 'edit',
      label: topic.value.editing ? i18n.t('publish.action.save') : i18n.t('publish.action.edit'),
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
      label: topic.value.recommended ? i18n.t('publish.action.unrecommend') : i18n.t('publish.action.recommend'),
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

function select(command) {
  open.value = false
  handleCommand(command)
}

function onClickOutside(e) {
  if (!root.value) return
  const path = e.composedPath ? e.composedPath() : (e.path || [])
  if (path.length) {
    if (!path.includes(root.value)) open.value = false
  }
  else {
    if (!root.value.contains(e.target)) open.value = false
  }
}

onMounted(() => window.addEventListener('click', onClickOutside))
onBeforeUnmount(() => window.removeEventListener('click', onClickOutside))

async function handleCommand(command) {
  if (command === 'edit') {
    toggleEditing()
  }
  else if (command === 'delete') {
    deleteTopic()
  }
  else if (command === 'recommend') {
    toggleRecommended()
  }
  else if (command === 'pin') {
    togglePinned()
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
    await useHttpPostForm('/api/users/forbidden', {
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
    useHttpDelete(`/api/topics/${topic.value.id}`).then(() => {
      useMsg({
        message: i18n.t('message.delete_success'),
        onClose() {
          useLinkTo('/')
        },
      })
    }).catch((e) => {
      useMsgError(i18n.t('message.delete_failure', { error: (e.message || e) }))
    })
  })
}

function toggleEditing() {
  if (topic.value.editing) {
    const action = i18n.t('publish.action.edit')
    useConfirm(i18n.t('dialog.message.confirm_action_post', { action })).then(() => {
      emit('onSwitchEditMode')
    })
  }
  else {
    emit('onSwitchEditMode')
  }
}

function toggleRecommended() {
  const action = topic.value.recommended ? i18n.t('publish.action.unrecommend') : i18n.t('publish.action.recommend')
  console.log(topic.value.recommended)
  useConfirm(i18n.t('dialog.message.confirm_action_post', { action })).then(function () {
    useHttpPatchForm(`/api/topics/${topic.value.slug}`, {
      body: { recommended: !topic.value.recommended },
    }).then(() => {
      topic.value.recommended = !topic.value.recommended
      emit('update:modelValue', topic.value)
      useMsgSuccess({ message: i18n.t('message.action_success', { action }) })
    }).catch((e) => {
      useMsgError(i18n.t('message.action_failure', { action, error: e.message || e }))
    })
  })
}

function togglePinned() {
  const action = topic.value.pinned ? i18n.t('publish.action.unpin') : i18n.t('publish.action.pin')
  useConfirm(i18n.t('dialog.message.confirm_action_post', { action })).then(function () {
    useHttpPatchForm(`/api/topics/${topic.value.slug}`, {
      body: { pinned: !topic.value.pinned },
    }).then(() => {
      topic.value.pinned = !topic.value.pinned
      emit('update:modelValue', topic.value)
      useMsgSuccess({ message: i18n.t('message.action_success', { action: action }) })
    }).catch((e) => {
      useMsgError(i18n.t('message.action_failure', { action: action, error: e.message || e }))
    })
  })
}
</script>
