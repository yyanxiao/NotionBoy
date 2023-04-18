import { Button } from "@/components/ui/button";

import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet";

import { List } from "lucide-react";
import { SideBarComponent } from "./sidebar";

export function MobileSideBarComponent() {
	return (
		<Sheet>
			<SheetTrigger asChild>
				<Button variant="ghost" className="px-2">
					<List />
				</Button>
			</SheetTrigger>
			<SheetContent position="left" size="content" className="h-screen">
				<SideBarComponent />
			</SheetContent>
		</Sheet>
	);
}
