// @ts-check

import mdx from '@astrojs/mdx';
import sitemap from '@astrojs/sitemap';
import { defineConfig } from 'astro/config';

import vue from '@astrojs/vue';
import { siteCopy } from './src/content/copy/site.ts';

// https://astro.build/config
export default defineConfig({
    site: siteCopy.seo.siteUrl,
    integrations: [mdx(), sitemap(), vue()],
    vite: {
        optimizeDeps: {
            exclude: ['marked'],
        },
        server: {
            proxy: {
                '/api': {
                    target: 'http://localhost:8080',
                    changeOrigin: true,
                },
            },
        },
    },
});
