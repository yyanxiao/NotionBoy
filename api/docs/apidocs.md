# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [model/common.proto](#model_common-proto)
    - [ErrorObject](#servicev1-ErrorObject)
    - [ErrorResponse](#servicev1-ErrorResponse)
    - [GenerateApiKeyResponse](#servicev1-GenerateApiKeyResponse)
    - [GenerateWechatQRCodeResponse](#servicev1-GenerateWechatQRCodeResponse)
    - [GenrateTokenRequest](#servicev1-GenrateTokenRequest)
    - [GenrateTokenResponse](#servicev1-GenrateTokenResponse)
    - [OAuthCallbackRequest](#servicev1-OAuthCallbackRequest)
    - [OAuthProvider](#servicev1-OAuthProvider)
    - [OAuthURLRequest](#servicev1-OAuthURLRequest)
    - [OAuthURLResponse](#servicev1-OAuthURLResponse)

    - [NumericEnum](#servicev1-NumericEnum)

- [model/conversation.proto](#model_conversation-proto)
    - [Conversation](#servicev1-Conversation)
    - [ConversationWithoutMessages](#servicev1-ConversationWithoutMessages)
    - [CreateConversationRequest](#servicev1-CreateConversationRequest)
    - [CreateMessageRequest](#servicev1-CreateMessageRequest)
    - [DeleteConversationRequest](#servicev1-DeleteConversationRequest)
    - [DeleteMessageRequest](#servicev1-DeleteMessageRequest)
    - [GetConversationRequest](#servicev1-GetConversationRequest)
    - [GetMessageRequest](#servicev1-GetMessageRequest)
    - [ListConversationsRequest](#servicev1-ListConversationsRequest)
    - [ListConversationsResponse](#servicev1-ListConversationsResponse)
    - [ListMessagesRequest](#servicev1-ListMessagesRequest)
    - [ListMessagesResponse](#servicev1-ListMessagesResponse)
    - [Message](#servicev1-Message)
    - [UpdateConversationRequest](#servicev1-UpdateConversationRequest)

- [model/product.proto](#model_product-proto)
    - [CreateProductRequest](#servicev1-CreateProductRequest)
    - [DeleteProductRequest](#servicev1-DeleteProductRequest)
    - [GetProductRequest](#servicev1-GetProductRequest)
    - [ListProductsRequest](#servicev1-ListProductsRequest)
    - [ListProductsResponse](#servicev1-ListProductsResponse)
    - [Product](#servicev1-Product)
    - [UpdateProductRequest](#servicev1-UpdateProductRequest)

- [model/order.proto](#model_order-proto)
    - [CreateOrderRequest](#servicev1-CreateOrderRequest)
    - [DeleteOrderRequest](#servicev1-DeleteOrderRequest)
    - [GetOrderRequest](#servicev1-GetOrderRequest)
    - [ListOrdersRequest](#servicev1-ListOrdersRequest)
    - [ListOrdersResponse](#servicev1-ListOrdersResponse)
    - [Order](#servicev1-Order)
    - [PayOrderRequest](#servicev1-PayOrderRequest)
    - [PayOrderResponse](#servicev1-PayOrderResponse)
    - [UpdateOrderRequest](#servicev1-UpdateOrderRequest)

- [server.proto](#server-proto)
    - [CheckStatusResponse](#servicev1-CheckStatusResponse)

    - [Service](#servicev1-Service)

- [Scalar Value Types](#scalar-value-types)



<a name="model_common-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## model/common.proto



<a name="servicev1-ErrorObject"></a>

### ErrorObject



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| code | [int32](#int32) |  |  |
| message | [string](#string) |  |  |






<a name="servicev1-ErrorResponse"></a>

### ErrorResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| correlationId | [string](#string) |  |  |
| error | [ErrorObject](#servicev1-ErrorObject) |  |  |






<a name="servicev1-GenerateApiKeyResponse"></a>

### GenerateApiKeyResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| api_key | [string](#string) |  |  |






<a name="servicev1-GenerateWechatQRCodeResponse"></a>

### GenerateWechatQRCodeResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| url | [string](#string) |  |  |
| qrcode | [string](#string) |  |  |






<a name="servicev1-GenrateTokenRequest"></a>

### GenrateTokenRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| magicCode | [string](#string) |  |  |
| qrcode | [string](#string) |  |  |






<a name="servicev1-GenrateTokenResponse"></a>

### GenrateTokenResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| token | [string](#string) |  |  |
| type | [string](#string) |  |  |
| expiry | [string](#string) |  |  |






<a name="servicev1-OAuthCallbackRequest"></a>

### OAuthCallbackRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| code | [string](#string) |  |  |
| state | [string](#string) |  |  |






<a name="servicev1-OAuthProvider"></a>

### OAuthProvider



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| url | [string](#string) |  |  |






<a name="servicev1-OAuthURLRequest"></a>

### OAuthURLRequest







<a name="servicev1-OAuthURLResponse"></a>

### OAuthURLResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| providers | [OAuthProvider](#servicev1-OAuthProvider) | repeated |  |








<a name="servicev1-NumericEnum"></a>

### NumericEnum
NumericEnum is one or zero.

| Name | Number | Description |
| ---- | ------ | ----------- |
| ZERO | 0 | ZERO means 0 |
| ONE | 1 | ONE means 1 |










<a name="model_conversation-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## model/conversation.proto



<a name="servicev1-Conversation"></a>

### Conversation



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | id uuid of the conversation |
| created_at | [string](#string) |  |  |
| updated_at | [string](#string) |  |  |
| instruction | [string](#string) |  |  |
| title | [string](#string) |  |  |
| messages | [Message](#servicev1-Message) | repeated |  |
| token_usage | [int32](#int32) |  |  |






<a name="servicev1-ConversationWithoutMessages"></a>

### ConversationWithoutMessages



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | id uuid of the conversation |
| created_at | [string](#string) |  | user_id uuid of the user |
| updated_at | [string](#string) |  |  |
| instruction | [string](#string) |  |  |
| title | [string](#string) |  |  |






<a name="servicev1-CreateConversationRequest"></a>

### CreateConversationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| instruction | [string](#string) |  | instruction instruction of the conversation, if this is empty, the conversation will be created with a default instruction |
| title | [string](#string) |  |  |






<a name="servicev1-CreateMessageRequest"></a>

### CreateMessageRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| conversation_id | [string](#string) |  | conversation_id uuid of the conversation |
| model | [string](#string) |  |  |
| request | [string](#string) |  | request request of the message |






<a name="servicev1-DeleteConversationRequest"></a>

### DeleteConversationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | id uuid of the conversation |






<a name="servicev1-DeleteMessageRequest"></a>

### DeleteMessageRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | id uuid of the message |
| conversation_id | [string](#string) |  | conversation_id uuid of the conversation |






<a name="servicev1-GetConversationRequest"></a>

### GetConversationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | id uuid of the conversation |






<a name="servicev1-GetMessageRequest"></a>

### GetMessageRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | id uuid of the message |
| conversation_id | [string](#string) |  | conversation_id uuid of the conversation |






<a name="servicev1-ListConversationsRequest"></a>

### ListConversationsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| limit | [int32](#int32) |  |  |
| offset | [int32](#int32) |  |  |






<a name="servicev1-ListConversationsResponse"></a>

### ListConversationsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| conversations | [ConversationWithoutMessages](#servicev1-ConversationWithoutMessages) | repeated |  |






<a name="servicev1-ListMessagesRequest"></a>

### ListMessagesRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| conversation_id | [string](#string) |  | conversation_id uuid of the conversation |
| limit | [int32](#int32) |  |  |
| offset | [int32](#int32) |  |  |






<a name="servicev1-ListMessagesResponse"></a>

### ListMessagesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| messages | [Message](#servicev1-Message) | repeated |  |






<a name="servicev1-Message"></a>

### Message



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | id uuid of the message |
| conversation_id | [string](#string) |  | conversation_id uuid of the conversation |
| request | [string](#string) |  | user_id uuid of the user |
| response | [string](#string) |  |  |
| created_at | [string](#string) |  |  |
| updated_at | [string](#string) |  |  |
| token_usage | [int32](#int32) |  |  |
| model | [string](#string) |  |  |






<a name="servicev1-UpdateConversationRequest"></a>

### UpdateConversationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | id uuid of the conversation |
| instruction | [string](#string) |  | instruction instruction of the conversation, if this is empty, the conversation will be created with a default instruction |
| title | [string](#string) |  |  |















<a name="model_product-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## model/product.proto



<a name="servicev1-CreateProductRequest"></a>

### CreateProductRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| price | [float](#float) |  |  |
| token | [int64](#int64) |  |  |
| storage | [int64](#int64) |  |  |
| description | [string](#string) |  |  |






<a name="servicev1-DeleteProductRequest"></a>

### DeleteProductRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |






<a name="servicev1-GetProductRequest"></a>

### GetProductRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |






<a name="servicev1-ListProductsRequest"></a>

### ListProductsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| limit | [int32](#int32) |  |  |
| offset | [int32](#int32) |  |  |






<a name="servicev1-ListProductsResponse"></a>

### ListProductsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| products | [Product](#servicev1-Product) | repeated |  |






<a name="servicev1-Product"></a>

### Product



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| created_at | [string](#string) |  |  |
| updated_at | [string](#string) |  |  |
| deleted | [string](#string) |  |  |
| uuid | [string](#string) |  |  |
| name | [string](#string) |  |  |
| price | [float](#float) |  |  |
| token | [int64](#int64) |  |  |
| storage | [int64](#int64) |  |  |
| description | [string](#string) |  |  |






<a name="servicev1-UpdateProductRequest"></a>

### UpdateProductRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| name | [string](#string) |  |  |
| price | [float](#float) |  |  |
| token | [int64](#int64) |  |  |
| storage | [int64](#int64) |  |  |
| description | [string](#string) |  |  |















<a name="model_order-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## model/order.proto



<a name="servicev1-CreateOrderRequest"></a>

### CreateOrderRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| product_id | [string](#string) |  |  |
| note | [string](#string) |  |  |






<a name="servicev1-DeleteOrderRequest"></a>

### DeleteOrderRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |






<a name="servicev1-GetOrderRequest"></a>

### GetOrderRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| is_check_payment | [bool](#bool) |  |  |






<a name="servicev1-ListOrdersRequest"></a>

### ListOrdersRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| status | [string](#string) |  |  |
| limit | [int32](#int32) |  |  |
| offset | [int32](#int32) |  |  |






<a name="servicev1-ListOrdersResponse"></a>

### ListOrdersResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| orders | [Order](#servicev1-Order) | repeated |  |






<a name="servicev1-Order"></a>

### Order



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| created_at | [string](#string) |  |  |
| updated_at | [string](#string) |  |  |
| deleted | [string](#string) |  |  |
| uuid | [string](#string) |  |  |
| user_id | [string](#string) |  |  |
| price | [float](#float) |  |  |
| status | [string](#string) |  |  |
| note | [string](#string) |  |  |
| payment_method | [string](#string) |  |  |
| product | [Product](#servicev1-Product) |  |  |






<a name="servicev1-PayOrderRequest"></a>

### PayOrderRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| payment_method | [string](#string) |  |  |






<a name="servicev1-PayOrderResponse"></a>

### PayOrderResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| status | [string](#string) |  |  |
| qrcode | [string](#string) |  |  |






<a name="servicev1-UpdateOrderRequest"></a>

### UpdateOrderRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| note | [string](#string) |  |  |















<a name="server-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## server.proto



<a name="servicev1-CheckStatusResponse"></a>

### CheckStatusResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| status | [google.rpc.Status](#google-rpc-Status) |  |  |












<a name="servicev1-Service"></a>

### Service


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Status | [.google.protobuf.Empty](#google-protobuf-Empty) | [CheckStatusResponse](#servicev1-CheckStatusResponse) |  |
| GenrateToken | [GenrateTokenRequest](#servicev1-GenrateTokenRequest) | [GenrateTokenResponse](#servicev1-GenrateTokenResponse) | GenrateToken generates a token for the user. using api key in the header. |
| OAuthProviders | [OAuthURLRequest](#servicev1-OAuthURLRequest) | [OAuthURLResponse](#servicev1-OAuthURLResponse) | get all Oauth providers |
| OAuthCallback | [OAuthCallbackRequest](#servicev1-OAuthCallbackRequest) | [GenrateTokenResponse](#servicev1-GenrateTokenResponse) | AuthCallback callback for oauth, will generate a token for the user |
| GenerateApiKey | [.google.protobuf.Empty](#google-protobuf-Empty) | [GenerateApiKeyResponse](#servicev1-GenerateApiKeyResponse) | GenerateApiKey generate a new api key for the user |
| DeleteApiKey | [.google.protobuf.Empty](#google-protobuf-Empty) | [.google.protobuf.Empty](#google-protobuf-Empty) | DeleteApiKey delete the api key for the user |
| GenerateWechatQRCode | [.google.protobuf.Empty](#google-protobuf-Empty) | [GenerateWechatQRCodeResponse](#servicev1-GenerateWechatQRCodeResponse) |  |
| CreateConversation | [CreateConversationRequest](#servicev1-CreateConversationRequest) | [Conversation](#servicev1-Conversation) |  |
| UpdateConversation | [UpdateConversationRequest](#servicev1-UpdateConversationRequest) | [Conversation](#servicev1-Conversation) | UpdateConversation update the conversation |
| GetConversation | [GetConversationRequest](#servicev1-GetConversationRequest) | [Conversation](#servicev1-Conversation) |  |
| ListConversations | [ListConversationsRequest](#servicev1-ListConversationsRequest) | [ListConversationsResponse](#servicev1-ListConversationsResponse) |  |
| DeleteConversation | [DeleteConversationRequest](#servicev1-DeleteConversationRequest) | [.google.protobuf.Empty](#google-protobuf-Empty) |  |
| CreateMessage | [CreateMessageRequest](#servicev1-CreateMessageRequest) | [Message](#servicev1-Message) stream |  |
| GetMessage | [GetMessageRequest](#servicev1-GetMessageRequest) | [Message](#servicev1-Message) |  |
| ListMessages | [ListMessagesRequest](#servicev1-ListMessagesRequest) | [ListMessagesResponse](#servicev1-ListMessagesResponse) |  |
| DeleteMessage | [DeleteMessageRequest](#servicev1-DeleteMessageRequest) | [.google.protobuf.Empty](#google-protobuf-Empty) |  |
| CreateOrder | [CreateOrderRequest](#servicev1-CreateOrderRequest) | [Order](#servicev1-Order) | CreateOrder create a new order |
| GetOrder | [GetOrderRequest](#servicev1-GetOrderRequest) | [Order](#servicev1-Order) | get order |
| ListOrders | [ListOrdersRequest](#servicev1-ListOrdersRequest) | [ListOrdersResponse](#servicev1-ListOrdersResponse) | list orders |
| DeleteOrder | [DeleteOrderRequest](#servicev1-DeleteOrderRequest) | [.google.protobuf.Empty](#google-protobuf-Empty) | delete order |
| UpdateOrder | [UpdateOrderRequest](#servicev1-UpdateOrderRequest) | [Order](#servicev1-Order) | update order |
| PayOrder | [PayOrderRequest](#servicev1-PayOrderRequest) | [PayOrderResponse](#servicev1-PayOrderResponse) | pay order |
| CreateProduct | [CreateProductRequest](#servicev1-CreateProductRequest) | [Product](#servicev1-Product) | CRUUD for products |
| GetProduct | [GetProductRequest](#servicev1-GetProductRequest) | [Product](#servicev1-Product) |  |
| ListProducts | [ListProductsRequest](#servicev1-ListProductsRequest) | [ListProductsResponse](#servicev1-ListProductsResponse) |  |
| DeleteProduct | [DeleteProductRequest](#servicev1-DeleteProductRequest) | [.google.protobuf.Empty](#google-protobuf-Empty) |  |
| UpdateProduct | [UpdateProductRequest](#servicev1-UpdateProductRequest) | [Product](#servicev1-Product) |  |





## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |
