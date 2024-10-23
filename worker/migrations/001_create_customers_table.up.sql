create table customers (
  customer_uuid uuid primary key,
  "name" varchar(255) not null,
  current_balance decimal(8,2) not null default 0
)