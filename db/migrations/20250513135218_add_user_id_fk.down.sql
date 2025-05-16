-- Remove the foreign key and column from orders table
ALTER TABLE `orders`
    DROP FOREIGN KEY `user_id_fk`,
    DROP COLUMN `user_id`;

-- Remove the system user only if it exists and no other tables reference it
DELETE FROM `users` 
WHERE id = 9999 
  AND email = 'system@system.com';