-- Create "accounts" table
CREATE TABLE `accounts`
(
    `id`                 bigint                                                             NOT NULL AUTO_INCREMENT,
    `created_at`         timestamp                                                          NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`         timestamp                                                          NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted`            bool                                                               NOT NULL DEFAULT 0,
    `uuid`               char(36)                                                           NULL,
    `user_id`            varchar(255)                                                       NOT NULL,
    `user_type`          enum ('wechat','telegram','github','google','twitter','microsoft') NULL     DEFAULT "wechat",
    `database_id`        varchar(255)                                                       NULL,
    `access_token`       varchar(255)                                                       NULL,
    `notion_user_id`     varchar(255)                                                       NULL,
    `notion_user_email`  varchar(255)                                                       NULL,
    `is_latest_schema`   bool                                                               NOT NULL DEFAULT 0,
    `is_openai_api_user` bool                                                               NOT NULL DEFAULT 0,
    `openai_api_key`     varchar(255)                                                       NULL,
    `api_key`            char(36)                                                           NULL,
    `is_admin`           bool                                                               NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `account_user_id_user_type` (`user_id`, `user_type`),
    UNIQUE INDEX `uuid` (`uuid`)
) CHARSET utf8mb4
  COLLATE utf8mb4_0900_ai_ci;
-- Create "chat_histories" table
CREATE TABLE `chat_histories`
(
    `id`               bigint       NOT NULL AUTO_INCREMENT,
    `created_at`       timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`       timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted`          bool         NOT NULL DEFAULT 0,
    `user_id`          bigint       NOT NULL,
    `conversation_idx` bigint       NOT NULL,
    `conversation_id`  char(36)     NOT NULL,
    `message_id`       varchar(255) NULL,
    `message_idx`      bigint       NULL,
    `request`          longtext     NULL,
    `response`         longtext     NULL,
    `token_usage`      bigint       NULL,
    PRIMARY KEY (`id`)
) CHARSET utf8mb4
  COLLATE utf8mb4_0900_ai_ci;
-- Create "conversations" table
CREATE TABLE `conversations`
(
    `id`          bigint       NOT NULL AUTO_INCREMENT,
    `created_at`  timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`  timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted`     bool         NOT NULL DEFAULT 0,
    `uuid`        char(36)     NOT NULL,
    `user_id`     char(36)     NOT NULL,
    `instruction` longtext     NULL,
    `title`       varchar(255) NULL,
    PRIMARY KEY (`id`),
    INDEX `conversation_user_id_created_at` (`user_id`, `created_at`),
    UNIQUE INDEX `uuid` (`uuid`)
) CHARSET utf8mb4
  COLLATE utf8mb4_0900_ai_ci;
-- Create "orders" table
CREATE TABLE `orders`
(
    `id`           bigint                                                                          NOT NULL AUTO_INCREMENT,
    `created_at`   timestamp                                                                       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`   timestamp                                                                       NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted`      bool                                                                            NOT NULL DEFAULT 0,
    `uuid`         char(36)                                                                        NOT NULL,
    `user_id`      char(36)                                                                        NOT NULL,
    `product_id`   char(36)                                                                        NOT NULL,
    `price`        double                                                                          NOT NULL,
    `status`       enum ('Unpaid','Paying','Paid','Processing','Cancelled','Refunded','Completed') NOT NULL DEFAULT "Unpaid",
    `note`         longtext                                                                        NULL,
    `payment_info` longtext                                                                        NULL,
    PRIMARY KEY (`id`),
    INDEX `order_user_id_created_at` (`user_id`, `created_at`),
    UNIQUE INDEX `uuid` (`uuid`)
) CHARSET utf8mb4
  COLLATE utf8mb4_0900_ai_ci;
-- Create "products" table
CREATE TABLE `products`
(
    `id`          bigint       NOT NULL AUTO_INCREMENT,
    `created_at`  timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`  timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted`     bool         NOT NULL DEFAULT 0,
    `uuid`        char(36)     NOT NULL,
    `name`        varchar(255) NOT NULL DEFAULT "Free",
    `description` longtext     NOT NULL,
    `price`       double       NOT NULL,
    `token`       bigint       NOT NULL DEFAULT 10000,
    `storage`     bigint       NOT NULL DEFAULT 100,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uuid` (`uuid`)
) CHARSET utf8mb4
  COLLATE utf8mb4_0900_ai_ci;
-- Create "quota" table
CREATE TABLE `quota`
(
    `id`         bigint       NOT NULL AUTO_INCREMENT,
    `created_at` timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted`    bool         NOT NULL DEFAULT 0,
    `user_id`    bigint       NOT NULL,
    `plan`       varchar(255) NOT NULL,
    `reset_time` timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `token`      bigint       NOT NULL,
    `token_used` bigint       NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `quota_user_id` (`user_id`)
) CHARSET utf8mb4
  COLLATE utf8mb4_0900_ai_ci;
-- Create "wechat_session" table
CREATE TABLE `wechat_session`
(
    `id`            bigint       NOT NULL AUTO_INCREMENT,
    `created_at`    timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`    timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted`       bool         NOT NULL DEFAULT 0,
    `session`       blob         NOT NULL,
    `dummy_user_id` varchar(255) NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `dummy_user_id` (`dummy_user_id`)
) CHARSET utf8mb4
  COLLATE utf8mb4_0900_ai_ci;
-- Create "conversation_messages" table
CREATE TABLE `conversation_messages`
(
    `id`                                 bigint    NOT NULL AUTO_INCREMENT,
    `created_at`                         timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`                         timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted`                            bool      NOT NULL DEFAULT 0,
    `uuid`                               char(36)  NOT NULL,
    `user_id`                            char(36)  NOT NULL,
    `conversation_id`                    char(36)  NOT NULL,
    `request`                            longtext  NULL,
    `response`                           longtext  NULL,
    `token_usage`                        bigint    NULL,
    `conversation_conversation_messages` bigint    NULL,
    PRIMARY KEY (`id`),
    INDEX `conversation_messages_conversations_conversation_messages` (`conversation_conversation_messages`),
    INDEX `conversationmessage_user_id_conversation_id_created_at` (`user_id`, `conversation_id`, `created_at`),
    UNIQUE INDEX `uuid` (`uuid`),
    CONSTRAINT `conversation_messages_conversations_conversation_messages` FOREIGN KEY (`conversation_conversation_messages`) REFERENCES `conversations` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL
) CHARSET utf8mb4
  COLLATE utf8mb4_0900_ai_ci;
