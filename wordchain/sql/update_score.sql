UPDATE scores
SET
    score = $2
WHERE id = $1;
