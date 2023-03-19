import { siteConfig } from "@/config/site";

import { Home, LogOut, Plus } from "lucide-react";
import Link from "next/link";
import { useContext } from "react";
import { AuthLoginButton } from "../auth";
import { Button, buttonVariants } from "../ui/button";
import ConversationListComponent from "./conversations";
import { ChatSettings } from "./settings";
import { ChatContext } from "@/lib/states/chat-context";
import { Separator } from "../ui/separator";
import { Icons } from "../icons";

export function SideBarComponent() {
	const {
		conversations,
		setConversations,
		selectedConversation,
		setSelectedConversation,
		handleCreateConversation,
	} = useContext(ChatContext);

	return (
		<div className="bg-gray-100 text-gray-800 h-full overflow-hidden">
			<div className="relative flex flex-col justify-between h-full ">
				<div className="sticky top-0 left-0 flex flex-row items-center justify-center">
					<Button
						variant="ghost"
						className="self-center"
						onClick={handleCreateConversation}
					>
						<Plus />
					</Button>
					<ChatSettings />
				</div>

				<ConversationListComponent />

				<div className="sticky bottom-0 left-0 container mx-auto flex flex-col items-start p-2 bg-gray-400 rounded-lg">
					<div className="w-full h-10 flex flex-row items-center p-2">
						<Home />
						<Link className="px-2" href={siteConfig.links.home}>
							Notionboy
						</Link>
					</div>
					<div className="w-full h-10 flex flex-row items-center p-2">
						<LogOut />
						<AuthLoginButton />
					</div>
					<Separator />
					<div className="flex flex-row my-2">
						<Link
							href={siteConfig.links.github}
							target="_blank"
							rel="noreferrer"
						>
							<div
								className={buttonVariants({
									size: "sm",
									variant: "ghost",
									className:
										"text-slate-700 dark:text-slate-400",
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
									className:
										"text-slate-700 dark:text-slate-400",
								})}
							>
								<Icons.twitter className="h-5 w-5 fill-current" />
								<span className="sr-only">Twitter</span>
							</div>
						</Link>
					</div>
				</div>
			</div>
		</div>
	);
}
