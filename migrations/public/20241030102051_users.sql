-- +goose ENVSUB ON
-- +goose Up
-- +goose StatementBegin
CREATE TABLE "public"."users" (
	"id" serial4 NOT NULL,
	"name" varchar(255) NOT NULL,
	"email" varchar(255) NOT NULL,
	"created_at" timestamp NULL,
	"deleted" bool not NULL default false,
	"ref_key" uuid NULL,
	"exp_date_ref_key" timestamp NULL,
	CONSTRAINT "users_pkey" PRIMARY KEY ("id"),
    "password" varchar(500) not NULL,
    "salt" varchar(500) not NULL
);
CREATE  INDEX "idx_id_email" ON "public"."users" USING btree ("id", "email");
CREATE UNIQUE INDEX "idx_uuid" ON "public"."users" USING btree ("ref_key");

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table "public"."users";
-- +goose StatementEnd
-- +goose ENVSUB OFF
