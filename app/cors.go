package app

import (
  "skfw/papaya"
  "skfw/papaya/pigeon/cors"
)

func ManageControlResourceShared(pn papaya.NetImpl) error {

  manageConsumers, err := cors.ManageConsumersNew()
  if err != nil {

    return err
  }

  manageConsumers.Grant("http://localhost")
  manageConsumers.Grant("http://localhost:5173") // vite - react
  manageConsumers.Grant("http://localhost:8000")
  manageConsumers.Grant("https://academy.skfw.net") // secure deploy
  manageConsumers.Grant("https://skfw.net")         // secure deploy
  manageConsumers.Grant("https://frontend-react-git-staging-academade.vercel.app")
  manageConsumers.Grant("https://academade.vercel.app")

  pn.Use(cors.MakeMiddlewareForManageConsumers(manageConsumers))

  return nil
}
