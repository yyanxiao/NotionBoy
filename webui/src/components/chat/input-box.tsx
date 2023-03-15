import { Loader2, Send } from "lucide-react";
import { useState } from "react";
import { Button } from "../ui/button";

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
		// <div className="fixed bottom-8  h-fit container mx-auto">
		<div className="h-fit container mx-auto">
			<div className="relative">
				<textarea
					placeholder="Type a message..."
					className="w-full p-2 border-sky-200 border-2 disabled:opacity-50"
					onChange={(e) => setInputValue(e.target.value)}
					value={inputValue}
					disabled={isLoading}
					rows={2}
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
					className="absolute right-1 mt-4 p-1 disabled:opacity-50"
					onClick={handleMessageSend}
				>
					{isLoading || isSending ? (
						<Loader2 className="animate-spin" />
					) : (
						<Send />
					)}
				</Button>
			</div>
		</div>
	);
}
