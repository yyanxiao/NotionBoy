import { Button } from "@/components/ui/button";

import { Plus } from "lucide-react";
import { useContext } from "react";

import { ChatSettings } from "./settings";
import { SideSheetComponent } from "./mobile-sidebar";
import { ChatContext } from "@/lib/states/chat-context";

export default function MobileChatHeader() {
	const {
		conversations,
		setConversations,
		selectedConversation,
		setSelectedConversation,
		handleCreateConversation,
	} = useContext(ChatContext);

	return (
		<div className="flex flex-row items-center justify-between mx-2">
			<SideSheetComponent />

			<div className="">{selectedConversation?.title}</div>
			<div className="flex flex-row">
				<Button
					variant="ghost"
					size="sm"
					onClick={handleCreateConversation}
				>
					<Plus />
				</Button>
				<ChatSettings />
			</div>
		</div>
	);
}
