import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
	Sheet,
	SheetContent,
	SheetDescription,
	SheetHeader,
	SheetTitle,
	SheetTrigger,
} from "@/components/ui/sheet";
import { DefaultInstruction, Instruction } from "@/config/prompts";
import { Conversation } from "@/lib/pb/model/conversation.pb";
import { Settings2 } from "lucide-react";
import { useEffect, useState } from "react";
import { Textarea } from "../ui/textarea";

import { useToast } from "@/hooks/use-toast";
import { Service } from "@/lib/pb/server.pb";
import { InstructionSelectComponent } from "./instruction-select";
type ChatSettingsProps = {
	conversations: Conversation[];
	selectedConversation: Conversation | undefined;
	onSelectConversation: (conversation: Conversation | undefined) => void;
	onSetConversations: (conversations: Conversation[]) => void;
};

export function ChatSettings({
	conversations,
	selectedConversation,
	onSelectConversation,
	onSetConversations,
}: ChatSettingsProps) {
	const { toast } = useToast();
	const [instruction, setInstruction] = useState<Instruction>({
		title: DefaultInstruction.title,
		instruction: DefaultInstruction.instruction,
		instructioncn: DefaultInstruction.instructioncn,
	} as Instruction);
	const [title, setTitle] = useState<string>(instruction.title);
	const [instructionStr, setInstructionStr] = useState<string>(
		instruction.instruction
	);

	useEffect(() => {
		setTitle(instruction.title);
		setInstructionStr(instruction.instruction);
	}, [instruction]);

	const handleUpdateConversation = () => {
		const conversation = {
			...selectedConversation,
			instruction: instruction.instruction,
			title: instruction.title,
		} as Conversation;
		Service.UpdateConversation(conversation)
			.then((resp) => {
				onSelectConversation(conversation);
				onSetConversations(
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
		<Sheet>
			<SheetTrigger asChild>
				<Button variant="ghost">
					<Settings2 />
				</Button>
			</SheetTrigger>
			<SheetContent position="right" size="xl" className="bg-gray-100">
				<SheetHeader>
					<SheetTitle>Change Settings</SheetTitle>
					<SheetDescription>
						Change settings for the conversation.
					</SheetDescription>
				</SheetHeader>
				<div className="flex flex-col rounded-lg p-4">
					<InstructionSelectComponent
						instruction={instruction}
						setInstruction={setInstruction}
					/>

					<div className="flex flex-row w-full items-center my-2">
						<Label className="w-20">Title:</Label>
						<Input
							type="text"
							placeholder="Title"
							value={title}
							onChange={(e) => setTitle(e.target.value)}
						/>
					</div>
					<div className="flex flex-row w-full items-center my-2">
						<Label className="w-20">Instraction:</Label>
						<Textarea
							className="h-36"
							placeholder="Instraction"
							value={instructionStr}
							onChange={(e) => setInstructionStr(e.target.value)}
						/>
					</div>
					<Button type="submit" onClick={handleUpdateConversation}>
						Save
					</Button>
				</div>
			</SheetContent>
		</Sheet>
	);
}
