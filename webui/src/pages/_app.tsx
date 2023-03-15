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
			</Head>

			<Layout key={router.asPath}>
				<Component {...pageProps} />
			</Layout>
		</>
	);
}
