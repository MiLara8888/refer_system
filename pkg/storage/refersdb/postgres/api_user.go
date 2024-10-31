package postgres

import (
	"context"
	"errors"
	"refers_rest/pkg/sqltemp"
	storage "refers_rest/pkg/storage/refersdb"
	"time"

	"github.com/jmoiron/sqlx"
)


// Проверка токена и возвращение токена
func (db *RefersDB) IsAuch(ctx context.Context, token string) (*storage.UserSerializer, error) {
	var (
		user = &storage.UserSerializer{}
	)
	r, ok, err := storage.TokenValid(token, db.config.SecretKey)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("токен не валиден")
	}
	user.ID, err = db.GetUserID(ctx, r.Email)
	if err != nil {
		return nil, err
	}
	if user.ID == 0 {
		return nil, errors.New("токен не валиден")
	}
	user.Email = r.Email
	return user, nil
}

func (db *RefersDB) GetUserID(ctx context.Context, email string) (int, error) {
	var (
		id int
	)
	sql := `select u.id
	from users u where u.email = :EMAIL and u.deleted=false`

	sqlt, err := sqltemp.Template("user_id", sql, struct{}{})
	if err != nil {
		return id, err
	}

	names := map[string]any{
		"EMAIL": email,
	}

	sqlt, args, err := sqlx.Named(sqlt, names)
	if err != nil {
		return id, err
	}

	sqlt = db.db.Rebind(sqlt)

	err = db.db.QueryRowContext(ctx, sqlt, args...).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}

func (db *RefersDB) UserIs(ctx context.Context, email string) (bool, error) {
	var (
		res bool
	)
	sql := `select :EMAIL in (select u.email
		from users u where u.email = :EMAIL) as res`

	sqlt, err := sqltemp.Template("user_is", sql, struct{}{})
	if err != nil {
		return res, err
	}

	names := map[string]any{
		"EMAIL": email,
	}

	sqlt, args, err := sqlx.Named(sqlt, names)
	if err != nil {
		return res, err
	}

	sqlt = db.db.Rebind(sqlt)

	err = db.db.QueryRowContext(ctx, sqlt, args...).Scan(&res)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (db *RefersDB) UserSave(ctx context.Context, data *storage.UserSerializer) (int, error) {
	var (
		id          int
		paswordHesh string
		salt        string
		err         error
	)
	_, paswordHesh, salt, err = storage.GenPassword([]byte(data.Password), nil)
	if err != nil {
		return id, err
	}

	sql := `insert into users ("name", email, created_at, "password", salt)
	values (:EMAIL, :EMAIL, :CREATED, :PASSWORD, :SALT)
	returning id`

	sqlt, err := sqltemp.Template("save_user", sql, struct{}{})
	if err != nil {
		return id, err
	}

	names := map[string]any{
		"SALT":     salt,
		"PASSWORD": paswordHesh,
		"EMAIL":    data.Email,
		"CREATED":  time.Now(),
	}

	sqlt, args, err := sqlx.Named(sqlt, names)
	if err != nil {
		return id, err
	}

	sqlt = db.db.Rebind(sqlt)

	err = db.db.QueryRowContext(ctx, sqlt, args...).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}

func (db *RefersDB) GetPasswordByLogin(ctx context.Context, email string) (*storage.UserPasswordSerializer, error) {
	var (
		password = &storage.UserPasswordSerializer{}
	)

	sql := `select u.id, u."password", u.salt
			from users u
			where u.email = :EMAIL and u.deleted = false `

	sqlt, err := sqltemp.Template("get_pass", sql, struct{}{})
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

	err = db.db.QueryRowContext(ctx, sqlt, args...).Scan(&password.ID, &password.Password, &password.Salt)
	if err != nil {
		return nil, err
	}

	return password, nil
}
