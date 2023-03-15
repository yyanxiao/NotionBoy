import {
	Inspect,
	List,
	Mail,
	Menu,
	MessageSquare,
	Plus,
	PlusCircle,
	Settings,
	Settings2,
	SidebarOpen,
	Trash2,
	UserPlus,
} from "lucide-react";
import { v4 as uuidv4 } from "uuid";
import { Button } from "@/components/ui/button";
import {
	Conversation,
	DeleteConversationRequest,
} from "@/lib/pb/model/conversation.pb";
import {
	Sheet,
	SheetContent,
	SheetDescription,
	SheetHeader,
	SheetTitle,
	SheetTrigger,
} from "@/components/ui/sheet";
import { Service } from "@/lib/pb/server.pb";
import { useToast } from "@/hooks/use-toast";
import {
	Popover,
	PopoverContent,
	PopoverTrigger,
} from "@radix-ui/react-popover";
import { useState } from "react";
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuLabel,
	DropdownMenuPortal,
	DropdownMenuSeparator,
	DropdownMenuSub,
	DropdownMenuSubContent,
	DropdownMenuSubTrigger,
	DropdownMenuTrigger,
} from "@radix-ui/react-dropdown-menu";
import {
	DefaultInstruction,
	Instruction,
	InstructionList,
	InstructionMap,
} from "@/config/prompts";
import {
	HoverCard,
	HoverCardContent,
	HoverCardTrigger,
} from "@radix-ui/react-hover-card";
import { Textarea } from "../ui/textarea";

type ConversationListProps = {
	conversations: Conversation[];
	selectedConversation: Conversation | undefined;
	onSelectConversation: (conversation: Conversation | undefined) => void;
	onSetConversations: (conversations: Conversation[]) => void;
};

export default function ConversationList({
	conversations,
	selectedConversation,
	onSelectConversation,
	onSetConversations,
}: ConversationListProps) {
	const { toast } = useToast();
	const [instruction, setInstruction] =
		useState<Instruction>(DefaultInstruction);

	const handleDeleteConversation = (conversation: Conversation) => {
		Service.DeleteConversation({
			id: conversation.id,
		} as DeleteConversationRequest)
			.then(() => {
				onSetConversations(
					conversations.filter((c) => c.id !== conversation.id)
				);
				const newCovId =
					conversations.length > 0 ? conversations[0] : undefined;
				onSelectConversation(newCovId);
				toast({
					title: "Success",
					description: `Conversation ${conversation.id} deleted`,
				});
			})
			.catch((error) => {
				toast({
					variant: "destructive",
					title: "Error",
					description: `Failed to delete conversation, ${JSON.stringify(
						error
					)}`,
				});
			});
	};

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
						className="border-2 border-stone-400 rounded-md px-2"
						onSelect={() => handleSelectInstruction(instruction)}
					>
						{instruction.title}
					</DropdownMenuItem>
				</HoverCardTrigger>
				<HoverCardContent>
					<div className="prose bg-stone-100 text-black text-sm p-2 rounded-lg">
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
							<div className="prose bg-stone-100 text-black text-sm p-2 rounded-lg">
								<p>Select your instruction</p>
							</div>
						</HoverCardContent>
					</HoverCard>
				</DropdownMenuTrigger>

				<DropdownMenuContent className="rounded-md border max-h-96  bg-stone-50 border-slate-100 p-1 shadow-md dark:border-slate-800  dark:text-slate-400 w-56 overflow-auto">
					<DropdownMenuItem className="border-2 bg-blue-200 rounded-md px-2">
						<strong>{`Selected: ${instruction?.title}`}</strong>
					</DropdownMenuItem>
					{dropdownMenuItem(DefaultInstruction)}
					{InstructionList.map(({ key, data }) => {
						return (
							<DropdownMenuSub key={key}>
								<DropdownMenuSubTrigger className="border-2 border-stone-400 rounded-md px-2">
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

	const listConversationsComponent = () => {
		return (
			<div className="flex flex-col bg-stone-100 p-1 rounded-3xl">
				<Button
					variant="ghost"
					className="self-center mb-2"
					onClick={handleCreateConversation}
				>
					<Plus />
				</Button>
				{conversations.map((conversation) => {
					return (
						<div
							className="flex items-center justify-start m-1 text-sm rounded-lg border border-stone-600 px-1"
							key={conversation.id}
						>
							<Button
								variant="ghost"
								key={conversation.id}
								className={`my-1 flex-1 ${
									selectedConversation?.id === conversation.id
										? "bg-blue-400"
										: ""
								} }`}
								onClick={() =>
									onSelectConversation(conversation)
								}
							>
								{conversation.title || conversation.id}
							</Button>
							<Button
								onClick={() =>
									handleDeleteConversation(conversation)
								}
							>
								<Trash2 />
							</Button>
						</div>
					);
				})}
			</div>
		);
	};

	return (
		<div className="flex justify-between h-10">
			<Popover>
				<PopoverTrigger>
					<List />
				</PopoverTrigger>
				<PopoverContent className="m-2">
					{listConversationsComponent()}
				</PopoverContent>
			</Popover>
			<div className="self-center">{selectedConversation?.title}</div>
			<div className="flex flex-row">
				<Button
					variant="ghost"
					size="sm"
					onClick={handleCreateConversation}
				>
					<Plus />
				</Button>
				{setInstructionComponent()}
				{/* <Popover>
					<PopoverTrigger className="px-2">
						<Settings />
					</PopoverTrigger>
					<PopoverContent className="bg-stone-100 w-96 h-96 p-1 rounded-3xl">
					</PopoverContent>
				</Popover> */}
			</div>
		</div>
	);
}
