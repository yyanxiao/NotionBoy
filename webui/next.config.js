/** @type {import('next').NextConfig} */
const nextConfig = {
	reactStrictMode: true,
	images: {
		unoptimized: true,
		domains: ["images.unsplash.com", "image.lexica.art"],
	},
	// basePath: "/web",
	// env: {},
	// async rewrites() {
	// 	return [
	// 		{
	// 			source: "/web/v1/:path*",
	// 			destination: "http://localhost:8001/v1/:path*",
	// 		},
	// 	];
	// },
};

module.exports = nextConfig;
