package util

import (
  "bytes"
  "fmt"
  "github.com/gofiber/fiber/v2"
  "golang.org/x/image/draw"
  "image"
  "image/jpeg"
  "image/png"
  "io"
  "mime/multipart"
  "net/http"
  "os"
  "reflect"
  "skfw/papaya/koala/kornet"
  "skfw/papaya/koala/mapping"
  "skfw/papaya/koala/pp"
  "strings"
)

func ValueToInt(value any) int {

  val := pp.KIndirectValueOf(value)

  if val.IsValid() {

    ty := val.Type()

    switch ty.Kind() {

    case reflect.Float64:

      return int(mapping.KValueToFloat(value))
    }

    return int(mapping.KValueToInt(value))
  }

  return 0
}

func ValueToArrayStr(data any) []string {

  var temp []string
  temp = make([]string, 0)

  val := pp.KIndirectValueOf(data)

  if val.IsValid() {

    ty := val.Type()

    switch ty.Kind() {
    case reflect.Array, reflect.Slice:

      // loop - validity - casting

      for i := 0; i < val.Len(); i++ {

        elem := val.Index(i)

        vElem := pp.KIndirectValueOf(elem)

        if vElem.IsValid() {

          tyElem := vElem.Type()

          switch tyElem.Kind() {

          case reflect.String:

            temp = append(temp, vElem.String())
          }
        }
      }
    }
  }

  return temp
}

func SafePathName(name string) string {

  return strings.Map(func(r rune) rune {
    if 65 <= r && r <= 90 {

      return r + 32 // to lower
    }
    if 97 <= r && r <= 122 {

      return r // keep lower
    }
    if r == 32 {

      return 45 // replace space with minus
    }
    if r == 46 {

      return r // keep dot
    }
    return -1
  }, name)
}

func ResizeImageX256(data []byte, format string) ([]byte, error) {

  format, _ = strings.CutPrefix(format, ".") // maybe use a extension

  var err error
  var img image.Image

  if img, _, err = image.Decode(bytes.NewReader(data)); err != nil {

    return nil, err
  }

  dst := image.NewRGBA(image.Rect(0, 0, 256, 256))

  draw.NearestNeighbor.Scale(dst, dst.Rect, img, img.Bounds(), draw.Over, nil)

  var buf bytes.Buffer

  switch format {
  case "png":

    if err = png.Encode(&buf, dst); err != nil {

      return nil, err
    }
    break

  case "jpe", "jpeg", "jpg":

    if err = jpeg.Encode(&buf, dst, nil); err != nil {

      return nil, err
    }
    break

  default:

    return nil, fmt.Errorf("unsupported format: %s", format)
  }

  return buf.Bytes(), nil
}

func SaveImageX256(key string, format string, output string) fiber.Handler {

  return func(ctx *fiber.Ctx) error {

    var err error
    var header *multipart.FileHeader
    var file multipart.File
    var buff []byte

    if header, err = ctx.FormFile(key); err != nil {

      return ctx.Status(http.StatusInternalServerError).JSON(kornet.Msg("unable to get header information from form-data", true))
    }

    if file, err = header.Open(); err != nil {

      return ctx.Status(http.StatusInternalServerError).JSON(kornet.Msg("unable to open image file", true))
    }

    defer file.Close()

    if buff, err = io.ReadAll(file); err != nil {

      return ctx.Status(http.StatusInternalServerError).JSON(kornet.Msg("unable to read image file", true))
    }

    if buff, err = ResizeImageX256(buff, format); err != nil {

      return ctx.Status(http.StatusInternalServerError).JSON(kornet.Msg("unable to resize image file", true))
    }

    if err = os.WriteFile(output, buff, 0644); err != nil {

      return ctx.Status(http.StatusInternalServerError).JSON(kornet.Msg("unable to save image file", true))
    }

    return ctx.Status(http.StatusCreated).JSON(kornet.Msg("upload image file successfully", false))
  }
}

func ResizeImage(data []byte, format string) ([]byte, error) {

  format, _ = strings.CutPrefix(format, ".") // maybe use a extension

  var err error
  var img image.Image

  if img, _, err = image.Decode(bytes.NewReader(data)); err != nil {

    return nil, err
  }

  ImageScale := .8

  dst := image.NewRGBA(image.Rect(0, 0, int(float64(img.Bounds().Max.X)*ImageScale), int(float64(img.Bounds().Max.Y)*ImageScale)))

  draw.NearestNeighbor.Scale(dst, dst.Rect, img, img.Bounds(), draw.Over, nil)

  var buf bytes.Buffer

  switch format {
  case "png":

    if err = png.Encode(&buf, dst); err != nil {

      return nil, err
    }
    break

  case "jpe", "jpeg", "jpg":

    if err = jpeg.Encode(&buf, dst, nil); err != nil {

      return nil, err
    }
    break

  default:

    return nil, fmt.Errorf("unsupported format: %s", format)
  }

  return buf.Bytes(), nil
}

func SaveImage(key string, format string, output string) fiber.Handler {

  return func(ctx *fiber.Ctx) error {

    var err error
    var header *multipart.FileHeader
    var file multipart.File
    var buff []byte

    if header, err = ctx.FormFile(key); err != nil {

      return ctx.Status(http.StatusInternalServerError).JSON(kornet.Msg("unable to get header information from form-data", true))
    }

    if file, err = header.Open(); err != nil {

      return ctx.Status(http.StatusInternalServerError).JSON(kornet.Msg("unable to open image file", true))
    }

    defer file.Close()

    if buff, err = io.ReadAll(file); err != nil {

      return ctx.Status(http.StatusInternalServerError).JSON(kornet.Msg("unable to read image file", true))
    }

    if buff, err = ResizeImage(buff, format); err != nil {

      return ctx.Status(http.StatusInternalServerError).JSON(kornet.Msg("unable to resize image file", true))
    }

    if err = os.WriteFile(output, buff, 0644); err != nil {

      return ctx.Status(http.StatusInternalServerError).JSON(kornet.Msg("unable to save image file", true))
    }

    return ctx.Status(http.StatusCreated).JSON(kornet.Msg("upload image file successfully", false))
  }
}
