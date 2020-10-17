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
	l := log.New(log.TypeServer)

	var ctx context.Context

	l.Debug().Log(ctx, "Proov")
	l.Info().Log(ctx, "Proov")
	e := errors.New("Viga")
	l.Error().WithError(e).Log(ctx, "Proov")
	l.Info().WithString("sõne", "arv1", 10).Log(ctx, "Proov")
	l.Info().WithJSON("JSON-väärtus", "a").Log(ctx, "Proov")
	l.Info().WithJSON("JSON-väärtus", "c").Log(ctx, "Proov 2")

}
