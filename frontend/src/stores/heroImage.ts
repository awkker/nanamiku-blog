import { atom, computed } from 'nanostores'

import { homeAssetsSettings } from './homeAssets'

// `heroImages` 是一个派生 store：
// 它不单独存数据，而是直接从首页素材设置里取出当前可用的背景图列表。
export const heroImages = computed(homeAssetsSettings, (settings) => settings.heroImages)

// `heroIndex` 则是“当前显示哪一张”的纯前端状态。
export const heroIndex = atom(0)

homeAssetsSettings.listen((settings) => {
  // 当后台把图片列表改短时，当前索引可能越界。
  // 这里监听原始列表变化，保证索引始终落在合法范围内。
  if (settings.heroImages.length === 0) {
    heroIndex.set(0)
    return
  }

  if (heroIndex.get() >= settings.heroImages.length) {
    heroIndex.set(0)
  }
})

export function shuffleHeroImage() {
  const currentHeroImages = heroImages.get()
  if (currentHeroImages.length <= 1) {
    // 没有图或只有一张图时，不需要随机，固定停在第 0 张即可。
    heroIndex.set(0)
    return
  }

  const current = heroIndex.get()
  let next: number
  do {
    // 用 do...while 的原因是：我们希望至少抽一次，
    // 并且尽量避免连续两次还是同一张图。
    next = Math.floor(Math.random() * currentHeroImages.length)
  } while (next === current && currentHeroImages.length > 1)
  heroIndex.set(next)
}
