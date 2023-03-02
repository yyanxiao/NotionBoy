# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [model/common.proto](#model_common-proto)
    - [ErrorObject](#servicev1-ErrorObject)
    - [ErrorResponse](#servicev1-ErrorResponse)
    - [GenerateApiKeyResponse](#servicev1-GenerateApiKeyResponse)
    - [GenrateTokenResponse](#servicev1-GenrateTokenResponse)

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






<a name="servicev1-GenrateTokenResponse"></a>

### GenrateTokenResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| token | [string](#string) |  |  |
| type | [string](#string) |  |  |
| expiry | [string](#string) |  |  |








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
| messages | [Message](#servicev1-Message) | repeated |  |






<a name="servicev1-ConversationWithoutMessages"></a>

### ConversationWithoutMessages



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | id uuid of the conversation |
| created_at | [string](#string) |  | user_id uuid of the user |
| updated_at | [string](#string) |  |  |






<a name="servicev1-CreateConversationRequest"></a>

### CreateConversationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| instruction | [string](#string) |  | instruction instruction of the conversation, if this is empty, the conversation will be created with a default instruction |






<a name="servicev1-CreateMessageRequest"></a>

### CreateMessageRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| conversation_id | [string](#string) |  | conversation_id uuid of the conversation |
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
| GenrateToken | [.google.protobuf.Empty](#google-protobuf-Empty) | [GenrateTokenResponse](#servicev1-GenrateTokenResponse) | GenrateToken generates a token for the user. using api key in the header. |
| GenerateApiKey | [.google.protobuf.Empty](#google-protobuf-Empty) | [GenerateApiKeyResponse](#servicev1-GenerateApiKeyResponse) | GenerateApiKey generate a new api key for the user |
| DeleteApiKey | [.google.protobuf.Empty](#google-protobuf-Empty) | [.google.protobuf.Empty](#google-protobuf-Empty) | DeleteApiKey delete the api key for the user |
| CreateConversation | [CreateConversationRequest](#servicev1-CreateConversationRequest) | [Conversation](#servicev1-Conversation) |  |
| GetConversation | [GetConversationRequest](#servicev1-GetConversationRequest) | [Conversation](#servicev1-Conversation) |  |
| ListConversations | [ListConversationsRequest](#servicev1-ListConversationsRequest) | [ListConversationsResponse](#servicev1-ListConversationsResponse) |  |
| DeleteConversation | [DeleteConversationRequest](#servicev1-DeleteConversationRequest) | [.google.protobuf.Empty](#google-protobuf-Empty) |  |
| CreateMessage | [CreateMessageRequest](#servicev1-CreateMessageRequest) | [Message](#servicev1-Message) |  |
| GetMessage | [GetMessageRequest](#servicev1-GetMessageRequest) | [Message](#servicev1-Message) |  |
| ListMessages | [ListMessagesRequest](#servicev1-ListMessagesRequest) | [ListMessagesResponse](#servicev1-ListMessagesResponse) |  |
| DeleteMessage | [DeleteMessageRequest](#servicev1-DeleteMessageRequest) | [.google.protobuf.Empty](#google-protobuf-Empty) |  |





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
