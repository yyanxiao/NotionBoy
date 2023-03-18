import { Loader2, Send } from "lucide-react";
import { useState } from "react";
import { Button } from "../ui/button";
import { Textarea } from "../ui/textarea";

type ChatInputBoxProps = {
	onSendMessage: (message: string) => void;
	isLoading: boolean;
};

export function ChatInputBox({ onSendMessage, isLoading }: ChatInputBoxProps) {
	const [inputValue, setInputValue] = useState<string>("");
	const [isSending, setIsSending] = useState(false);
	const handleMessageSend = () => {
		setIsSending(true);
		onSendMessage(inputValue);
		setInputValue("");
		setIsSending(false);
	};

	const isEmptyInput = () => inputValue.trim() === "";

	return (
		<div className="relative m-2">
			<Textarea
				placeholder="Type a message..."
				className="w-full disabled:opacity-50 h-10 md:h-20"
				onChange={(e) => setInputValue(e.target.value)}
				value={inputValue}
				disabled={isLoading}
				rows={1}
				onKeyDown={(e) => {
					if (e.key === "Enter" && !isEmptyInput()) {
						e.preventDefault();
						handleMessageSend();
					}
				}}
			/>

			<Button
				variant="ghost"
				disabled={isLoading || isEmptyInput()}
				className="absolute right-0 bottom-0"
				onClick={handleMessageSend}
			>
				{isLoading || isSending ? (
					<Loader2 className="animate-spin" />
				) : (
					<Send />
				)}
			</Button>
		</div>
	);
}
