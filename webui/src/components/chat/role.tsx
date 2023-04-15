import { Button } from "@/components/ui/button";
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useToast } from "@/hooks/use-toast";

import { Prompt } from "@/lib/pb/model/common.pb";
import { Conversation } from "@/lib/pb/model/conversation.pb";
import { Service } from "@/lib/pb/server.pb";
import { ChatContext } from "@/lib/states/chat-context";
import { useContext, useEffect, useState } from "react";

const defaultPrompt = {
	act: "ChatGPT",
	prompt: "You are ChatGPT, a large language model trained by OpenAI. Follow the user's instructions carefully. Respond using markdown.",
} as Prompt;

export function RoleDialog() {
	const {
		conversations,
		setConversations,
		selectedConversation,
		setSelectedConversation,
		handleCreateConversation,
	} = useContext(ChatContext);

	const { toast } = useToast();

	const [prompts, setPrompts] = useState<Prompt[]>([]);
	const [filteredPrompts, setFilteredPrompts] = useState<Prompt[]>([]);
	const [selectedPrompt, setSelectedPrompt] = useState<Prompt>(defaultPrompt);
	const [searchValue, setSearchValue] = useState<string>("");

	const [isOpen, setIsOpen] = useState(false);

	useEffect(() => {
		if (prompts.length == 0) {
			const localPrompts = localStorage.getItem("prompts");
			if (localPrompts) {
				const prompts = JSON.parse(localPrompts);
				setPrompts(prompts);
			} else {
				Service.ListPrompts({})
					.then((res) => {
						const prompts = [
							defaultPrompt,
							...(res?.prompts || []),
						];
						setPrompts(prompts);
						localStorage.setItem(
							"prompts",
							JSON.stringify(prompts)
						);
					})
					.catch((err) => {
						console.error(err);
					});
			}
		} else {
			setFilteredPrompts(prompts);
		}
	}, [prompts]);

	useEffect(() => {
		setFilteredPrompts(
			prompts.filter((p) => {
				const regExp = new RegExp(searchValue, "i");
				return regExp.test(p.act as string);
			})
		);
	}, [searchValue]);

	const handleUpdateConversation = () => {
		const conversation = {
			...selectedConversation,
			instruction: selectedPrompt?.prompt,
			title: selectedPrompt?.act,
		} as Conversation;
		Service.UpdateConversation(conversation)
			.then(() => {
				setSelectedConversation(conversation);
				setConversations(
					conversations.map((c) => {
						if (c.id === conversation.id) {
							return conversation;
						}
						return c;
					})
				);
			})
			.catch((error) => {
				toast({
					variant: "destructive",
					title: "Error",
					description: `Failed to update conversation, ${JSON.stringify(
						error
					)}`,
				});
			});
	};

	return (
		<Dialog open={isOpen} onOpenChange={setIsOpen}>
			<DialogTrigger asChild>
				<Button
					variant="outline"
					size="sm"
					className="flex flex-row w-full space-x-1"
				>
					<span>🎭 选择角色</span>
				</Button>
			</DialogTrigger>
			<DialogContent className="w-3/4 h-3/4">
				<DialogHeader>
					<DialogTitle>角色 🎭 扮演</DialogTitle>
					<DialogDescription>
						你期望 ChatGPT 以哪个角色来回复你的消息？
					</DialogDescription>
				</DialogHeader>
				<div>
					<Input
						type="text"
						placeholder="Search Role"
						value={searchValue}
						onChange={(e) => {
							setSearchValue(e.target.value);
						}}
					/>
				</div>
				<div className="grid w-full grid-cols-1 gap-2 p-1 overflow-auto md:grid-cols-2">
					{filteredPrompts.map((p) => (
						<Button
							variant={"outline"}
							size="lg"
							className={`h-20 text-start ${
								selectedPrompt.prompt == p.prompt
									? "bg-[#3da9fc]"
									: ""
							} `}
							key={p.prompt}
							value={p.act}
							onClick={(e) => {
								setSelectedPrompt(p);
							}}
						>
							{p.act}
						</Button>
					))}
				</div>
				{selectedPrompt && (
					<div className="h-24 p-2 overflow-auto prose border-2 border-gray-200">
						<Label>{selectedPrompt?.prompt}</Label>
					</div>
				)}

				<DialogFooter>
					<Button
						type="submit"
						onClick={() => {
							handleUpdateConversation();
							setIsOpen(false);
						}}
					>
						确认
					</Button>
				</DialogFooter>
			</DialogContent>
		</Dialog>
	);
}
