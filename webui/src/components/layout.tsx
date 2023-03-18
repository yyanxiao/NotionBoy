import { Toaster } from "./ui/toaster";

export const Layout = ({ children }: { children: React.ReactNode }) => {
	return (
		<div className="min-h-screen max-w-8xl sm:px-6 md:px-8">
			<Toaster />
			{children}
		</div>
	);
};
