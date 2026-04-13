<template>
  <span class="tabular-nums text-[11px] tracking-wide text-[#39c5bb]/90">
    {{ displayText }}
  </span>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { api } from '../../lib/api'

interface WeatherData {
  temp: string
  feels_like: string
  humidity: string
  desc: string
  icon: string
  wind_speed: string
  location: string
}

interface Props {
  location?: string
}

const props = defineProps<Props>()

// 后端返回的是语义化天气图标 key，这里再映射成本地要展示的 Unicode 符号。
const iconMap: Record<string, string> = {
  sunny: '\u2600',
  partly_cloudy: '\u26C5',
  cloudy: '\u2601',
  light_rain: '\u2602',
  rain: '\u2602',
  snow: '\u2744',
  thunderstorm: '\u26A1',
}

const displayText = ref('--\u00B0C')

async function loadWeather() {
  try {
    // 如果外部传了 location，就拼到查询参数里；
    // 否则交给后端自行决定默认地点。
    const query = props.location ? `?location=${encodeURIComponent(props.location)}` : ''
    const data = await api.get<WeatherData>(`/weather${query}`)
    const icon = iconMap[data.icon] || '\u2601'
    displayText.value = `${icon} ${data.temp}\u00B0C`
  } catch {
    // 天气信息不是首页核心流程，失败时直接静默降级成占位值即可。
    displayText.value = '--\u00B0C'
  }
}

watch(() => props.location, () => {
  // 后台修改地点配置后，组件会自动重新拉天气。
  void loadWeather()
})

onMounted(async () => {
  await loadWeather()
})
</script>
