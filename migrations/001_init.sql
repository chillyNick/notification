-- +goose Up
CREATE TABLE mail (
    user_id int PRIMARY KEY,
    mail varchar(128) not null
);

INSERT INTO mail(user_id, mail)
VALUES
    (123, 'test@test.com');

-- +goose Down
DROP TABLE mail;
