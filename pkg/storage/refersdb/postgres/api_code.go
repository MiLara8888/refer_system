package postgres

import (
	"context"
	"refers_rest/pkg/sqltemp"
	storage "refers_rest/pkg/storage/refersdb"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// получить код по email
func (db *RefersDB) GetCode(ctx context.Context, email string) (*storage.RefCodeSerializer, error) {
	var (
		res = &storage.RefCodeSerializer{}
	)

	sql := `select u.id as user_id, u.ref_key, u.exp_date_ref_key
			from users u
			where u.email = :EMAIL and u.deleted=false
			`

	sqlt, err := sqltemp.Template("get_code", sql, struct{}{})
	if err != nil {
		return nil, err
	}

	names := map[string]any{
		"EMAIL": email,
	}

	sqlt, args, err := sqlx.Named(sqlt, names)
	if err != nil {
		return nil, err
	}

	sqlt = db.db.Rebind(sqlt)

	err = db.db.QueryRowContext(ctx, sqlt, args...).Scan(&res.UserID, &res.Code, &res.ExpDate)
	if err != nil {
		return nil, err
	}

	return res, nil
}




// удаление свего ref_code
func (db *RefersDB) RefCodeDelete(ctx context.Context, user *storage.UserSerializer) error {

	sql := `UPDATE users SET ref_key  = null, exp_date_ref_key = null
	WHERE email = :EMAIL and deleted = false and id = :ID
			returning ref_key, exp_date_ref_key, id`
	sqlt, err := sqltemp.Template("update_code", sql, struct{}{})
	if err != nil {
		return err
	}

	names := map[string]any{
		"EMAIL": user.Email,
		"ID":    user.ID,
	}

	sqlt, args, err := sqlx.Named(sqlt, names)
	if err != nil {
		return err
	}

	sqlt = db.db.Rebind(sqlt)

	_, err = db.db.ExecContext(ctx, sqlt, args...)
	if err != nil {
		return err
	}

	return nil
}

// обновление ref_code
func (db *RefersDB) RefCodeUpdate(ctx context.Context, user *storage.UserSerializer, exp_day int) (*storage.RefCodeSerializer, error) {
	var (
		res = &storage.RefCodeSerializer{}
	)

	uuid := uuid.NewString()

	sql := `UPDATE users SET ref_key  = :UUID, exp_date_ref_key = :EXP_REF_KEY
			WHERE email = :EMAIL and deleted = false and id = :ID
			returning ref_key, exp_date_ref_key, id`

	sqlt, err := sqltemp.Template("update_code", sql, struct{}{})
	if err != nil {
		return nil, err
	}

	names := map[string]any{
		"UUID":        uuid,
		"EMAIL":       user.Email,
		"ID":          user.ID,
		"EXP_REF_KEY": time.Now().Add(time.Duration(((24 * time.Hour) * time.Duration(exp_day)))),
	}

	sqlt, args, err := sqlx.Named(sqlt, names)
	if err != nil {
		return nil, err
	}

	sqlt = db.db.Rebind(sqlt)

	err = db.db.QueryRowContext(ctx, sqlt, args...).Scan(&res.Code, &res.ExpDate, &res.UserID)
	if err != nil {
		return nil, err
	}

	return res, nil
}
