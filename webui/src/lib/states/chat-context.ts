import { createContext, Dispatch, SetStateAction } from "react";
import { Conversation, Message } from "../pb/model/conversation.pb";

interface ChatContextProps {
	conversations: Conversation[];
	setConversations: Dispatch<SetStateAction<Conversation[]>>;
	selectedConversation: Conversation | undefined;
	setSelectedConversation: Dispatch<SetStateAction<Conversation>>;
	handleCreateConversation: () => void;
}

export const ChatContext = createContext({} as ChatContextProps);

interface MessageContextProps {
	selectedConversation: Conversation;
	isLoading: boolean;
	messages: Message[] | undefined;
	model: string;
	temperature: number;
	maxTokens: number;
	setModel: Dispatch<SetStateAction<string>>;
	setTemperature: Dispatch<SetStateAction<number>>;
	setMaxTokens: Dispatch<SetStateAction<number>>;
	onMessageSend: (
		message: string,
		model: string,
		temperature: number,
		maxTokens: number
	) => void;
	onMessageUpdate: (
		message: Message,
		model: string,
		temperature: number,
		maxTokens: number
	) => void;
	onMessageDelete: (conversationID: string, messageId: string) => void;
}

export const MessageContext = createContext({} as MessageContextProps);
