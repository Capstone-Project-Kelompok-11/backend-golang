package controllers

import (
	"skfw/papaya/bunny/swag"
	"skfw/papaya/koala/kornet"
	"skfw/papaya/koala/mapping"
)

func AnonymController(router swag.SwagRouterImpl) {

	router.Get("/", &mapping.KMap{
		"request":   nil,
		"responses": swag.OkJSON(&kornet.Result{}),
	}, func(ctx *swag.SwagContext) error {

		return nil
	})
}
