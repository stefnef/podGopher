package main

import (
	"context"
	"database/sql"
	"log"
	repositoryDistribution "podGopher/adapter/outbound/repository/postgres/distribution"
	repositoryEpisode "podGopher/adapter/outbound/repository/postgres/episode"
	"podGopher/adapter/outbound/repository/postgres/migration"
	repositoryRSS "podGopher/adapter/outbound/repository/postgres/rss"
	repositoryShow "podGopher/adapter/outbound/repository/postgres/show"
	repositoryUser "podGopher/adapter/outbound/repository/postgres/user"
	"podGopher/core/domain/service/distribution"
	"podGopher/core/domain/service/episode"
	"podGopher/core/domain/service/rss"
	"podGopher/core/domain/service/show"
	"podGopher/core/domain/service/user"
	"podGopher/core/port/inbound"
	"podGopher/env"
	"podGopher/integration/web"
	"podGopher/integration/web/auth"

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
	var adminAuth = app.createAuth()
	app.router = web.NewRouter(portMap, adminAuth)
}

func (app *App) createAuth() auth.AdminAuth {
	return auth.NewAdminAuth(env.AdminUser.GetValue(), env.AdminPassword.GetValue())
}

func (app *App) Start() {
	log.Fatal(app.router.Run(":3000"))
}

func (app *App) Stop() {
	_ = app.db.Close()
	app.ctx.Done()
}

func (app *App) createPortMap() inbound.PortMap {
	var showRepository = repositoryShow.NewPostgresShowRepository(app.db)
	var episodeRepository = repositoryEpisode.NewPostgresEpisodeRepository(app.db)
	var distributionRepository = repositoryDistribution.NewPostgresDistributionRepository(app.db)
	var rssRepository = repositoryRSS.NewPostgresRSSRepository(app.db)
	var userRepository = repositoryUser.NewPostgresUserRepository(app.db)

	var createShowPort = show.NewCreateShowService(showRepository)
	var getShowPort = show.NewGetShowService(showRepository)
	var createEpisodePort = episode.NewCreateEpisodeService(showRepository, episodeRepository)
	var getEpisodePort = episode.NewGetEpisodeService(showRepository, episodeRepository)
	var createDistributionPort = distribution.NewCreateDistributionService(showRepository, distributionRepository)
	var getDistributionPort = distribution.NewGetDistributionService(showRepository, distributionRepository)
	var updateDistributionPort = distribution.NewUpdateDistributionService(showRepository, episodeRepository, distributionRepository, distributionRepository)
	var getRSSPort = rss.NewGetRSSService(rssRepository)
	var createUserPort = user.NewCreateUserService(showRepository, userRepository)

	return inbound.PortMap{
		inbound.CreateShow:         createShowPort,
		inbound.GetShow:            getShowPort,
		inbound.GetAllShows:        getShowPort,
		inbound.CreateEpisode:      createEpisodePort,
		inbound.GetEpisode:         getEpisodePort,
		inbound.CreateDistribution: createDistributionPort,
		inbound.GetDistribution:    getDistributionPort,
		inbound.UpdateDistribution: updateDistributionPort,
		inbound.GetRSS:             getRSSPort,
		inbound.CreateUser:         createUserPort,
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
