import { MainNav } from "@/components/main-nav";

import { AuthLoginButton } from "./auth";

export function SiteHeader() {
	return (
		<header className="flex items-center justify-between md:sticky top-0 w-full  border-b border-b-slate-200 bg-white dark:border-b-slate-700 dark:bg-slate-900">
			<MainNav />
			<AuthLoginButton />
		</header>
	);
}
