<template>
  <div
    class="w-full flex items-center px-4 py-3 bg-gray-800/70 border border-gray-700 rounded-lg text-white placeholder-gray-400 transition-colors transition-shadow relative focus-within:ring-2 focus-within:ring-purple-600/50">
    <input id="tags" v-model="tags" name="tags" type="hidden" />

    <!-- selected tags -->
    <div class="flex flex-wrap gap-2 mr-3">
      <div
        v-for="tag in tags"
        :key="tag"
        class="h-8 px-3 bg-gray-700/60 text-gray-300 rounded-lg text-sm whitespace-nowrap flex items-center gap-1">
        <span class="text-sm">#{{ tag }}</span>
        <button type="button" class="text-gray-400 hover:text-red-500 w-3 h-3 flex items-center justify-center"
          @click="clickRemoveTag(tag)">
          <Icon name="Fa7SolidTimes" class="w-3 h-3"></Icon>
        </button>
      </div>
    </div>

    <input
      ref="tagInput"
      v-model="inputTag"
      :placeholder="$t('publish.add_some_tags')"
      class="flex-1 bg-transparent border-none p-0 m-0 focus:outline-none text-sm placeholder-gray-400"
      type="text"
      @keydown.delete="removeTag"
      @keydown.enter="addTag"
      @keydown.;="addTag"
      @keydown.,="addTag" />
  </div>
</template>

<script setup>
const props = defineProps({
  modelValue: {
    type: Array,
    default() {
      return []
    },
  },
})

const maxTagCount = 4
const maxWordCount = 15

const tagInput = ref(null)
const tags = ref(props.modelValue || [])
const inputTag = ref('')

const emits = defineEmits(['update:modelValue'])

function removeTag(event) {
  if (event.currentTarget && event.currentTarget.value) {
    return
  }
  const selectionStart = tagInput.value && tagInput.value.selectionStart
  if (!inputTag.value || selectionStart === 0) {
    tags.value.splice(tags.value.length - 1, 1)
    emits('update:modelValue', tags.value)
  }
}

function clickRemoveTag(tag) {
  if (!tag) return
  const index = tags.value.indexOf(tag)
  if (index !== -1) {
    tags.value.splice(index, 1)
    emits('update:modelValue', tags.value)
  }
}

function addTag(event) {
  if (event) {
    event.stopPropagation()
    event.preventDefault()
  }
  addTagName(inputTag.value)
}

function addTagName(tagName) {
  if (!tagName) {
    return false
  }
  if (tags.value && tags.value.length >= maxTagCount) {
    return false
  }
  if (tagName.length > maxWordCount) {
    return false
  }
  if (tags.value && tags.value.includes(tagName)) {
    return false
  }
  tags.value.push(tagName)
  inputTag.value = ''
  emits('update:modelValue', tags.value)
  return true
}
</script>
