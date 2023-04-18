import { useToast } from "@/hooks/use-toast";
import { Message } from "@/lib/pb/model/conversation.pb";
import { MessageContext } from "@/lib/states/chat-context";
import { parseDateTime } from "@/lib/utils";

import "highlight.js/styles/github.css";
import { Bot, Check, Copy, Edit2, Trash2, User } from "lucide-react";

import { useContext, useEffect, useRef, useState } from "react";
import { Button } from "../ui/button";
import { Input } from "../ui/input";

import { MarkdownComponent } from "./markdown";

export default function ChatWindow() {
	const {
		selectedConversation,
		isLoading,
		messages,
		model,
		temperature,
		maxTokens,
		setModel,
		setTemperature,
		setMaxTokens,
		onMessageSend,
		onMessageUpdate,
		onMessageDelete,
	} = useContext(MessageContext);

	const [isEditing, setIsEditing] = useState<boolean>(false);
	const [inputValue, setInputValue] = useState<string>("");
	const [selectedMessage, setSelectedMessage] = useState<Message>();

	const messagesEndRef = useRef<HTMLDivElement | null>(null);
	const { toast } = useToast();
	useEffect(() => {
		if (messagesEndRef.current) {
			messagesEndRef.current.scrollTo({
				top: messagesEndRef.current.scrollHeight,
				behavior: "smooth",
			});
		}
	}, [messages, messagesEndRef]);

	useEffect(() => {
		if (selectedMessage) {
			setInputValue(selectedMessage.request as string);
		}
	}, [selectedMessage]);

	const messageComponents = (message: Message, isResponse: boolean) => {
		const md = () => {
			if (isResponse && message.response) {
				return <MarkdownComponent text={message.response} />;
			} else if (message.request) {
				if (isEditing && selectedMessage?.id === message.id) {
					return (
						<Input
							className="w-full disabled:opacity-50"
							onChange={(e) => setInputValue(e.target.value)}
							value={inputValue}
							disabled={isLoading}
							onBlur={() => {
								setIsEditing(false);
							}}
							onKeyDown={(e) => {
								// using shift + enter to send a message
								// using enter to create a new line
								if (
									e.key === "Enter" &&
									inputValue.trim() != ""
								) {
									e.preventDefault();
									handleMessageEdit(message);
								} else if (e.key === "Escape") {
									setIsEditing(false);
								}
							}}
						/>
					);
				}
				return <MarkdownComponent text={message.request} />;
			} else {
				return null;
			}
		};
		const icon = () => {
			return (
				<div className="w-8 h-8">{isResponse ? <Bot /> : <User />}</div>
			);
		};

		const copyMessage = () => {
			return (
				<div className="absolute top-0 right-0">
					<Button
						variant="ghost"
						className="p-0 mx-1"
						size={"sm"}
						onClick={() => {
							navigator.clipboard.writeText(
								isResponse
									? (message.response as string)
									: (message.request as string)
							);
							toast({
								title: "Copied to clipboard",
								variant: "default",
							});
						}}
					>
						<Copy size={18} />
					</Button>
				</div>
			);
		};

		const deleteMessage = () => {
			return (
				<div className="absolute top-0 right-6">
					<Button
						variant="ghost"
						className="p-0 mx-1"
						size={"sm"}
						onClick={() => {
							onMessageDelete(
								selectedConversation.id as string,
								message.id as string
							);
						}}
					>
						<Trash2 size={18} />
					</Button>
				</div>
			);
		};
		const handleMessageEdit = (message: Message) => {
			const newMessage = {
				...message,
				request: inputValue,
				response: "",
			} as Message;
			onMessageUpdate(newMessage, model, temperature, maxTokens);
			setIsEditing(false);
		};

		const editMessage = (message: Message) => {
			if (isResponse) {
				return;
			}
			if (isEditing && selectedMessage?.id === message.id) {
				return (
					<div className="absolute top-0 right-12">
						<Button
							variant="ghost"
							className="p-0 mx-1"
							size={"sm"}
							onClick={() => handleMessageEdit(message)}
						>
							<Check size={18} />
						</Button>
					</div>
				);
			}
			return (
				<div className="absolute top-0 right-12">
					<Button
						variant="ghost"
						className="p-0 mx-1"
						size={"sm"}
						onClick={() => {
							setIsEditing(true);
							setSelectedMessage(message);
						}}
					>
						<Edit2 size={18} />
					</Button>
				</div>
			);
		};

		return (
			<div
				className={`flex flex-row items-start justify-start w-full rounded-lg p-2 `}
				key={`${message.id}-resp`}
			>
				{icon()}
				<div
					className={`relative flex flex-col w-full text-sm rounded-lg ${
						isResponse ? "bg-sky-100" : "bg-blue-100"
					}`}
				>
					<div className="px-2 my-1">
						<strong>{isResponse ? "Bot" : "User"}</strong>
						<span className="px-2">
							{parseDateTime(message.updatedAt as string)}
						</span>
						{isResponse && message.tokenUsage && (
							<strong className="px-2">
								Token: {message.tokenUsage}
							</strong>
						)}
					</div>
					{copyMessage()}
					{deleteMessage()}
					{editMessage(message)}
					{md()}
				</div>
			</div>
		);
	};

	const renderMessages = () => {
		if (messages && messages.length > 0) {
			return messages.map((message) => {
				return (
					<div key={message.id} className="flex flex-col">
						{message.request && (
							<div>{messageComponents(message, false)}</div>
						)}
						{message.response && (
							<div>{messageComponents(message, true)}</div>
						)}
					</div>
				);
			});
		}
	};

	return (
		<div ref={messagesEndRef} className="flex flex-col flex-1 mt-2">
			<div className="box-border flex flex-col overflow-auto">
				{renderMessages()}
			</div>
		</div>
	);
}
