import { useToast } from "@/hooks/use-toast";
import {
	Conversation,
	CreateMessageRequest,
	Message,
	UpdateMessageRequest,
} from "@/lib/pb/model/conversation.pb";
import { useEffect, useState } from "react";
import { v4 as uuidv4 } from "uuid";

import { Service } from "@/lib/pb/server.pb";
import { useRouter } from "next/router";

import ChatWindow from "@/components/chat/chat-window";
import { ChatInputBox } from "@/components/chat/input-box";
import { siteConfig } from "@/config/site";
import { currentTime, isLogin } from "@/lib/utils";

import MobileChatHeader from "@/components/chat/mobile-chat-header";
import { SideBarComponent } from "@/components/chat/sidebar";
import { DefaultInstruction } from "@/config/prompts";
import { ChatContext, MessageContext } from "@/lib/states/chat-context";

interface ChatResponse {
	message: string;
}

export default function Chat() {
	const [conversations, setConversations] = useState<Conversation[]>([]);
	const [selectedConversation, setSelectedConversation] =
		useState<Conversation>(newConversation());
	const [conversationMessagesMap, setConversationMessagesMap] = useState<
		Map<string, Message[]>
	>(new Map());
	const [messages, setMessages] = useState<Message[]>([]);
	const [isLoading, setIsLoading] = useState(false);
	const [model, setModel] = useState<string>("gpt-3.5-turbo");
	const [temperature, setTemperature] = useState<number>(1);
	const [maxTokens, setMaxTokens] = useState<number>(1000);

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
			Service.ListConversations({
				limit: 100,
			})
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
		if (conversationMessagesMap.has(selectedConversation.id as string)) {
			setMessages(
				conversationMessagesMap.get(
					selectedConversation.id as string
				) as []
			);
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
		Service.ListMessages({
			conversationId: selectedConversation.id,
			limit: 100,
		})
			.then((response) => {
				if (response.messages === undefined) {
					return;
				}
				if (selectedConversation.id) {
					conversationMessagesMap.set(
						selectedConversation.id,
						response.messages
					);
					setMessages(response.messages);
					setConversationMessagesMap(
						new Map(conversationMessagesMap)
					);
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

	const handlerMessage = (
		messages: Message[],
		message: Message,
		newMessage: Message,
		resp: ChatResponse
	) => {
		if (newMessage.response === undefined) {
			return;
		}
		if (resp.message != newMessage.response) {
			resp.message += newMessage.response;
		}
		message.response = resp.message;
		if (newMessage.updatedAt != "0001-01-01T00:00:00Z") {
			message.updatedAt = newMessage.updatedAt;
		} else {
			message.updatedAt = currentTime();
		}
		// set id to uuid of DB record
		if (newMessage.id != "00000000-0000-0000-0000-000000000000") {
			message.id = newMessage.id;
		}
		message.tokenUsage = newMessage.tokenUsage;
		handleAddMessae(messages, message);

		// check if selected conversation in conversations
		// if not exists, add it to conversations as it is a new conversation
		const conversation = conversations.find(
			(c) => c.id === selectedConversation.id
		);
		if (conversation === undefined) {
			setConversations([selectedConversation, ...conversations]);
		}
	};

	const handleAddMessae = (messages: Message[], newMessage: Message) => {
		messages = messages.filter((m) => m.id !== newMessage.id);
		const newMessages = [...messages, newMessage];
		// 		console.log(`messages: ${JSON.stringify(messages)},
		// \nnewMessage: ${JSON.stringify(newMessage)},
		// \nnewMessages: ${JSON.stringify(newMessages)},
		// 			`);
		setMessages(newMessages);
		conversationMessagesMap.set(
			selectedConversation.id as string,
			newMessages
		);
		setConversationMessagesMap(new Map(conversationMessagesMap));
	};

	const handleMessageUpdate = (
		message: Message,
		model: string,
		temperature: number,
		maxTokens: number
	) => {
		const request = message.request;
		if (request === undefined || request === "") {
			return;
		}
		setIsLoading(true);
		const updateMessageRequest = {
			conversationId: message.conversationId,
			id: message.id,
			model: model,
			request: request,
			temperature: temperature,
			maxTokens: maxTokens,
		} as UpdateMessageRequest;
		const newMessage = {
			...message,
			model: model,
			temperature: temperature,
			maxTokens: maxTokens,
			response: "",
		};

		// filter out the message with create at newer than the message to be updated
		const filteredMessages =
			messages?.filter(
				(m) =>
					new Date(m.createdAt as string).getTime() <
					new Date(message.createdAt as string).getTime()
			) || [];
		console.log("filteredMessages", filteredMessages);
		handleAddMessae(filteredMessages, newMessage);
		const resp = {
			message: "",
		} as ChatResponse;
		// send message request to server
		Service.UpdateMessage(updateMessageRequest, (msg) => {
			handlerMessage(filteredMessages, newMessage, msg, resp);
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

	const handleMessageSend = (
		request: string,
		model: string,
		temperature: number,
		maxTokens: number
	) => {
		if (request === undefined || request === "") {
			return;
		}
		setIsLoading(true);
		const createMessageRequest = {
			conversationId: selectedConversation?.id,
			request: request,
			model: model,
			temperature: temperature,
			maxTokens: maxTokens,
		} as CreateMessageRequest;

		let message: Message = {
			id: uuidv4(),
			conversationId: selectedConversation?.id as string,
			request: request,
			createdAt: currentTime(),
		};
		// add message request to messageMap
		handleAddMessae(messages, message);

		const resp = {
			message: "",
		} as ChatResponse;
		// send message request to server
		Service.CreateMessage(createMessageRequest, (msg) => {
			handlerMessage(messages, message, msg, resp);
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

	const handleMessageDelete = (conversationId: string, messageId: string) => {
		setIsLoading(true);
		Service.DeleteMessage({
			conversationId: conversationId,
			id: messageId,
		})
			.then(() => {
				const messages = conversationMessagesMap
					.get(conversationId)
					?.filter((m) => m.id !== messageId);
				setMessages(messages as Message[]);
				conversationMessagesMap.set(
					conversationId,
					messages as Message[]
				);
				setConversationMessagesMap(new Map(conversationMessagesMap));
				toast({
					variant: "default",
					title: "DeleteMessage success",
				});
			})
			.catch((error) => {
				toast({
					variant: "destructive",
					title: "DeleteMessage error",
					description: JSON.stringify(error),
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
			<MessageContext.Provider
				value={{
					selectedConversation: selectedConversation as Conversation,
					isLoading,
					messages: messages,
					model: model,
					setModel: setModel,
					temperature: temperature,
					setTemperature: setTemperature,
					maxTokens: maxTokens,
					setMaxTokens: setMaxTokens,
					onMessageSend: handleMessageSend,
					onMessageUpdate: handleMessageUpdate,
					onMessageDelete: handleMessageDelete,
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
									<ChatWindow />
								</div>
								<div className="sticky bottom-0 bg-[#fffffe]">
									<ChatInputBox />
								</div>
							</div>
						</div>
					</div>
				</div>
			</MessageContext.Provider>
		</ChatContext.Provider>
	);
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
