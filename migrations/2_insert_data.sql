-- +goose Up
-- +goose StatementBegin
INSERT INTO role (name)
VALUES ('Администратор'),
       ('Менеджер'),
       ('Разработчик'),
       ('Тестировщик'),
       ('Дизайнер');

INSERT INTO employee (name, role_id)
VALUES ('Иванов Петр', (select id from role where name = 'Администратор')),
       ('Сидорова Анна', (select id from role where name = 'Менеджер')),
       ('Петров Алексей', (select id from role where name = 'Разработчик')),
       ('Козлова Елена', (select id from role where name = 'Разработчик')),
       ('Смирнов Дмитрий', (select id from role where name = 'Тестировщик')),
       ('Федорова Ольга', (select id from role where name = 'Дизайнер')),
       ('Николаев Игорь', (select id from role where name = 'Менеджер')),
       ('Васильева Мария', (select id from role where name = 'Разработчик'));
-- +goose StatementEnd


