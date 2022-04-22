BEGIN TRANSACTION;
CREATE TABLE IF NOT EXISTS "players" (
	"id"        TEXT NOT NULL PRIMARY KEY UNIQUE,
	"username"	TEXT NOT NULL UNIQUE,
    "avatar"    TEXT,
    "sex"       TEXT,
    "email"     TEXT UNIQUE
);
CREATE TABLE IF NOT EXISTS "stats" (
    "id"            INTEGER NOT NULL PRIMARY KEY UNIQUE,
    "loss_count"    INTEGER,
    "win_count"     INTEGER,
    "duration"      INTEGER
);
CREATE TABLE IF NOT EXISTS "pid_sid" (
    "pid"   TEXT NOT NULL PRIMARY KEY UNIQUE,
    "sid"   INTEGER NOT NULL UNIQUE
);
COMMIT;
