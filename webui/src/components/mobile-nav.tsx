import * as React from "react";
import Link from "next/link";

import { siteConfig } from "@/config/site";

import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Icons } from "./icons";
import { Menu } from "lucide-react";

export function MobileNav() {
	return (
		<div className="md:hidden mx-1">
			<DropdownMenu>
				<DropdownMenuTrigger>
					<Menu />
				</DropdownMenuTrigger>
				<DropdownMenuContent>
					<DropdownMenuItem>
						<Link href="/" className="flex items-center">
							<Icons.home className="mr-2 h-4 w-4" /> Home
						</Link>
					</DropdownMenuItem>
					<DropdownMenuItem>
						<Link
							href={siteConfig.links.chatgpt}
							className="flex items-center"
						>
							<Icons.message className="mr-2 h-4 w-4" /> Chat
						</Link>
					</DropdownMenuItem>
				</DropdownMenuContent>
			</DropdownMenu>
		</div>
	);
}
