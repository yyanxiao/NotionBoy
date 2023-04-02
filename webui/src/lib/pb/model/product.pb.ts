/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/
export type Product = {
  id?: string
  createdAt?: string
  updatedAt?: string
  deleted?: string
  uuid?: string
  name?: string
  price?: number
  token?: string
  storage?: string
  description?: string
}

export type CreateProductRequest = {
  name?: string
  price?: number
  token?: string
  storage?: string
  description?: string
}

export type ListProductsRequest = {
  limit?: number
  offset?: number
}

export type ListProductsResponse = {
  products?: Product[]
}

export type GetProductRequest = {
  id?: string
}

export type DeleteProductRequest = {
  id?: string
}

export type UpdateProductRequest = {
  id?: string
  name?: string
  price?: number
  token?: string
  storage?: string
  description?: string
}
