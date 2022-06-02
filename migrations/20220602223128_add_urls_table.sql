-- +goose Up
create table urls (
    url text unique not null,
    short_url varchar(10) unique not null
);

-- +goose Down
drop table urls;