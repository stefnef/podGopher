package main

import (
	"context"
	"database/sql"
	"log"
	repositoryDistribution "podGopher/adapter/outbound/repository/postgres/distribution"
	repositoryEpisode "podGopher/adapter/outbound/repository/postgres/episode"
	"podGopher/adapter/outbound/repository/postgres/migration"
	repositoryShow "podGopher/adapter/outbound/repository/postgres/show"
	"podGopher/core/domain/service/distribution"
	"podGopher/core/domain/service/episode"
	"podGopher/core/domain/service/show"
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
	var showRepository = repositoryShow.NewPostgresShowRepository(app.db)
	var episodeRepository = repositoryEpisode.NewPostgresEpisodeRepository(app.db)
	var distributionRepository = repositoryDistribution.NewPostgresDistributionRepository(app.db)

	var createShowPort = show.NewCreateShowService(showRepository)
	var getShowPort = show.NewGetShowService(showRepository)
	var createEpisodePort = episode.NewCreateEpisodeService(showRepository, episodeRepository)
	var getEpisodePort = episode.NewGetEpisodeService(showRepository, episodeRepository)
	var createDistributionPort = distribution.NewCreateDistributionService(showRepository, distributionRepository)

	return inbound.PortMap{
		inbound.CreateShow:         createShowPort,
		inbound.GetShow:            getShowPort,
		inbound.CreateEpisode:      createEpisodePort,
		inbound.GetEpisode:         getEpisodePort,
		inbound.CreateDistribution: createDistributionPort,
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
