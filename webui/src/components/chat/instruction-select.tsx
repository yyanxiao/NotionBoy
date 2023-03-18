import {
	DefaultInstruction,
	Instruction,
	InstructionList,
} from "@/config/prompts";

import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from "@/components/ui/select";

import { useEffect, useState } from "react";

type ConversationListProps = {
	instruction: Instruction;
	setInstruction: (instruction: Instruction) => void;
};

export function InstructionSelectComponent({
	instruction,
	setInstruction,
}: ConversationListProps) {
	const [selectedList, setSelectedList] = useState<Instruction[]>([
		DefaultInstruction,
	]);
	const [title, setTitle] = useState<string>(DefaultInstruction.title);

	const [category, setCategory] = useState<string>("default");

	useEffect(() => {
		const list = InstructionList.find((item) => item.key === category);
		if (list) {
			setSelectedList(list.data);
		}
	}, [category]);

	useEffect(() => {
		const instruction = selectedList.find((item) => item.title === title);
		if (instruction) {
			setInstruction(instruction);
		}
	}, [title]);

	return (
		<div className="flex flex-col">
			<div className="flex flex-row w-full">
				<Select onValueChange={setCategory}>
					<SelectTrigger className="w-full">
						<SelectValue placeholder="Instruction categories..." />
					</SelectTrigger>
					<SelectContent>
						<SelectItem key="disabled" value="" disabled>
							Select Instruction Category
						</SelectItem>
						{InstructionList.map((item) => (
							<SelectItem key={item.key} value={item.key}>
								{item.key}
							</SelectItem>
						))}
					</SelectContent>
				</Select>
			</div>

			<div className="flex flex-row">
				<Select onValueChange={setTitle}>
					<SelectTrigger className="w-full">
						<SelectValue placeholder="Select Instruction" />
					</SelectTrigger>
					<SelectContent>
						<SelectItem key="disabled" value="" disabled>
							Select Instruction
						</SelectItem>
						{selectedList.map((item) => (
							<SelectItem key={item.title} value={item.title}>
								{item.title}
							</SelectItem>
						))}
					</SelectContent>
				</Select>
			</div>
		</div>
	);
}
