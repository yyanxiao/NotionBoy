import { siteConfig } from "@/config/site";
import { OAuthCallbackRequest } from "@/lib/pb/model/common.pb";
import { Service } from "@/lib/pb/server.pb";
import Cookies from "js-cookie";
import { useRouter } from "next/router";
import { useEffect, useState } from "react";

export default function Callback() {
	const [isValid, setIsValid] = useState(false);
	const [isLoading, setIsLoading] = useState(true);
	const [errorMessage, setErrorMessage] = useState("");

	const router = useRouter();

	useEffect(() => {
		const req = {
			code: router.query.code as string,
			state: router.query.state as string,
		} as OAuthCallbackRequest;
		if (!router.isReady) {
			setIsLoading(true);
			return;
		}

		console.log(`req is in useEffect ${JSON.stringify(req)}`);
		Service.OAuthCallback(req)
			.then((res) => {
				console.log(`get token: ${res}`);
				if (res.token) {
					Cookies.set("token", res.token, {
						expires: new Date(res.expiry as string),
					});
					Cookies.set("tokenExpire", res.expiry as string, {
						expires: new Date(res.expiry as string),
					});
				}
				setIsValid(true);
				router.push(siteConfig.links.home);
				setIsLoading(false);
			})
			.catch((err) => {
				console.log("Get token failed", err);
				setErrorMessage(err.message);
				setIsLoading(false);
			});
	}, [router.isReady]);

	return (
		<div className="container mx-auto p-8 ">
			<article className="prose prose-sm sm:prose lg:prose-lg xl:prose-2xl">
				{!isLoading && isValid && <h1>Redirecting...</h1>}
				{!isLoading && !isValid && (
					<div>
						<h1>Login Failed</h1>
						<div>
							<p className="text-red-500">
								<strong>Error Message:</strong> {errorMessage}
							</p>
						</div>
					</div>
				)}
				{isLoading && <h1>Validating...</h1>}
			</article>
		</div>
	);
}
