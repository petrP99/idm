CREATE TABLE IF NOT EXISTS role
(
    id         bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name       text        NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS employee
(
    id         bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name       text        NOT NULL,
    role_id    bigint REFERENCES role (id),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

INSERT INTO role (name)
VALUES ('Администратор'),
       ('Менеджер'),
       ('Разработчик'),
       ('Тестировщик'),
       ('Дизайнер');

INSERT INTO employee (name, role_id)
VALUES ('Иванов Петр', 1),
       ('Сидорова Анна', 2),
       ('Петров Алексей', 3),
       ('Козлова Елена', 3)


