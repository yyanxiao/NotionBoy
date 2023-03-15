import { MainNav } from "@/components/main-nav";

import { Auth } from "./auth";
import { MobileNav } from "./mobile-nav";

export function SiteHeader() {
	return (
		<header className="sticky top-0 w-full  border-b border-b-slate-200 bg-white dark:border-b-slate-700 dark:bg-slate-900">
			<div className="flex h-12 items-center justify-between">
				<MainNav />
				<MobileNav />
				<Auth />
			</div>
		</header>
	);
}
