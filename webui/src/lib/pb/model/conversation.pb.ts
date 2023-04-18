/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/
export type Message = {
  id?: string
  conversationId?: string
  request?: string
  response?: string
  createdAt?: string
  updatedAt?: string
  tokenUsage?: number
  model?: string
}

export type ConversationWithoutMessages = {
  id?: string
  createdAt?: string
  updatedAt?: string
  instruction?: string
  title?: string
}

export type Conversation = {
  id?: string
  createdAt?: string
  updatedAt?: string
  instruction?: string
  title?: string
  messages?: Message[]
  tokenUsage?: number
}

export type ListConversationsRequest = {
  limit?: number
  offset?: number
}

export type ListConversationsResponse = {
  conversations?: ConversationWithoutMessages[]
}

export type GetConversationRequest = {
  id?: string
}

export type CreateConversationRequest = {
  instruction?: string
  title?: string
}

export type UpdateConversationRequest = {
  id?: string
  instruction?: string
  title?: string
}

export type DeleteConversationRequest = {
  id?: string
}

export type CreateMessageRequest = {
  conversationId?: string
  model?: string
  request?: string
  temperature?: number
  maxTokens?: number
}

export type UpdateMessageRequest = {
  conversationId?: string
  id?: string
  model?: string
  request?: string
  temperature?: number
  maxTokens?: number
}

export type ListMessagesRequest = {
  conversationId?: string
  limit?: number
  offset?: number
}

export type ListMessagesResponse = {
  messages?: Message[]
}

export type GetMessageRequest = {
  id?: string
  conversationId?: string
}

export type DeleteMessageRequest = {
  id?: string
  conversationId?: string
}
