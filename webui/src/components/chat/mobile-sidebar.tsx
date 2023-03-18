import { Button } from "@/components/ui/button";

import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet";

import { Conversation } from "@/lib/pb/model/conversation.pb";

import { List } from "lucide-react";
import { SideBarComponent } from "./sidebar";

type SideSheetProps = {
	conversations: Conversation[];
	selectedConversation: Conversation | undefined;
	onSelectConversation: (conversation: Conversation | undefined) => void;
	onSetConversations: (conversations: Conversation[]) => void;
};

export function SideSheetComponent({
	conversations,
	selectedConversation,
	onSelectConversation,
	onSetConversations,
}: SideSheetProps) {
	return (
		<Sheet>
			<SheetTrigger asChild>
				<Button variant="ghost" className="px-2">
					<List />
				</Button>
			</SheetTrigger>
			<SheetContent
				position="left"
				size="content"
				className="bg-gray-100 text-gray-800 h-screen"
			>
				<SideBarComponent
					conversations={conversations}
					selectedConversation={selectedConversation}
					onSelectConversation={onSelectConversation}
					onSetConversations={onSetConversations}
				/>
			</SheetContent>
		</Sheet>
	);
}
