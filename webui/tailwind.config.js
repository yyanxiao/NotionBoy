const { fontFamily } = require("tailwindcss/defaultTheme");

/** @type {import('tailwindcss').Config} */
module.exports = {
	darkMode: ["class", '[data-theme="dark"]'],
	content: ["./src/**/*.{ts,tsx}"],
	theme: {
		// https://www.happyhues.co/palettes/3
		// colors: {
		// 	my: {
		// 		primary: "#fffffe",
		// 		secondary: "#90b4ce",
		// 		tertiary: "#ef4565",
		// 		highlight: "#3da9fc",
		// 		stroke: "#2a2a2a",
		// 	},
		// },
		extend: {
			fontFamily: {
				sans: ["var(--font-sans)", ...fontFamily.sans],
			},
			keyframes: {
				"accordion-down": {
					from: { height: 0 },
					to: { height: "var(--radix-accordion-content-height)" },
				},
				"accordion-up": {
					from: { height: "var(--radix-accordion-content-height)" },
					to: { height: 0 },
				},
			},
			animation: {
				"accordion-down": "accordion-down 0.2s ease-out",
				"accordion-up": "accordion-up 0.2s ease-out",
			},
		},
	},
	plugins: [
		require("@tailwindcss/typography"),
		require("tailwindcss-animate"),
		require("@tailwindcss/forms"),
	],
};
