import { useEffect, useState } from "react";

import { SiteHeader } from "@/components/site-header";
import { Button } from "@/components/ui/button";
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from "@/components/ui/dialog";
import { siteConfig } from "@/config/site";
import { useToast } from "@/hooks/use-toast";
import { Order, PayOrderConfig } from "@/lib/pb/model/order.pb";
import { Product } from "@/lib/pb/model/product.pb";
import { Service } from "@/lib/pb/server.pb";
import { isLogin } from "@/lib/utils";
import { Loader2 } from "lucide-react";
import Head from "next/head";
import { useRouter } from "next/router";
import { QRCodeSVG } from "qrcode.react";
import WeChatPayment from "@/components/WeChatPayment";
import Script from "next/script";

export default function Example() {
	const [productId, setProductID] = useState<string>("");
	const [product, setProduct] = useState<Product>();

	const [isLoading, setIsLoading] = useState<boolean>(false);

	const [wxpayQrcode, setWxpayQrcode] = useState<string>();
	const [wxpayConfig, setWxpayConfig] = useState<PayOrderConfig>();

	const { toast } = useToast();
	const router = useRouter();

	useEffect(() => {
		if (!isLogin()) {
			router.push(siteConfig.links.login);
			return;
		}
	}, []);

	useEffect(() => {
		if (router.isReady) {
			const { product } = router.query;
			if (product) {
				setProductID(product.toString());
			}
		}
	}, [router]);

	// show product for the order
	useEffect(() => {
		if (productId) {
			Service.GetProduct({ id: productId })
				.then((resp) => {
					setProduct(resp);
				})
				.catch((err) => {
					toast({
						variant: "destructive",
						title: "Get product error",
						description: JSON.stringify(err),
					});
				});
		}
	}, [productId]);

	const handleCheckout = async () => {
		// create order
		let order: Order;
		try {
			order = await Service.CreateOrder({ productId: productId });
			if (!order) {
				return;
			}
		} catch (err) {
			toast({
				variant: "destructive",
				title: "Create order error",
				description: JSON.stringify(err),
			});
			return;
		}

		try {
			const resp = await Service.PayOrder({ id: order.id });
			if (resp.qrcode) {
				setWxpayQrcode(resp.qrcode);
			} else if (resp.config) {
				setWxpayConfig(resp.config);
			} else {
				toast({
					variant: "destructive",
					title: "Get QRcdode error",
					description: "Do not get qrcode",
				});
			}
		} catch (err) {
			toast({
				variant: "destructive",
				title: "pay order error",
				description: JSON.stringify(err),
			});
		}

		const interval = setInterval(() => {
			if (order) {
				Service.GetOrder({ id: order.id }).then((resp) => {
					if (resp.status == "Completed") {
						clearInterval(interval);
						setIsLoading(false);
						toast({
							variant: "default",
							title: "Success pay for order",
							description: "Order success",
						});
						router.push("/");
					} else if (
						resp.status == "Paid" ||
						resp.status == "Processing"
					) {
						setIsLoading(true);
					}
				});
			}
		}, 1000 * 5);

		const timeout = setTimeout(() => {
			clearInterval(interval);
			setIsLoading(false);
		}, 1000 * 60 * 5);
		return () => {
			setIsLoading(false);
			clearInterval(interval);
			clearTimeout(timeout);
		};
	};

	const payWithWechat = () => {
		if (window.navigator.userAgent.indexOf("MicroMessenger") != -1) {
			return (
				<>
					<Dialog>
						<DialogTrigger onClick={handleCheckout} asChild>
							<Button className="w-full">微信支付</Button>
						</DialogTrigger>
						<DialogContent>
							<DialogHeader>
								<DialogTitle>微信支付</DialogTitle>
								<DialogDescription>
									请稍后，正在为您跳转到微信支付
								</DialogDescription>
							</DialogHeader>
							<div className="inline-flex items-center justify-center w-full ">
								{loaddingShow()}
								{wxpayConfig && (
									<WeChatPayment
										cfg={wxpayConfig}
										onSuccess={() => {
											() => setIsLoading(true);
										}}
									></WeChatPayment>
								)}
							</div>
						</DialogContent>
					</Dialog>
				</>
			);
		} else {
			return (
				<Dialog>
					<DialogTrigger onClick={handleCheckout} asChild>
						<Button className="w-full">微信支付</Button>
					</DialogTrigger>
					<DialogContent>
						<DialogHeader>
							<DialogTitle>扫码支付</DialogTitle>
							<DialogDescription>
								请扫描下方二维码进行进行支付
							</DialogDescription>
						</DialogHeader>
						<div className="inline-flex items-center justify-center w-full ">
							{loaddingShow()}
							{wxpayQrcode && (
								<QRCodeSVG
									className="w-1/3"
									size={192}
									value={wxpayQrcode}
								/>
							)}
						</div>
					</DialogContent>
				</Dialog>
			);
		}
	};

	const loaddingShow = () => {
		if (isLoading) {
			return (
				<div className="fixed inset-0 z-10 flex items-center justify-center bg-black bg-opacity-50">
					<Loader2
						className="w-10 h-10 text-white animate-spin"
						strokeWidth="2"
					/>
				</div>
			);
		}
	};

	const showOrder = () => {
		if (product) {
			return (
				<>
					<div className="bg-gray-50">
						<div className="max-w-2xl p-10 mx-auto my-24 bg-white border border-gray-200 rounded-lg shadow-sm ">
							<div className="flex flex-col flex-1 m-12">
								<div className="flex">
									<div className="flex-1 min-w-0">
										<h4 className="text-lg">
											{product.name}
										</h4>
										<ul></ul>

										<p className="mt-1 text-gray-500">
											{product.description}
										</p>
										<ul>
											<li className="mt-1 text-gray-500">
												- {product.token} OpenAI Token
											</li>
											<li className="mt-1 text-gray-500">
												- {product.storage} MB Storage
											</li>
										</ul>
									</div>
								</div>

								<div className="flex justify-end flex-1 pt-2">
									<p className="mt-1 font-medium text-gray-900">
										总价 ¥ {product.price}
									</p>
								</div>
							</div>
							{payWithWechat()}
						</div>
					</div>
				</>
			);
		}
	};
	return (
		<>
			<Script src="https://res.wx.qq.com/open/js/jweixin-1.6.0.js"></Script>
			<SiteHeader />
			{/* {loaddingShow()} */}
			{showOrder()}
		</>
	);
}
