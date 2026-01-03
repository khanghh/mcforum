<template>
  <div
    class="relative flex flex-col rounded-lg shadow-lg text-white overflow-hidden transform transition-all duration-300 pointer-events-auto w-full max-w-sm"
    :class="[bgColor]"
    @mouseenter="onMouseEnter"
    @mouseleave="onMouseLeave">
    <!-- Main Content -->
    <div class="flex items-center p-4">
      <div class="flex-shrink-0 mr-3">
        <!-- Icons based on type -->
        <svg v-if="toast.type === 'success'" class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
          <path fill-rule="evenodd"
            d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
            clip-rule="evenodd" />
        </svg>
        <svg v-else-if="toast.type === 'error'" class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
          <path fill-rule="evenodd"
            d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
            clip-rule="evenodd" />
        </svg>
        <svg v-else-if="toast.type === 'warning'" class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
          <path fill-rule="evenodd"
            d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z"
            clip-rule="evenodd" />
        </svg>
        <svg v-else class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
          <path fill-rule="evenodd"
            d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z"
            clip-rule="evenodd" />
        </svg>
      </div>
      <div class="flex-grow">
        <p v-if="toast.title" class="font-bold text-sm">{{ toast.title }}</p>
        <p class="font-medium text-sm">{{ toast.message }}</p>
      </div>
      <button class="flex-shrink-0 ml-3 text-white hover:text-gray-200 transition-colors" @click="dismiss">
        <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
          <path fill-rule="evenodd"
            d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
            clip-rule="evenodd" />
        </svg>
      </button>
    </div>

    <!-- Progress Bar -->
    <div class="h-1 w-full bg-white/30 rounded-b-lg overflow-hidden relative">
      <div
        class="h-full w-full origin-left transition-transform linear duration-[50ms]"
        :class="[progressColor]"
        :style="{ transform: `scaleX(${progress})` }"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onBeforeUnmount, ref } from 'vue';
import { useToast, type Toast } from '@/composables/useToast';

const props = defineProps<{
  toast: Toast
}>();

const { removeToast, pauseToast, resumeToast } = useToast();

const progress = ref(1);
let progressInterval: number;
let isPaused = false;

const bgColor = computed(() => {
  switch (props.toast.type) {
    case 'success': return 'bg-green-500';
    case 'error': return 'bg-red-500';
    case 'warning': return 'bg-yellow-500';
    case 'info': return 'bg-blue-500';
    default: return 'bg-blue-500';
  }
});

const progressColor = computed(() => {
  switch (props.toast.type) {
    case 'success': return 'bg-green-200';
    case 'error': return 'bg-red-200';
    case 'warning': return 'bg-yellow-200';
    case 'info': return 'bg-blue-200';
    default: return 'bg-blue-200';
  }
});

const dismiss = () => {
  removeToast(props.toast.id);
};

const onMouseEnter = () => {
  isPaused = true;
  pauseToast(props.toast.id);
};

const onMouseLeave = () => {
  isPaused = false;
  resumeToast(props.toast.id);
};

onMounted(() => {
  if (!props.toast.duration) return;

  const updateProgress = () => {
    if (isPaused) return;

    // We need to calculate progress based on remaining time in the toast object
    // But useToast updates remainingTime only on pause.
    // So we should rely on local calculation synced with start time.
    // Actually, useToast.ts sets startTime when added or resumed.

    const now = Date.now();
    const startTime = props.toast.startTime || now;
    const duration = props.toast.duration || 5000;

    // If we just resumed, the startTime was reset to now, and duration is effectively the remaining time?
    // Let's check useToast.ts:
    // resumeToast: toast.startTime = Date.now(); setTimeout(..., toast.remainingTime)
    // So if resumed, the "total duration" for this segment is toast.remainingTime.
    // This makes the math tricky if we want a smooth bar from 0 to 100% of the *original* duration.
    // However, the visual bar usually just represents "time left / total time".
    // If we pause and resume, the bar should continue from where it left off.

    // Let's try a simpler approach:
    // The bar scale is remainingTime / originalDuration.
    // But useToast doesn't track originalDuration separately if we overwrite duration?
    // useToast has `duration` (original) and `remainingTime`.
    // When resumed, `startTime` is reset, but `duration` (original) is kept?
    // No, `useToast` uses `toast.remainingTime` for the new timeout.
    // It doesn't seem to change `toast.duration`.

    // So:
    // Current Remaining = toast.remainingTime - (Date.now() - toast.startTime)
    // Progress = Current Remaining / toast.duration

    // Wait, `useToast` only updates `remainingTime` when PAUSED.
    // While running, `remainingTime` is static (the value at start of this segment).
    // So:
    // Active Remaining = props.toast.remainingTime - (now - props.toast.startTime!)

    const elapsedSinceStart = now - startTime;
    const currentRemaining = (props.toast.remainingTime || duration) - elapsedSinceStart;

    progress.value = Math.max(0, currentRemaining / (props.toast.duration || 5000));
  };

  progressInterval = window.setInterval(updateProgress, 50);
});

onBeforeUnmount(() => {
  if (progressInterval) clearInterval(progressInterval);
});
</script>
