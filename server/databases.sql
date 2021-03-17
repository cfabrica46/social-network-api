CREATE TABLE If NOT EXISTS users(
    id INTEGER PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    password  TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS friends(
    id1 INTERGER NOT NULL,
    id2 INTERGER NOT NULL,

    FOREIGN KEY(id1)
        REFERENCES users(id)
            ON UPDATE CASCADE
            ON DELETE CASCADE,

    FOREIGN KEY(id2)
        REFERENCES users(id)
            ON UPDATE CASCADE
            ON DELETE CASCADE
);

INSERT INTO users(username, password)
    VALUES
        ('cfabrica46','01234'),
        ('arturo','12345'),
        ('luis', 'lolsito123');

INSERT INTO friends(id1,id2)
    VALUES
        (1,2),
        (1,3);