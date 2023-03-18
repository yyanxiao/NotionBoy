import { SiteHeader } from "@/components/site-header";

export default function Home() {
	return (
		<div>
			<SiteHeader />
			<div className="flex-grow container mx-auto prose p-8">
				<h1>Welcome to NotionBoy</h1>
			</div>
		</div>
	);
}
