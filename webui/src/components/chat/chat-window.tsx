import { Conversation, Message } from "@/lib/pb/model/conversation.pb";
import { Bot, User } from "lucide-react";
import { marked } from "marked";
import { useEffect, useRef } from "react";

type Props = {
	messages: Message[] | undefined;
	selectedConversation: Conversation;
};

function parseMarkdown(text: string) {
	return { __html: marked(text) };
}

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
		<div ref={messagesEndRef} className="flex-1 flex flex-col">
			{props.selectedConversation ? (
				<div className="flex-1 w-full flex flex-col overflow-auto box-border rounded-lg p-4">
					{props.messages?.map((message) => {
						return (
							<div key={message.id}>
								{message.request !== undefined &&
									message.request !== "" && (
										<div
											className="w-full flex flex-row justify-start items-center"
											key={`${message.id}-req`}
										>
											<User className="flex-none w-8 h-8 " />
											<div
												className="prose dark:prose-invert lg:prose-xl text-start bg-blue-200 m-1 p-2 rounded-xl overflow-auto"
												dangerouslySetInnerHTML={parseMarkdown(
													message.request
												)}
											></div>
										</div>
									)}
								{message.response !== undefined &&
									message.response !== "" && (
										<div
											className=" w-full flex flex-row justify-end items-center"
											key={`${message.id}-resp`}
										>
											<div
												className="prose dark:prose-invert lg:prose-xl text-start bg-green-200 p-2 rounded-xl overflow-auto"
												dangerouslySetInnerHTML={parseMarkdown(
													message.response
												)}
											></div>
											<Bot className="flex-none w-8 h-8" />
										</div>
									)}
							</div>
						);
					})}
				</div>
			) : (
				<div className="flex-1 w-full flex flex-col overflow-auto box-border rounded-lg p-4">
					<p>No conversation selected</p>
				</div>
			)}
		</div>
	);
}
