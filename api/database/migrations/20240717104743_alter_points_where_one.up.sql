-- +migrate Up

-- alter points column for old data to make 0 where points are 1
-- as in old functionality it was putting 1 where points were not provided, in new we want 0

UPDATE questions
SET "points" = 0 
where "points" = 1;