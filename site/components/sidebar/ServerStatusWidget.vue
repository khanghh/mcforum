<template>
  <div class="gaming-card rounded-xl p-4">
    <div class="text-lg font-bold text-green-400 mb-4 flex items-center gaming-title">
      <Icon name="TablerServer" class="mr-2" /> {{ $t('widgets.server_status') }}
    </div>
    <div class="space-y-3">
      <div class="flex justify-between items-center p-2 rounded bg-yellow-500/10">
        <span class="text-sm text-gray-400">Server IP</span>
        <span class="font-bold text-amber-400 gaming-title">play.mineviet.com</span>
      </div>
      <div class="flex justify-between items-center p-2 rounded bg-green-500/10">
        <span class="text-sm text-gray-400">Status</span>
        <span class="px-2 py-1 bg-green-500 text-white text-xs font-bold rounded gaming-title">
          {{ mcStatus.status.toUpperCase() }}
        </span>
      </div>
      <div class="flex justify-between items-center p-2 rounded bg-blue-500/10">
        <span class="text-sm text-gray-400">Online</span>
        <span class="font-bold text-green-400 gaming-title">{{ mcStatus.playersOnline }}</span>
      </div>
      <div class="flex justify-between items-center p-2 rounded bg-purple-500/10">
        <span class="text-sm text-gray-400">Version</span>
        <span class="font-bold text-purple-400 gaming-title">{{ mcStatus.version }}</span>
      </div>
    </div>
    <a href="https://play.mineviet.com"
      title="Xem trạng thái và tham gia máy chủ MineViet"
      class="w-full mt-4 px-4 py-3
             bg-gradient-to-r from-green-600 to-emerald-600
             text-white rounded-lg font-bold gaming-title
             hover:scale-105 transition-transform
             flex items-center justify-center">
      <Icon name="TablerPlayerPlayFilled" class="mr-2" /> {{ $t('widgets.server_status.join_server') }}
    </a>
  </div>
</template>

<script setup lang="ts">
const route = useRoute()
const api = useApi()


const mcVersions: Record<number, string> = {
  774: "1.21.11",
  773: "1.21.10",
  772: "1.21.8",
  771: "1.21.6",
  770: "1.21.5",
  769: "1.21.4",
  768: "1.21.3",
  767: "1.21.1",
  766: "1.20.6",
  765: "1.20.4",
  764: "1.20.2",
  763: "1.20.1",
  762: "1.19.4",
  761: "1.19.3",
  760: "1.19.2",
  759: "1.19",
  758: "1.18.2",
  757: "1.18.1",
  756: "1.17.1",
  755: "1.17",
  754: "1.16.5",
  753: "1.16.3",
  751: "1.16.2",
  736: "1.16.1",
  735: "1.16",
  578: "1.15.2",
  575: "1.15.1",
  573: "1.15",
  498: "1.14.4",
  477: "1.14",
  404: "1.13.2",
  393: "1.13",
  340: "1.12.2",
  316: "1.11.2",
  210: "1.10.2",
  110: "1.9.4",
  47: "1.8.9",
  5: "1.7.10",
  4: "1.7.2",
}

function getMCVersion(protocol: number): string {
  return mcVersions[protocol] || "Unknown"
}

const { data: mcStatus } = await useAsyncData(
  'mc-status',
  () => api.getMCStatus().then(data => {
    data.version = getMCVersion(data.protocol)
    return data
  }).catch(() => { }),
  {
    default: () => ({
      status: 'offline',
      version: 'Unknown',
      online: 0,
      playersOnline: 0,
      playersMax: 0,
    }),
    watch: [() => route.fullPath],
  }
)

</script>
