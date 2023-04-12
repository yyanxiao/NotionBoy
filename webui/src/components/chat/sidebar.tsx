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
		<div className="h-full overflow-hidden text-[#fffffe] bg-[#094067] rounded-lg my-2 md:my-0">
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
				<Separator />

				<ConversationListComponent />
				<Separator />
				<div className="container sticky bottom-0 left-0 flex flex-col items-center px-6 py-2 mx-auto">
					<div className="flex flex-row items-center w-full h-10 p-2">
						<Home />
						<Link className="p-2" href={siteConfig.links.home}>
							Notionboy
						</Link>
					</div>
					<div className="flex flex-row items-center w-full h-10 p-2">
						<LogOut />
						<AuthLoginButton />
					</div>

					<div className="flex flex-row w-full">
						<Link
							href={siteConfig.links.github}
							target="_blank"
							rel="noreferrer"
						>
							<div
								className={buttonVariants({
									size: "sm",
									variant: "ghost",
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
