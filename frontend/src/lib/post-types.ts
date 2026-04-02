export interface PostTag {
  name: string
  slug: string
}

export interface PostSummary {
  id: string
  slug: string
  title: string
  excerpt: string
  hero_image_url: string
  category: string
  published_at?: string
  view_count: number
  like_count: number
  comment_count: number
  created_at: string
  tags?: PostTag[]
}

export interface PostDetail extends PostSummary {
  content_markdown: string
  status: string
  updated_at: string
  liked: boolean
}

export interface HeadingItem {
  depth: number
  slug: string
  text: string
}

export function toPostUrl(slug: string): string {
  return `/blog/${encodeURIComponent(slug)}`
}
