import { Button } from "@/components/ui/button";

import { Plus } from "lucide-react";
import { useContext } from "react";

import { ChatSettings } from "./settings";
import { MobileSideBarComponent } from "./mobile-sidebar";
import { ChatContext } from "@/lib/states/chat-context";
import { RoleDialog } from "./role";

export default function MobileChatHeader() {
	const {
		conversations,
		setConversations,
		selectedConversation,
		setSelectedConversation,
		handleCreateConversation,
	} = useContext(ChatContext);

	return (
		<div className="flex flex-row items-center justify-between text-[#fffffe] bg-[#094067] rounded-lg">
			<MobileSideBarComponent />

			<div className="">{selectedConversation?.title}</div>
			<div className="flex flex-row">
				<Button
					variant="ghost"
					size="sm"
					onClick={handleCreateConversation}
				>
					<Plus />
				</Button>
			</div>
		</div>
	);
}
