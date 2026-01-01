<template>
  <div class="relative">
    <!-- Cover Photo -->
    <div
      class="h-48 sm:h-64 bg-gradient-to-r from-purple-900 via-indigo-900 to-blue-900 rounded-2xl overflow-hidden relative group">
      <img
        :src="coverImageSrc || 'https://images.unsplash.com/photo-1511512578047-dfb367046420?auto=format&fit=crop&w=1600&q=80'"
        class="w-full h-full object-cover opacity-80">
      <!-- Animated overlay -->
      <div class="absolute inset-0 bg-gradient-to-t from-purple-900/60 via-transparent to-transparent"></div>
      <!-- Upload button (visible on hover for profile owner) -->
      <input ref="coverInput" type="file" accept="image/*" class="sr-only" @change="onCoverFileChange" />
      <button v-if="isSelf" @click="handleCoverClick" type="button"
        class="absolute top-3 right-3 opacity-0 group-hover:opacity-100 transition-opacity text-sm font-medium text-white px-3 py-1 rounded-md bg-black/40 backdrop-blur-sm">
        <Icon name="TablerCloudUpload" class="mr-1" />
        {{ $t('profile.actions.change_cover') }}
      </button>
    </div>

    <div class="px-6 -mt-12 sm:-mt-16 relative z-10 flex items-end justify-between gap-4">
      <div class="flex flex-col sm:flex-row items-center sm:items-end gap-4 w-full">
        <div class="relative flex-shrink-0 p-1 mx-auto sm:mx-0">
          <AvatarDecorated :user="user" :is-self="isSelf">
            <AvatarEdit v-if="isSelf" v-model="user.avatar" :username="user.username"
              class="w-full h-full object-cover"
              size="120" />
            <Avatar v-else :src="user.avatar" :username="user.username"
              class="w-full h-full object-cover"
              size="120" />
          </AvatarDecorated>
        </div>
        <div class="pb-2 w-full sm:w-auto text-center sm:text-left">
          <div class="flex items-center gap-2 flex-wrap justify-center sm:justify-start">
            <h1
              class="text-2xl sm:text-3xl font-bold text-white gaming-title drop-shadow-[0_0_15px_rgba(139,92,246,0.8)]">
              {{ user.username }}
            </h1>
          </div>
          <div class="flex items-center gap-2 mt-1 justify-center sm:justify-start w-full">
            <div
              class="gaming-card min-w-0 w-full sm:w-auto px-3 py-1 rounded-lg backdrop-blur-sm flex flex-wrap items-center group relative md:inline-flex md:flex-nowrap break-all"
              :class="{ invisible: !user.statusMessage && !isSelf }">

              <span v-if="!isEditing"
                :class="['text-sm text-purple-200 font-medium pr-4 min-w-0 w-full md:w-auto break-words break-all whitespace-normal md:line-clamp-1', user.statusMessage ? '' : 'italic']">
                <Icon name="Fa7SolidQuoteLeft" class="text-purple-400 text-sm mr-1" />
                {{ user.statusMessage || 'Say something...' }}
              </span>

              <div v-if="isEditing" class="relative inline-block">
                <input ref="statusInput" v-model="editingStatus"
                  @input="updateWidth"
                  @keyup.enter="saveStatus"
                  @keyup.escape="cancelStatusMsgEditing"
                  class="bg-transparent text-sm text-purple-200 font-medium outline-none min-w-0"
                  placeholder="Say something..." maxlength="80" />
                <span ref="measureSpan" aria-hidden="true"
                  style="position:absolute; visibility:hidden; white-space:pre; font-size:0.875rem; font-weight:500;">
                </span>
              </div>

              <div v-if="isSelf" class="absolute right-2 opacity-0 group-hover:opacity-100 transition-opacity">
                <button v-if="isEditing" @click="saveStatus" type="button"
                  class="text-purple-300 hover:text-purple-100">
                  <Icon name="Fa7SolidCheck" class="text-xs" />
                </button>
                <button v-else @click="startStatusMsgEditing" type="button"
                  class="text-purple-300 hover:text-purple-100">
                  <Icon name="Fa7SolidPencil" class="text-xs" />
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div v-if="currentUser" class="hidden md:flex items-center gap-2 mb-3 whitespace-nowrap">
        <button v-if="!isSelf"
          :class="isFollowing
            ? 'px-4 py-2 border-2 border-purple-500/50 text-purple-300 rounded-lg font-bold flex items-center gaming-title text-sm hover:bg-purple-500/10 transition-colors'
            : 'px-4 py-2 bg-gradient-to-r from-purple-600 to-pink-600 text-white rounded-lg font-bold flex items-center neon-border gaming-title text-sm hover:scale-105 transition-transform'"
          :aria-label="isFollowing ? 'Unfollow user' : 'Follow user'"
          @click="isFollowing ? handleUnfollow() : handleFollow()">
          <Icon :name="isFollowing ? 'Fa7SolidUserMinus' : 'Fa7SolidUserPlus'" class="mr-2" />
          {{ isFollowing ? 'UNFOLLOW' : 'FOLLOW' }}
        </button>
        <nuxt-link v-if="isSelf" to="/users/me/profile"
          class="px-4 py-2 border-2 border-purple-500/50 text-purple-300 rounded-lg font-bold flex items-center gaming-title text-sm hover:bg-purple-500/10 transition-colors">
          <Icon name="Fa7SolidGear" class="mr-2" /> EDIT PROFILE
        </nuxt-link>
      </div>
    </div>

    <!-- Mobile Settings Button -->
    <div v-if="isSelf" class="absolute top-4 right-4 md:hidden">
      <nuxt-link to="/users/me/profile"
        class="gaming-card text-purple-300 px-4 py-2 rounded-lg text-sm font-medium backdrop-blur-sm block">
        <Icon name="Fa7SolidGear" />
      </nuxt-link>
    </div>
  </div>
</template>

<script setup>
import AvatarDecorated from '../AvatarDecorated.vue'
import { useMsgError } from '@/composables/useMsg'
import { nextTick, ref } from 'vue'

const userStore = useUserStore()
const dialog = useConfirmDialog()
const api = useApi()

const props = defineProps({
  user: {
    type: Object,
    required: true,
  },
})

const isFollowing = ref(false)
const isEditing = ref(false)
const editingStatus = ref('')
const statusInput = ref()
const measureSpan = ref()
const { user: currentUser } = storeToRefs(userStore)
const isSelf = computed(() => currentUser?.value && currentUser?.value?.id === props.user.id)
const coverImageSrc = ref(props.user.backgroundImage)
const avatar = computed({
  get: () => props.user.avatar ?? '',
  set: v => {
    if (!props.user) return
    props.user.avatar = v
  },
})

function onAvatarUpdate(val) {
  avatar.value = val
  props.user.avatar = val
  if (currentUser?.value && currentUser.value.id === props.user.id) {
    userStore.user.avatar = val
  }
}

if (!isSelf.value) {
  isFollowing.value = await api.isFollowing(props.user.id).catch(() => false)
}

const handleFollow = async () => {
  try {
    await api.followUser(props.user.username)
    isFollowing.value = true
  }
  catch (e) {
    console.error(e)
  }
}

const handleUnfollow = async () => {
  try {
    await api.unfollowUser(props.user.username)
    isFollowing.value = false
  }
  catch (e) {
    console.error(e)
  }
}

function startStatusMsgEditing() {
  isEditing.value = true
  editingStatus.value = props.user.statusMessage || ''
  updateWidth()
}

function cancelStatusMsgEditing() {
  isEditing.value = false
  editingStatus.value = ''
}

function updateWidth() {
  nextTick(() => {
    const span = measureSpan.value
    const input = statusInput.value
    if (!span || !input) return
    span.textContent = editingStatus.value || input.placeholder || ''
    const width = span.offsetWidth
    input.style.width = `${Math.max(60, width + 16)}px`
  })
}

watch(isEditing, async (val) => {
  if (val) {
    await nextTick()
    updateWidth()
    const input = statusInput.value
    if (input) {
      input.focus()
      const len = (editingStatus.value || '').length
      try { input.setSelectionRange(len, len) } catch (e) { }
    }
  }
})

watch(editingStatus, () => updateWidth())

async function saveStatus() {
  if (editingStatus.value === props.user.statusMessage) {
    isEditing.value = false
    return
  }
  try {
    await api.setStatusMessage(editingStatus.value)
    props.user.statusMessage = editingStatus.value
    isEditing.value = false
  } catch (e) {
    const errMsg = e?.response?.data?.error?.message || e.message || 'Unknown error'
    useMsgError('Status update failed')
    return
  }
}

// Cover upload handlers
const coverInput = ref(null)
const handleCoverClick = () => {
  if (coverInput?.value) coverInput.value.click()
}

let objectUrl = null
async function onCoverFileChange(e) {
  const file = e.target.files && e.target.files[0]
  if (!file) return
  if (objectUrl) URL.revokeObjectURL(objectUrl)
  objectUrl = URL.createObjectURL(file)
  await api.uploadCover(file).then((res) => {
    props.user.backgroundImage = res.coverImage
    coverImageSrc.value = objectUrl
  }).catch((err) => {
    const errMsg = err?.data?.error?.message || err.message || 'Unknown error'
    dialog.show({
      title: $t('dialog.title.error_occurred'),
      message: errMsg,
      variant: 'warning',
    })
  })
  e.target.value = ''
}
</script>

<style scoped>
.neon-border {
  box-shadow: 0 0 10px rgba(139, 92, 246, 0.5), 0 0 20px rgba(139, 92, 246, 0.3);
}
</style>
