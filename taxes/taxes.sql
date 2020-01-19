CREATE TABLE IF NOT EXISTS taxes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email varchar(32) UNIQUE NOT NULL,
    amount REAL
);

INSERT INTO taxes(email, amount) VALUES ("a@a.com", 250.55), ("b@a.com", 1000);