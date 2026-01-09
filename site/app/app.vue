<template>
  <NuxtLoadingIndicator color="#27ae60" :height="5" />
  <NuxtLayout>
    <NuxtPage />
  </NuxtLayout>
</template>


<script setup lang="ts">

useSeoMeta({
  title: 'Minecraft Việt Nam – MineViet Network',
  ogType: 'website',
  ogUrl: 'https://mineviet.com',
  ogTitle: 'MineViet Minecraft Server – Chơi SkyBlock, Factions, Survival và Mini Game',
  ogSiteName: 'MineViet Network',
  ogImage: "https://mineviet.com/images/mineviet.png",
  description: 'MineViet Network - Máy chủ Minecraft Việt Nam với nhiều chế độ chơi: SkyBlock, Survival, Factions, BedWars và nhiều mini game. Tham gia ngay hôm nay!',
  ogDescription: 'MineViet Network - Máy chủ Minecraft Việt Nam với nhiều chế độ chơi: SkyBlock, Survival, Factions, BedWars và nhiều mini game. Tham gia ngay hôm nay!',
  twitterCard: 'summary_large_image',
  twitterSite: '@minevietnetwork',
  twitterTitle: 'MineViet Minecraft Server – Chơi SkyBlock, Factions, Survival và Mini Game',
  twitterDescription: 'MineViet Network - Máy chủ Minecraft Việt Nam với nhiều chế độ chơi: SkyBlock, Survival, Factions, BedWars và nhiều mini game. Tham gia ngay hôm nay!',
  twitterImage: "https://mineviet.com/images/mineviet.png",
})

</script>
