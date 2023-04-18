-- Modify "conversation_messages" table
ALTER TABLE `conversation_messages` ADD COLUMN `model` varchar(255) NOT NULL DEFAULT "gpt-3.5-turbo";
-- Modify "conversations" table
ALTER TABLE `conversations` ADD COLUMN `token_usage` bigint NULL;
