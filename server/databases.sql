CREATE TABLE IF NOT EXISTS users(
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

CREATE TABLE IF NOT EXISTS posts(
    id INTEGER PRIMARY KEY,
    content TEXT NOT NULL,
    date TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS users_posts(
    user_ID INTEGER NOT NULL,
    post_ID INTEGER NOT NULL,

    FOREIGN KEY(user_ID)
        REFERENCES users(id)
            ON UPDATE CASCADE
            ON DELETE CASCADE,
    
    FOREIGN KEY(post_ID)
        REFERENCES posts(id)
            ON UPDATE CASCADE
            ON DELETE CASCADE

);

INSERT INTO users(username, password)
    VALUES
        ('cfabrica46',  '01234'),
        ('arturo',      '12345'),
        ('luis',        'lolsito123'),
        ('jose',        'asdasd');

INSERT INTO friends(id1,id2)
    VALUES
        (1,2),
        (1,3);

INSERT INTO posts(content,date)
    VALUES
        ("HOLIIIIII",   datetime('now','localtime')),
        ("HOLA :D",     datetime('now','localtime')),
        ("UWU",         datetime('now','localtime')),
        (">:V",         datetime('now','localtime')),
        ("owo",         datetime('now','localtime'));

INSERT INTO users_posts(user_ID,post_ID)
    VALUES
        (1,1),
        (1,3),
        (2,2),
        (3,4),
        (4,5);