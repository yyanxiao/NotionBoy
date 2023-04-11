/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

export enum NumericEnum {
  ZERO = "ZERO",
  ONE = "ONE",
}

export type ErrorResponse = {
  correlationId?: string
  error?: ErrorObject
}

export type ErrorObject = {
  code?: number
  message?: string
}

export type GenrateTokenRequest = {
  magicCode?: string
  qrcode?: string
}

export type GenrateTokenResponse = {
  token?: string
  type?: string
  expiry?: string
}

export type GenerateApiKeyResponse = {
  apiKey?: string
}

export type OAuthCallbackRequest = {
  code?: string
  state?: string
}

export type OAuthURLRequest = {
}

export type OAuthProvider = {
  name?: string
  url?: string
}

export type OAuthURLResponse = {
  providers?: OAuthProvider[]
}

export type GenerateWechatQRCodeResponse = {
  url?: string
  qrcode?: string
}

export type Prompt = {
  act?: string
  prompt?: string
}

export type ListPromptsResponse = {
  prompts?: Prompt[]
}
