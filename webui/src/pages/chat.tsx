import { useState, useEffect } from "react";
import { useToast } from "@/hooks/use-toast";
import { v4 as uuidv4 } from "uuid";
import {
	Conversation,
	CreateMessageRequest,
	Message,
} from "@/lib/pb/model/conversation.pb";

import { Service } from "@/lib/pb/server.pb";
import Cookies from "js-cookie";
import { useRouter } from "next/router";
import ConversationList from "@/components/chat/conversation-list";
import ChatWindow from "@/components/chat/chat-window";
import { ChatInputBox } from "@/components/chat/input-box";
import { isLogin } from "@/lib/utils";
import { siteConfig } from "@/config/site";

export default function Chat() {
	const [conversations, setConversations] = useState<Conversation[]>([]);
	const [selectedConversation, setSelectedConversation] =
		useState<Conversation>();
	const [messageMap, setMessageMap] = useState<Map<string, Message[]>>(
		new Map()
	);
	const [isLoading, setIsLoading] = useState(false);

	const router = useRouter();
	const { toast } = useToast();

	useEffect(() => {
		// check if user is logged in
		if (!isLogin()) {
			router.push(siteConfig.links.login);
			return;
		}

		// list conversatons when pages loads and there are no conversations
		if (conversations.length == 0) {
			Service.ListConversations({})
				.then((response) => {
					console.log("ListConversations", response);
					if (response.conversations === undefined) {
						return;
					}
					setConversations(response.conversations);
				})
				.catch((error) => {
					toast({
						variant: "destructive",
						title: "ListConversations error",
						description: JSON.stringify(error),
					});
				});
		}
		// create conversation on page load
		const conversation = {
			id: uuidv4(),
			createdAt: new Date().toISOString(),
			updatedAt: new Date().toISOString(),
		} as Conversation;
		setConversations([conversation, ...conversations]);

		// select conversation on page load
		setSelectedConversation(conversation);
	}, []);

	// handle conversation selection with messages
	useEffect(() => {
		if (selectedConversation === undefined) {
			return;
		}
		console.log("selectedConversation", selectedConversation);
		console.log("messageMap", messageMap);
		if (messageMap.has(selectedConversation.id as string)) {
			return;
		}
		setIsLoading(true);
		Service.ListMessages({ conversationId: selectedConversation.id })
			.then((response) => {
				console.log("ListMessages", response);
				if (response.messages === undefined) {
					return;
				}
				if (selectedConversation.id) {
					messageMap.set(selectedConversation.id, response.messages);
					setMessageMap(new Map(messageMap));
				}
			})
			.catch((error) => {
				toast({
					variant: "destructive",
					title: "ListMessage error",
					description: JSON.stringify(error),
				});
			})
			.finally(() => {
				setIsLoading(false);
			});
	}, [selectedConversation]);

	const handleSelectConversation = (
		conversation: Conversation | undefined
	) => {
		if (conversation === undefined) {
			return;
		}
		setSelectedConversation(conversation);
		if (messageMap && messageMap.has(conversation.id as string)) {
			return;
		}
		setIsLoading(true);
		Service.ListMessages({ conversationId: conversation.id })
			.then((response) => {
				if (response.messages === undefined) {
					toast({
						variant: "destructive",
						description: "No messages for the conversation",
					});
					return;
				}
				if (response.messages.length == 0) {
					toast({
						description: "No messages for the conversation",
					});
					return;
				}

				if (conversation.id) {
					messageMap.set(conversation.id, response.messages);
					setMessageMap(new Map(messageMap));
				}
			})
			.catch((error) => {
				toast({
					variant: "default",
					title: "ListMessages Error",
					description: error.message,
				});
			})
			.finally(() => {
				setIsLoading(false);
			});
	};

	const handleMessageSend = (request: string) => {
		setIsLoading(true);
		if (request === undefined || request === "") {
			return;
		}

		const createMessageRequest = {
			conversationId: selectedConversation?.id,
			request: request,
		} as CreateMessageRequest;

		let message: Message = {
			id: uuidv4(),
			conversationId: selectedConversation?.id as string,
			request: request,
		};
		// add message request to messageMap
		addMessageToMessageMap(
			selectedConversation?.id as string,
			message,
			messageMap
		);
		setMessageMap(new Map(messageMap));

		// send message request to server
		Service.CreateMessage(createMessageRequest)
			.then((response) => {
				message.response = response.response;
				message.createdAt = response.createdAt;
				message.updatedAt = response.updatedAt;
				message.tokenUsage = response.tokenUsage;
			})
			.catch((error) => {
				toast({
					variant: "destructive",
					title: "Send Message Failed",
					description: error.message,
				});
			})
			.finally(() => {
				setIsLoading(false);
			});

		addMessageToMessageMap(
			selectedConversation?.id as string,
			message,
			messageMap
		);
		setMessageMap(new Map(messageMap));
	};

	return (
		<div className="flex-grow container mx-auto flex flex-col ">
			<ConversationList
				conversations={conversations}
				selectedConversation={selectedConversation}
				onSelectConversation={handleSelectConversation}
				onSetConversations={setConversations}
			/>
			<ChatWindow
				messages={messageMap.get(selectedConversation?.id as string)}
				selectedConversation={selectedConversation as Conversation}
			/>
			<ChatInputBox
				onSendMessage={handleMessageSend}
				isLoading={isLoading}
			/>
		</div>
	);
}

function addMessageToMessageMap(
	conversationId: string,
	message: Message,
	messageMap: Map<string, Message[]>
) {
	const messages = messageMap.get(conversationId);
	if (messages === undefined) {
		messageMap.set(conversationId, [message]);
	} else {
		if (messages.length > 0) {
			const latestMessage = messages[messages.length - 1];
			if (latestMessage.id === message.id) {
				messages[messages.length - 1] = message;
			} else {
				messages.push(message);
			}
		} else {
			messages.push(message);
		}

		messageMap.set(conversationId, messages);
	}
}
