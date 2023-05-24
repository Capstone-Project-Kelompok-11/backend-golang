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
  "math/rand"
  "mime/multipart"
  "net/http"
  "net/url"
  "os"
  "reflect"
  "skfw/papaya/koala/kio"
  "skfw/papaya/koala/kornet"
  "skfw/papaya/koala/mapping"
  "skfw/papaya/koala/pp"
  "skfw/papaya/koala/tools/posix"
  "strings"
  "time"
)

var ImageLimitSize = 1024 * 1024 * 2    // 2MB
var FileLimitSize = 1024 * 1024 * 2     // 2MB
var DocumentLimitSize = 1024 * 1024 * 2 // 2MB

func HexRand(size int) string {

  temp := ""
  digits := strings.Split("0123456789abcdef", "")
  randomize := rand.New(rand.NewSource(time.Now().UnixNano()))
  k := 1

  min := 0
  max := len(digits) - 1

  for i := 0; i < size; i++ {

    k = randomize.Intn(max-min+1) + min

    temp += digits[k]
  }

  return temp
}

func CheckAvailableFileName(src string) bool {

  return !kio.KFileNew(src).IsExist()
}

func GenUniqFileNameOutput(dir string, filename string) (string, string) { // filename, output

  var output string

  for {

    // re-generate with randomize
    filename = HexRand(7) + "." + filename
    output = posix.KPathNew(dir).JoinStr(filename)

    if !CheckAvailableFileName(output) {
      continue
    }

    break
  }

  return filename, output
}

func SourceNetOrEmpty(src string) bool {

  if src != "" {

    return strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://") || strings.HasPrefix(src, "ftp://")
  }

  return true // string is empty
}

func RemoveExtensionFromFileName(filename string) string {

  if filename != "" {

    tokens := strings.Split(filename, ".")
    return strings.Join(tokens[:len(tokens)-1], ".")
  }

  // passing
  return filename
}

func SafeParseSearchAndSortOrder(search string, sort string) (string, string, error) {

  var err error
  search = strings.TrimSpace(search)
  sort = strings.TrimSpace(sort)

  if search, err = url.QueryUnescape(search); err != nil {

    return search, sort, err
  }

  search = "%" + strings.ReplaceAll(search, " ", "%") + "%"
  switch strings.ToLower(sort) {
  case "asc", "ascending":
    sort = "ASC"
    break
  case "desc", "descending":
    sort = "DESC"
    break
  default:
    sort = "ASC"
    break
  }

  return search, sort, nil
}

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

func ValueToInt64(value any) int64 {

  val := pp.KIndirectValueOf(value)

  if val.IsValid() {

    ty := val.Type()

    switch ty.Kind() {

    case reflect.Float64:

      return int64(mapping.KValueToFloat(value))
    }

    return mapping.KValueToInt(value)
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
    if 48 <= r && r <= 57 {

      return r // number only
    }
    if 65 <= r && r <= 90 {

      return r + 32 // to lower
    }
    if 97 <= r && r <= 122 {

      return r // keep lower
    }
    switch r {
    case 32: // replace <space> with "-"
      return 45
    case 45, 46, 95: // "-", ".", "_"
      return r
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

    size := len(buff)

    if size > ImageLimitSize {

      return ctx.Status(http.StatusBadRequest).JSON(kornet.Msg("image file is too big than "+kornet.ReprByte(uint64(ImageLimitSize)), true))
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

func SaveFile(key string, output string) fiber.Handler {

  // no fix, limit reachable, not work for stream or large file

  return func(ctx *fiber.Ctx) error {

    var err error
    var header *multipart.FileHeader
    var file multipart.File
    var buff []byte

    if header, err = ctx.FormFile(key); err != nil {

      return ctx.Status(http.StatusInternalServerError).JSON(kornet.Msg("unable to get header information from form-data", true))
    }

    if file, err = header.Open(); err != nil {

      return ctx.Status(http.StatusInternalServerError).JSON(kornet.Msg("unable to open file", true))
    }

    defer file.Close()

    if buff, err = io.ReadAll(file); err != nil {

      return ctx.Status(http.StatusInternalServerError).JSON(kornet.Msg("unable to read file", true))
    }

    size := len(buff)

    if size > FileLimitSize {

      return ctx.Status(http.StatusBadRequest).JSON(kornet.Msg("file is too big than "+kornet.ReprByte(uint64(FileLimitSize)), true))
    }

    if err = os.WriteFile(output, buff, 0644); err != nil {

      return ctx.Status(http.StatusInternalServerError).JSON(kornet.Msg("unable to save file", true))
    }

    return ctx.Status(http.StatusCreated).JSON(kornet.Msg("upload file successfully", false))
  }
}

func SaveDocument(key string, output string) fiber.Handler {

  return func(ctx *fiber.Ctx) error {

    var err error
    var header *multipart.FileHeader
    var file multipart.File
    var buff []byte

    if header, err = ctx.FormFile(key); err != nil {

      return ctx.Status(http.StatusInternalServerError).JSON(kornet.Msg("unable to get header information from form-data", true))
    }

    if file, err = header.Open(); err != nil {

      return ctx.Status(http.StatusInternalServerError).JSON(kornet.Msg("unable to open document", true))
    }

    defer file.Close()

    if buff, err = io.ReadAll(file); err != nil {

      return ctx.Status(http.StatusInternalServerError).JSON(kornet.Msg("unable to read document", true))
    }

    size := len(buff)

    if size > DocumentLimitSize {

      return ctx.Status(http.StatusBadRequest).JSON(kornet.Msg("document file is too big than "+kornet.ReprByte(uint64(DocumentLimitSize)), true))
    }

    if err = os.WriteFile(output, buff, 0644); err != nil {

      return ctx.Status(http.StatusInternalServerError).JSON(kornet.Msg("unable to save document", true))
    }

    return ctx.Status(http.StatusCreated).JSON(kornet.Msg("upload document successfully", false))
  }
}
