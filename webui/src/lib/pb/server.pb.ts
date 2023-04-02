/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as fm from "./fetch.pb"
import * as GoogleProtobufEmpty from "./google/protobuf/empty.pb"
import * as GoogleRpcStatus from "./google/rpc/status.pb"
import * as Servicev1Common from "./model/common.pb"
import * as Servicev1Conversation from "./model/conversation.pb"
import * as Servicev1Order from "./model/order.pb"
import * as Servicev1Product from "./model/product.pb"
export type CheckStatusResponse = {
  status?: GoogleRpcStatus.Status
}

export class Service {
  static Status(req: GoogleProtobufEmpty.Empty, initReq?: fm.InitReq): Promise<CheckStatusResponse> {
    return fm.fetchReq<GoogleProtobufEmpty.Empty, CheckStatusResponse>(`/v1/status?${fm.renderURLSearchParams(req, [])}`, {...initReq, method: "GET"})
  }
  static GenrateToken(req: Servicev1Common.GenrateTokenRequest, initReq?: fm.InitReq): Promise<Servicev1Common.GenrateTokenResponse> {
    return fm.fetchReq<Servicev1Common.GenrateTokenRequest, Servicev1Common.GenrateTokenResponse>(`/v1/auth/token`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static OAuthProviders(req: Servicev1Common.OAuthURLRequest, initReq?: fm.InitReq): Promise<Servicev1Common.OAuthURLResponse> {
    return fm.fetchReq<Servicev1Common.OAuthURLRequest, Servicev1Common.OAuthURLResponse>(`/v1/auth/providers?${fm.renderURLSearchParams(req, [])}`, {...initReq, method: "GET"})
  }
  static OAuthCallback(req: Servicev1Common.OAuthCallbackRequest, initReq?: fm.InitReq): Promise<Servicev1Common.GenrateTokenResponse> {
    return fm.fetchReq<Servicev1Common.OAuthCallbackRequest, Servicev1Common.GenrateTokenResponse>(`/v1/auth/callback`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static GenerateApiKey(req: GoogleProtobufEmpty.Empty, initReq?: fm.InitReq): Promise<Servicev1Common.GenerateApiKeyResponse> {
    return fm.fetchReq<GoogleProtobufEmpty.Empty, Servicev1Common.GenerateApiKeyResponse>(`/v1/auth/apikey`, {...initReq, method: "POST"})
  }
  static DeleteApiKey(req: GoogleProtobufEmpty.Empty, initReq?: fm.InitReq): Promise<GoogleProtobufEmpty.Empty> {
    return fm.fetchReq<GoogleProtobufEmpty.Empty, GoogleProtobufEmpty.Empty>(`/v1/auth/apikey`, {...initReq, method: "DELETE"})
  }
  static GenerateWechatQRCode(req: GoogleProtobufEmpty.Empty, initReq?: fm.InitReq): Promise<Servicev1Common.GenerateWechatQRCodeResponse> {
    return fm.fetchReq<GoogleProtobufEmpty.Empty, Servicev1Common.GenerateWechatQRCodeResponse>(`/v1/auth/wechat/qrcode`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static CreateConversation(req: Servicev1Conversation.CreateConversationRequest, initReq?: fm.InitReq): Promise<Servicev1Conversation.Conversation> {
    return fm.fetchReq<Servicev1Conversation.CreateConversationRequest, Servicev1Conversation.Conversation>(`/v1/conversations`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static UpdateConversation(req: Servicev1Conversation.UpdateConversationRequest, initReq?: fm.InitReq): Promise<Servicev1Conversation.Conversation> {
    return fm.fetchReq<Servicev1Conversation.UpdateConversationRequest, Servicev1Conversation.Conversation>(`/v1/conversations/${req["id"]}`, {...initReq, method: "PUT", body: JSON.stringify(req, fm.replacer)})
  }
  static GetConversation(req: Servicev1Conversation.GetConversationRequest, initReq?: fm.InitReq): Promise<Servicev1Conversation.Conversation> {
    return fm.fetchReq<Servicev1Conversation.GetConversationRequest, Servicev1Conversation.Conversation>(`/v1/conversations/${req["id"]}?${fm.renderURLSearchParams(req, ["id"])}`, {...initReq, method: "GET"})
  }
  static ListConversations(req: Servicev1Conversation.ListConversationsRequest, initReq?: fm.InitReq): Promise<Servicev1Conversation.ListConversationsResponse> {
    return fm.fetchReq<Servicev1Conversation.ListConversationsRequest, Servicev1Conversation.ListConversationsResponse>(`/v1/conversations?${fm.renderURLSearchParams(req, [])}`, {...initReq, method: "GET"})
  }
  static DeleteConversation(req: Servicev1Conversation.DeleteConversationRequest, initReq?: fm.InitReq): Promise<GoogleProtobufEmpty.Empty> {
    return fm.fetchReq<Servicev1Conversation.DeleteConversationRequest, GoogleProtobufEmpty.Empty>(`/v1/conversations/${req["id"]}`, {...initReq, method: "DELETE"})
  }
  static CreateMessage(req: Servicev1Conversation.CreateMessageRequest, entityNotifier?: fm.NotifyStreamEntityArrival<Servicev1Conversation.Message>, initReq?: fm.InitReq): Promise<void> {
    return fm.fetchStreamingRequest<Servicev1Conversation.CreateMessageRequest, Servicev1Conversation.Message>(`/v1/conversations/${req["conversationId"]}/messages`, entityNotifier, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static GetMessage(req: Servicev1Conversation.GetMessageRequest, initReq?: fm.InitReq): Promise<Servicev1Conversation.Message> {
    return fm.fetchReq<Servicev1Conversation.GetMessageRequest, Servicev1Conversation.Message>(`/v1/conversations/${req["conversationId"]}/messages/${req["id"]}?${fm.renderURLSearchParams(req, ["conversationId", "id"])}`, {...initReq, method: "GET"})
  }
  static ListMessages(req: Servicev1Conversation.ListMessagesRequest, initReq?: fm.InitReq): Promise<Servicev1Conversation.ListMessagesResponse> {
    return fm.fetchReq<Servicev1Conversation.ListMessagesRequest, Servicev1Conversation.ListMessagesResponse>(`/v1/conversations/${req["conversationId"]}/messages?${fm.renderURLSearchParams(req, ["conversationId"])}`, {...initReq, method: "GET"})
  }
  static DeleteMessage(req: Servicev1Conversation.DeleteMessageRequest, initReq?: fm.InitReq): Promise<GoogleProtobufEmpty.Empty> {
    return fm.fetchReq<Servicev1Conversation.DeleteMessageRequest, GoogleProtobufEmpty.Empty>(`/v1/conversations/${req["conversationId"]}/messages/${req["id"]}`, {...initReq, method: "DELETE"})
  }
  static CreateOrder(req: Servicev1Order.CreateOrderRequest, initReq?: fm.InitReq): Promise<Servicev1Order.Order> {
    return fm.fetchReq<Servicev1Order.CreateOrderRequest, Servicev1Order.Order>(`/v1/orders`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static GetOrder(req: Servicev1Order.GetOrderRequest, initReq?: fm.InitReq): Promise<Servicev1Order.Order> {
    return fm.fetchReq<Servicev1Order.GetOrderRequest, Servicev1Order.Order>(`/v1/orders/${req["id"]}?${fm.renderURLSearchParams(req, ["id"])}`, {...initReq, method: "GET"})
  }
  static ListOrders(req: Servicev1Order.ListOrdersRequest, initReq?: fm.InitReq): Promise<Servicev1Order.ListOrdersResponse> {
    return fm.fetchReq<Servicev1Order.ListOrdersRequest, Servicev1Order.ListOrdersResponse>(`/v1/orders?${fm.renderURLSearchParams(req, [])}`, {...initReq, method: "GET"})
  }
  static DeleteOrder(req: Servicev1Order.DeleteOrderRequest, initReq?: fm.InitReq): Promise<GoogleProtobufEmpty.Empty> {
    return fm.fetchReq<Servicev1Order.DeleteOrderRequest, GoogleProtobufEmpty.Empty>(`/v1/orders/${req["id"]}`, {...initReq, method: "DELETE"})
  }
  static UpdateOrder(req: Servicev1Order.UpdateOrderRequest, initReq?: fm.InitReq): Promise<Servicev1Order.Order> {
    return fm.fetchReq<Servicev1Order.UpdateOrderRequest, Servicev1Order.Order>(`/v1/orders/${req["id"]}`, {...initReq, method: "PATCH", body: JSON.stringify(req, fm.replacer)})
  }
  static PayOrder(req: Servicev1Order.PayOrderRequest, initReq?: fm.InitReq): Promise<Servicev1Order.PayOrderResponse> {
    return fm.fetchReq<Servicev1Order.PayOrderRequest, Servicev1Order.PayOrderResponse>(`/v1/orders/${req["id"]}/pay`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static CreateProduct(req: Servicev1Product.CreateProductRequest, initReq?: fm.InitReq): Promise<Servicev1Product.Product> {
    return fm.fetchReq<Servicev1Product.CreateProductRequest, Servicev1Product.Product>(`/internal/v1/products`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static GetProduct(req: Servicev1Product.GetProductRequest, initReq?: fm.InitReq): Promise<Servicev1Product.Product> {
    return fm.fetchReq<Servicev1Product.GetProductRequest, Servicev1Product.Product>(`/internal/v1/products/${req["id"]}?${fm.renderURLSearchParams(req, ["id"])}`, {...initReq, method: "GET"})
  }
  static ListProducts(req: Servicev1Product.ListProductsRequest, initReq?: fm.InitReq): Promise<Servicev1Product.ListProductsResponse> {
    return fm.fetchReq<Servicev1Product.ListProductsRequest, Servicev1Product.ListProductsResponse>(`/internal/v1/products?${fm.renderURLSearchParams(req, [])}`, {...initReq, method: "GET"})
  }
  static DeleteProduct(req: Servicev1Product.DeleteProductRequest, initReq?: fm.InitReq): Promise<GoogleProtobufEmpty.Empty> {
    return fm.fetchReq<Servicev1Product.DeleteProductRequest, GoogleProtobufEmpty.Empty>(`/internal/v1/products/${req["id"]}`, {...initReq, method: "DELETE"})
  }
  static UpdateProduct(req: Servicev1Product.UpdateProductRequest, initReq?: fm.InitReq): Promise<Servicev1Product.Product> {
    return fm.fetchReq<Servicev1Product.UpdateProductRequest, Servicev1Product.Product>(`/internal/v1/products/${req["id"]}`, {...initReq, method: "PATCH", body: JSON.stringify(req, fm.replacer)})
  }
}
