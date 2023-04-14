import { useToast } from "@/hooks/use-toast";
import { Conversation, Message } from "@/lib/pb/model/conversation.pb";
import { parseDateTime } from "@/lib/utils";

import "highlight.js/styles/github.css";
import { Bot, Copy, Trash2, User } from "lucide-react";

import { useEffect, useRef } from "react";
import { Button } from "../ui/button";

import { MarkdownComponent } from "./markdown";

type Props = {
	messages: Message[] | undefined;
	selectedConversation: Conversation;
	onMessageDelete: (conversationId: string, messageId: string) => void;
};

export default function ChatWindow(props: Props) {
	const messagesEndRef = useRef<HTMLDivElement | null>(null);
	const { toast } = useToast();
	useEffect(() => {
		if (messagesEndRef.current) {
			messagesEndRef.current.scrollTo({
				top: messagesEndRef.current.scrollHeight,
				behavior: "smooth",
			});
		}
	}, [props.messages, messagesEndRef]);

	const messageComponents = (message: Message, isResponse: boolean) => {
		const md = () => {
			if (isResponse && message.response) {
				return <MarkdownComponent text={message.response} />;
			} else if (message.request) {
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
							props.onMessageDelete(
								props.selectedConversation.id as string,
								message.id as string
							);
						}}
					>
						<Trash2 size={18} />
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
					{md()}
				</div>
			</div>
		);
	};

	const renderMessages = () => {
		if (props.messages && props.messages.length > 0) {
			return props.messages.map((message) => {
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
