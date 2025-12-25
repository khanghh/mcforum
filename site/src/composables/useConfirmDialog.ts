import { ref, toRefs, reactive } from 'vue'

type ConfirmDialogOptions = {
  title?: string
  message?: string
  confirmText?: string
  cancelText?: string
  onConfirm?: () => void
  onCancel?: () => void
  variant?: 'info' | 'warning'
  icon?: string
}

const visible = ref(false)
const title = ref('')
const message = ref('')
const confirmText = ref('OK')
const cancelText = ref('Cancel')
const variant = ref<'info' | 'warning'>('info')
const icon = ref<string | null | undefined>(null)

let resolver: ((value: boolean) => void) | null = null

function hide(result = false) {
  visible.value = false
  if (resolver) {
    resolver(result)
    resolver = null
  }
}

function confirm() {
  hide(true)
}

function cancel() {
  hide(false)
}

/**
 * Convenience API: show({ title, message, onConfirm, onCancel })
 * returns a Promise<boolean> and calls callbacks accordingly.
 */
function show(opts: ConfirmDialogOptions) {
  const { onConfirm, onCancel, ...rest } = opts

  if (visible.value && resolver) {
    resolver(false)
  }

  title.value = rest.title ?? ''
  message.value = rest.message ?? ''
  confirmText.value = rest.confirmText ?? 'OK'
  cancelText.value = rest.cancelText ?? 'Cancel'
  variant.value = rest.variant ?? 'info'
  icon.value = rest.icon

  visible.value = true

  return new Promise<boolean>((resolve) => {
    resolver = (v: boolean) => {
      resolve(v)
      resolver = null
      if (v) {
        onConfirm?.()
      } else {
        onCancel?.()
      }
    }
  })
}

const state = reactive({
  visible,
  title,
  message,
  confirmText,
  cancelText,
  variant,
  icon,
})

export function useConfirmDialog() {
  return {
    show,
    hide,
    confirm,
    cancel,
    state: toRefs(state),
  }
}
