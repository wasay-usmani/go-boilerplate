
-- +migrate Up
create schema if not exists go_boilerplate CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci;
create user if not exists 'go_boilerplate'@'localhost' identified by 'go_boilerplate';
create user if not exists 'go_boilerplate'@'%' identified by 'go_boilerplate';
GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, ALTER, REFERENCES ON go_boilerplate.* to 'go_boilerplate'@'localhost';
GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, ALTER, REFERENCES ON go_boilerplate.* to 'go_boilerplate'@'%';

-- +migrate Down
revoke all on 'go_boilerplate'.* from 'go_boilerplate'@'%';
revoke all on 'go_boilerplate'.* from 'go_boilerplate'@'localhost';
drop user if exists 'go_boilerplate'@'%';
drop user if exists 'go_boilerplate'@'localhost';
drop schema if exists go_boilerplate;
