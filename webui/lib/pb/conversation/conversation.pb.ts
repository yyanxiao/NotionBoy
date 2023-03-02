/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/
export type Message = {
  id?: string
  conversationId?: string
  userId?: string
  request?: string
  response?: string
  createdAt?: string
  updatedAt?: string
}

export type ConversationWithoutMessages = {
  id?: string
  userId?: string
  createdAt?: string
  updatedAt?: string
}

export type ConversationConversationWithoutMessages = {
  id?: string
  userId?: string
  createdAt?: string
  updatedAt?: string
}

export type Conversation = {
  conversation?: ConversationConversationWithoutMessages
  messages?: string[]
}

export type ListConversationsResponse = {
  conversations?: Conversation[]
}

export type ListConversationsRequest = {
  limit?: number
  offset?: number
}

export type GetConversationRequest = {
  id?: string
}

export type CreateConversationRequest = {
  request?: string
}

export type UpdateConversationRequest = {
  id?: string
  userId?: string
  conversationId?: string
  messageId?: string
  request?: string
  response?: string
}

export type DeleteConversationRequest = {
  id?: string
}

export type DeleteConversationResponse = {
  success?: boolean
}
