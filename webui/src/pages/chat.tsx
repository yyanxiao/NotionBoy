import { useState, useEffect } from "react";
import { useToast } from "@/hooks/use-toast";
import { v4 as uuidv4 } from "uuid";
import {
	Conversation,
	CreateMessageRequest,
	Message,
} from "@/lib/pb/model/conversation.pb";

import { Service } from "@/lib/pb/server.pb";
import { useRouter } from "next/router";

import ChatWindow from "@/components/chat/chat-window";
import { ChatInputBox } from "@/components/chat/input-box";
import { currentTime, isLogin } from "@/lib/utils";
import { siteConfig } from "@/config/site";

import { SideBarComponent } from "@/components/chat/sidebar";
import MobileChatHeader from "@/components/chat/mobile-chat-header";
import { DefaultInstruction } from "@/config/prompts";
import { ChatContext } from "@/lib/states/chat-context";

export default function Chat() {
	const [conversations, setConversations] = useState<Conversation[]>([]);
	const [selectedConversation, setSelectedConversation] =
		useState<Conversation>(newConversation());
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

		if (selectedConversation === undefined) {
			// create conversation on page load
			const conversation = newConversation();
			setConversations([conversation, ...conversations]);

			// select conversation on page load
			setSelectedConversation(conversation);
		} else {
			// check if conversation is already in the list
			const conversation = conversations.find(
				(c) => c.id === selectedConversation.id
			);
			if (conversation === undefined) {
				setConversations([selectedConversation, ...conversations]);
			}
		}
	}, []);

	// handle conversation selection with messages
	useEffect(() => {
		if (selectedConversation === undefined) {
			return;
		}
		// router.asPath = `/chat/${selectedConversation.id}`;
		if (messageMap.has(selectedConversation.id as string)) {
			return;
		}

		// check if selected conversation in conversations
		// if not exists, add it to conversations as it is a new conversation
		const conversation = conversations.find(
			(c) => c.id === selectedConversation.id
		);
		if (conversation === undefined) {
			setConversations([selectedConversation, ...conversations]);
		}

		setIsLoading(true);
		Service.ListMessages({ conversationId: selectedConversation.id })
			.then((response) => {
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

	const handleMessageSend = (request: string, model: string) => {
		setIsLoading(true);
		if (request === undefined || request === "") {
			return;
		}

		const createMessageRequest = {
			conversationId: selectedConversation?.id,
			request: request,
			model: model,
		} as CreateMessageRequest;

		console.log("createMessageRequest", createMessageRequest);

		let message: Message = {
			id: uuidv4(),
			conversationId: selectedConversation?.id as string,
			request: request,
			createdAt: currentTime(),
		};
		// add message request to messageMap
		addMessageToMessageMap(
			selectedConversation?.id as string,
			message,
			messageMap
		);
		setMessageMap(new Map(messageMap));
		let fullMessage = "";
		// send message request to server
		Service.CreateMessage(createMessageRequest, (msg) => {
			if (msg.response === undefined) {
				return;
			}
			if (fullMessage != msg.response) {
				fullMessage += msg.response;
			}
			message.response = fullMessage;
			if (msg.createdAt != "0001-01-01T00:00:00Z") {
				message.createdAt = msg.createdAt;
			}
			if (msg.updatedAt != "0001-01-01T00:00:00Z") {
				message.updatedAt = msg.updatedAt;
			} else {
				message.updatedAt = currentTime();
			}
			message.tokenUsage = msg.tokenUsage;
			addMessageToMessageMap(
				selectedConversation?.id as string,
				message,
				messageMap
			);
			setMessageMap(new Map(messageMap));
			// check if selected conversation in conversations
			// if not exists, add it to conversations as it is a new conversation
			const conversation = conversations.find(
				(c) => c.id === selectedConversation.id
			);
			if (conversation === undefined) {
				setConversations([selectedConversation, ...conversations]);
			}
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
	};

	const handleCreateConversation = () => {
		const conversation = newConversation();
		setSelectedConversation(conversation);
		setConversations([conversation, ...conversations]);
	};

	return (
		<ChatContext.Provider
			value={{
				conversations,
				setConversations,
				selectedConversation,
				setSelectedConversation,
				handleCreateConversation,
			}}
		>
			<div>
				<div className="hidden lg:block fixed left-0 top-0 bottom-0 w-[19.5rem]">
					<SideBarComponent />
				</div>
				<div className="lg:pl-[19.5rem] bg-[#fffffe]">
					<div className="flex flex-col h-full max-w-6xl mx-auto">
						<div className="relative flex flex-col min-h-screen">
							<div className="sticky top-0 left-0 h-10 rounded-sm lg:hidden ">
								<MobileChatHeader />
							</div>
							<div className="flex-grow">
								<ChatWindow
									messages={messageMap.get(
										selectedConversation?.id as string
									)}
									selectedConversation={
										selectedConversation as Conversation
									}
								/>
							</div>
							<div className="sticky bottom-0 bg-[#fffffe]">
								<ChatInputBox
									onSendMessage={handleMessageSend}
									isLoading={isLoading}
								/>
							</div>
						</div>
					</div>
				</div>
			</div>
		</ChatContext.Provider>
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

function newConversation(): Conversation {
	return {
		id: uuidv4(),
		instruction: DefaultInstruction.instruction,
		title: DefaultInstruction.title,
		createdAt: new Date().toISOString(),
		updatedAt: new Date().toISOString(),
	} as Conversation;
}
