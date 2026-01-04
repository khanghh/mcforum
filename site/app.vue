<template>
  <NuxtLoadingIndicator color="#27ae60" :height="5" />
  <NuxtLayout>
    <NuxtPage />
  </NuxtLayout>
</template>


<script setup lang="ts">

useSeoMeta({
  title: 'MineViet Network - Server Minecraft Việt Nam Hàng Đầu',
  ogTitle: 'MineViet Minecraft Server - Chơi SkyBlock, Factions, Survival và Mini Game',
  ogSiteName: 'MineViet Network',
  ogImage: "https://mineviet.com/images/mineviet.png",
  description: 'MineViet Network - Máy chủ Minecraft Việt Nam với nhiều chế độ chơi: SkyBlock, Factions, Prison, BedWars, SkyWars và nhiều mini game. Tham gia ngay hôm nay!',
  ogDescription: 'Tham gia MineViet Network — máy chủ Minecraft Việt Nam với SkyBlock, Factions, Survival và Mini game, cộng đồng năng động và sự kiện thường xuyên.',
  twitterCard: 'summary_large_image',
})

</script>
