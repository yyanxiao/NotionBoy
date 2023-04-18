/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as Servicev1Product from "./product.pb"
export type Order = {
  id?: string
  createdAt?: string
  updatedAt?: string
  deleted?: string
  uuid?: string
  userId?: string
  price?: number
  status?: string
  note?: string
  paymentMethod?: string
  product?: Servicev1Product.Product
}

export type CreateOrderRequest = {
  productId?: string
  note?: string
}

export type ListOrdersRequest = {
  status?: string
  limit?: number
  offset?: number
}

export type ListOrdersResponse = {
  orders?: Order[]
}

export type GetOrderRequest = {
  id?: string
  isCheckPayment?: boolean
}

export type DeleteOrderRequest = {
  id?: string
}

export type UpdateOrderRequest = {
  id?: string
  note?: string
}

export type PayOrderRequest = {
  id?: string
  paymentMethod?: string
}

export type PayOrderConfig = {
  timestamp?: string
  nonceStr?: string
  prePayId?: string
  signType?: string
  package?: string
  paySign?: string
  appId?: string
}

export type PayOrderResponse = {
  status?: string
  qrcode?: string
  config?: PayOrderConfig
}
