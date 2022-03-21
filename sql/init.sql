create table BEER
(
    id       serial primary key,
    name     text           not null,
    brewery  text           not null,
    country  text           not null,
    currency text           not null,
    price    numeric(10, 6) not null
);