import { Layout } from "@/components/layout";
import "@/styles/globals.css";
import type { AppProps } from "next/app";
import Head from "next/head";
import { useRouter } from "next/router";

export default function App({
	Component,
	pageProps: { session, ...pageProps },
}: AppProps) {
	const router = useRouter();
	return (
		<>
			<Head>
				<title>NotionBoy</title>
				<meta
					name="viewport"
					// This code sets the initial scale to 1 and disables user scaling, which should prevent the auto-scaling issue on Safari for iPhone.
					content="width=device-width, initial-scale=1, maximum-scale=1, user-scalable=0"
				/>
			</Head>

			<Layout key={router.asPath}>
				<Component {...pageProps} />
			</Layout>
		</>
	);
}
