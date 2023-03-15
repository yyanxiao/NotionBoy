import { cn } from "@/lib/utils";
import { SiteFooter } from "./site-footer";
import { SiteHeader } from "./site-header";
import { Toaster } from "./ui/toaster";

export const Layout = ({ children }: { children: React.ReactNode }) => {
	return (
		<div className={cn("min-h-screen container mx-auto flex flex-col ")}>
			<SiteHeader />
			<Toaster />
			{children}
			<SiteFooter />
		</div>
	);
};
