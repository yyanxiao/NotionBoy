import { TOKEN, TOKEN_EXPIRE } from "@/components/auth";
import { Button } from "@/components/ui/button";
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { siteConfig } from "@/config/site";
import { useToast } from "@/hooks/use-toast";
import {
	GenrateTokenRequest,
	OAuthProvider,
	OAuthURLRequest,
} from "@/lib/pb/model/common.pb";
import { Service } from "@/lib/pb/server.pb";
import Cookies from "js-cookie";
import { marked } from "marked";
import Link from "next/link";
import { useRouter } from "next/router";
import { useEffect, useState } from "react";

export default function SignIn() {
	const router = useRouter();
	const { toast } = useToast();
	const [isLoggedIn, setIsLoggedIn] = useState(false);
	const [magicCode, setMagicCode] = useState("");
	const [providers, setProviders] = useState<OAuthProvider[]>([]);
	const [wechatQRCodeURL, setWechatQRCodeURL] = useState("");

	useEffect(() => {
		checkLogin();
		if (isLoggedIn) {
			return;
		}

		Service.OAuthProviders({} as OAuthURLRequest).then((res) => {
			setProviders(res.providers as OAuthProvider[]);
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

	const handleLoginWithWechatQRCode = () => {
		let wechatQRCode = "";

		Service.GenerateWechatQRCode({} as OAuthURLRequest)
			.then((res) => {
				wechatQRCode = res.qrcode as string;
				setWechatQRCodeURL(res.url as string);
			})
			.catch((err) => {
				console.log("GenerateWechatQRCode error: ", err);
			});
		const interval = setInterval(() => {
			console.log(
				`wechatQRCode: ${wechatQRCode}, wechatQRCodeURL: ${wechatQRCodeURL}`
			);
			Service.GenrateToken({
				qrcode: wechatQRCode,
			} as GenrateTokenRequest)
				.then((resp) => {
					clearInterval(interval);
					Cookies.set(TOKEN, resp.token as string);
					Cookies.set(TOKEN_EXPIRE, resp.expiry as string);
					setIsLoggedIn(true);
					router.push(siteConfig.links.home);
				})
				.catch((err) => {
					console.log("GenrateToken error: ", err);
				});
		}, 5000);

		setTimeout(() => {
			clearInterval(interval);
			toast({
				variant: "destructive",
				title: "Sign In error",
				description: "Sign In with QRCode timeout",
			});
		}, 1000 * 60 * 5);
	};

	if (isLoggedIn) {
		router.push("/");
		return (
			<div className="container flex flex-col items-center flex-grow p-8 mx-auto">
				<div className="m-4 prose">
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
					title: "Sign In with MagicCode error",
					description: JSON.stringify(err),
				});
			});
	};

	const loginWithWechatQRCode = () => {
		return (
			<Dialog>
				<DialogTrigger
					onClick={() => handleLoginWithWechatQRCode()}
					asChild
				>
					<Button className="w-full">微信扫码登录</Button>
				</DialogTrigger>
				<DialogContent>
					<DialogHeader>
						<DialogTitle>请扫描下方二维码进行登录</DialogTitle>
						<DialogDescription>
							需要先关注微信公众号
						</DialogDescription>
					</DialogHeader>
					<div className="bg-indigo-300">
						<img
							className="object-cover w-full h-full "
							src={wechatQRCodeURL}
						></img>
					</div>
				</DialogContent>
			</Dialog>
		);
	};

	const loginWithMagicCode = () => {
		return (
			<Dialog>
				<DialogTrigger
					onClick={() => handleLoginWithWechatQRCode()}
					asChild
				>
					<Button className="w-full">使用 MagicCode 登录</Button>
				</DialogTrigger>
				<DialogContent>
					<DialogHeader>
						<DialogTitle className="prose">
							使用微信或者 Telegram 获取 MagicCode
						</DialogTitle>
						<DialogDescription
							className="prose"
							dangerouslySetInnerHTML={{
								__html: marked(
									`微信搜索 **NotionBoy** 关注并回复 \`/MagicCode\`
									<br />
									OR
									<br/>
									Telegram 关注 [**NotionBoy**](https://t.me/TheNotionBoyBot) 并回复 \`/MagicCode\`
									`
								),
							}}
						></DialogDescription>
					</DialogHeader>
					<div className="flex items-center w-full max-w-sm space-x-2">
						<Input
							type="text"
							placeholder="Sign In with MagicCode"
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
				</DialogContent>
			</Dialog>
		);
	};

	return (
		<div className="container flex flex-col items-center flex-grow p-8 mx-auto">
			<div className="m-4 prose">
				<h1>Please Login</h1>
			</div>
			<div className="container flex flex-col items-center max-w-sm mx-auto space-y-2 bg-white">
				{providers.map((provider) => {
					if (provider.url === undefined) {
						return (
							<Button disabled key={provider.name}>
								Sign In with {provider.name}
							</Button>
						);
					}
					return (
						<Link
							href={provider.url}
							key={provider.name}
							className="w-full"
						>
							<Button className="w-full">
								Sign In with {provider.name}
							</Button>
						</Link>
					);
				})}
				{loginWithWechatQRCode()}
				{loginWithMagicCode()}
			</div>
		</div>
	);
}
