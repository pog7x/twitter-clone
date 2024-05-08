import { defineConfig, loadEnv } from 'vite';
import vue from '@vitejs/plugin-vue';
import vitePluginRequire from 'vite-plugin-require';
import { fileURLToPath, URL } from 'url';

export default ({ mode }) => {
	const env = loadEnv(mode, process.cwd(), '');

	return defineConfig({
		plugins: [vue(), vitePluginRequire.default()],
		resolve: {
			alias: {
				'@': fileURLToPath(new URL('./src', import.meta.url)),
			},
		},
		server: {
			port: parseInt(env.VITE_PORT),
			host: env.VITE_HOST,
		},
	});
};
