package main

import (
	"amalfi/pkg/router"

	"github.com/rs/zerolog/log"
	// "github.com/BurntSushi/toml"
	// "github.com/ethereum/go-ethereum/common"
	// "github.com/ethereum/go-ethereum/crypto"
	// "github.com/rs/zerolog"
	// "github.com/rs/zerolog/log"
)

// Load in all the reserves

func main() {

	router, err := router.NewRouter()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize bot")
	}

	router.Route()
	// // I need to refactor this to use freaking interfaces as a CFMM
}
