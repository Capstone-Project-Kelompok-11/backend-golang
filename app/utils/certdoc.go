package util

import (
  "bytes"
  "github.com/h2non/bimg"
  "github.com/signintech/gopdf"
  "os"
)

func TextCenter(pdf *gopdf.GoPdf, text string, pageSize gopdf.Rect) error {

  var err error
  var width float64
  if width, err = pdf.MeasureTextWidth(text); err != nil {
    return err
  }
  pdf.SetX(pageSize.W/2 - width/2)
  if err = pdf.Text(text); err != nil {
    return err
  }
  return nil
}

func TextMiddle(pdf *gopdf.GoPdf, text string, pageSize gopdf.Rect) error {

  var err error
  var width float64
  if width, err = pdf.MeasureTextWidth(text); err != nil {
    return err
  }
  pdf.SetXY(pageSize.W/2-width/2, pageSize.H/2)
  if err = pdf.Text(text); err != nil {
    return err
  }
  return nil
}

func GenerateCertificateInCaches(title string, context string) (string, error) {

  var err error

  arialFontPath := "assets/fonts/arial.ttf"
  templateCertificateDocument := "templates/documents/certificates/cert.pdf"
  cacheDir := "assets/public/caches"

  pdf := gopdf.GoPdf{}
  config := gopdf.Config{PageSize: *gopdf.PageSizeA4Landscape}
  pdf.Start(config)
  pdf.AddPage()

  tpl1 := pdf.ImportPage(templateCertificateDocument, 1, "/MediaBox")

  pdf.UseImportedTemplate(tpl1, 0, 0, config.PageSize.W, config.PageSize.H)

  if err = pdf.AddTTFFont("arial", arialFontPath); err != nil {

    return "", err
  }

  if err = pdf.SetFont("arial", "", 18); err != nil {

    return "", err
  }

  pdf.SetTextColor(0, 0, 0)

  pdf.SetY(164)

  if err = TextCenter(&pdf, title, config.PageSize); err != nil {

    return "", err
  }

  if err = TextMiddle(&pdf, context, config.PageSize); err != nil {

    return "", err
  }

  buff := bytes.NewBuffer([]byte{})

  if err = pdf.Write(buff); err != nil {
    panic(err)
  }

  _, output := GenUniqFileNameOutput(cacheDir, "certificate")

  if err = os.WriteFile(output+".pdf", buff.Bytes(), 0644); err != nil {

    return "", err
  }

  var data []byte

  data, err = bimg.NewImage(buff.Bytes()).Convert(bimg.JPEG)
  if err != nil {
    panic(err)
  }

  if err = os.WriteFile(output+".jpg", data, 0644); err != nil {

    return "", err
  }

  return output, nil
}
