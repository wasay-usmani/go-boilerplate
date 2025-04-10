
-- +migrate Up
create schema if not exists go_boilerplate CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci;
create user if not exists 'goboiler'@'localhost' identified by 'goboiler';
create user if not exists 'goboiler'@'%' identified by 'goboiler';
GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, ALTER, REFERENCES ON go_boilerplate.* to 'goboiler'@'localhost';
GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, ALTER, REFERENCES ON go_boilerplate.* to 'goboiler'@'%';

-- +migrate Down
revoke all on goboiler.* from 'goboiler'@'%';
revoke all on goboiler.* from 'goboiler'@'localhost';
drop user if exists 'goboiler'@'%';
drop user if exists 'goboiler'@'localhost';
drop schema if exists go_boilerplate;
