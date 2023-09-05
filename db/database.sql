CREATE TABLE "follows" (
  "following_user_id" integer,
  "followed_user_id" integer,
  "created_at" timestamp
);

CREATE TABLE "posts" (
  "id" integer PRIMARY KEY,
  "title" varchar,
  "body" text,
  "user_id" int,
  "status" varchar,
  "created_at" timestamp
);

CREATE TABLE "products" (
  "id" integer PRIMARY KEY,
  "sku" varchar,
  "title" varchar,
  "description" text,
  "images" varchar[],
  "videos" varchar[],
  "vendors" int UNIQUE,
  "content" int UNIQUE,
  "comments" int UNIQUE,
  "review" int[],
  "tags" varchar[],
  "seo" json
);

CREATE TABLE "contents" (
  "id" integer PRIMARY KEY,
  "body" text,
  "image_gallery" varchar[]
);

CREATE TABLE "users" (
  "id" integer PRIMARY KEY,
  "username" varchar,
  "role" varchar,
  "created_at" timestamp
);

CREATE TABLE "addresses" (
  "id" integer PRIMARY KEY,
  "overal" text,
  "province" varchar,
  "city" varchar,
  "homtown" varchar,
  "street" varchar,
  "pluck" int,
  "title" varchar
);

CREATE TABLE "clients" (
  "id" integer PRIMARY KEY,
  "username" varchar,
  "name" varchar,
  "family" varchar,
  "mobile" varchar,
  "email" varchar,
  "address" integer[],
  "interests" integer[]
);

CREATE TABLE "vendors_admins" (
  "id" integer PRIMARY KEY,
  "name" varchar,
  "mobile" varchar,
  "email" varchar,
  "address" integer[]
);

CREATE TABLE "super_admin" (
  "id" integer PRIMARY KEY,
  "name" varchar,
  "mobile" varchar,
  "email" varchar,
  "address" integer[]
);

CREATE TABLE "vendors" (
  "id" integer PRIMARY KEY,
  "name" varchar,
  "address" integer[],
  "phone" varchar[],
  "created_at" timestamp,
  "updated_at" timestamp,
  "deleted_at" timestamp
);

CREATE TABLE "comments" (
  "id" integer PRIMARY KEY,
  "user_id" integer,
  "user_name" varchar,
  "product_id" integer,
  "created_at" timestamp,
  "comment" text
);

COMMENT ON COLUMN "posts"."body" IS 'Content of the post';

COMMENT ON COLUMN "contents"."body" IS 'body of content';

ALTER TABLE "posts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "vendors" ADD FOREIGN KEY ("id") REFERENCES "products" ("vendors");

ALTER TABLE "contents" ADD FOREIGN KEY ("id") REFERENCES "products" ("content");

ALTER TABLE "comments" ADD FOREIGN KEY ("id") REFERENCES "products" ("comments");
