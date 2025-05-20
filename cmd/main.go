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

	userHandler := helpers.NewBuildUser(db)
	authHandler := helpers.NewBuildAuth(db)
	forgotPasswordHandler := helpers.NewBuildForgotPassword(db)
	questionHandler := helpers.NewBuildQuestion(db)

	routes.NewRoute(authHandler, userHandler, forgotPasswordHandler, questionHandler)
}
