import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import vitePluginRequire from 'vite-plugin-require';
import { fileURLToPath, URL } from 'url';

export default defineConfig({
	plugins: [vue(), vitePluginRequire.default()],
	resolve: {
		alias: [
			{
				find: '@',
				replacement: fileURLToPath(new URL('./src', import.meta.url)),
			},
		],
	},
});
