<template>
  <NuxtLoadingIndicator color="#27ae60" :height="5" />
  <NuxtLayout>
    <NuxtPage />
  </NuxtLayout>
</template>


<script setup lang="ts">

useSeoMeta({
  title: 'MineViet Network - Server Minecraft Việt Nam Hàng Đầu',
  ogType: 'website',
  ogUrl: 'https://mineviet.com',
  ogTitle: 'MineViet Minecraft Server - Chơi SkyBlock, Factions, Survival và Mini Game',
  ogSiteName: 'MineViet Network',
  ogImage: "https://mineviet.com/images/banner_1200x630.png",
  description: 'MineViet Network - Máy chủ Minecraft Việt Nam với nhiều chế độ chơi: SkyBlock, Factions, Prison, BedWars, SkyWars và nhiều mini game. Tham gia ngay hôm nay!',
  ogDescription: 'Tham gia MineViet Network — máy chủ Minecraft Việt Nam với SkyBlock, Factions, Survival và Mini game, cộng đồng năng động và sự kiện thường xuyên.',
  twitterCard: 'summary_large_image',
  twitterSite: '@minevietnetwork',
  twitterTitle: 'MineViet Minecraft Server - Chơi SkyBlock, Factions, Survival và Mini Game',
  twitterDescription: 'Tham gia MineViet Network — máy chủ Minecraft Việt Nam với SkyBlock, Factions, Survival và Mini game, cộng đồng năng động và sự kiện thường xuyên.',
  twitterImage: "https://mineviet.com/images/banner_1200x630.png",
})

</script>
