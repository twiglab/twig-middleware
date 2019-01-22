/*
expamle:

package main

import (
	"time"

	"github.com/twiglab/twig"
	"github.com/twiglab/twig-middleware/ratelimiter"
	"golang.org/x/time/rate"
)

func main() {
	web := twig.TODO()
	limiter := rate.NewLimiter(rate.Every(10), 100)

	web.Pre(ratelimiter.New(limiter))

	web.Config().
		Get("/hello", twig.HelloTwig).
		Done()

	web.Start()

	twig.Signal(twig.Graceful(web, 10*time.Second))
}

*/
package ratelimiter

import (
	"github.com/twiglab/twig"
	"github.com/twiglab/twig/middleware"
)

type Allower interface {
	Allow() bool
}

type Conifg struct {
	Skipper middleware.Skipper
	Allower Allower
}

var DefaultConfig = Conifg{
	Skipper: middleware.DefaultSkipper,
}

func New(allower Allower) twig.MiddlewareFunc {
	config := DefaultConfig
	config.Allower = allower
	return NewWithConifg(config)
}

func NewWithConifg(config Conifg) twig.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = middleware.DefaultSkipper
	}
	if config.Allower == nil {
		panic("Limiter is nil")
	}

	return func(next twig.HandlerFunc) twig.HandlerFunc {
		return func(c twig.Ctx) error {
			if config.Skipper(c) || config.Allower.Allow() {
				return next(c)
			}
			return twig.ErrTooManyRequests
		}
	}
}
