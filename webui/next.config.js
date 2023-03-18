/** @type {import('next').NextConfig} */
const nextConfig = {
	reactStrictMode: true,
	images: {
		unoptimized: true,
		domains: ["images.unsplash.com", "image.lexica.art"],
	},
	basePath: "/web",
};

module.exports = nextConfig;
