package helpers

import "github.com/rs/zerolog/log"

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatal().Err(err).Msg(msg)
	}
}
