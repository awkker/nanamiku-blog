export type GifPackKey = 'miku' | 'nanamichiaki'

export interface GifEmoteItem {
  id: string
  pack: GifPackKey
  path: string
  src: string
  token: string
  label: string
}

export type RichContentSegment =
  | { type: 'text'; key: string; value: string }
  | { type: 'gif'; key: string; path: string; src: string; label: string }

const mikuFiles = Array.from({ length: 20 }, (_, index) => `${String(index + 1).padStart(2, '0')}.gif`)
const nanamichiakiFiles = [
  '上吊.gif',
  '上工.gif',
  '举牌.gif',
  '口水.gif',
  '启动.gif',
  '哼哼.gif',
  '喜欢.gif',
  '喝茶.gif',
  '对手指.gif',
  '想要.gif',
  '慌.gif',
  '戳戳.gif',
  '拉泪.gif',
  '灵光.gif',
  '点赞.gif',
  '生气.gif',
  '疑惑.gif',
  '端详.gif',
  '震惊.gif',
  '鼓掌.gif',
]

const GIF_TOKEN_PATTERN = /\[gif:([^\]\s]+)\]/g

function normalizeGifPath(path: string): string {
  try {
    return decodeURIComponent(path)
  } catch {
    return path
  }
}

function toGifSrc(path: string): string {
  return `/gif/${path.split('/').map((segment) => encodeURIComponent(segment)).join('/')}`
}

function toGifToken(path: string): string {
  return `[gif:${path}]`
}

function buildPack(pack: GifPackKey, fileNames: string[]): GifEmoteItem[] {
  return fileNames.map((fileName) => {
    const path = `${pack}/${fileName}`
    return {
      id: path,
      pack,
      path,
      src: toGifSrc(path),
      token: toGifToken(path),
      label: fileName.replace(/\.gif$/i, ''),
    }
  })
}

export const GIF_EMOTE_PACKS: Record<GifPackKey, GifEmoteItem[]> = {
  miku: buildPack('miku', mikuFiles),
  nanamichiaki: buildPack('nanamichiaki', nanamichiakiFiles),
}

export const GIF_PACK_ORDER: GifPackKey[] = ['miku', 'nanamichiaki']

const allowedGifPaths = new Set(
  Object.values(GIF_EMOTE_PACKS)
    .flat()
    .map((item) => item.path),
)

export function appendGifToken(content: string, path: string): string {
  const normalizedPath = normalizeGifPath(path)
  if (!allowedGifPaths.has(normalizedPath)) return content

  const token = toGifToken(normalizedPath)
  if (!content.trim()) return token
  return /\s$/.test(content) ? `${content}${token}` : `${content}\n${token}`
}

export function parseRichContentWithGif(content: string): RichContentSegment[] {
  const text = content || ''
  const segments: RichContentSegment[] = []
  let cursor = 0
  let segmentId = 0

  GIF_TOKEN_PATTERN.lastIndex = 0
  let match = GIF_TOKEN_PATTERN.exec(text)

  while (match) {
    const raw = match[0]
    const tokenPath = normalizeGifPath(match[1] || '')

    if (match.index > cursor) {
      segments.push({
        type: 'text',
        key: `text-${segmentId++}`,
        value: text.slice(cursor, match.index),
      })
    }

    if (allowedGifPaths.has(tokenPath)) {
      segments.push({
        type: 'gif',
        key: `gif-${segmentId++}`,
        path: tokenPath,
        src: toGifSrc(tokenPath),
        label: tokenPath.split('/').pop()?.replace(/\.gif$/i, '') || tokenPath,
      })
    } else {
      segments.push({
        type: 'text',
        key: `text-${segmentId++}`,
        value: raw,
      })
    }

    cursor = match.index + raw.length
    match = GIF_TOKEN_PATTERN.exec(text)
  }

  if (cursor < text.length) {
    segments.push({
      type: 'text',
      key: `text-${segmentId++}`,
      value: text.slice(cursor),
    })
  }

  return segments
}
