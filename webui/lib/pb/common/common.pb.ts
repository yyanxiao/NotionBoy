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
