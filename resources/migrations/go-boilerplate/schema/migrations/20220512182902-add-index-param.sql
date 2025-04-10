
-- +migrate Up
GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, ALTER, INDEX, REFERENCES ON go_boilerplate.* to 'goboiler'@'localhost';
GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, ALTER, INDEX, REFERENCES ON go_boilerplate.* to 'goboiler'@'%';

-- +migrate Down
GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, ALTER, REFERENCES ON go_boilerplate.* to 'goboiler'@'localhost';
GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, ALTER, REFERENCES ON go_boilerplate.* to 'goboiler'@'%';
