<template>
  <component :is="iconMap[name]" weight="fill" size="1.2em" v-bind="$attrs" class="svg-inline"
    :class="{ 'animate-spin': spin }" />
</template>

<script setup lang="ts">
import {
  PhChatCircle,
  PhArrowUp,
  PhBookmark,
  PhSignOut,
  PhCalendarPlus,
  PhChartLine,
  PhCheckCircle,
  PhClock,
  PhCloudArrowUp,
  PhGear,
  PhChatText,
  PhCrosshair,
  PhCrown,
  PhCube,
  PhEnvelope,
  PhWarning,
  PhEye,
  PhFunnel,
  PhFire,
  PhFlag,
  PhGameController,
  PhHeart,
  PhHouse,
  PhImage,
  PhInfo,
  PhMapPin,
  PhMedal,
  PhPalette,
  PhPlay,
  PhPushPin,
  PhPlus,
  PhQuestion,
  PhHardDrive,
  PhShare,
  PhShield,
  PhShoppingCart,
  PhArrowRight,
  PhSpinner,
  PhStar,
  PhList,
  PhThumbsUp,
  PhCaretDown,
  PhX,
  PhTrophy,
  PhLightning,
  PhTarget,
  PhUser,
  PhUserCircle,
  PhUserPlus,
  PhUserMinus,
  PhUsers,
  PhBell,
  PhDiscordLogo,
  PhTwitterLogo,
  PhYoutubeLogo,
  PhTwitchLogo,
  PhHardDrives,
  PhAcorn,
  PhActivity,
  PhPuzzlePiece,
  PhBomb,
  PhTrash,
} from "@phosphor-icons/vue"

const iconMap = {
  'arrow-up': PhArrowUp,
  'bookmark': PhBookmark,
  'calendar-plus': PhCalendarPlus,
  'chart-line': PhChartLine,
  'check-circle': PhCheckCircle,
  'clock': PhClock,
  'cloud-upload-alt': PhCloudArrowUp,
  'cog': PhGear,
  'comment-alt': PhChatText,
  'comments': PhChatText,
  'crosshairs': PhCrosshair,
  'crown': PhCrown,
  'cube': PhCube,
  'envelope': PhEnvelope,
  'exclamation-triangle': PhWarning,
  'eye': PhEye,
  'filter': PhFunnel,
  'fire': PhFire,
  'flag': PhFlag,
  'gamepad': PhGameController,
  'heart': PhHeart,
  'home': PhHouse,
  'image': PhImage,
  'info-circle': PhInfo,
  'map-marker-alt': PhMapPin,
  'medal': PhMedal,
  'palette': PhPalette,
  'play': PhPlay,
  'thumbtack': PhPushPin,
  'plus': PhPlus,
  'question-circle': PhQuestion,
  'hdd': PhHardDrive,
  'hard-drive': PhHardDrive,
  'share-alt': PhShare,
  'share': PhShare,
  'shield-alt': PhShield,
  'shopping-cart': PhShoppingCart,
  'arrow-right': PhArrowRight,
  'spinner': PhSpinner,
  'star': PhStar,
  'list': PhList,
  'th-list': PhList,
  'thumbs-up': PhThumbsUp,
  'times': PhX,
  'trophy': PhTrophy,
  'bolt': PhLightning,
  'bullseye': PhTarget,
  'user': PhUser,
  'user-circle': PhUserCircle,
  'user-plus': PhUserPlus,
  'user-minus': PhUserMinus,
  'users': PhUsers,
  'bell': PhBell,
  'caret-down': PhCaretDown,
  'chevron-down': PhCaretDown,
  'comment': PhChatCircle,
  'comment-dots': PhChatCircle,
  'sign-out-alt': PhSignOut,
  'discord': PhDiscordLogo,
  'twitter': PhTwitterLogo,
  'youtube': PhYoutubeLogo,
  'twitch': PhTwitchLogo,
  'server': PhHardDrives,
  'ban': PhActivity,
  'puzzle-piece': PhPuzzlePiece,
  'bomb': PhBomb,
  'trash': PhTrash,
  'acorn': PhAcorn,
}

type IconName = keyof typeof iconMap

const props = defineProps<{
  name: IconName
  spin?: boolean
}>()
</script>

<style scoped>
.svg-inline {
  box-sizing: content-box;
  display: var(--fa-display, inline-block);
  overflow: visible;
}

.animate-spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }

  to {
    transform: rotate(360deg);
  }
}
</style>