-- First ensure we have at least one user to reference
INSERT IGNORE INTO `users` (`id`, `name`, `email`, `password`, `is_admin`, `created_at`, `updated_at`) 
VALUES (9999, 'system_user', 'system@system.com', 'temp_password', false, NOW(), NOW());

-- Add the user_id column initially as nullable
ALTER TABLE `orders` ADD COLUMN `user_id` int;

-- Update existing orders to reference the system user
UPDATE `orders` SET `user_id` = 9999 WHERE `user_id` IS NULL;

-- Now make the column NOT NULL and add the foreign key
ALTER TABLE `orders` 
    MODIFY COLUMN `user_id` int NOT NULL,
    ADD CONSTRAINT `user_id_fk` FOREIGN KEY (`user_id`) 
        REFERENCES `users` (`id`);