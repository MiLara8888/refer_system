package postgres

import (
	"context"
	"fmt"
	"refers_rest/pkg/sqltemp"
	storage "refers_rest/pkg/storage/refersdb"
	"time"

	"github.com/jmoiron/sqlx"
)

// созание подписчика на код
func (db *RefersDB) CreateReferals(ctx context.Context, code, email string) (*storage.RefSubsirbeSerializer, error) {
	var (
		res   = &storage.RefSubsirbeSerializer{}
		refer = &storage.RefCodeSerializer{}
	)

	//получаю рефера
	sql := `select u.id as user_id, u.ref_key as code
			from users u
			where u.ref_key = :CODE`

	sqlt, err := sqltemp.Template("get_refer", sql, struct{}{})
	if err != nil {
		return nil, err
	}

	names := map[string]any{
		"CODE": code,
	}

	sqlt, args, err := sqlx.Named(sqlt, names)
	if err != nil {
		return nil, err
	}

	sqlt = db.db.Rebind(sqlt)

	err = db.db.QueryRowContext(ctx, sqlt, args...).Scan(&refer.UserID, &refer.Code)
	if err != nil {
		return nil, err
	}

	sql = `insert into referals ( email, created_at, ref_key, user_id)
	values (:EMAIL, :CREATED, :CODE, :USER)
	returning id, email, ref_key`

	sqlt, err = sqltemp.Template("save_referals", sql, struct{}{})
	if err != nil {
		return nil, err
	}

	names = map[string]any{
		"EMAIL":   email,
		"CREATED": time.Now(),
		"CODE":    refer.Code,
		"USER":    refer.UserID,
	}

	sqlt, args, err = sqlx.Named(sqlt, names)
	if err != nil {
		return nil, err
	}

	sqlt = db.db.Rebind(sqlt)

	err = db.db.QueryRowContext(ctx, sqlt, args...).Scan(&res.ID, &res.Email, &res.RefKey)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	res.UserId = *refer

	return res, nil
}

func (db *RefersDB) ReferalsList(ctx context.Context, id int) ([]*storage.RefSubsirbeSerializer, error) {
	var (
		res = []*storage.RefSubsirbeSerializer{}
	)

	names := map[string]any{
		"ID": id,
	}

	sql := `select r.id, r.email, r.ref_key
	from referals r
	where r.user_id = :ID`

	sqlt, err := sqltemp.Template("referal_list", sql, nil)
	if err != nil {
		return nil, err
	}

	sqlt, args, err := sqlx.Named(sqlt, names)
	if err != nil {
		return nil, err
	}

	sqlt, args, err = sqlx.In(sqlt, args...)
	if err != nil {
		return nil, err
	}
	sqlt = db.db.Rebind(sqlt)

	rows, err := db.db.QueryxContext(ctx, sqlt, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		s := &storage.RefSubsirbeSerializer{}
		err = rows.Scan(&s.ID, &s.Email, &s.RefKey)
		if err != nil {
			return nil, err
		}
		res = append(res, s)
	}

	return res, nil
}
