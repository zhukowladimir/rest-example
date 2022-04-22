BEGIN TRANSACTION;
CREATE TABLE IF NOT EXISTS "players" (
	"id"        TEXT NOT NULL PRIMARY KEY UNIQUE,
	"username"	TEXT NOT NULL UNIQUE,
    "avatar"    TEXT,
    "sex"       TEXT,
    "email"     TEXT UNIQUE
);
COMMIT;
