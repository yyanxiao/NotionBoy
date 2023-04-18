import { Button } from "@/components/ui/button";
import { DefaultInstruction } from "@/config/prompts";
import { useToast } from "@/hooks/use-toast";
import {
	Conversation,
	DeleteConversationRequest,
} from "@/lib/pb/model/conversation.pb";
import { Service } from "@/lib/pb/server.pb";
import { ChatContext } from "@/lib/states/chat-context";
import { currentTime, parseDateTime } from "@/lib/utils";
import { Check, Edit2, Trash2 } from "lucide-react";
import { useContext, useEffect, useRef, useState } from "react";
import { v4 as uuidv4 } from "uuid";
import { Input } from "../ui/input";

export default function ConversationListComponent() {
	const {
		conversations,
		setConversations,
		selectedConversation,
		setSelectedConversation,
	} = useContext(ChatContext);
	const { toast } = useToast();
	const inputRef = useRef<HTMLInputElement>(null);

	const [isRenaming, setIsRenaming] = useState<boolean>(false);
	const [title, setTitle] = useState<string>("");

	useEffect(() => {
		if (selectedConversation) {
			setTitle(selectedConversation.title as string);
		}
	}, [selectedConversation]);

	const handleDeleteConversation = (conversation: Conversation) => {
		Service.DeleteConversation({
			id: conversation.id,
		} as DeleteConversationRequest)
			.then(() => {
				setConversations(
					conversations.filter((c) => c.id !== conversation.id)
				);

				if (conversations.length > 0) {
					setSelectedConversation(conversations[0]);
				} else {
					const newCov = {
						id: uuidv4(),
						instruction: DefaultInstruction.instruction,
						title: DefaultInstruction.title,
					} as Conversation;
					setSelectedConversation(newCov);
				}

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

	const handleRenameConversation = (conversation: Conversation) => {
		if (title != conversation.title && title.length > 0) {
			Service.UpdateConversation({
				id: conversation.id,
				instruction: conversation.instruction,
				title: title,
			})
				.then((newConversation) => {
					setSelectedConversation(newConversation);
					conversations.filter((c) => c.id !== conversation.id);
					setConversations(
						conversations.map((c) => {
							if (c.id === conversation.id) {
								return newConversation;
							}
							return c;
						})
					);
					toast({
						title: "Success",
						description: `Conversation ${conversation.id} renamed`,
					});
				})
				.catch((error) => {
					toast({
						variant: "destructive",
						title: "Error",
						description: `Failed to rename conversation, ${JSON.stringify(
							error
						)}`,
					});
				});
		}
		setIsRenaming(false);
	};

	const handleEdit = () => {
		setIsRenaming(true);
		console.log(inputRef);
		if (inputRef.current) {
			inputRef.current.focus();
		}
	};

	const showIcon = (conversation: Conversation) => {
		const edit = () => {
			if (isRenaming && conversation.id === selectedConversation?.id) {
				return (
					<Button
						size={"sm"}
						variant="ghost"
						onClick={() => handleRenameConversation(conversation)}
					>
						<Check size={14} />
					</Button>
				);
			} else {
				return (
					<Button size={"sm"} variant="ghost" onClick={handleEdit}>
						<Edit2 size={14} />
					</Button>
				);
			}
		};
		return (
			<>
				{edit()}
				{conversation.id === selectedConversation?.id && (
					<Button
						size={"sm"}
						variant="ghost"
						onClick={() => handleDeleteConversation(conversation)}
					>
						<Trash2 size={14} />
					</Button>
				)}
			</>
		);
	};

	const showConversation = (conversation: Conversation) => {
		if (isRenaming && conversation.id === selectedConversation?.id) {
			return (
				<Input
					type="text"
					ref={inputRef}
					className="flex-1 border-0 border-b-2 rounded-none focus:ring-0 focus:border-[#3da9fc]"
					value={title}
					minLength={1}
					onBlur={() => {
						setIsRenaming(false);
					}}
					onKeyDown={(e) => {
						if (e.key === "Enter" && title.length > 0) {
							handleRenameConversation(conversation);
						}
					}}
					onChange={(e) => {
						setTitle(e.target.value);
					}}
				/>
			);
		} else {
			return (
				<Button
					variant="ghost"
					key={conversation.id}
					className={`my-1 lg:w-60 w-48 border-0 hover:bg-[#3da9fc] ${
						selectedConversation?.id === conversation.id
							? "bg-[#3da9fc]"
							: ""
					} }`}
					onClick={() => setSelectedConversation(conversation)}
				>
					<div className="text-xs">
						<p>{conversation.title || conversation.id}</p>
						<p>
							{conversation.createdAt
								? parseDateTime(conversation.createdAt)
								: currentTime()}
						</p>
					</div>
				</Button>
			);
		}
	};

	return (
		<div className="flex flex-col w-full h-full overflow-auto">
			{conversations.map((conversation) => {
				return (
					<div className="px-1" key={conversation.id}>
						<div
							className="flex items-center px-1 m-1 text-sm rounded-lg"
							key={conversation.id}
						>
							{showConversation(conversation)}
							{conversation.id === selectedConversation?.id &&
								showIcon(conversation)}
						</div>
					</div>
				);
			})}
		</div>
	);
}
