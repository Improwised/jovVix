-- +migrate Up

-- fill for old data columns

UPDATE questions
SET "type" = (
    CASE WHEN answers::jsonb @> '[1,2,3,4]'::jsonb OR answers::jsonb @> '[1,2,3,4,5]'::jsonb
         THEN 2
         ELSE 1
    END
)
where "type" is null;