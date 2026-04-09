import { atom, computed } from 'nanostores'

import { homeAssetsSettings } from './homeAssets'

export const heroImages = computed(homeAssetsSettings, (settings) => settings.heroImages)

export const heroIndex = atom(0)

homeAssetsSettings.listen((settings) => {
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
    heroIndex.set(0)
    return
  }

  const current = heroIndex.get()
  let next: number
  do {
    next = Math.floor(Math.random() * currentHeroImages.length)
  } while (next === current && currentHeroImages.length > 1)
  heroIndex.set(next)
}
