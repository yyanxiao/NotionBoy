import { Conversation, Message } from "@/lib/pb/model/conversation.pb";

import "highlight.js/styles/github.css";
import { Bot, User } from "lucide-react";

import { useEffect, useRef } from "react";
import { MarkdownComponent } from "./markdown";
type Props = {
	messages: Message[] | undefined;
	selectedConversation: Conversation;
};

export default function ChatWindow(props: Props) {
	const messagesEndRef = useRef<HTMLDivElement | null>(null);
	useEffect(() => {
		if (messagesEndRef.current) {
			messagesEndRef.current.scrollTo({
				top: messagesEndRef.current.scrollHeight,
				behavior: "smooth",
			});
		}
	}, [props.messages, messagesEndRef]);
	return (
		<div ref={messagesEndRef} className="flex flex-col flex-1">
			{props.selectedConversation ? (
				<div className="box-border flex flex-col flex-1 w-full p-4 overflow-auto rounded-lg">
					{props.messages?.map((message) => {
						return (
							<div key={message.id}>
								{message.request !== undefined &&
									message.request !== "" && (
										<div
											className="flex flex-row items-center justify-end w-full my-2"
											key={`${message.id}-req`}
										>
											<div className="p-2 mx-2 overflow-auto prose bg-blue-200 dark:prose-invert lg:prose-lg text-start rounded-xl">
												<MarkdownComponent
													text={message.request}
												/>
											</div>
											<User className="flex-none w-8 h-8 " />
										</div>
									)}
								{message.response !== undefined &&
									message.response !== "" && (
										<div
											className="flex flex-row items-center justify-start w-full my-2 "
											key={`${message.id}-resp`}
										>
											<Bot className="flex-none w-8 h-8" />
											<div className="p-2 mx-2 overflow-auto prose bg-green-200 dark:prose-invert lg:prose-lg text-start rounded-xl">
												<MarkdownComponent
													text={message.response}
												/>
											</div>
										</div>
									)}
							</div>
						);
					})}
				</div>
			) : (
				<div className="box-border flex flex-col flex-1 w-full p-4 overflow-auto rounded-lg">
					<p>No conversation selected</p>
				</div>
			)}
		</div>
	);
}
