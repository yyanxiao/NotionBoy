// WeChatPayment.tsx

import { PayOrderConfig } from "@/lib/pb/model/order.pb";
import React, { useEffect } from "react";

interface WeChatPaymentProps {
	cfg: PayOrderConfig | undefined;
	onSuccess: () => void;
}

const WeChatPayment = ({ cfg, onSuccess }: WeChatPaymentProps) => {
	useEffect(() => {
		const onBridgeReady = () => {
			if (!cfg) return;
			WeixinJSBridge?.invoke(
				"getBrandWCPayRequest",
				{
					appId: cfg.appId, //公众号ID，由商户传入
					timeStamp: cfg.timestamp, //时间戳，自1970年以来的秒数
					nonceStr: cfg.nonceStr, //随机串
					package: cfg.package,
					signType: cfg.signType, //微信签名方式：
					paySign: cfg.paySign, //微信签名
				},
				(res) => {
					if (res.err_msg === "get_brand_wcpay_request:ok") {
						// 使用以上方式判断前端返回,微信团队郑重提示：
						// res.err_msg 将在用户支付成功后返回 ok，但并不保证它绝对可靠。
						onSuccess();
					}
				}
			);
		};

		if (typeof WeixinJSBridge === "undefined") {
			document.addEventListener(
				"WeixinJSBridgeReady",
				onBridgeReady,
				false
			);
		} else {
			onBridgeReady();
		}
	}, []);

	return <div></div>;
};

export default WeChatPayment;
