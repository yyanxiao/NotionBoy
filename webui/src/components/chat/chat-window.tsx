import { useToast } from "@/hooks/use-toast";
import { Conversation, Message } from "@/lib/pb/model/conversation.pb";
import { parseDateTime } from "@/lib/utils";

import "highlight.js/styles/github.css";
import { Bot, Copy, User } from "lucide-react";

import { useEffect, useRef } from "react";
import { Button } from "../ui/button";

import { MarkdownComponent } from "./markdown";
import { RoleDialog } from "./role";

type Props = {
	messages: Message[] | undefined;
	selectedConversation: Conversation;
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
					<div className="absolute top-0 right-0">
						<Button
							variant="ghost"
							className="py-0"
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
		return (
			<div className="fixed h-48 md:w-96 top-1/3 left-1/4 lg:left-1/2">
				<div className="grid justify-between grid-cols-2 gap-2">
					<div className="bg-[#abd1c6] rounded-lg w-full">
						<RoleDialog />
					</div>
					{/* <div className="bg-[#abd1c6] rounded-lg w-full">
						<UploadFileComponent />
					</div> */}
				</div>
			</div>
		);
	};

	return (
		<div ref={messagesEndRef} className="flex flex-col flex-1 mt-2">
			<div className="box-border flex flex-col overflow-auto">
				{renderMessages()}
			</div>
		</div>
	);
}
