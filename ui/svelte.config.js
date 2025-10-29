import adapter from '@sveltejs/adapter-auto';
//import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	preprocess: vitePreprocess(),

	compilerOptions: {
		experimental: {
			async: true
		}
	},

	kit: {
		adapter: adapter()
	}
	/*
  	kit: {
      adapter: adapter({
        pages: 'build',
        assets: 'build',
        fallback: 'index.html',
      }),
      paths: {
        base: '', // leave default
      },
      prerender: {
        entries: ['*'],
      },
    },
	*/
};

export default config;
