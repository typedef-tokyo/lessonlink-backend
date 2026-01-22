-- Create "data_campuses" table
CREATE TABLE `data_campuses` (
  `id` int NOT NULL AUTO_INCREMENT,
  `campus` varchar(16) NOT NULL,
  `campus_name` varchar(32) NOT NULL,
  `order_index` int NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `campus` (`campus`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "sys_sessions" table
CREATE TABLE `sys_sessions` (
  `session_id` varchar(128) NOT NULL,
  `user_id` int NOT NULL,
  `value` text NOT NULL,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`session_id`),
  INDEX `user_id` (`user_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "tbl_schedule_invisible_rooms" table
CREATE TABLE `tbl_schedule_invisible_rooms` (
  `id` int NOT NULL AUTO_INCREMENT,
  `schedule_id` int NOT NULL,
  `room_index` int NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "data_lessons" table
CREATE TABLE `data_lessons` (
  `id` int NOT NULL AUTO_INCREMENT,
  `campus` varchar(16) NOT NULL,
  `name` varchar(32) NOT NULL,
  `duration` int NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `campus` (`campus`, `name`),
  CONSTRAINT `data_lessons_ibfk_1` FOREIGN KEY (`campus`) REFERENCES `data_campuses` (`campus`) ON UPDATE RESTRICT ON DELETE RESTRICT
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "data_rooms" table
CREATE TABLE `data_rooms` (
  `id` int NOT NULL AUTO_INCREMENT,
  `campus` varchar(16) NOT NULL,
  `room_index` int NOT NULL,
  `name` varchar(32) NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `campus` (`campus`, `room_index`),
  CONSTRAINT `data_rooms_ibfk_1` FOREIGN KEY (`campus`) REFERENCES `data_campuses` (`campus`) ON UPDATE RESTRICT ON DELETE RESTRICT
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "data_roles" table
CREATE TABLE `data_roles` (
  `id` int NOT NULL AUTO_INCREMENT,
  `role_key` varchar(16) NOT NULL,
  `role_name` varchar(16) NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `role_key` (`role_key`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "tbl_users" table
CREATE TABLE `tbl_users` (
  `id` int NOT NULL AUTO_INCREMENT,
  `role_key` varchar(16) NOT NULL,
  `user_name` varchar(64) NOT NULL,
  `password` text NOT NULL,
  `name` varchar(64) NOT NULL,
  `update_user_id` int NOT NULL,
  `delete_flag` int NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `role_key` (`role_key`),
  INDEX `update_user_id` (`update_user_id`),
  UNIQUE INDEX `user_name` (`user_name`),
  CONSTRAINT `tbl_users_ibfk_1` FOREIGN KEY (`role_key`) REFERENCES `data_roles` (`role_key`) ON UPDATE RESTRICT ON DELETE RESTRICT,
  CONSTRAINT `tbl_users_ibfk_2` FOREIGN KEY (`update_user_id`) REFERENCES `tbl_users` (`id`) ON UPDATE RESTRICT ON DELETE RESTRICT
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "tbl_schedules" table
CREATE TABLE `tbl_schedules` (
  `id` int NOT NULL AUTO_INCREMENT,
  `campus` varchar(16) NOT NULL,
  `title` varchar(64) NOT NULL,
  `history_index` int NOT NULL,
  `start_time` int NOT NULL,
  `end_time` int NOT NULL,
  `create_user` int NOT NULL,
  `last_update_user` int NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `campus` (`campus`),
  INDEX `create_user` (`create_user`),
  INDEX `last_update_user` (`last_update_user`),
  CONSTRAINT `tbl_schedules_ibfk_1` FOREIGN KEY (`campus`) REFERENCES `data_campuses` (`campus`) ON UPDATE RESTRICT ON DELETE RESTRICT,
  CONSTRAINT `tbl_schedules_ibfk_2` FOREIGN KEY (`create_user`) REFERENCES `tbl_users` (`id`) ON UPDATE RESTRICT ON DELETE RESTRICT,
  CONSTRAINT `tbl_schedules_ibfk_3` FOREIGN KEY (`last_update_user`) REFERENCES `tbl_users` (`id`) ON UPDATE RESTRICT ON DELETE RESTRICT
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "tbl_schedule_items" table
CREATE TABLE `tbl_schedule_items` (
  `id` int NOT NULL AUTO_INCREMENT,
  `schedule_id` int NOT NULL,
  `history_index` int NOT NULL,
  `lesson_id` int NOT NULL,
  `identifier` varchar(36) NOT NULL,
  `duration` int NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `lesson_id` (`lesson_id`),
  INDEX `schedule_id` (`schedule_id`),
  UNIQUE INDEX `schedule_id_2` (`schedule_id`, `history_index`, `identifier`),
  CONSTRAINT `tbl_schedule_items_ibfk_1` FOREIGN KEY (`schedule_id`) REFERENCES `tbl_schedules` (`id`) ON UPDATE RESTRICT ON DELETE RESTRICT,
  CONSTRAINT `tbl_schedule_items_ibfk_2` FOREIGN KEY (`lesson_id`) REFERENCES `data_lessons` (`id`) ON UPDATE RESTRICT ON DELETE RESTRICT
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "tbl_schedule_room_items" table
CREATE TABLE `tbl_schedule_room_items` (
  `id` int NOT NULL AUTO_INCREMENT,
  `schedule_id` int NOT NULL,
  `history_index` int NOT NULL,
  `item_tag` varchar(32) NOT NULL,
  `lesson_id` int NOT NULL,
  `identifier` varchar(36) NOT NULL,
  `duration` int NOT NULL,
  `start_time_hour` int NOT NULL,
  `start_time_minutes` int NOT NULL,
  `end_time_hour` int NOT NULL,
  `end_time_minutes` int NOT NULL,
  `room_index` int NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `lesson_id` (`lesson_id`),
  INDEX `schedule_id` (`schedule_id`),
  UNIQUE INDEX `schedule_id_2` (`schedule_id`, `history_index`, `identifier`),
  CONSTRAINT `tbl_schedule_room_items_ibfk_1` FOREIGN KEY (`schedule_id`) REFERENCES `tbl_schedules` (`id`) ON UPDATE RESTRICT ON DELETE RESTRICT
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
