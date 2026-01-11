package main

import (
	"github.com/migueldor/gator/internal/config"
	"github.com/migueldor/gator/internal/database"
)

type state struct {
	db     *database.Queries
	config *config.Config
}
