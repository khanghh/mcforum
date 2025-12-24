<template>
  <component :is="iconMap[name]" v-bind="$attrs" class="svg-inline" />
</template>

<script setup lang="ts">
import {
  PhArrowUp,
  PhBookmark,
  PhCalendarPlus,
  PhChartLine,
  PhCheckCircle,
  PhClock,
  PhCloudArrowUp,
  PhGear,
  PhChatText,
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
  PhX,
  PhTrophy,
  PhUser,
  PhUserCircle,
  PhUserPlus,
  PhUsers,
  PhBell,
} from "@phosphor-icons/vue"

const iconMap = {
  'PhArrowUp': PhArrowUp,
  'PhBookmark': PhBookmark,
  'PhCalendarPlus': PhCalendarPlus,
  'PhChartLine': PhChartLine,
  'PhCheckCircle': PhCheckCircle,
  'PhClock': PhClock,
  'PhCloudArrowUp': PhCloudArrowUp,
  'PhGear': PhGear,
  'PhChatText': PhChatText,
  'PhCrown': PhCrown,
  'PhCube': PhCube,
  'PhEnvelope': PhEnvelope,
  'PhWarning': PhWarning,
  'PhEye': PhEye,
  'PhFunnel': PhFunnel,
  'PhFire': PhFire,
  'PhFlag': PhFlag,
  'PhGameController': PhGameController,
  'PhHeart': PhHeart,
  'PhHouse': PhHouse,
  'PhImage': PhImage,
  'PhInfo': PhInfo,
  'PhMapPin': PhMapPin,
  'PhMedal': PhMedal,
  'PhPalette': PhPalette,
  'PhPlay': PhPlay,
  'PhPlus': PhPlus,
  'PhQuestion': PhQuestion,
  'PhHardDrive': PhHardDrive,
  'PhShare': PhShare,
  'PhShield': PhShield,
  'PhShoppingCart': PhShoppingCart,
  'PhArrowRight': PhArrowRight,
  'PhSpinner': PhSpinner,
  'PhStar': PhStar,
  'PhList': PhList,
  'PhThumbsUp': PhThumbsUp,
  'PhX': PhX,
  'PhTrophy': PhTrophy,
  'PhUser': PhUser,
  'PhUserCircle': PhUserCircle,
  'PhUserPlus': PhUserPlus,
  'PhUsers': PhUsers,
  'PhBell': PhBell,
} as const

type IconName = keyof typeof iconMap

const props = defineProps<{ name: IconName }>()

// All other attributes and event listeners are forwarded via $attrs
</script>

<style scoped>
.svg-inline {
  box-sizing: content-box;
  display: var(--fa-display, inline-block);
  height: 1em;
  overflow: visible;
  vertical-align: -0.125em;
  width: var(--fa-width, 1.25em);
}
</style>