interface SiteConfig {
	name: string;
	description: string;
	links: {
		twitter: string;
		github: string;
		chatgpt: string;
		login: string;
		home: string;
		authCallback: string;
		price: string;
		order: string;
		profile: string;
	};
	authPages: string[];
}

function buildPath(path: string): string {
	if (process.env.NODE_ENV === "development") {
		return path;
	}
	return path + ".html";
}

export const siteConfig: SiteConfig = {
	name: "NotionBoy",
	description:
		"NotionBoy is a note app base on Notion. It's a web app, you can use it in your browser.",
	links: {
		twitter: "https://twitter.com/LiuVaayne",
		github: "https://github.com/vaayne/NotionBoy",
		chatgpt: buildPath("/chat"),
		login: buildPath("/login"),
		home: "/",
		authCallback: buildPath("/authcallback"),
		price: buildPath("/price"),
		order: buildPath("/order"),
		profile: buildPath("/profile"),
	},
	authPages: [buildPath("/chat"), buildPath("/order"), buildPath("/user")],
};
