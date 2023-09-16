/**
  This is the SQL script that will be used to initialize the database schema.
  We will evaluate you based on how well you design your database.
  1. How you design the tables.
  2. How you choose the data types and keys.
  3. How you name the fields.
  In this assignment we will use PostgreSQL as the database.
  */

CREATE TABLE sawitpro_user (
  id serial PRIMARY KEY,
  full_name text NOT NULL,
  phone_number text NOT NULL,
  "password" text not null,
  created_time timestamp NOT NULL default now(),
  updated_time timestamp,
  successful_login_count int not null default 0,
  CONSTRAINT sawitpro_user_phone_number_uniquekey UNIQUE (phone_number)
);

INSERT INTO sawitpro_user (full_name, phone_number, "password") VALUES ('name1', '+6281234567890', 'password1');
INSERT INTO sawitpro_user (full_name, phone_number, "password") VALUES ('name2', '+6289876543210', 'password2');