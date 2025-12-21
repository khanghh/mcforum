<template>
  <div class="dialog-component">
    <div
      class="dialog-mask"
      :class="{ visible: visible }"
      :style="{ zIndex: zIndex }" />
    <transition
      appear
      enter-active-class="animate__animated animate__fadeIn"
      leave-active-class="animate__animated animate__fadeOut"
      mode="out-in">
      <div
        v-if="visible"
        class="dialog-wrapper"
        :style="{ zIndex: zIndex + 1 }">
        <div
          class="dialog-content"
          :style="{
            width: dialogContentWidth,
            maxWidth: dialogContentMaxWidth,
          }">
          <div class="dialog-header">
            <div class="dialog-title">
              {{ title }}
            </div>
            <div class="dialog-close">
              <icon name="X" @click="close" />
            </div>
          </div>
          <div class="dialog-main">
            <slot />
          </div>
          <div
            class="dialog-footer"
            :style="{ justifyContent: btnsCenter ? 'center' : 'flex-end' }">
            <el-button v-if="cancelBtnVisible" @click="cancel">
              {{ $t('dialog.button.cancel') }}
            </el-button>
            <el-button v-if="okBtnVisible" type="primary" @click="ok">
              {{ okBtnText }}
            </el-button>
          </div>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
const i18n = useI18n()
const emits = defineEmits(['update:visible', 'ok', 'cancel'])
defineExpose({
  show,
  close,
})

const props = defineProps({
  visible: {
    type: Boolean,
    default: false,
  },
  // 宽度
  width: {
    type: Number,
    default: 520, // 0表示100%
  },
  // 最大宽度
  maxWidth: {
    type: Number,
    default: 0, // 0表示不设置
  },
  zIndex: {
    type: Number,
    default: 1001,
  },
  title: {
    type: String,
    default: '',
  },
  btnsCenter: {
    type: Boolean,
    default: false,
  },
  cancelBtnVisible: {
    type: Boolean,
    default: true,
  },
  okBtnVisible: {
    type: Boolean,
    default: true,
  },
  okBtnText: {
    type: String,
    default: '',
  },
})

const okBtnText = computed(() => props.okBtnText || i18n.t('dialog.button.confirm'))

const dialogContentWidth = computed(() => {
  if (props.width > 0) {
    return `${props.width}px`
  }
  // return 'auto'
  return '100%'
})
const dialogContentMaxWidth = computed(() => {
  if (props.maxWidth > 0) {
    return `${props.maxWidth}px`
  }
  return ''
})

function show() {
  emits('update:visible', true)
}
function close() {
  emits('update:visible', false)
}
function ok() {
  emits('ok')
}
function cancel() {
  emits('cancel')
  close()
}
</script>

<style lang="scss" scoped>
.dialog-mask {
  // transition: all 2s;
  position: fixed;
  top: 0;
  right: 0;
  bottom: 0;
  left: 0;
  overflow: auto;
  background: #000000;
  opacity: 0.45;
  display: none;

  &.visible {
    display: block;
  }
}

.dialog-wrapper {
  // transition: all 2s;
  position: fixed;
  top: 0;
  right: 0;
  bottom: 0;
  left: 0;

  display: flex;
  align-items: center;
  justify-content: center;

  .dialog-content {
    position: relative;
    margin: 0 auto;
    // margin-top: 15vh;
    background: #ffffff;
    border-radius: 8px;

    padding: 24px;

    .dialog-header {
      display: flex;
      align-items: center;
      justify-content: space-between;

      .dialog-title {
        font-size: 16px;
        font-weight: 500;
        color: var(--text-color);
      }

      .dialog-close {
        cursor: pointer;
        padding: 0 0 0 10px;

        svg:hover {
          color: red;
        }
      }
    }

    .dialog-main {
      padding: 12px 0;
    }

    .dialog-footer {
      display: flex;
      align-items: center;
      // justify-content: flex-end;
      column-gap: 12px;

      .chaitin-btn {
        width: 78px;
      }
    }
  }
}
</style>
