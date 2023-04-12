import { ClassValue, clsx } from "clsx";
import Cookies from "js-cookie";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
	return twMerge(clsx(inputs));
}
export function formatDate(input: string | number): string {
	const date = new Date(input);
	return date.toLocaleDateString("en-US", {
		month: "long",
		day: "numeric",
		year: "numeric",
	});
}

export function absoluteUrl(path: string) {
	return `${process.env.NEXT_PUBLIC_APP_URL}${path}`;
}

export function isLogin() {
	const token = Cookies.get("token");
	return token != undefined;
}

export function currentTime() {
	return new Date().toISOString().replace(/T/, " ").replace(/\..+/, "");
}

export function parseDateTime(input: string) {
	return new Date(input).toLocaleString();
}
