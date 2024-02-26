CREATE TABLE "users"(
    "id" bigserial primary key,
    "username" varchar(255) NOT NULL,
    "email" varchar(255) NOT NULL,
    "password" varchar(255) NOT NULL
)