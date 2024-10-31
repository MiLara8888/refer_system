package migrator


import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	settings "refers_rest/pkg/settings"

	_ "github.com/jackc/pgx"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

type CommandType string

const (
	CMD_CREATE_SCHEMA CommandType = "schema"
	CMD_CREATE        CommandType = "create"
	CMD_VERSION       CommandType = "version"
	CMD_STATUS        CommandType = "status"
	CMD_RESET         CommandType = "reset"
	CMD_REDO          CommandType = "redo"
	CMD_UP            CommandType = "up"
	CMD_DOWN          CommandType = "down"
)

type Options struct {
	*settings.Config
	Command   CommandType
	Dir       string
	Args      []string
	EnvFile   string
	MigSchema string
}

func load(dir string) ([]string, error) {
	ret := make([]string, 0, 5)
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && path != dir {
			ret = append(ret, info.Name())
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return ret, nil
}


// генератор базы данных
func Migrate(opts ...func(*Options)) error {
	var (
		err     error
		options = &Options{}
		config  *settings.Config
	)

	for _, fn := range opts {
		fn(options)
	}

	if options.Dir == "" {
		return errors.New("укажите директорию с файлами")
	}

	if options.MigSchema == "" {
		return errors.New("не указана схема базы данных")
	}

	switch options.EnvFile {
	case "":
		config, err = settings.InitEnv()
		if err != nil {
			return err
		}
	default:
		config, err = settings.InitEnv(options.EnvFile)
		if err != nil {
			return err
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10*time.Second))
	defer cancel()

	conn, err := sql.Open("pgx", config.DB.UrlPostgres())
	if err != nil {
		return err
	}

	defer func() {
		err := conn.Close()
		if err != nil {
			log.Panic(err)
		}
	}()

	err = conn.PingContext(ctx)
	if err != nil {
		return err
	}
	dir := path.Join(options.Dir, options.MigSchema)
	os.Setenv("SCHEMA_DB", options.MigSchema)

	log.Println(options.MigSchema)
	if options.Command == CMD_CREATE_SCHEMA {
		_, err = conn.ExecContext(context.Background(), fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS "%s"`, strings.ToLower(options.MigSchema)))
		if err != nil {
			return err
		}

		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.Mkdir(dir, os.ModePerm); err != nil {
				log.Fatal(err)
			}
		}

		return nil
	}
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	_, err = conn.ExecContext(context.Background(), fmt.Sprintf(`SET search_path TO "%s"`, strings.ToLower(options.MigSchema)))
	if err != nil {
		return err
	}

	err = goose.RunContext(ctx, string(options.Command), conn, dir, options.Args...)
	if err != nil {
		return err
	}

	return nil
}
