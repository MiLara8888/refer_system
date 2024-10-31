/* для работы с гусём */
package main

import (
	"log"
	"os"
	"runtime"

	"github.com/urfave/cli"
	"refers_rest/pkg/migrator"
	settings "refers_rest/pkg/settings"

	_ "github.com/jackc/pgx"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var FLAGS = []cli.Flag{
	cli.StringFlag{
		Name:  "dir,d",
		Usage: "migrate directory",
	},
	&cli.StringFlag{
		Name:     "schema, s",
		Usage:    "имя схемы",
		Required: true,
	},
}
var (
	err error
	// настройка подключения
	config *settings.Config
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)

}

func main() {
	app := &cli.App{
		Name:  "news migrator db",
		Usage: "news migrator db",
		Commands: []cli.Command{
			{
				Name:  "schema",
				Usage: "откат изменений в db",
				Flags: FLAGS,
				Action: func(c *cli.Context) error {
					dir := c.String("dir")
					schema := c.String("schema")
					err := migrator.Migrate(func(o *migrator.Options) {
						o.Command = migrator.CMD_CREATE_SCHEMA
						o.MigSchema = schema
						o.Dir = dir
					})
					if err != nil {
						log.Println(err)
						return err
					}
					return nil
				},
			},
			{
				Name:  "down",
				Usage: "откат изменений в db",
				Flags: FLAGS,
				Action: func(c *cli.Context) error {
					dir := c.String("dir")
					schema := c.String("schema")
					err := migrator.Migrate(func(o *migrator.Options) {
						o.Command = migrator.CMD_DOWN
						o.MigSchema = schema
						o.Dir = dir
					})
					if err != nil {
						log.Println(err)
						return err
					}
					return nil
				},
			},
			{
				Name:  "up",
				Usage: "внести посл. изменения",
				Flags: FLAGS,
				Action: func(c *cli.Context) error {
					dir := c.String("dir")
					schema := c.String("schema")
					err := migrator.Migrate(func(o *migrator.Options) {
						o.Command = migrator.CMD_UP
						o.MigSchema = schema
						o.Dir = dir
					})
					if err != nil {
						log.Println(err)
						return err
					}
					return nil
				},
			},
			{
				Name:  "redo",
				Usage: "",
				Flags: FLAGS,
				Action: func(c *cli.Context) error {
					dir := c.String("dir")
					schema := c.String("schema")
					err := migrator.Migrate(func(o *migrator.Options) {
						o.Command = migrator.CMD_REDO
						o.MigSchema = schema
						o.Dir = dir
					})
					if err != nil {
						log.Println(err)
						return err
					}
					return nil
				},
			},
			{
				Name:  "reset",
				Usage: "",
				Flags: FLAGS,
				Action: func(c *cli.Context) error {
					dir := c.String("dir")
					schema := c.String("schema")
					err := migrator.Migrate(func(o *migrator.Options) {
						o.Command = migrator.CMD_RESET
						o.MigSchema = schema
						o.Dir = dir
					})
					if err != nil {
						log.Println(err)
						return err
					}
					return nil
				},
			},
			{
				Name:  "status",
				Usage: "",
				Flags: FLAGS,
				Action: func(c *cli.Context) error {
					dir := c.String("dir")
					schema := c.String("schema")
					err := migrator.Migrate(func(o *migrator.Options) {
						o.Command = migrator.CMD_STATUS
						o.MigSchema = schema
						o.Dir = dir
					})
					if err != nil {
						log.Println(err)
						return err
					}
					return nil
				},
			},
			{
				Name:  "version",
				Usage: "",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "dir,d",
						Usage: "migrate directory",
					},
				},
				Action: func(c *cli.Context) error {
					dir := c.String("dir")
					schema := c.String("schema")
					err := migrator.Migrate(func(o *migrator.Options) {
						o.Command = migrator.CMD_VERSION
						o.MigSchema = schema
						o.Dir = dir
					})
					if err != nil {
						log.Println(err)
						return err
					}
					return nil
				},
			},
			{
				Name:  "create",
				Usage: "создание миграций",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "dir,d",
						Usage: "migrate directory",
					},
					&cli.StringFlag{
						Name:     "name, n",
						Usage:    "имя файла",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "schema, s",
						Usage:    "имя схемы",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					dir := c.String("dir")
					schema := c.String("schema")
					err := migrator.Migrate(func(o *migrator.Options) {
						o.Command = migrator.CMD_CREATE
						o.Dir = dir
						o.MigSchema = schema
						o.Args = []string{c.String("name"), "sql"}
					})
					if err != nil {
						log.Println(err)
						return err
					}
					return nil
				},
			},
		},
	}
	app.Run(os.Args)
}
