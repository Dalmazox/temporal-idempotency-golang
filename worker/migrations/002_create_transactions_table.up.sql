create table transactions (
  transaction_uuid uuid primary key,
  customer_uuid uuid references customers(customer_uuid) not null,
  amount decimal(8,2) not null,
  created_at timestamptz not null default current_timestamp
)