package main

import (
	"bytes"
	"github.com/signintech/gopdf"
	"github.com/signintech/pdft"
	"os"
)

func main() {

	var err error
	var pt pdft.PDFt

	if err = pt.Open("templates/documents/certificates/cert.pdf"); err != nil {
		panic(err)
	}

	if err = pt.AddFont("arial", "assets/fonts/arial.ttf"); err != nil {
		panic(err)
	}

	if err = pt.SetFont("arial", "", 24); err != nil {
		panic(err)
	}

	if err = pt.Insert("SAMSUDIN", 1, 0, 0, -1, -1, gopdf.Middle); err != nil {
		panic(err)
	}

	buff := bytes.NewBuffer(nil)

	if err = pt.SaveTo(buff); err != nil {
		panic(err)
	}

	if err = os.WriteFile("exam.pdf", buff.Bytes(), 0644); err != nil {
		return
	}
}
