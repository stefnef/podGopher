package main

import (
	"context"
	"database/sql"
	"log"
	repository2 "podGopher/adapter/outbound/repository/postgres"
	"podGopher/adapter/outbound/repository/postgres/migration"
	"podGopher/core/domain/service"
	"podGopher/core/port/inbound"
	"podGopher/env"
	"podGopher/integration/web"

	"github.com/gin-gonic/gin"
	postgresClient "gocloud.dev/postgres"
)

func main() {
	var app = NewApp("env/.env")

	defer app.Stop()

	app.Start()
}

type App struct {
	ctx    context.Context
	db     *sql.DB
	router *gin.Engine
}

func loadEnvironment(filename string) {
	if err := env.Load(filename); err != nil {
		log.Fatal(err)
	}
}

func NewApp(environmentFilePath string) *App {
	loadEnvironment(environmentFilePath)
	var app = &App{
		context.Background(),
		nil,
		nil,
	}
	app.createSqlDb()

	app.startMigration()
	app.createWebRouter()

	return app
}

func (app *App) createWebRouter() {
	var portMap = app.createPortMap()
	app.router = web.NewRouter(portMap)
}

func (app *App) Start() {
	log.Fatal(app.router.Run(":3000"))
}

func (app *App) Stop() {
	app.db.Close()
	app.ctx.Done()
}

func (app *App) createPortMap() inbound.PortMap {
	var showRepository = repository2.NewPostgresShowRepository(app.db)
	var episodeRepository = repository2.NewPostgresEpisodeRepository(app.db)
	var createShowPort = service.NewCreateShowService(showRepository)
	var getShowPort = service.NewGetShowService(showRepository)
	var createEpisodePort = service.NewCreateEpisodeService(showRepository, episodeRepository)
	return inbound.PortMap{
		inbound.CreateShow:    createShowPort,
		inbound.GetShow:       getShowPort,
		inbound.CreateEpisode: createEpisodePort,
	}
}

func (app *App) createSqlDb() {
	dsn := migration.GetPostgresConnectionString()
	db, err := postgresClient.Open(app.ctx, dsn)
	if err != nil {
		log.Fatal(err)
	}
	app.db = db
}

func (app *App) startMigration() {
	dbMigration, err := migration.NewMigration()
	if err != nil {
		log.Fatal(err)
	}
	if err := dbMigration.Migrate(); err != nil {
		log.Printf("WARNING on migration: %s", err)
	}
}
