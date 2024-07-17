-- +migrate Down

-- to again make points to 1 where points are 0, in case the up migration make changes were not desired

UPDATE questions
SET "points" = 1 
where "points" = 0;