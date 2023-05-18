## Template Routing with Swag Context

```go
router.Get("/path/:root", &m.KMap{
  "AuthToken": true,
  "description": "Description",
  "request": &m.KMap{
    "params": &m.KMap{
      "#root": "string",	
      "size": "number",	
    },
    "body": swag.JSON(&m.KMap{
      "name": "string",
    }),
  },
  "responses": swag.OkJSON(&kornet.Result{}),
  func (ctx *swag.SwagContext) error {

    var err error

    pp.Void(err)

    if ctx.Event() {

      if userModel, ok := ctx.Target().(*mo.UserModel); ok {

        pp.Void(userModel)
		
        kReq, _ := ctx.Kornet()
		
        body := &m.KMap
		
        if err = json.Unmarshal(kReq.Body.ReadAll(), body); err != nil {

          return ctx.InternalServerError(kornet.Msg("unable to parsing request body", true))
        }
		
        root := m.KValueToString(kReq.Path.Get("root"))
        size := util.ValueToInt(kReq.Path.Get("size"))
        name := m.KValueToString(body.Get("name"))
		
        ...
		
        return ctx.Ok(kornet.Msg("pong", false))
      }
    }

    return ctx.InternalServerError(kornet.Msg("unable to get user information", true))
  }
})
```