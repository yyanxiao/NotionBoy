import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from "@/components/ui/select";
import { Loader2, Send } from "lucide-react";
import { useState } from "react";
import { Button } from "../ui/button";
import { Label } from "../ui/label";
import { Textarea } from "../ui/textarea";

type ChatInputBoxProps = {
	onSendMessage: (message: string, model: string) => void;
	isLoading: boolean;
};

interface Model {
	name: string;
	value: string;
}

const models: Model[] = [
	{
		name: "GPT-3Dot5Turbo",
		value: "gpt-3.5-turbo",
	},
	{
		name: "GPT4",
		value: "gpt-4",
	},
	// {
	// 	name: "GPT432K",
	// 	value: "gpt-4-32k",
	// },
];

export function ChatInputBox({ onSendMessage, isLoading }: ChatInputBoxProps) {
	const [inputValue, setInputValue] = useState<string>("");
	const [isSending, setIsSending] = useState(false);
	const [model, setModel] = useState<string>("gpt-3.5-turbo");
	const handleMessageSend = () => {
		setIsSending(true);
		onSendMessage(inputValue, model);
		setInputValue("");
		setIsSending(false);
	};

	const isEmptyInput = () => inputValue.trim() === "";

	return (
		<div className="relative flex flex-col items-center m-2">
			<div className="flex flex-row items-center space-x-2">
				<Label>
					<strong className="text-gray-700">Model:</strong>
				</Label>
				<Select value={model} onValueChange={setModel}>
					<SelectTrigger className="w-[180px]">
						<SelectValue placeholder="Model" />
					</SelectTrigger>
					<SelectContent>
						{models.map((model) => (
							<SelectItem key={model.value} value={model.value}>
								{model.name}
							</SelectItem>
						))}
					</SelectContent>
				</Select>
			</div>
			<Textarea
				placeholder="Type a message and send with Shift+Enter ..."
				className="w-full h-10 disabled:opacity-50 md:h-20"
				onChange={(e) => setInputValue(e.target.value)}
				value={inputValue}
				disabled={isLoading}
				rows={1}
				onKeyDown={(e) => {
					// using shift + enter to send a message
					// using enter to create a new line
					if (e.key === "Enter" && e.shiftKey && !isEmptyInput()) {
						e.preventDefault();
						handleMessageSend();
					}
				}}
			/>

			<Button
				variant="ghost"
				disabled={isLoading || isEmptyInput()}
				className="absolute bottom-0 right-0"
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
