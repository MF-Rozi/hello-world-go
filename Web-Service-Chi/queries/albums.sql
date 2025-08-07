-- name: GetAlbums :many
SELECT id, title, artist, price
FROM albums
ORDER BY id
LIMIT $1
OFFSET
    $2;

-- name: GetAlbumByID :one
SELECT id, title, artist, price FROM albums WHERE id = $1;

-- name: CreateAlbum :one
INSERT INTO
    albums (title, artist, price)
VALUES ($1, $2, $3) RETURNING id,
    title,
    artist,
    price;

-- name: UpdateAlbum :one
UPDATE albums
SET
    title = $2,
    artist = $3,
    price = $4
WHERE
    id = $1 RETURNING id,
    title,
    artist,
    price;

-- name: DeleteAlbum :exec
DELETE FROM albums WHERE id = $1;