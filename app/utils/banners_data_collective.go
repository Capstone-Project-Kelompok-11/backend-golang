package util

import (
	"lms/app/models"
	"skfw/papaya/bunny/swag"
	m "skfw/papaya/koala/mapping"
	"skfw/papaya/koala/tools/posix"
)

func BannersDataCollective(ctx *swag.SwagContext, data []models.Banners) []m.KMapImpl {

	//var err error
	res := make([]m.KMapImpl, 0)

	//var URL *url.URL

	//if URL, err = url.Parse(ctx.BaseURL()); err != nil {
	//
	//  URL = &url.URL{}
	//}

	imagePub := posix.KPathNew("/public/image")
	//imagePub := posix.KPathNew("/api/v1/public/image")

	for _, banner := range data {

		if banner.Src != "" {

			//URL.Path = imagePub.Copy().JoinStr(banner.Src)
			//URL.RawPath = URL.Path
			//
			//banner.Src = URL.String()

			banner.Src = imagePub.Copy().JoinStr(banner.Src)
		}

		mm := &m.KMap{
			"id":  banner.ID,
			"alt": banner.Alt,
			"src": banner.Src,
		}

		res = append(res, mm)
	}

	return res
}
