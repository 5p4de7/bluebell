CREATE TABLE `user` (
            `id` bigint(20) NOT NULL AUTO_INCREMENT,
            `user_id` bigint(20) NOT NULL,
            `username` varchar(64) collate utf8mb4_general_ci not null ,
            `password` varchar(64) collate utf8mb4_general_ci not null ,
            `email` varchar(64) collate utf8mb4_general_ci ,
            `gender` tinyint(4) NOT NULL DEFAULT '0',
            `create_time` timestamp not null default current_timestamp,
            `update_time` timestamp not null default current_timestamp
                        ON UPDATE current_timestamp,
    primary key (`id`),
    unique key `idx_username` (`username`) using  BTREE,
    unique key `idx_user_id` (`user_id`) using  BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 collate=utf8mb4_general_ci;