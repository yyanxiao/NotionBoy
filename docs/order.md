# 订单支付流程

1. 用户前端页面点击`保存订单`按钮，调用`/v1/orders`接口，将订单信息保存到数据库中，返回订单号。
2. 用户前端页面点击支付按钮，调用`/v1/orders/{order_id}/pay`接口，将订单号传递给后端
3. 后端收到支付的请求
   1. 更改订单状态为`Paying`
   2. 调用微信支付的 prepay 接口，获取到支付二维码返回给前端
4. 前端收到支付二维码，展示给用户，用户扫码支付
5. 前端同时轮询后端的 `/v1/orders/{id}` 接口并且带上 `is_check_payment` 参数
6. 后端接收到 is_check_payment 的请求后
   1. 查询订单状态，如果是 Paying，则调用微信支付的接口，查询支付是否完成
      1. 如果支付完成，更新 Order 状态为 Paid，异步去执行 Process Order 的操作
      2. 如果支付为完成，直接正常返回前端
   2. 如果 Order 状态是 Paid，异步执行 Process Order 的操作
      1. 更新 Order 状态为 Processing
      2. 更新用户购买的服务
      3. 更新完成之后更改 Order 的状态为 Completed
7. 如果前端轮询后端，得到 Completed 状态的 Order，则展示给用户，订单已完成
8. 跳转用户页面，查看以购买的内容
