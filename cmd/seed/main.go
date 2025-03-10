package main

import (
	"context"
	"log"
	"os"
	"runtime/debug"
	"time"

	"github.com/iota-uz/iota-sdk/modules"
	"github.com/iota-uz/iota-sdk/modules/bichat"
	"github.com/iota-uz/iota-sdk/modules/core"
	"github.com/iota-uz/iota-sdk/modules/core/domain/aggregates/user"
	coreseed "github.com/iota-uz/iota-sdk/modules/core/seed"
	"github.com/iota-uz/iota-sdk/modules/crm"
	"github.com/iota-uz/iota-sdk/modules/finance"
	"github.com/iota-uz/iota-sdk/modules/hrm"
	"github.com/iota-uz/iota-sdk/modules/warehouse"
	"github.com/iota-uz/iota-sdk/modules/website"
	"github.com/iota-uz/iota-sdk/pkg/application"
	"github.com/iota-uz/iota-sdk/pkg/composables"
	"github.com/iota-uz/iota-sdk/pkg/configuration"
	"github.com/iota-uz/iota-sdk/pkg/eventbus"

	"github.com/jackc/pgx/v5/pgxpool"
)

func pgxPool() *pgxpool.Pool {
	conf := configuration.Use()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	pool, err := pgxpool.New(ctx, conf.Database.Opts)
	if err != nil {
		panicWithStack(err)
	}
	return pool
}

func panicWithStack(err error) {
	errorWithStack := string(debug.Stack()) + "\n\nError: " + err.Error()
	panic(errorWithStack)
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			configuration.Use().Unload()
			log.Println(r)
			os.Exit(1)
		}
	}()

	conf := configuration.Use()
	ctx := context.Background()
	pool := pgxPool()
	app := application.New(pool, eventbus.NewEventPublisher(conf.Logger()))
	if err := modules.Load(app, modules.BuiltInModules...); err != nil {
		panicWithStack(err)
	}
	app.RegisterNavItems(core.NavItems...)
	app.RegisterNavItems(bichat.NavItems...)
	app.RegisterNavItems(hrm.NavItems...)
	app.RegisterNavItems(finance.NavItems...)
	app.RegisterNavItems(warehouse.NavItems...)
	app.RegisterNavItems(crm.NavItems...)
	app.RegisterNavItems(website.NavItems...)
	tx, err := pool.Begin(ctx)
	if err != nil {
		panicWithStack(err)
	}
	seeder := application.NewSeeder()
	seeder.Register(
		coreseed.CreateCurrencies,
		coreseed.CreatePermissions,
		coreseed.UserSeedFunc("test@gmail.com", "TestPass123!", user.UILanguageEN),
	)
	if err := seeder.Seed(composables.WithTx(ctx, tx), app); err != nil {
		panicWithStack(err)
	}
	if err := tx.Commit(ctx); err != nil {
		panicWithStack(err)
	}
}
