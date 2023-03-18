import { siteConfig } from "@/config/site";
import { useToast } from "@/hooks/use-toast";
import { Service } from "@/lib/pb/server.pb";

import Cookies from "js-cookie";
import { useRouter } from "next/router";
import { useEffect, useState } from "react";
import { Button } from "./ui/button";

export const TOKEN = "token";
export const TOKEN_EXPIRE = "tokenExpire";

const PROTECTED_PATHS = siteConfig.authPages;

export function AuthLoginButton() {
	const [token, setToken] = useState<string | null>(null);
	// add a const to show if the user is logged in or not

	const [errorMessage, setErrorMessage] = useState("");

	const [isShowLoginButton, setIsShowLoginButton] = useState(true);
	const router = useRouter();
	const { toast } = useToast();

	useEffect(() => {
		if (router.isReady) {
			if (router.pathname == siteConfig.links.login) {
				setIsShowLoginButton(false);
			}
			const token = Cookies.get(TOKEN);
			if (token) {
				setToken(token);
			}
			handleProtectedPage();
		}
	}, [router.isReady]);

	useEffect(() => {
		if (errorMessage) {
			toast({
				variant: "destructive",
				title: "Auth error",
				description: errorMessage,
			});
		}
	}, [errorMessage]);

	const handleProtectedPage = () => {
		for (const path of PROTECTED_PATHS) {
			if (router.pathname.startsWith(path)) {
				const tokenExpire = Cookies.get(TOKEN_EXPIRE);
				const diff = Date.parse(tokenExpire as string) - Date.now();

				if (diff < 0) {
					Cookies.remove(TOKEN);
					Cookies.remove(TOKEN_EXPIRE);
				} else if (diff < 3600) {
					// refresh token
					Service.GenrateToken({})
						.then((res) => {
							if (res.token) {
								Cookies.set(TOKEN, res.token, {
									expires: new Date(res.expiry as string),
								});
								Cookies.set(
									TOKEN_EXPIRE,
									res.expiry as string,
									{
										expires: new Date(res.expiry as string),
									}
								);
								setToken(res.token);
							} else {
								setErrorMessage("Get token failed");
								redirectToLogin();
							}
						})
						.catch((err) => {
							setErrorMessage(`Get token failed ${err}`);
							redirectToLogin();
						});
				} else {
					setToken(Cookies.get(TOKEN) as string);
				}
			}
		}
	};

	const handleSignOut = () => {
		Cookies.remove(TOKEN);
		Cookies.remove(TOKEN_EXPIRE);
		router.push(siteConfig.links.home);
		router.reload();
	};

	const redirectToLogin = () => {
		router.push(siteConfig.links.login);
	};

	const loginButton = () => {
		if (isShowLoginButton) {
			return token ? (
				<Button variant="ghost" size="sm" onClick={handleSignOut}>
					Logout
				</Button>
			) : (
				<Button variant="ghost" size="sm" onClick={redirectToLogin}>
					Login
				</Button>
			);
		}
	};
	return <>{loginButton()}</>;
}
