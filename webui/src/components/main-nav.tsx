import * as React from "react";
import Link from "next/link";

import { siteConfig } from "@/config/site";

import {
	NavigationMenu,
	NavigationMenuItem,
	NavigationMenuLink,
	NavigationMenuList,
	navigationMenuTriggerStyle,
} from "@/components/ui/navigation-menu";

export function MainNav() {
	const link = (url: string, name: string) => {
		return (
			<Link href={url} legacyBehavior passHref>
				<NavigationMenuLink className={navigationMenuTriggerStyle()}>
					<span className="font-bold">{name}</span>
				</NavigationMenuLink>
			</Link>
		);
	};

	return (
		<div className="flex">
			<NavigationMenu>
				<NavigationMenuList>
					<NavigationMenuItem>
						{link(siteConfig.links.home, siteConfig.name)}
						{link(siteConfig.links.chatgpt, "Chat")}
						{link(siteConfig.links.price, "Pricing")}
					</NavigationMenuItem>
				</NavigationMenuList>
			</NavigationMenu>
		</div>
	);
}
