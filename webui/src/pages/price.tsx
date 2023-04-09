import { SiteHeader } from "@/components/site-header";
import { siteConfig } from "@/config/site";
import { useToast } from "@/hooks/use-toast";
import { Product } from "@/lib/pb/model/product.pb";
import { Service } from "@/lib/pb/server.pb";
import { Check } from "lucide-react";
import Link from "next/link";
import { useEffect, useState } from "react";

export default function PricePage() {
	const [products, setProducts] = useState<Product[]>([]);

	const { toast } = useToast();

	useEffect(() => {
		// fetch products
		Service.ListProducts({})
			.then((res) => {
				if (res.products) {
					res.products.sort(
						(a, b) => (b.price as number) - (a.price as number)
					);

					setProducts(res.products);
				} else {
					toast({
						variant: "destructive",
						title: "Errors",
						description: "Fetch Products error",
					});
				}
			})
			.catch((err) => {
				toast({
					variant: "destructive",
					title: "Errors",
					description: JSON.stringify(err),
				});
			});
	}, []);

	const planComponent = (p: Product) => {
		const features = [`${p.token} OpenAI Token`, `${p.storage}MB Storage`];

		return (
			<div
				key={p.id}
				className="max-w-2xl mx-auto mt-16 rounded-3xl ring-1 ring-gray-200 sm:mt-20 lg:mx-0 lg:flex lg:max-w-none"
			>
				<div className="p-8 sm:p-10 lg:flex-auto">
					<h3 className="text-2xl font-bold tracking-tight text-gray-900">
						{p.name}
					</h3>
					<p className="mt-6 text-base leading-7 text-gray-600">
						{p.description}
					</p>
					<div className="flex items-center mt-10 gap-x-4">
						<h4 className="flex-none text-sm font-semibold leading-6 text-indigo-600">
							What’s included
						</h4>
						<div className="flex-auto h-px bg-gray-100" />
					</div>
					<ul
						role="list"
						className="grid grid-cols-1 gap-4 mt-8 text-sm leading-6 text-gray-600 sm:grid-cols-2 sm:gap-6"
					>
						{features.map((feature) => (
							<li key={feature} className="flex gap-x-3">
								<Check className="flex-none w-5 h-6 text-indigo-600" />
								{feature}
							</li>
						))}
					</ul>
				</div>
				<div className="p-2 -mt-2 lg:mt-0 lg:w-full lg:max-w-md lg:flex-shrink-0">
					<div className="py-10 text-center rounded-2xl bg-gray-50 ring-1 ring-inset ring-gray-900/5 lg:flex lg:flex-col lg:justify-center lg:py-16">
						<div className="max-w-xs px-8 mx-auto">
							<p className="text-base font-semibold text-gray-600"></p>
							<p className="flex items-baseline justify-center mt-6 gap-x-2">
								<span className="text-3xl font-bold tracking-tight text-gray-900">
									¥ {p.price} / 月
								</span>
							</p>
							<Link
								href={`${siteConfig.links.order}?product=${p.id}`}
								className={`block w-full px-3 py-2 mt-10 text-sm font-semibold text-center text-white bg-indigo-600 rounded-md shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 ${
									p.price === 0
										? "opacity-50 cursor-not-allowed"
										: ""
								}}`}
							>
								购买
							</Link>
						</div>
					</div>
				</div>
			</div>
		);
	};

	return (
		<>
			<SiteHeader />

			<div className="py-24 bg-white sm:py-32">
				<div className="px-6 mx-auto max-w-7xl lg:px-8">
					<div className="max-w-2xl mx-auto sm:text-center">
						<h2 className="text-3xl font-bold tracking-tight text-gray-900 sm:text-4xl">
							NotionBoy 会员计划
						</h2>
						<p className="mt-6 text-lg leading-8 text-gray-600">
							NotionBoy 的会员主要包含两个部分:
							一个是服务，一个是用量. 服务是指 NotionBoy
							所提供的所有服务，例如 Chat，Note， Zlib 等。
							用量是指第三方的服务用量，比如保存图片所用的存储，使用
							OpenAI 的 API 消耗的 Token。
							<strong>
								所有的会员都能使用 NotionBoy
								的所有服务，哪怕不付费，不同的是付费会员有更多的用量。
							</strong>
							<Link
								href="https://www.notion.so/vaayne/2a4906985bc647209a96069bfebfdc3d"
								target={"_blank"}
								className={`inline-flex  self-center w-64 px-3 py-2 mt-10 text-sm font-semibold items-center text-white bg-indigo-600 rounded-md shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600`}
							>
								详细的介绍请查看账户与会员 FAQ
							</Link>
						</p>
					</div>

					{products.map((p) => planComponent(p))}
				</div>
			</div>
		</>
	);
}
