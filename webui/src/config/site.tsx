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
	};
	authPages: string[];
}

export const siteConfig: SiteConfig = {
	name: "NotionBoy",
	description:
		"NotionBoy is a note app base on Notion. It's a web app, you can use it in your browser.",
	links: {
		twitter: "https://twitter.com/LiuVaayne",
		github: "https://github.com/vaayne/NotionBoy",
		chatgpt: "/web/chat.html",
		login: "/web/login.html",
		home: "/web/",
		authCallback: "/web/auth-callback.html",
	},
	authPages: ["/web/chat.html"],
};
