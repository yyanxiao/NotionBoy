import { TOKEN, TOKEN_EXPIRE } from "@/components/auth";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { siteConfig } from "@/config/site";
import { useToast } from "@/hooks/use-toast";
import { GenrateTokenRequest, OAuthURLRequest } from "@/lib/pb/model/common.pb";
import { Service } from "@/lib/pb/server.pb";
import Cookies from "js-cookie";
import Link from "next/link";
import { useRouter } from "next/router";
import { useEffect, useState } from "react";

interface Provider {
	id: string;
	name: string;
	url: string | undefined;
}

export default function SignIn() {
	const router = useRouter();
	const { toast } = useToast();
	const [validProviders, setValidProviders] = useState<Provider[]>([]);
	const [isLoggedIn, setIsLoggedIn] = useState(false);
	const [magicCode, setMagicCode] = useState("");
	const providers: Provider[] = [
		{
			id: "github",
			name: "Github",
		} as Provider,
	];

	useEffect(() => {
		checkLogin();
		if (isLoggedIn) {
			return;
		}
		// for each provider, fetch the url from the backend
		providers.forEach(async (provider) => {
			const res = await Service.OAuthURL({
				provider: provider.id,
			} as OAuthURLRequest);
			provider.url = res.url;
			setValidProviders([...validProviders, provider]);
		});
	}, []);

	const checkLogin = () => {
		const oldToken = Cookies.get(TOKEN);
		if (oldToken == undefined) {
			setIsLoggedIn(false);
			return;
		}
		setIsLoggedIn(true);
	};

	if (isLoggedIn) {
		return (
			<div className="flex-grow container mx-auto flex flex-col p-8 items-center">
				<div className="prose m-4">
					<h1>You are already logged in!</h1>
				</div>
			</div>
		);
	}

	const handleLoginWithMagicCode = (code: string) => {
		Service.GenrateToken({ magicCode: code } as GenrateTokenRequest)
			.then((resp) => {
				Cookies.set(TOKEN, resp.token as string);
				Cookies.set(TOKEN_EXPIRE, resp.expiry as string);
				setIsLoggedIn(true);
				router.push(siteConfig.links.home);
			})
			.catch((err) => {
				toast({
					variant: "destructive",
					title: "Login with MagicCode error",
					description: JSON.stringify(err),
				});
			});
	};

	return (
		<div className="flex-grow container mx-auto flex flex-col p-8 items-center">
			<div className="prose m-4">
				<h1>Please Login</h1>
			</div>
			<div className="container mx-auto max-w-sm flex flex-col bg-white p-4 items-center">
				{validProviders.map((provider) => {
					if (provider.url === undefined) {
						return (
							<Button disabled key={provider.id}>
								Sign In with {provider.name}
							</Button>
						);
					}
					return (
						<Link
							href={provider.url}
							key={provider.id}
							className="w-full m-2"
						>
							<Button className="w-full">
								Sign In with {provider.name}
							</Button>
						</Link>
					);
				})}

				<div className="flex w-full max-w-sm items-center space-x-2">
					<Input
						type="text"
						placeholder="Login with MagicCode"
						value={magicCode}
						onChange={(e) => setMagicCode(e.target.value)}
					/>
					<Button
						type="submit"
						onClick={() => handleLoginWithMagicCode(magicCode)}
					>
						Login
					</Button>
				</div>
			</div>
		</div>
	);
}
