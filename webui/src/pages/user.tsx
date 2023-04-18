import { SiteHeader } from "@/components/site-header";

export default function UserPage() {
	return (
		<>
			<SiteHeader />
			<div className="container flex flex-col items-center mx-auto my-10">
				<div className="prose">
					<h1>User Page</h1>
				</div>
			</div>
		</>
	);
}
