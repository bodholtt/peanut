const colors = require('tailwindcss/colors')
/** @type {import('tailwindcss').Config} */
export default {
	content: ['./src/**/*.{astro,html,js,jsx,md,mdx,svelte,ts,tsx,vue}'],
	theme: {
		extend: {},
		colors: {
			twcolor: colors.amber,
			white: colors.white,
			black: colors.black,
			red: colors.red,
			transparent: 'transparent',
			gray: colors.gray,
		},
	},
	plugins: [],
}
