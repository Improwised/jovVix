-- +migrate Up
alter table users 
add column img_key text default '';