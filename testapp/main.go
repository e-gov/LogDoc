/*
Rakendus testapp on LogDoc testimiseks kasutatav rakendus.
*/
package main

import (
	"context"
	"errors"

	"github.com/e-gov/LogDoc/testapp/log"
)

func main() {
	log := log.New(log.TypeServer)

	var ctx context.Context

	log.Debug().Log(ctx, "Proov")
	log.Info().Log(ctx, "Proov")
	e := errors.New("Viga")
	log.Error().WithError(e).Log(ctx, "Proov")
	log.Info().WithString("sõne", "arv1a", 10).Log(ctx, "Proov")
	log.Info().WithJSON("JSON-väärtus", "a").Log(ctx, "Proov")
	log.Info().WithJSON("JSON-väärtus", "c").Log(ctx, "Proov 2")

	log.Debug().
		WithString("method", method).
		WithString("url", httpReq.URL).
		WithString("contentLength", httpReq.ContentLength).
		Log(ctx, "request")

	log.Info().
		WithStringf("ptr", "%p", client).
		WithStringf("server", "%s:%d", hp.host, hp.port).
		Log(ctx, "connected")

}
