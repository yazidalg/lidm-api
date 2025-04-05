package main

import (
	"github.com/yazidalg/lidm_backend/internal/config"
	"github.com/yazidalg/lidm_backend/internal/helpers"
	"github.com/yazidalg/lidm_backend/internal/routes"
)

func main() {
	config.LoadEnv()
	db := config.ConnectDB()
	config.MigrateDb(db)

	auth := helpers.NewBuildAuth(db)

	routes.NewRoute(auth)
}
