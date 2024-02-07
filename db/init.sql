CREATE table clients (
  id int primary key,
  c_limit int not null,
  balance int not null  
);

ALTER TABLE clients ADD CONSTRAINT check_balance_positive CHECK (abs(balance) <= c_limit);

CREATE table transactions (
  id int auto_increment primary key, 
  client_id int not null,
  value int not null,
  type varchar(1) not null,
  description varchar(10) not null,
  create_at TIMESTAMP not null default NOW()
);

insert into clients (id, c_limit, balance) 
values
  (1,100000, 0),
  (2,80000, 0),
  (3,1000000, 0),
  (4,10000000, 0),
  (5,500000, 0);

