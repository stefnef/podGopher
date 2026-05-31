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

type repositoryOutAdapters struct {
	showRepository         *repositoryShow.PostgresShowOutAdapter
	episodeRepository      *repositoryEpisode.PostgresEpisodeOutAdapter
	distributionRepository *repositoryDistribution.PostgresDistributionOutAdapter
	rssRepository          *repositoryRSS.PostgresRssOutAdapter
	userRepository         *repositoryUser.PostgresUserOutAdapter
}

func initRepositories(app *App) repositoryOutAdapters {
	return repositoryOutAdapters{
		showRepository:         repositoryShow.NewPostgresShowRepository(app.db),
		episodeRepository:      repositoryEpisode.NewPostgresEpisodeRepository(app.db),
		distributionRepository: repositoryDistribution.NewPostgresDistributionRepository(app.db),
		rssRepository:          repositoryRSS.NewPostgresRSSRepository(app.db),
		userRepository:         repositoryUser.NewPostgresUserRepository(app.db),
	}
}

type servicePorts struct {
	createShowPort         *show.CreateShowService
	getShowPort            *show.GetShowService
	createEpisodePort      *episode.CreateEpisodeService
	getEpisodePort         *episode.GetEpisodeService
	createDistributionPort *distribution.CreateDistributionService
	getDistributionPort    *distribution.GetDistributionService
	updateDistributionPort *distribution.UpdateDistributionService
	getRSSPort             *rss.GetRSSService
	createUserPort         *user.CreateUserService
	assignUserPort         *user.AssignUserService
}

func (s *servicePorts) initShows(repos repositoryOutAdapters) {
	s.createShowPort = show.NewCreateShowService(repos.showRepository)
	s.getShowPort = show.NewGetShowService(repos.showRepository)
}

func (s *servicePorts) initEpisodes(repos repositoryOutAdapters) {
	s.createEpisodePort = episode.NewCreateEpisodeService(repos.showRepository, repos.episodeRepository)
	s.getEpisodePort = episode.NewGetEpisodeService(repos.showRepository, repos.episodeRepository)
}

func (s *servicePorts) initDistributions(repos repositoryOutAdapters) {
	s.createDistributionPort = distribution.NewCreateDistributionService(repos.showRepository, repos.distributionRepository)
	s.getDistributionPort = distribution.NewGetDistributionService(repos.showRepository, repos.distributionRepository)
	s.updateDistributionPort = distribution.NewUpdateDistributionService(repos.showRepository, repos.episodeRepository, repos.distributionRepository, repos.distributionRepository)
}

func (s *servicePorts) initRSS(repos repositoryOutAdapters) {
	s.getRSSPort = rss.NewGetRSSService(repos.rssRepository)
}

func (s *servicePorts) initUsers(repos repositoryOutAdapters) {
	s.createUserPort = user.NewCreateUserService(repos.userRepository, nil)
	s.assignUserPort = user.NewAssignUserService(repos.showRepository, repos.userRepository, repos.userRepository)
}

func initServices(app *App) *servicePorts {
	var repos = initRepositories(app)

	var services = new(servicePorts)
	services.initShows(repos)
	services.initEpisodes(repos)
	services.initDistributions(repos)
	services.initRSS(repos)
	services.initUsers(repos)
	return services
}

func (app *App) createPortMap() inbound.PortMap {
	var services = initServices(app)

	return inbound.PortMap{
		inbound.CreateShow:         services.createShowPort,
		inbound.GetShow:            services.getShowPort,
		inbound.GetAllShows:        services.getShowPort,
		inbound.CreateEpisode:      services.createEpisodePort,
		inbound.GetEpisode:         services.getEpisodePort,
		inbound.CreateDistribution: services.createDistributionPort,
		inbound.GetDistribution:    services.getDistributionPort,
		inbound.UpdateDistribution: services.updateDistributionPort,
		inbound.GetRSS:             services.getRSSPort,
		inbound.CreateUser:         services.createUserPort,
		inbound.AssignUser:         services.assignUserPort,
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
