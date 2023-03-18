import { Separator } from "@/components/ui/separator";
import { DefaultInstruction, Instruction } from "@/config/prompts";
import { siteConfig } from "@/config/site";

import { Conversation } from "@/lib/pb/model/conversation.pb";

import { Home, LogOut, Plus } from "lucide-react";
import Link from "next/link";
import { useState } from "react";
import { AuthLoginButton } from "../auth";
import { Button } from "../ui/button";
import ConversationListComponent from "./conversations";
import { ChatSettings } from "./settings";
import { v4 as uuidv4 } from "uuid";

type ConversationListProps = {
	conversations: Conversation[];
	selectedConversation: Conversation | undefined;
	onSelectConversation: (conversation: Conversation | undefined) => void;
	onSetConversations: (conversations: Conversation[]) => void;
};

export function SideBarComponent({
	conversations,
	selectedConversation,
	onSelectConversation,
	onSetConversations,
}: ConversationListProps) {
	const [instruction, setInstruction] =
		useState<Instruction>(DefaultInstruction);

	const handleCreateConversation = () => {
		const conversation = {
			id: uuidv4(),
			createdAt: new Date().toISOString(),
			updatedAt: new Date().toISOString(),
			instruction: instruction.instruction,
			title: instruction.title,
		} as Conversation;
		onSelectConversation(conversation);
		onSetConversations([conversation, ...conversations]);
	};

	return (
		<div className="bg-gray-100 text-gray-800 h-screen overflow-hidden">
			<div className="relative flex flex-col justify-between h-full ">
				<div className="sticky top-0 left-0 flex flex-row items-center justify-center">
					<Button
						variant="ghost"
						className="self-center"
						onClick={handleCreateConversation}
					>
						<Plus />
					</Button>
					<ChatSettings
						conversations={conversations}
						selectedConversation={selectedConversation}
						onSelectConversation={onSelectConversation}
						onSetConversations={onSetConversations}
					/>
				</div>
				<ConversationListComponent
					conversations={conversations}
					selectedConversation={selectedConversation}
					onSelectConversation={onSelectConversation}
					onSetConversations={onSetConversations}
				/>
				{/* <Separator /> */}
				<div className="sticky bottom-0 left-0 container mx-auto flex flex-col items-start p-2 bg-gray-300">
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
				</div>
			</div>
		</div>
	);
}
