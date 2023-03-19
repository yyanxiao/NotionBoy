import { Button } from "@/components/ui/button";
import { useToast } from "@/hooks/use-toast";
import {
	Conversation,
	DeleteConversationRequest,
} from "@/lib/pb/model/conversation.pb";
import { Service } from "@/lib/pb/server.pb";
import { ChatContext } from "@/lib/states/chat-context";
import { v4 as uuidv4 } from "uuid";
import { Trash2 } from "lucide-react";
import { useContext } from "react";
import { DefaultInstruction } from "@/config/prompts";

export default function ConversationListComponent() {
	const {
		conversations,
		setConversations,
		selectedConversation,
		setSelectedConversation,
	} = useContext(ChatContext);
	const { toast } = useToast();

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

	return (
		<div className="flex flex-col h-full w-full overflow-auto">
			{conversations.map((conversation) => {
				return (
					<div
						className="flex items-center justify-start m-1 text-sm rounded-lg px-1 border-2 border-gray-400"
						key={conversation.id}
					>
						<Button
							variant="ghost"
							key={conversation.id}
							className={`my-1 flex-1 border-0 ${
								selectedConversation?.id === conversation.id
									? "bg-blue-400"
									: ""
							} }`}
							onClick={() =>
								setSelectedConversation(conversation)
							}
						>
							<div className="text-xs">
								<p>{conversation.title || conversation.id}</p>
								<p>{conversation.createdAt}</p>
							</div>
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
}
