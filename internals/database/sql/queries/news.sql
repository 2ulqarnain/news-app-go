-- name: CreateTableNews :exec
CREATE TABLE IF NOT EXISTS news (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    news_url TEXT NOT NULL,
    image_url TEXT NOT NULL,
    slug TEXT NOT NULL,
    created_on TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER created_on_trigger
BEFORE UPDATE OF created_on ON news FOR EACH ROW
WHEN NEW.created_on IS NOT OLD.created_on
BEGIN
   SELECT RAISE(ABORT,'Error: The created_on column is read-only and cannot be modified.');
END;


-- name: CreateNews :exec
INSERT INTO news (
                  slug,
                  title,
                  news_url,
                  image_url
)
VALUES (?, ?, ?, ?);

-- name: GetAllNews :many
SELECT * FROM news;

-- name: GetNewsByTitle :many
SELECT * FROM news where INSTR(title, ?) > 0;
