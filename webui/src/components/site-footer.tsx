import { siteConfig } from "@/config/site";
import { Icons } from "@/components/icons";
import Link from "next/link";
import { buttonVariants } from "./ui/button";

export function SiteFooter() {
	return (
		<footer className="container mx-auto border-b border-b-slate-200 bg-white dark:border-b-slate-700 dark:bg-slate-900">
			<div className="flex flex-row items-center justify-between px-2 border-t border-t-slate-200 py-1 dark:border-t-slate-700 md:flex-row md:py-0">
				<p className="text-center text-sm leading-loose text-slate-600 dark:text-slate-400 md:text-center">
					Built by{" "}
					<a
						href={siteConfig.links.twitter}
						target="_blank"
						rel="noreferrer"
						className="font-medium underline underline-offset-4"
					>
						Vaayne
					</a>
				</p>
				<nav className="flex items-center space-x-1">
					<Link
						href={siteConfig.links.github}
						target="_blank"
						rel="noreferrer"
					>
						<div
							className={buttonVariants({
								size: "sm",
								variant: "ghost",
								className: "text-slate-700 dark:text-slate-400",
							})}
						>
							<Icons.gitHub className="h-5 w-5" />
							<span className="sr-only">GitHub</span>
						</div>
					</Link>
					<Link
						href={siteConfig.links.twitter}
						target="_blank"
						rel="noreferrer"
					>
						<div
							className={buttonVariants({
								size: "sm",
								variant: "ghost",
								className: "text-slate-700 dark:text-slate-400",
							})}
						>
							<Icons.twitter className="h-5 w-5 fill-current" />
							<span className="sr-only">Twitter</span>
						</div>
					</Link>
				</nav>
			</div>
		</footer>
	);
}
