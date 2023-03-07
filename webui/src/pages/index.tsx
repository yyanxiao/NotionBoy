import Cookies from "js-cookie";

import { useRouter } from "next/router";

export default function Home() {
	const router = useRouter();
	const token = router.query.token;
	if (token) {
		localStorage.setItem("token", token as string);
		Cookies.set("token", token as string);
		router.push("/web");
	}

	return (
		<div className="container mx-auto">
			<h1 className="text-3xl font-bold underline">Hello world!</h1>
			<a className="link link-primary" href="/web/chat.html">
				NotionBoy Chat
			</a>
		</div>
	);
}
