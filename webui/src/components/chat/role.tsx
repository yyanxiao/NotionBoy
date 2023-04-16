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
import { Switch } from "@/components/ui/switch";
import { useToast } from "@/hooks/use-toast";
import {
	CreatePromptRequest,
	DeletePromptRequest,
	Prompt,
	UpdatePromptRequest,
} from "@/lib/pb/model/common.pb";
import { CreateConversationRequest } from "@/lib/pb/model/conversation.pb";
import { Service } from "@/lib/pb/server.pb";
import { ChatContext } from "@/lib/states/chat-context";

import { Edit2, Trash2 } from "lucide-react";
import { useContext, useEffect, useState } from "react";
import { v4 as uuidv4 } from "uuid";
import { Separator } from "../ui/separator";
import { Textarea } from "../ui/textarea";

const defaultPrompt = {
	id: uuidv4(),
	act: "ChatGPT",
	prompt: "You are ChatGPT, a large language model trained by OpenAI. Follow the user's instructions carefully. Respond using markdown.",
	isCustom: false,
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
	const [customPrompts, setCustomPrompts] = useState<Prompt[]>([]);
	const [filteredPrompts, setFilteredPrompts] = useState<Prompt[]>([]);
	const [selectedPrompt, setSelectedPrompt] = useState<Prompt>(defaultPrompt);
	const [searchValue, setSearchValue] = useState<string>("");

	const [roleName, setRoleName] = useState<string>("");
	const [roleDescription, setRoleDescription] = useState<string>("");

	const [isCustomRole, setIsCustomRole] = useState<boolean>(false);

	const [isOpen, setIsOpen] = useState(false);
	const [isOpenEditRole, setIsOpenEditRole] = useState(false);

	useEffect(() => {
		// get custom prompts
		Service.ListPrompts({
			isCustom: true,
		})
			.then((res) => {
				if (res?.prompts) {
					setCustomPrompts(res.prompts);
				}
			})
			.catch((err) => {
				toast({
					variant: "destructive",
					title: "获取自定义角色失败",
					description: JSON.stringify(err),
				});
			});

		// get default prompts
		const localPrompts = localStorage.getItem("prompts");
		if (localPrompts) {
			const prompts = JSON.parse(localPrompts);
			setPrompts(prompts);
		} else {
			Service.ListPrompts({
				isCustom: false,
			})
				.then((res) => {
					const prompts = [defaultPrompt, ...(res?.prompts || [])];
					setPrompts(prompts);
					localStorage.setItem("prompts", JSON.stringify(prompts));
				})
				.catch((err) => {
					toast({
						variant: "destructive",
						title: "获取默认角色失败",
						description: JSON.stringify(err),
					});
				});
		}
	}, []);

	useEffect(() => {
		let allPrompts = [...prompts, ...customPrompts];
		if (isCustomRole) {
			allPrompts = customPrompts;
		}
		setFilteredPrompts(
			allPrompts.filter((p) => {
				const regExp = new RegExp(searchValue, "i");
				return regExp.test(p.act as string);
			})
		);
	}, [prompts, customPrompts, searchValue, isCustomRole]);

	const handleCreateConversationWithRole = () => {
		const createConversationRequest = {
			instruction: selectedPrompt?.prompt,
			title: selectedPrompt?.act,
		} as CreateConversationRequest;
		Service.CreateConversation(createConversationRequest)
			.then((conversation) => {
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

	const handleCreateCustomRole = (act: string, prompt: string) => {
		const createPromptRequest = {
			act: act,
			prompt: prompt,
		} as CreatePromptRequest;
		Service.CreatePrompt(createPromptRequest)
			.then((prompt) => {
				setCustomPrompts([...customPrompts, prompt]);
				setIsOpenEditRole(false);
				toast({
					variant: "default",
					title: "success",
					description: "新建 role 成功",
				});
			})
			.catch((error) => {
				toast({
					variant: "destructive",
					title: "Error",
					description: `Failed to create custom prompt, ${JSON.stringify(
						error
					)}`,
				});
			});
	};

	const handleUpdateCustomRole = (
		id: string,
		act: string,
		prompt: string
	) => {
		const updatePromptRequest = {
			id: id,
			act: act,
			prompt: prompt,
		} as UpdatePromptRequest;
		Service.UpdatePrompt(updatePromptRequest)
			.then((prompt) => {
				setCustomPrompts(
					customPrompts.map((p) => {
						if (p.id === prompt.id) {
							return prompt;
						}
						return p;
					})
				);
				setIsOpenEditRole(false);
				toast({
					variant: "default",
					title: "success",
					description: "更新 role 成功",
				});
			})
			.catch((error) => {
				toast({
					variant: "destructive",
					title: "Error",
					description: `Failed to update custom prompt, ${JSON.stringify(
						error
					)}`,
				});
			});
	};

	const handleDeleteCustomRole = (id: string) => {
		Service.DeletePrompt({ id } as DeletePromptRequest)
			.then(() => {
				setCustomPrompts(customPrompts.filter((p) => p.id !== id));
				toast({
					variant: "default",
					title: "success",
					description: "删除 role 成功",
				});
			})
			.catch((error) => {
				toast({
					variant: "destructive",
					title: "Error",
					description: `Failed to delete custom prompt, ${JSON.stringify(
						error
					)}`,
				});
			});
	};

	const editRoleDetail = () => {
		return (
			<div className="flex flex-col space-y-2">
				<Separator />
				<div className="grid w-full  items-center gap-1.5">
					<Label>角色名称</Label>
					<Input
						type="text"
						id="role-name"
						placeholder="Role name"
						value={roleName}
						onChange={(e) => setRoleName(e.target.value)}
					/>
				</div>

				<div className="grid w-full items-center gap-1.5">
					<Label>角色定义</Label>
					<Textarea
						id="role-prompt"
						placeholder="Role prompt"
						className="h-24"
						value={roleDescription}
						onChange={(e) => setRoleDescription(e.target.value)}
					/>
				</div>
			</div>
		);
	};

	const editCustomRoleButton = (role: Prompt | undefined) => {
		if (role) {
			return (
				<Button
					type="button"
					onClick={() => {
						handleUpdateCustomRole(
							role?.id as string,
							roleName,
							roleDescription
						);
					}}
				>
					确认
				</Button>
			);
		}
		return (
			<Button
				type="button"
				onClick={() => {
					handleCreateCustomRole(roleName, roleDescription);
				}}
			>
				确认
			</Button>
		);
	};

	const editCustomRole = (role: Prompt | undefined) => {
		return (
			<Dialog
				open={isOpenEditRole}
				onOpenChange={() => {
					setIsOpenEditRole(!isOpenEditRole);
				}}
			>
				<DialogTrigger
					onClick={() => {
						if (role) {
							setRoleName(role.act as string);
							setRoleDescription(role.prompt as string);
						} else {
							setRoleName("");
							setRoleDescription("");
						}
					}}
				>
					{role ? (
						<div className="absolute top-0 right-0 inline-flex items-center justify-center px-2 bg-transparent rounded-md h-9 hover:bg-slate-100 dark:hover:bg-slate-800 dark:text-slate-100 dark:hover:text-slate-100">
							<Edit2 size={16} />
						</div>
					) : (
						<div className="inline-flex items-center justify-center px-2 text-white rounded-md h-9 bg-slate-900 hover:bg-slate-700 dark:bg-slate-50 dark:text-slate-900">
							新建角色
						</div>
					)}
				</DialogTrigger>
				<DialogContent>
					<DialogHeader>
						<DialogTitle>编辑角色 🎭</DialogTitle>
						<DialogDescription>
							请自定义角色的名称和定义
						</DialogDescription>
						{editRoleDetail()}
					</DialogHeader>
					<DialogFooter>{editCustomRoleButton(role)}</DialogFooter>
				</DialogContent>
			</Dialog>
		);
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
				<div className="grid items-center grid-cols-2 gap-2 justify-items-center">
					<div className="flex items-center space-x-2">
						<Label htmlFor="myroles">我的角色</Label>
						<Switch
							id="myroles"
							onCheckedChange={() =>
								setIsCustomRole(!isCustomRole)
							}
						/>
					</div>
					<div className="inline-flex items-center justify-center px-2 rounded-md h-9 hover:bg-white">
						{editCustomRole(undefined)}
					</div>
				</div>
				<Input
					type="text"
					placeholder="Search Role"
					className=""
					value={searchValue}
					onChange={(e) => {
						setSearchValue(e.target.value);
					}}
				/>
				<div className="grid w-full grid-cols-1 gap-2 p-1 overflow-auto md:grid-cols-2">
					{filteredPrompts.map((p) => (
						<div
							key={p.prompt}
							className={`relative h-20 rounded-lg ${
								selectedPrompt.prompt == p.prompt
									? "bg-[#3da9fc]"
									: ""
							}`}
						>
							<Button
								variant={"outline"}
								size="lg"
								className="w-full h-20"
								value={p.act}
								onClick={(e) => {
									setSelectedPrompt(p);
								}}
							>
								{p.act}
							</Button>
							{p.isCustom &&
								selectedPrompt.id == p.id &&
								editCustomRole(p)}
							{p.isCustom && selectedPrompt.id == p.id && (
								<Button
									variant="ghost"
									size="sm"
									className="absolute top-0 right-8"
									onClick={() =>
										handleDeleteCustomRole(p.id as string)
									}
								>
									<Trash2 size={16} />
								</Button>
							)}
						</div>
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
							handleCreateConversationWithRole();
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
