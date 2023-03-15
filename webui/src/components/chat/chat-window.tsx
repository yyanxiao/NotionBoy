import { Conversation, Message } from "@/lib/pb/model/conversation.pb";
import { marked } from "marked";

type Props = {
	messages: Message[] | undefined;
	selectedConversation: Conversation;
};

function parseMarkdown(text: string) {
	return { __html: marked(text) };
}

export default function ChatWindow(props: Props) {
	return (
		<div className="flex-grow container mx-auto flex flex-col items-center justify-center overflow-hidden bg-white border-2 border-b-0 rounded-lg border-blue-200">
			{props.selectedConversation ? (
				<div className="flex-1 w-full flex flex-col overflow-auto box-border rounded-lg p-4">
					{props.messages?.map((message) => {
						return (
							<div key={message.id}>
								{message.request !== undefined &&
									message.request !== "" && (
										<div
											className=" w-full flex flex-row"
											key={`${message.id}-req`}
										>
											<div
												className="prose text-start bg-blue-200 m-1 p-2 rounded-xl"
												dangerouslySetInnerHTML={parseMarkdown(
													message.request
												)}
											></div>
											<div className="flex-1"></div>
										</div>
									)}
								{message.response !== undefined &&
									message.response !== "" && (
										<div
											className=" w-full flex flex-row"
											key={`${message.id}-resp`}
										>
											<div className="flex-1"></div>
											<div
												className="prose text-start bg-green-200 p-2 rounded-xl"
												dangerouslySetInnerHTML={parseMarkdown(
													message.response
												)}
											></div>
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
