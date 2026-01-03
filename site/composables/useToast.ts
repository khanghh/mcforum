import { ref, readonly } from 'vue'

export type ToastType = 'success' | 'error' | 'warning' | 'info'

export interface Toast {
  id: string
  type: ToastType
  title?: string
  message: string
  duration?: number
  startTime?: number
  remainingTime?: number
}

interface ToastOptions {
  title?: string
  duration?: number
}

const toasts = ref<Toast[]>([])
const toastTimeouts = new Map<string, number>()

function generateId(): string {
  return `toast-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`
}

function addToast(type: ToastType, message: string, options: ToastOptions = {}) {
  const id = generateId()
  const duration = options.duration ?? 5000

  const toast: Toast = {
    id,
    type,
    title: options.title,
    message,
    duration,
    startTime: Date.now(),
    remainingTime: duration
  }

  toasts.value.push(toast)

  // Auto remove after duration
  if (duration > 0) {
    const timeoutId = window.setTimeout(() => {
      removeToast(id)
    }, duration)
    toastTimeouts.set(id, timeoutId)
  }

  return id
}

function removeToast(id: string) {
  const index = toasts.value.findIndex(t => t.id === id)
  if (index !== -1) {
    toasts.value.splice(index, 1)

    // Clear timeout if exists
    const timeoutId = toastTimeouts.get(id)
    if (timeoutId) {
      clearTimeout(timeoutId)
      toastTimeouts.delete(id)
    }
  }
}

function pauseToast(id: string) {
  const toast = toasts.value.find(t => t.id === id)
  if (toast && toast.duration && toast.duration > 0 && toast.remainingTime && toast.remainingTime > 0) {
    // Clear existing timeout
    const timeoutId = toastTimeouts.get(id)
    if (timeoutId) {
      clearTimeout(timeoutId)
      toastTimeouts.delete(id)
    }

    // Calculate remaining time
    const elapsed = Date.now() - (toast.startTime || Date.now())
    toast.remainingTime -= elapsed
  }
}

function resumeToast(id: string) {
  const toast = toasts.value.find(t => t.id === id)
  if (toast && toast.duration && toast.duration > 0 && toast.remainingTime && toast.remainingTime > 0) {
    // Reset start time
    toast.startTime = Date.now()

    // Set new timeout with remaining time
    const timeoutId = window.setTimeout(() => {
      removeToast(id)
    }, toast.remainingTime)
    toastTimeouts.set(id, timeoutId)
  }
}

function clearAll() {
  toasts.value = []
  toastTimeouts.forEach(timeoutId => clearTimeout(timeoutId))
  toastTimeouts.clear()
}

export function useToast() {
  return {
    toasts: readonly(toasts),
    addToast,
    removeToast,
    pauseToast,
    resumeToast,
    clearAll,
    success: (message: string, options?: ToastOptions) => addToast('success', message, options),
    error: (message: string, options?: ToastOptions) => addToast('error', message, options),
    warning: (message: string, options?: ToastOptions) => addToast('warning', message, options),
    info: (message: string, options?: ToastOptions) => addToast('info', message, options)
  }
}
