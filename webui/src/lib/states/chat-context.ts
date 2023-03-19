import { createContext, Dispatch, SetStateAction } from "react";
import { Conversation } from "../pb/model/conversation.pb";

interface ChatContextProps {
	conversations: Conversation[];
	setConversations: Dispatch<SetStateAction<Conversation[]>>;
	selectedConversation: Conversation | undefined;
	setSelectedConversation: Dispatch<SetStateAction<Conversation>>;
	handleCreateConversation: () => void;
}

export const ChatContext = createContext({} as ChatContextProps);
