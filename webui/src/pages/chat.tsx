import { useRef, useState, useEffect } from "react";
import { Menu, Send, Plus } from "lucide-react";
import { marked } from "marked";

import {
	Conversation,
	CreateMessageRequest,
	Message,
} from "@/lib/pb/model/conversation.pb";

import { v4 as uuidv4 } from "uuid";
import { Service } from "@/lib/pb/server.pb";

function parseMarkdown(text: string) {
	return { __html: marked(text) };
}

export default function Chat() {
	const [conversations, setConversations] = useState<Conversation[]>([]);
	const [currentConversation, setCurrentConversation] =
		useState<Conversation>({
			id: uuidv4(),
		} as Conversation);
	const [currentConversationId, setCurrentConversationId] =
		useState<string>("");
	const [currentMessages, setCurrentMessages] = useState<Message[]>([]);
	const [inputValue, setInputValue] = useState<string>("");
	const [isSidebarOpen, setIsSidebarOpen] = useState<boolean>(false);
	const [isLoading, setIsLoading] = useState(false);

	// on page load
	useEffect(() => {
		let conversations: Conversation[] = [];
		Service.ListConversations({})
			.then((response) => {
				console.log("ListConversations", response);
				if (response.conversations === undefined) {
					return;
				}
				conversations = response.conversations;
				setConversations(conversations);
			})
			.catch((error) => {
				console.log(error);
			});
		const conversation = {
			id: uuidv4(),
			createdAt: new Date().toISOString(),
			updatedAt: new Date().toISOString(),
		} as Conversation;
		setCurrentConversation(conversation);
		setConversations([...conversations, conversation]);
	}, []);

	// on currentConversation change, set currentConversationId and fetch messages
	useEffect(() => {
		if (currentConversation.id === undefined) {
			return;
		}
		console.log("currentConversation", JSON.stringify(currentConversation));
		setCurrentConversationId(currentConversation.id);
		if (currentConversation.messages === undefined) {
			handleListMessages(currentConversationId);
		}
	}, [currentConversation]);

	// create conversation on local, only create on server when user sends message
	const handleCreateConversation = () => {
		const conversation = {
			id: uuidv4(),
			createdAt: new Date().toISOString(),
			updatedAt: new Date().toISOString(),
			messages: [],
		} as Conversation;
		setCurrentConversation(conversation);
		setConversations([...conversations, conversation]);
	};

	const handleListMessages = (conversationId: string) => {
		if (conversationId === undefined || conversationId === "") {
			return;
		}
		Service.ListMessages({ conversationId: conversationId })
			.then((response) => {
				console.log(
					"handleConversationClick list messages: ",
					response
				);
				response.messages?.sort((a, b) => {
					if (a.createdAt === undefined) {
						return 1;
					} else if (b.createdAt === undefined) {
						return -1;
					}
					return (
						new Date(a.createdAt).getTime() -
						new Date(b.createdAt).getTime()
					);
				});
				// update current conversation and messages
				setCurrentMessages(response.messages ?? []);
				let conversation = currentConversation;
				conversation.messages = response.messages ?? [];
				setCurrentConversation(conversation);
			})
			.catch((error) => {
				console.log(
					"handleConversationClick list message error: ",
					error
				);
			});
	};

	const handleConversationClick = (conversation: Conversation) => {
		if (conversation.id != undefined) {
			conversations.forEach((c) => {
				if (c.id === conversation.id) {
					setCurrentConversation(c);
				}
			});
		}
	};

	const leftSideBar = () => {
		return (
			<div className="h-full flex flex-col max-w-lg bg-gray-400 text-white">
				<button
					key="button"
					className="btn"
					onClick={handleCreateConversation}
				>
					<Plus />
				</button>
				{conversations.map((conversation) => {
					return (
						<button
							key={conversation.id}
							className="btn hover:bg-slate-500 focus:bg-slate-800"
							value={currentConversationId}
							onClick={() =>
								handleConversationClick(conversation)
							}
						>
							<p>{conversation.id}</p>
						</button>
					);
				})}
			</div>
		);
	};

	const inputWithSendButton = () => {
		return (
			<div className="relative">
				<input
					type="text"
					placeholder="Type here"
					className="input input-bordered w-full p-2"
					onChange={(e) => setInputValue(e.target.value)}
					value={inputValue}
					onKeyDown={(e) => {
						if (e.key === "Enter") {
							handleMessageSend(inputValue);
						}
					}}
				/>
				<button
					className="absolute right-2 mt-2 p-1 disabled:opacity-50"
					disabled={isLoading}
					onClick={() => handleMessageSend(inputValue)}
				>
					<Send />
				</button>
			</div>
		);
	};

	const renderMessages = (messages: Message[]) => {
		return messages.map((message) => {
			return (
				<div key={message.id}>
					{message.request !== undefined &&
						message.request !== "" && (
							<div
								className="chat chat-start"
								key={`${message.id}-req`}
							>
								<div
									className="chat-bubble chat-bubble-accent prose text-sm"
									dangerouslySetInnerHTML={parseMarkdown(
										message.request
									)}
								></div>
							</div>
						)}
					{message.response !== undefined &&
						message.response !== "" && (
							<div
								className="chat chat-end"
								key={`${message.id}-resp`}
							>
								<div
									className="chat-bubble chat-bubble-info prose text-sm"
									dangerouslySetInnerHTML={parseMarkdown(
										message.response
									)}
								></div>
							</div>
						)}
				</div>
			);
		});
	};

	const handleMessageSend = (request: string | undefined) => {
		if (request === undefined || request === "") {
			return;
		}
		setIsLoading(true);
		console.log("handleMessageSend", currentConversation.id, request);
		setInputValue("");
		const createMessageRequest = {
			conversationId: currentConversation.id,
			request: request,
		} as CreateMessageRequest;

		let message: Message = {
			id: uuidv4(),
			conversationId: currentConversation.id,
			request: request,
		};
		setCurrentMessages([...currentMessages, message]);

		Service.CreateMessage(createMessageRequest)
			.then((response) => {
				setCurrentMessages([...currentMessages, response]);
			})
			.catch((error) => {
				// todo: handle error
				console.log(error);
			});
		setIsLoading(false);
		setInputValue("");
	};

	const handleMenuClick = () => {
		setIsSidebarOpen(!isSidebarOpen);
	};

	return (
		<div className="w-full h-screen flex flex-row overflow-hidden">
			<header>
				<title>NotionBoy Chat</title>
			</header>
			{isSidebarOpen && leftSideBar()}
			<div className="flex-1 w-full flex flex-col bg-gray-2">
				<div className="w-full h-12 bg-success dark:bg-black items-center justify-between flex flex-row rounded-lg">
					<button
						className="btn bg-success dark:bg-black border-0"
						onClick={handleMenuClick}
					>
						<Menu />
					</button>
					<div className="prose">
						<h1>NotionBoy Chat</h1>
					</div>
					<div></div>
				</div>

				<div className="flex-1 w-full flex flex-col overflow-auto box-border bg-stone-100 dark:bg-neutral-focus rounded-lg p-4">
					{renderMessages(currentMessages)}
				</div>
				{isLoading && <progress className="progress w-full"></progress>}
				<div className="w-full">{inputWithSendButton()}</div>
			</div>
		</div>
	);
}
