-- name: CreateBid :one
INSERT INTO bids (product_id, bidder_id, bid_amount)
VALUES ($1, $2, $3)
RETURNING *;

-- name: getBidsByProductId :many
SELECT * FROM bids
WHERE product_id = $1
ORDER BY created_at DESC;

-- name: GetHighestBidByProductId :one
SELECT *
FROM bids
WHERE product_id = $1
ORDER BY created_at DESC
LIMIT 1;
