<script setup>
import { computed } from 'vue'
import {
  CloudUpload,
  Trophy,
  X,
  ThumbsUp,
  MessageSquareMore,
  Plus,
  Trash2,
  LoaderCircle,
  Mail,
  UserRound,
  Star,
  Pencil,
  LogOut,
  Search,
  ImagePlus,
  BookOpenText,
  Ban,
  Lock,
  ShieldCheck,
  MessageSquareText,
  FileText,
  Undo2,
  Settings,
  LockOpen,
  SquarePen,
} from 'lucide-vue-next'

const icons = {
  CloudUpload,
  Trophy,
  X,
  ThumbsUp,
  MessageSquareMore,
  Plus,
  Trash2,
  LoaderCircle,
  Mail,
  UserRound,
  Star,
  Pencil,
  LogOut,
  Search,
  ImagePlus,
  BookOpenText,
  Ban,
  Lock,
  ShieldCheck,
  MessageSquareText,
  FileText,
  Undo2,
  Settings,
  LockOpen,
  SquarePen,
}

const props = defineProps({
  name: {
    type: String,
    required: true,
  },
  size: {
    type: String,
    default: '1em',
  },
  strokeWidth: {
    type: Number,
    default: 2,
  },
  color: String,
  filled: Boolean,
  defaultClass: String,
})

const icon = computed(() => icons[props.name])
const fill = computed(() => props.filled ? props.color || 'currentColor' : 'none')
</script>

<template>
  <component :is="icon" :size="size" :fill="fill" :color="color" :stroke-width="strokeWidth"
    :default-class="defaultClass" />
</template>
