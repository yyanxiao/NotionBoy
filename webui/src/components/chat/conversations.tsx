import { Button } from "@/components/ui/button";
import { useToast } from "@/hooks/use-toast";
import {
	Conversation,
	DeleteConversationRequest,
} from "@/lib/pb/model/conversation.pb";
import { Service } from "@/lib/pb/server.pb";

import { Trash2 } from "lucide-react";

type ConversationListProps = {
	conversations: Conversation[];
	selectedConversation: Conversation | undefined;
	onSelectConversation: (conversation: Conversation | undefined) => void;
	onSetConversations: (conversations: Conversation[]) => void;
};

export default function ConversationListComponent({
	conversations,
	selectedConversation,
	onSelectConversation,
	onSetConversations,
}: ConversationListProps) {
	const { toast } = useToast();

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

	return (
		<div className="flex-grow flex flex-col w-full overflow-auto">
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
							onClick={() => onSelectConversation(conversation)}
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
