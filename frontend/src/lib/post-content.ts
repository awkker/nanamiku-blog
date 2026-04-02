import type { HeadingItem } from './post-types'

let markedModule: Awaited<ReturnType<typeof import('marked')>> | null = null

export function escapeHtml(input: string): string {
  return input
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;')
    .replaceAll("'", '&#39;')
}

export function normalizeMarkdown(source: string): string {
  return source
    .replace(/```[\s\S]*?```/g, ' ')
    .replace(/`[^`]*`/g, ' ')
    .replace(/!\[[^\]]*]\([^)]*\)/g, ' ')
    .replace(/\[([^\]]+)]\([^)]*\)/g, '$1')
    .replace(/<\/?[^>]+(>|$)/g, ' ')
    .replace(/[#>*_~\-]/g, ' ')
}

export function getReadingStats(source: string) {
  const normalized = normalizeMarkdown(source)
  const latinWords = normalized.match(/[A-Za-z0-9]+/g) ?? []
  const cjkChars = normalized.match(/[\u3400-\u9fff]/g) ?? []
  const wordCount = latinWords.length + cjkChars.length

  return {
    wordCount,
    readingMinutes: Math.max(1, Math.ceil(wordCount / 260)),
  }
}

export function slugifyHeading(text: string): string {
  const normalized = text
    .trim()
    .toLowerCase()
    .normalize('NFKD')
    .replace(/[\u0300-\u036f]/g, '')
    .replace(/[^\p{Letter}\p{Number}\u4e00-\u9fff\s-]/gu, '')
    .replace(/\s+/g, '-')
    .replace(/-+/g, '-')
    .replace(/^-|-$/g, '')

  return normalized || 'section'
}

export function collectHeadings(
  markdown: string,
  markedLib: Awaited<ReturnType<typeof import('marked')>>,
): HeadingItem[] {
  const counters = new Map<string, number>()
  const collected: HeadingItem[] = []
  const tokens = markedLib.marked.lexer(markdown)

  const visit = (list: unknown[]) => {
    list.forEach((token) => {
      if (!token || typeof token !== 'object') {
        return
      }

      const item = token as {
        type?: string
        depth?: number
        text?: string
        tokens?: unknown[]
        items?: Array<{ tokens?: unknown[] }>
      }

      if (item.type === 'heading') {
        const rawText = item.text?.trim() || `章节 ${collected.length + 1}`
        const base = slugifyHeading(rawText)
        const nextIndex = (counters.get(base) ?? 0) + 1
        counters.set(base, nextIndex)
        const slug = nextIndex === 1 ? base : `${base}-${nextIndex}`
        collected.push({
          depth: Number(item.depth) || 2,
          slug,
          text: rawText,
        })
        return
      }

      if (item.type === 'blockquote' && Array.isArray(item.tokens)) {
        visit(item.tokens)
      }

      if (item.type === 'list' && Array.isArray(item.items)) {
        item.items.forEach((listItem) => {
          if (Array.isArray(listItem.tokens)) {
            visit(listItem.tokens)
          }
        })
      }
    })
  }

  visit(tokens)
  return collected
}

export async function renderPostMarkdown(markdown: string): Promise<{
  html: string
  headings: HeadingItem[]
}> {
  if (!markdown) {
    return {
      html: '',
      headings: [],
    }
  }

  try {
    if (!markedModule) {
      markedModule = await import('marked')
    }

    const headings = collectHeadings(markdown, markedModule)
    let headingCursor = 0

    const renderer = new markedModule.marked.Renderer()
    renderer.heading = function ({ tokens, depth }: { tokens: unknown[]; depth: number }) {
      const content = this.parser.parseInline(tokens)
      const fallbackSlug = `section-${headingCursor + 1}`
      const slug = headings[headingCursor]?.slug || fallbackSlug
      headingCursor += 1
      return `<h${depth} id="${slug}" data-heading-id="${slug}">${content}</h${depth}>`
    }

    const html = await markedModule.marked.parse(markdown, {
      gfm: true,
      breaks: false,
      renderer,
    })

    return {
      html: String(html),
      headings: headings.filter((item) => item.depth <= 3),
    }
  } catch {
    const readable = escapeHtml(markdown).replace(/\n/g, '<br>')
    return {
      html: `<p>${readable}</p>`,
      headings: [],
    }
  }
}
