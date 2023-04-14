import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from "@/components/ui/select";
import { Info, Loader2, Send, Settings2 } from "lucide-react";
import { useState } from "react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
	Popover,
	PopoverContent,
	PopoverTrigger,
} from "@/components/ui/popover";
import {
	Tooltip,
	TooltipContent,
	TooltipProvider,
	TooltipTrigger,
} from "@/components/ui/tooltip";
import { Message } from "@/lib/pb/model/conversation.pb";
import * as Form from "@radix-ui/react-form";
import { Slider } from "../ui/slider";
import { Textarea } from "../ui/textarea";
import { RoleDialog } from "./role";
type ChatInputBoxProps = {
	onSendMessage: (
		message: string,
		model: string,
		temperature: number,
		maxTokens: number
	) => void;
	isLoading: boolean;
	messages: Message[] | undefined;
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

export function ChatInputBox({
	onSendMessage,
	isLoading,
	messages,
}: ChatInputBoxProps) {
	const [inputValue, setInputValue] = useState<string>("");
	const [isSending, setIsSending] = useState<boolean>(false);
	const [model, setModel] = useState<string>("gpt-3.5-turbo");
	const [temperature, setTemperature] = useState<number>(1);
	const [maxTokens, setMaxTokens] = useState<number>(1000);
	const [isOpen, setIsOpen] = useState<boolean>(false);

	const handleMessageSend = () => {
		setIsSending(true);
		onSendMessage(inputValue, model, temperature, maxTokens);
		setInputValue("");
		setIsSending(false);
	};

	const isEmptyInput = () => inputValue.trim() === "";

	const tooltip = (title: string, desc: string) => {
		return (
			<div className="flex flex-row items-center space-x-1">
				<Label>{title}</Label>
				<TooltipProvider>
					<Tooltip defaultOpen={false}>
						<TooltipTrigger type="button">
							<Info size="14" />
						</TooltipTrigger>
						<TooltipContent className="w-48 prose">
							<p>{desc}</p>
						</TooltipContent>
					</Tooltip>
				</TooltipProvider>
			</div>
		);
	};

	const settings = () => {
		return (
			<Popover open={isOpen} onOpenChange={setIsOpen}>
				<PopoverTrigger asChild>
					<Button variant="outline">
						<Settings2 />
						<span className="sr-only">Open popover</span>
					</Button>
				</PopoverTrigger>
				<PopoverContent className="w-96">
					<Form.Root className="space-y-4">
						{selectModel()}
						{selectTemperature()}
						{selectMaxTokens()}
						<Button
							type="button"
							className="w-full"
							onClick={() => {
								if (maxTokens < 100 || maxTokens > 4000) {
									alert("MaxTokens 范围是 100 ~ 4000");
									return;
								}
								setIsOpen(false);
							}}
						>
							确定
						</Button>
					</Form.Root>
				</PopoverContent>
			</Popover>
		);
	};

	const selectModel = () => {
		return (
			<Form.Field
				name="Model"
				className="grid items-center grid-cols-3 gap-4"
			>
				<Form.Label>
					{tooltip(
						"Model:",
						"GPT-4 更加智能, 但是消耗的Token 是GPT-3.5 的 30 倍"
					)}
				</Form.Label>
				<Select value={model} onValueChange={setModel}>
					<SelectTrigger className="col-span-2">
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
			</Form.Field>
		);
	};

	const selectTemperature = () => {
		return (
			<Form.Field
				name="Temperature"
				className="grid items-center grid-cols-3 gap-4"
			>
				<Form.Label>
					{tooltip(
						`创造力: ${temperature / 2}`,
						`可选范围是 0 ~ 1, 默认为 0.5. 值越高, 生成的文本越有趣, 但是也越不可靠`
					)}
				</Form.Label>
				<Slider
					// defaultValue={[1]}
					value={[temperature]}
					max={2.0}
					step={0.2}
					onValueChange={(values) => setTemperature(values[0])}
					className="col-span-2"
				/>
			</Form.Field>
		);
	};

	const selectMaxTokens = () => {
		return (
			<Form.Field
				name="MaxTokens"
				className="grid items-center grid-cols-3 gap-4"
			>
				<Form.Label>
					{tooltip(
						"MaxTokens:",
						"ChatGPT 返回的最大 Token 数量, 默认为 1000, 可选范围是 100 ~ 4096"
					)}
				</Form.Label>
				<Form.Control type="number" asChild>
					<Input
						type="number"
						id="tokens"
						min={100}
						max={4096}
						placeholder="Max Tokens"
						className="col-span-2"
						value={maxTokens}
						onChange={(e) => setMaxTokens(parseInt(e.target.value))}
					/>
				</Form.Control>
			</Form.Field>
		);
	};

	return (
		<div className="relative flex flex-col items-center m-2 space-y-2">
			<div className="flex flex-row items-center justify-center space-x-2">
				{messages && messages.length == 0 && <RoleDialog />}
				{settings()}
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
