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
		<div className="h-full overflow-hidden text-gray-800 bg-gray-100">
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

				<div className="container sticky bottom-0 left-0 flex flex-col items-start p-2 mx-auto bg-gray-400 rounded-lg">
					<div className="flex flex-row items-center w-full h-10 p-2">
						<Home />
						<Link className="px-2" href={siteConfig.links.home}>
							Notionboy
						</Link>
					</div>
					<div className="flex flex-row items-center w-full h-10 p-2">
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
								<Icons.gitHub className="w-5 h-5" />
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
								<Icons.twitter className="w-5 h-5 fill-current" />
								<span className="sr-only">Twitter</span>
							</div>
						</Link>
					</div>
				</div>
			</div>
		</div>
	);
}
