INSERT INTO `data_campuses` (`id`, `campus`, `campus_name`, `order_index`, `created_at`, `updated_at`) VALUES
(1, 'shibuya', '渋谷', 1, NOW(), NOW()),
(2, 'shinjuku', '新宿', 2, NOW(), NOW()),
(3, 'ikebukuro', '池袋', 3, NOW(), NOW());


INSERT INTO `data_roles` (`id`, `role_key`, `role_name`, `created_at`, `updated_at`) VALUES
(1, 'owner', 'オーナー', NOW(), NOW()),
(2, 'editor', '編集者', NOW(), NOW()),
(3, 'viewer', '閲覧者', NOW(), NOW());


INSERT INTO `tbl_users` (`id`, `role_key`, `user_name`, `password`, `name`, `update_user_id`, `delete_flag`, `created_at`, `updated_at`) VALUES
(1, 'owner', 'admin@admin.com', '$argon2id$v=19$m=65536,t=1,p=4$vGqq6iPu5/iGGWoa0PjKBQ$iaWzePFKW9thZnjthrpNh6jUQEYjQUICmstVNnM6jnE', 'admin', 1, 0, NOW(), NOW()); -- パスワードはadmin
