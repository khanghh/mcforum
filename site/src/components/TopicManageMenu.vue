<template>
  <div ref="root" class="relative inline-block text-left">
    <button type="button"
      :class="['gap-2 px-2 py-1 rounded text-sm border border-purple-500/20 hover:bg-purple-600/20', open ? 'bg-purple-600/20' : 'text-purple-300']"
      aria-haspopup="true"
      :aria-label="$t('publish.manage')"
      @click="open = !open">
      <Icon name="Fa7SolidGear" class="mr-2" />
      <span class="whitespace-nowrap">{{ $t('publish.manage') }}</span>
    </button>

    <transition name="fade">
      <div v-if="open"
        class="origin-top-right absolute right-0 mt-2 w-40 rounded-md shadow-lg bg-white ring-1 ring-black ring-opacity-5 z-50">
        <div class="py-1">

          <button v-if="isOwnerOrAdmin && isPending"
            class="w-full text-left text-sm text-gray-700 px-3 py-2 hover:bg-gray-100 flex items-center"
            @click="open = false; approveTopic()">
            <Icon name="Fa7SolidCheck" class="mr-2" />
            {{ $t('publish.action.approve') }}
          </button>

          <button v-if="isOwnerOrAdmin && isPending"
            class="w-full text-left text-sm text-gray-700 px-3 py-2 hover:bg-gray-100 flex items-center"
            @click="open = false; rejectTopic()">
            <Icon name="Fa7SolidTimes" class="mr-2" />
            {{ $t('publish.action.reject') }}
          </button>

          <button v-if="isTopicOwner || isOwnerOrAdmin"
            class="w-full text-left text-sm text-gray-700 px-3 py-2 hover:bg-gray-100 flex items-center"
            @click="open = false; navigateTo(`/topics/edit/${topicId}`)">
            <Icon name="Fa7SolidEdit" class="mr-2" />
            {{ $t('publish.action.edit') }}
          </button>

          <button v-if="isOwnerOrAdmin && !isPending"
            class="w-full text-left text-sm text-gray-700 px-3 py-2 hover:bg-gray-100 flex items-center"
            @click="open = false; toggleRecommended()">
            <Icon :name="recommended ? 'PhSealCheckFill' : 'PhSealCheckBold'" class="mr-2" />
            {{ recommended ? $t('publish.action.unrecommend') : $t('publish.action.recommend') }}
          </button>

          <button v-if="isOwnerOrAdmin && !isPending"
            class="w-full text-left text-sm text-gray-700 px-3 py-2 hover:bg-gray-100 flex items-center"
            @click="open = false; togglePinned()">
            <Icon :name="pinned ? 'TablerPinFilled' : 'TablerPin'" class="mr-2" />
            {{ pinned ? $t('publish.action.unpin') : $t('publish.action.pin') }}
          </button>

          <button v-if="isTopicOwner || isOwnerOrAdmin"
            class="w-full text-left text-sm text-gray-700 px-3 py-2 hover:bg-gray-100 flex items-center"
            @click="open = false; deleteTopic()">
            <Icon name="Fa7SolidTrashCan" class="mr-2" />
            {{ $t('publish.action.delete') }}
          </button>

        </div>
      </div>
    </transition>
  </div>
</template>

<script setup lang="ts">
import { TopicStatus, type Topic } from '@/types'
const i18n = useI18n()
const api = useApi()
const dialog = useConfirmDialog()
const toast = useToast()

const props = defineProps({
  modelValue: {
    type: Object as PropType<Topic>,
    required: true,
  },
})

const topic = ref(props.modelValue)
const topicId = computed(() => {
  const parts = topic.value.slug.split('.')
  return parts.length > 1 ? parts[1] : topic.value.id?.toString() || ''
})

const emit = defineEmits(['update:modelValue'])

// dropdown state
const open = ref(false)
const root = ref<HTMLElement | null>(null)

const userStore = useUserStore()

const isOwnerOrAdmin = userIsOwner(userStore.user) || userIsAdmin(userStore.user)
const isTopicOwner = computed(() => userStore.user && userStore.user.id === topic.value.user.id)
const recommended = computed(() => topic.value.recommended)
const pinned = computed(() => topic.value.pinned)
const isPending = computed(() => topic.value.status === TopicStatus.PendingReview)

function onClickOutside(e: any) {
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

function approveTopic() {
  dialog.show({
    title: i18n.t('dialog.title.confirm_approve'),
    message: i18n.t('dialog.message.confirm_approve_post'),
    confirmText: i18n.t('dialog.button.confirm'),
    cancelText: i18n.t('dialog.button.cancel'),
    variant: 'warning',
    icon: 'Fa7SolidCheck',
    onConfirm() {
      return api.approveTopic(topic.value.slug)
        .then(() => {
          topic.value.status = TopicStatus.Active
          emit('update:modelValue', topic.value)
        }).catch((e) => {
          const errMsg = e.data?.error.message || e.message || e
          const msg = i18n.t('message.action_failure', { error: errMsg })
          toast.error(msg)
        })
    },
  })
}

function rejectTopic() {
  dialog.show({
    title: i18n.t('dialog.title.confirm_reject'),
    message: i18n.t('dialog.message.confirm_reject_post'),
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
