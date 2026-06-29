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