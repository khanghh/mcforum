<template>
  <div v-if="isTopicOwner || isOwner || isAdmin" ref="root" class="relative inline-block text-left">
    <button type="button"
      :class="['gap-2 px-2 py-1 rounded text-sm border border-purple-500/20 hover:bg-purple-600/20', open ? 'bg-purple-600/20' : 'text-purple-300']"
      aria-haspopup="true"
      :aria-expanded="String(open)"
      :aria-label="$t('publish.manage')"
      @click="open = !open">
      <Icon name="Fa7SolidGear" class="mr-2" />
      <span class="whitespace-nowrap">{{ $t('publish.manage') }}</span>
    </button>

    <transition name="fade">
      <div v-if="open"
        class="origin-top-right absolute right-0 mt-2 w-40 rounded-md shadow-lg bg-white ring-1 ring-black ring-opacity-5 z-50">
        <div class="py-1">
          <button v-if="isTopicOwner"
            class="w-full text-left text-sm text-gray-700 px-3 py-2 hover:bg-gray-100 flex items-center"
            @click="open = false; toggleEditing()">
            <Icon name="Fa7SolidFileEdit" class="mr-2" />
            {{ $t('publish.action.edit') }}
          </button>

          <button v-if="isTopicOwner || isOwner || isAdmin"
            class="w-full text-left text-sm text-gray-700 px-3 py-2 hover:bg-gray-100 flex items-center"
            @click="open = false; deleteTopic()">
            <Icon name="Fa7SolidTrashCan" class="mr-2" />
            {{ $t('publish.action.delete') }}
          </button>

          <button v-if="isOwner || isAdmin"
            class="w-full text-left text-sm text-gray-700 px-3 py-2 hover:bg-gray-100 flex items-center"
            @click="open = false; toggleRecommended()">
            <Icon :name="recommended ? 'PhSealCheckFill' : 'PhSealCheckBold'" class="mr-2" />
            {{ recommended ? $t('publish.action.unrecommend') : $t('publish.action.recommend') }}
          </button>

          <button v-if="isOwner || isAdmin"
            class="w-full text-left text-sm text-gray-700 px-3 py-2 hover:bg-gray-100 flex items-center"
            @click="open = false; togglePinned()">
            <Icon :name="pinned ? 'TablerPinFilled' : 'TablerPin'" class="mr-2" />
            {{ pinned ? $t('publish.action.unpin') : $t('publish.action.pin') }}
          </button>

          <button v-if="isOwner || isAdmin"
            class="w-full text-left text-sm text-gray-700 px-3 py-2 hover:bg-gray-100 flex items-center"
            @click="open = false; forbidden(7)">
            <Icon name="TablerBan" class="mr-2" />
            {{ $t('profile.actions.mute_7days') }}
          </button>

          <button v-if="isOwner"
            class="w-full text-left text-sm text-gray-700 px-3 py-2 hover:bg-gray-100 flex items-center"
            @click="open = false; forbidden(-1)">
            <Icon name="TablerBan" class="mr-2" />
            {{ $t('profile.actions.mute_permanent') }}
          </button>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
const i18n = useI18n()
const api = useApi()
const dialog = useConfirmDialog()
const toast = useToast()

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
const isTopicOwner = computed(() => !!userStore.user && !!topic.value && !!topic.value.user && userStore.user.id === topic.value.user.id)
const editing = computed(() => !!topic.value && !!topic.value.editing)
const recommended = computed(() => !!topic.value && !!topic.value.recommended)
const pinned = computed(() => !!topic.value && !!topic.value.pinned)

// Note: menu actions are called directly from the template now.

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
  dialog.show({
    title: i18n.t('dialog.title.confirm_delete'),
    message: i18n.t('dialog.message.confirm_delete_post'),
    confirmText: i18n.t('dialog.button.confirm'),
    cancelText: i18n.t('dialog.button.cancel'),
    variant: 'warning',
    icon: 'Fa7SolidTrashCan',
    onConfirm() {
      return api.deleteTopic(topic.value.slug)
        .then(() => {
          toast.success(i18n.t('message.delete_success'))
          navigateTo(topic.value?.forum?.slug ? `/forums/${topic.value.forum.slug}` : '/')
        }).catch((e) => {
          const errMsg = e.data?.error.message || e.message || e
          const msg = i18n.t('message.delete_failure', { error: errMsg })
          toast.error(msg)
        })
    },
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
  return api.setTopicFlags(topic.value.slug, { recommended: !topic.value.recommended })
    .then(() => {
      topic.value.recommended = !topic.value.recommended
      emit('update:modelValue', topic.value)
      toast.success(i18n.t('message.action_success', { action }))
    }).catch((e) => {
      const errMsg = e.data?.error.message || e.message || e
      const msg = i18n.t('message.action_failure', { action, error: errMsg })
      toast.error(msg)
    })
}

function togglePinned() {
  const action = topic.value.pinned ? i18n.t('publish.action.unpin') : i18n.t('publish.action.pin')
  return api.setTopicFlags(topic.value.slug, { pinned: !topic.value.pinned })
    .then(() => {
      topic.value.pinned = !topic.value.pinned
      emit('update:modelValue', topic.value)
      toast.success(i18n.t('message.action_success', { action }))
    }).catch((e) => {
      const errMsg = e.data?.error.message || e.message || e
      const msg = i18n.t('message.action_failure', { action, error: errMsg })
      toast.error(msg)
    })
}
</script>
