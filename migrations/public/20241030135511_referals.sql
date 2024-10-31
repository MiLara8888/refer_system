-- +goose ENVSUB ON
-- +goose Up
-- +goose StatementBegin
CREATE TABLE "public"."referals" (
	"id" serial4 NOT NULL,
	"email" varchar(255) NOT NULL,
	"created_at" timestamp NULL,
	"deleted" bool not NULL default false,
	"ref_key" uuid NULL,
    "user_id" int8 NULL,
    CONSTRAINT "user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users"("id") ON DELETE RESTRICT
);
CREATE  INDEX "idx_id_email_ref" ON "public"."referals"  ("id", "email");
CREATE UNIQUE INDEX "idx_uuid_ref" ON "public"."referals"  ("ref_key", "email") ;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table "public"."referals";
-- +goose StatementEnd
-- +goose ENVSUB OFF
