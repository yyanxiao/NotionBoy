import { Button } from "@/components/ui/button";
import {
	DefaultInstruction,
	Instruction,
	InstructionList,
} from "@/config/prompts";
import { useToast } from "@/hooks/use-toast";
import { Conversation } from "@/lib/pb/model/conversation.pb";
import { Service } from "@/lib/pb/server.pb";

import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuPortal,
	DropdownMenuSub,
	DropdownMenuSubContent,
	DropdownMenuSubTrigger,
	DropdownMenuTrigger,
} from "@radix-ui/react-dropdown-menu";
import {
	HoverCard,
	HoverCardContent,
	HoverCardTrigger,
} from "@radix-ui/react-hover-card";

import { Plus, Settings } from "lucide-react";
import { useState } from "react";
import { v4 as uuidv4 } from "uuid";
import { ChatSettings } from "./settings";
import { SideSheetComponent } from "./mobile-sidebar";

type ConversationListProps = {
	conversations: Conversation[];
	selectedConversation: Conversation | undefined;
	onSelectConversation: (conversation: Conversation | undefined) => void;
	onSetConversations: (conversations: Conversation[]) => void;
};

export default function MobileChatHeader({
	conversations,
	selectedConversation,
	onSelectConversation,
	onSetConversations,
}: ConversationListProps) {
	const { toast } = useToast();
	const [instruction, setInstruction] =
		useState<Instruction>(DefaultInstruction);

	const handleUpdateConversation = (instruction: Instruction) => {
		const conversation = {
			...selectedConversation,
			instruction: instruction.instruction,
			title: instruction.title,
		} as Conversation;
		Service.UpdateConversation(conversation)
			.then((resp) => {
				onSelectConversation(conversation);
				onSetConversations(
					conversations.map((c) => {
						if (c.id === conversation.id) {
							return conversation;
						}
						return c;
					})
				);
			})
			.catch((error) => {
				toast({
					variant: "destructive",
					title: "Error",
					description: `Failed to update conversation, ${JSON.stringify(
						error
					)}`,
				});
			});
	};

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

	const handleSelectInstruction = (instruction: Instruction) => {
		setInstruction(instruction);
		handleUpdateConversation(instruction);
	};

	const dropdownMenuItem = (instruction: Instruction) => {
		return (
			<HoverCard key={instruction.title}>
				<HoverCardTrigger>
					<DropdownMenuItem
						className="border-2 bg-white text-black rounded-md px-2"
						onSelect={() => handleSelectInstruction(instruction)}
					>
						{instruction.title}
					</DropdownMenuItem>
				</HoverCardTrigger>
				<HoverCardContent>
					<div className="prose bg-white text-black text-sm p-2 rounded-lg">
						<strong>EN:</strong>
						<p>{instruction.instruction}</p>
						<strong>中文:</strong>
						<p>{instruction.instructioncn}</p>
					</div>
				</HoverCardContent>
			</HoverCard>
		);
	};

	const setInstructionComponent = () => {
		return (
			<DropdownMenu>
				<DropdownMenuTrigger>
					<HoverCard>
						<HoverCardTrigger>
							<Settings />
						</HoverCardTrigger>
						<HoverCardContent>
							<div className="prose bg-white text-black text-sm p-2 rounded-lg">
								<p>Select your instruction</p>
							</div>
						</HoverCardContent>
					</HoverCard>
				</DropdownMenuTrigger>

				<DropdownMenuContent className="rounded-md border max-h-96  bg-white text-black border-white p-1 shadow-md dark:border-slate-800  dark:text-slate-400 w-56 overflow-auto">
					<DropdownMenuItem className="border-2 bg-white text-black rounded-md px-2">
						<strong>{`Selected: ${instruction?.title}`}</strong>
					</DropdownMenuItem>
					{dropdownMenuItem(DefaultInstruction)}
					{InstructionList.map(({ key, data }) => {
						return (
							<DropdownMenuSub key={key}>
								<DropdownMenuSubTrigger className="border-2 bg-white text-black rounded-md px-2">
									{key}
								</DropdownMenuSubTrigger>
								<DropdownMenuPortal>
									<DropdownMenuSubContent>
										{data?.map((v) => {
											return dropdownMenuItem(v);
										})}
									</DropdownMenuSubContent>
								</DropdownMenuPortal>
							</DropdownMenuSub>
						);
					})}
				</DropdownMenuContent>
			</DropdownMenu>
		);
	};

	return (
		<div className="flex flex-row items-center justify-between mx-2">
			<SideSheetComponent
				conversations={conversations}
				selectedConversation={selectedConversation}
				onSelectConversation={onSelectConversation}
				onSetConversations={onSetConversations}
			/>

			<div className="">{selectedConversation?.title}</div>
			<div className="flex flex-row">
				<Button
					variant="ghost"
					size="sm"
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
		</div>
	);
}
