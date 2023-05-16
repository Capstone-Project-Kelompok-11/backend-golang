package main

import (
	"lms/app"
	"skfw/papaya"
)

func main() {

	pn := papaya.NetNew()
	logger := pn.Logger()

	if err := app.App(pn); err != nil {

		logger.Error(err)
	}

	if err := pn.Close(); err != nil {

		logger.Error(err)
	}

	logger.Log("Shutdown ...")
}
