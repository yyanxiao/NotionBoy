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
	return (
		<div className="flex">
			<NavigationMenu>
				<NavigationMenuList>
					<NavigationMenuItem>
						<Link href="/" legacyBehavior passHref>
							<NavigationMenuLink
								className={navigationMenuTriggerStyle()}
							>
								<span className="font-bold">
									{siteConfig.name}
								</span>
							</NavigationMenuLink>
						</Link>
						<Link
							href={siteConfig.links.chatgpt}
							legacyBehavior
							passHref
						>
							<NavigationMenuLink
								className={navigationMenuTriggerStyle()}
							>
								<span className="font-bold">Chat</span>
							</NavigationMenuLink>
						</Link>
					</NavigationMenuItem>
				</NavigationMenuList>
			</NavigationMenu>
		</div>
	);
}
