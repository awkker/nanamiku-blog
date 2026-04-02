import rss from '@astrojs/rss';
import { SITE_DESCRIPTION, SITE_TITLE } from '../consts';
import { getPublishedPostSummaries } from '../lib/post-source';
import { toPostUrl } from '../lib/post-types';

export async function GET(context) {
	const { data: posts } = await getPublishedPostSummaries();
	return rss({
		title: SITE_TITLE,
		description: SITE_DESCRIPTION,
		site: context.site,
		items: posts.map((post) => ({
			title: post.title,
			description: post.excerpt,
			pubDate: new Date(post.published_at || post.created_at),
			link: toPostUrl(post.slug),
		})),
	});
}
