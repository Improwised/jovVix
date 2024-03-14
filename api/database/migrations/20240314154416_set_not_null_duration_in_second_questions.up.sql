-- +migrate Up

-- set default value as 5 second
update questions
set duration_in_seconds = 5
where duration_in_seconds is null;

-- set not null
ALTER TABLE questions 
ALTER COLUMN duration_in_seconds SET NOT NULL;

