// WeixinJSBridge.d.ts

interface WeixinJSBridge {
	invoke(
		eventName: string,
		data: Record<string, any>,
		callback: (res: Record<string, any>) => void
	): void;
	on(eventName: string, callback: (res: Record<string, any>) => void): void;
}

declare var WeixinJSBridge: WeixinJSBridge | undefined;
