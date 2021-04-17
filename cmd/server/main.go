package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/wufe/polo/pkg"
	"github.com/wufe/polo/pkg/background"
	"github.com/wufe/polo/pkg/background/queues"
	"github.com/wufe/polo/pkg/http/proxy"
	"github.com/wufe/polo/pkg/http/rest"
	"github.com/wufe/polo/pkg/http/routing"
	"github.com/wufe/polo/pkg/services"
	"github.com/wufe/polo/pkg/storage"
	"github.com/wufe/polo/pkg/utils"
)

func main() {
	environment := utils.DetectEnvironment()

	var mutexBuilder utils.MutexBuilder
	mutexBuilder = func() utils.RWLocker {
		return utils.GetMutex(environment)
	}

	// Configuration (.yml)
	configuration, applications := storage.LoadConfigurations(environment, mutexBuilder)

	// Instance
	existingInstance, _ := storage.DetectInstance(environment)
	if existingInstance == nil {
		storage.NewInstance(fmt.Sprint(configuration.Global.Port)).Persist(environment)
	} else {
		log.Infof("Detected existing instance on host %s", existingInstance.Host)
		return
	}

	// Storage
	folder := environment.GetExecutableFolder()
	database := storage.NewDB(folder)
	appStorage := storage.NewApplication(environment)
	sesStorage := storage.NewSession(database, environment)

	mediator := background.NewMediator(
		queues.NewSessionBuild(),
		queues.NewSessionDestroy(),
		queues.NewSessionFilesystem(),
		queues.NewSessionCleanup(),
		queues.NewSessionStart(),
		queues.NewSessionHealthCheck(),
		queues.NewApplicationInit(),
		queues.NewApplicationFetch(),
	)

	// Workers
	background.NewSessionBuildWorker(&configuration.Global, appStorage, sesStorage, mediator, mutexBuilder)
	background.NewSessionStartWorker(sesStorage, mediator)
	background.NewSessionCleanWorker(sesStorage, mediator)
	background.NewSessionFilesystemWorker(mediator)
	background.NewSessionDestroyWorker(mediator)
	background.NewSessionHealthcheckWorker(mediator)
	background.NewApplicationInitWorker(&configuration.Global, mediator)
	background.NewApplicationFetchWorker(sesStorage, mediator)

	// Services
	staticService := services.NewStaticService(environment)
	queryService := services.NewQueryService(environment, sesStorage, appStorage)
	requestService := services.NewRequestService(environment, sesStorage, appStorage, mediator)

	// HTTP
	proxy := proxy.NewHandler(environment)
	routing := routing.NewHandler(environment, proxy, sesStorage, appStorage, queryService, requestService, staticService)
	rest := rest.NewHandler(environment, staticService, routing, proxy, queryService, requestService)

	// Startup
	pkg.NewStartup(
		configuration,
		applications,
		rest,
		staticService,
		appStorage,
		sesStorage,
		mediator,
		mutexBuilder,
	).Start(&pkg.StartupOptions{
		WatchApplications: true,
		LoadSessionHelper: true,
		StartServer:       true,
	})

}
