package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"io"
	"lms/app/models"
	"net/http"
	"net/url"
	"os"
	"reflect"
	m "skfw/papaya/koala/mapping"
	"skfw/papaya/koala/pp"
	"skfw/papaya/pigeon"
	"skfw/papaya/pigeon/drivers/common"
	"skfw/papaya/pigeon/drivers/postgresql"
	"skfw/papaya/pigeon/easy"
	"skfw/papaya/pigeon/templates/basicAuth/util"
	"strconv"
	"testing"
	"time"
)

type Map map[string]any
type Handler func(req *http.Request) error
type Maps []Map

func Mapping(data Map) m.KMapImpl {

	mm := m.KMap(data)
	return &mm
}

func Do(method string, URL url.URL, payload Map, handler Handler) ([]byte, error) {

	var err error
	var req *http.Request
	var res *http.Response

	buff := []byte(m.KMapEncodeJSON(payload))
	body := bytes.NewReader(buff)

	if req, err = http.NewRequest(method, URL.String(), body); err != nil {

		return nil, err
	}

	req.Header.Set("Origin", "https://skfw.net")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 Edg/114.0.1823.37")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Encoding", "gzip,deflate")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", strconv.Itoa(len(buff)))

	if handler != nil {

		if err = handler(req); err != nil {

			return nil, err
		}
	}

	client := http.Client{Timeout: time.Second * 15}

	if res, err = client.Do(req); err != nil {

		return nil, err
	}

	fmt.Printf("Status: %s\n", res.Status)
	fmt.Printf("StatusCode: %d\n", res.StatusCode)

	defer func(Body io.ReadCloser) {

		if err = Body.Close(); err != nil {

			panic(err)
		}
	}(res.Body)

	var data []byte

	if data, err = io.ReadAll(res.Body); err != nil {

		return nil, err
	}

	return data, nil
}

func Req(method string, URL url.URL, payload Map, handler Handler) (*Map, error) {

	var err error
	var buff []byte

	if buff, err = Do(method, URL, payload, handler); err != nil {

		return nil, err
	}

	var data Map

	if err = json.Unmarshal(buff, &data); err != nil {

		return nil, err
	}

	return &data, nil
}

func Token(token string) Handler {

	return func(req *http.Request) error {

		req.Header.Set("Authorization", "Bearer "+token)

		return nil
	}
}

func ValueToInt(value any) int {

	val := pp.KIndirectValueOf(value)

	if val.IsValid() {

		ty := val.Type()

		switch ty.Kind() {

		case reflect.Float64:

			return int(m.KValueToFloat(value))
		}

		return int(m.KValueToInt(value))
	}

	return 0
}

func TestResetDatabase(t *testing.T) {

	var err error

	os.Setenv("DB_USERNAME", "user")
	os.Setenv("DB_PASSWORD", "1234")
	os.Setenv("DB_NAME", "academy")

	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_CHARSET", "utf8")
	os.Setenv("DB_TIMEZONE", "UTC")
	os.Setenv("DB_SECURE", "false")

	var conn common.DBConnectionImpl

	if conn, err = postgresql.DBConnectionNew(pigeon.InitLoadEnviron); err != nil {

		t.Error(err)
	}

	db := conn.GORM()

	db.AutoMigrate(&models.Users{}, &models.Sessions{})

	users := db.Model(&models.Users{})
	sessions := db.Model(&models.Sessions{})
	courses := db.Model(&models.Courses{})
	modules := db.Model(&models.Modules{})
	quizzes := db.Model(&models.Quizzes{})
	checkout := db.Model(&models.Checkout{})
	completionCourses := db.Model(&models.CompletionCourses{})
	completionModules := db.Model(&models.CompletionModules{})
	assignments := db.Model(&models.Assignments{})
	events := db.Model(&models.Events{})

	var tx *gorm.DB
	var password string

	users.Unscoped().Exec("DELETE FROM users")

	userId := easy.Idx(uuid.New())
	currentTime := time.Now().UTC()
	password, _ = util.HashPassword("User@1234")

	if tx = users.Exec("INSERT INTO users (id, username, email, password, created_at, updated_at, admin) VALUES (?, ?, ?, ?, ?, ?, ?)", userId, "user", "user@mail.co", password, currentTime, currentTime, false); tx.Error != nil {

		t.Error(tx.Error)
	}

	adminId := easy.Idx(uuid.New())
	password, _ = util.HashPassword("Admin@1234")

	if tx = users.Exec("INSERT INTO users (id, username, email, password, created_at, updated_at, admin) VALUES (?, ?, ?, ?, ?, ?, ?)", adminId, "admin", "admin@mail.co", password, currentTime, currentTime, true); tx.Error != nil {

		t.Error(tx.Error)
	}

	sessions.Unscoped().Exec("DELETE FROM sessions")

	courses.Unscoped().Exec("DELETE FROM courses")

	modules.Unscoped().Exec("DELETE FROM modules")

	quizzes.Unscoped().Exec("DELETE FROM quizzes")

	checkout.Unscoped().Exec("DELETE FROM checkout")

	completionCourses.Unscoped().Exec("DELETE FROM completion_courses")

	completionModules.Unscoped().Exec("DELETE FROM completion_modules")

	assignments.Unscoped().Exec("DELETE FROM assignments")

	events.Unscoped().Exec("DELETE FROM events")
}

var origin = url.URL{
	Scheme: "http",
	Host:   "localhost:8000",
}

var token, adminToken string

func TestAdminLogin(t *testing.T) {

	var err error
	var data *Map

	origin.Path = "/api/v1/users/login"
	origin.RawPath = origin.Path

	if data, err = Req(http.MethodPost, origin, Map{
		"username": "admin",
		"email":    "*",
		"password": "Admin@1234",
	}, nil); err != nil {

		t.Error("failed login")
	}

	mm := Mapping(*data)
	adminToken = m.KValueToString(mm.Get("data.token"))

	if m.KValueToBool(mm.Get("error")) {

		t.Error(m.KValueToString(mm.Get("message")))
	}

	t.Log("success login, token:", adminToken)
}

func TestLogin(t *testing.T) {

	var err error
	var data *Map

	origin.Path = "/api/v1/users/login"
	origin.RawPath = origin.Path

	if data, err = Req(http.MethodPost, origin, Map{
		"username": "user",
		"email":    "*",
		"password": "User@1234",
	}, nil); err != nil {

		t.Error("failed login")
	}

	mm := Mapping(*data)
	token = m.KValueToString(mm.Get("data.token"))

	if m.KValueToBool(mm.Get("error")) {

		t.Error(m.KValueToString(mm.Get("message")))
	}

	t.Log("success login, data:", mm.Get("data"))
}

func TestCatchAdminSessions(t *testing.T) {

	var err error
	var data *Map

	origin.Path = "/api/v1/users/sessions"
	origin.RawPath = origin.Path

	if data, err = Req(http.MethodGet, origin, Map{}, Token(adminToken)); err != nil {

		t.Error("failed catch admin sessions")
	}

	mm := Mapping(*data)

	if m.KValueToBool(mm.Get("error")) {

		t.Error(m.KValueToString(mm.Get("message")))
	}

	t.Log("success login, data:", mm.Get("data"))
}

func TestCatchAdminStat(t *testing.T) {

	var err error
	var data *Map

	origin.Path = "/api/v1/admin/stats"
	origin.RawPath = origin.Path

	if data, err = Req(http.MethodGet, origin, Map{}, Token(adminToken)); err != nil {

		t.Error("failed catch admin stat")
	}

	mm := Mapping(*data)

	if m.KValueToBool(mm.Get("error")) {

		t.Error(m.KValueToString(mm.Get("message")))
	}

	t.Log("success fetch stats, data:", mm.Get("data"))
}

var courseId string

func TestCreateCourse(t *testing.T) {

	var err error
	var data *Map

	origin.Path = "/api/v1/admin/course"
	origin.RawPath = origin.Path

	if data, err = Req(http.MethodPost, origin, Map{
		"description": "Life is Short, Like You",
		"price":       "1000",
		"level":       "intermediate",
		"name":        "Life is Short",
		"category":    "life,short",
	}, Token(adminToken)); err != nil {

		t.Error("failed create course")
	}

	mm := Mapping(*data)

	if m.KValueToBool(mm.Get("error")) {

		t.Error(m.KValueToString(mm.Get("message")))
	}

	courseId = m.KValueToString(mm.Get("data.id"))

	t.Log("success create course, data:", mm.Get("data"))
}

var moduleId string

func TestCreateModule(t *testing.T) {

	var err error
	var data *Map

	origin.Path = "/api/v1/admin/module"
	origin.RawPath = origin.Path
	origin.RawQuery = "id=" + courseId

	if data, err = Req(http.MethodPost, origin, Map{
		"description": "waste time with playful",
		"name":        "playful",
	}, Token(adminToken)); err != nil {

		t.Error("failed create module")
	}

	mm := Mapping(*data)

	if m.KValueToBool(mm.Get("error")) {

		t.Error(m.KValueToString(mm.Get("message")))
	}

	moduleId = m.KValueToString(mm.Get("data.id"))

	t.Log("success create module, data:", mm.Get("data"))
}

func TestCreateQuiz(t *testing.T) {

	var err error
	var data *Map

	origin.Path = "/api/v1/admin/module/quiz"
	origin.RawPath = origin.Path
	origin.RawQuery = "id=" + moduleId

	if data, err = Req(http.MethodPost, origin, Map{
		"quizzes": []m.KMapImpl{
			&m.KMap{
				"question": "what creator play for waste time in everyday?",
				"choices": []m.KMapImpl{
					&m.KMap{
						"text":  "mobile bang bang",
						"valid": false,
					},
					&m.KMap{
						"text":  "call of duty mobile",
						"valid": true,
					},
					&m.KMap{
						"text":  "watch anime",
						"valid": true,
					},
				},
			},
		},
	}, Token(adminToken)); err != nil {

		t.Error("failed create quiz")
	}

	mm := Mapping(*data)

	if m.KValueToBool(mm.Get("error")) {

		t.Error(m.KValueToString(mm.Get("message")))
	}

	t.Log("success create quiz, data:", mm.Get("data"))
}

func TestCheckoutCourse(t *testing.T) {

	var err error
	var data *Map

	origin.Path = "/api/v1/users/checkout"
	origin.RawPath = origin.Path
	origin.RawQuery = "id=" + courseId

	if data, err = Req(http.MethodPost, origin, Map{}, Token(token)); err != nil {

		t.Error("failed checkout course")
	}

	mm := Mapping(*data)

	if m.KValueToBool(mm.Get("error")) {

		t.Error(m.KValueToString(mm.Get("message")))
	}

	t.Log("success checkout course, data:", mm.Get("data"))
}

func TestCheckoutVerify(t *testing.T) {

	var err error
	var data *Map

	origin.Path = "/api/v1/users/checkout/verify"
	origin.RawPath = origin.Path
	origin.RawQuery = ""

	if data, err = Req(http.MethodPost, origin, Map{
		"payment_method": "debt-visa-card",
	}, Token(token)); err != nil {

		t.Error("failed checkout verify")
	}

	mm := Mapping(*data)

	if m.KValueToBool(mm.Get("error")) {

		t.Error(m.KValueToString(mm.Get("message")))
	}

	t.Log("success checkout verify, data:", mm.Get("data"))
}

func TestCheckoutHistory(t *testing.T) {

	var err error
	var data *Map

	origin.Path = "/api/v1/users/checkout/history"
	origin.RawPath = origin.Path
	origin.RawQuery = ""

	if data, err = Req(http.MethodGet, origin, Map{}, Token(token)); err != nil {

		t.Error("failed catch checkout history")
	}

	mm := Mapping(*data)

	if m.KValueToBool(mm.Get("error")) {

		t.Error(m.KValueToString(mm.Get("message")))
	}

	t.Log("success catch checkout history, data:", mm.Get("data"))
}

func TestCatchCourse(t *testing.T) {

	var err error
	var data *Map

	origin.Path = "/api/v1/users/course"
	origin.RawPath = origin.Path
	origin.RawQuery = "id=" + courseId

	if data, err = Req(http.MethodGet, origin, Map{}, Token(token)); err != nil {

		t.Error("failed catch course")
	}

	mm := Mapping(*data)

	if m.KValueToBool(mm.Get("error")) {

		t.Error(m.KValueToString(mm.Get("message")))
	}

	t.Log("success catch course, data:", mm.Get("data"))
}

func TestCatchQuiz(t *testing.T) {

	var err error
	var data *Map

	origin.Path = "/api/v1/users/quiz"
	origin.RawPath = origin.Path
	origin.RawQuery = "id=" + moduleId

	if data, err = Req(http.MethodGet, origin, Map{}, Token(token)); err != nil {

		t.Error("failed catch quiz")
	}

	mm := Mapping(*data)

	if m.KValueToBool(mm.Get("error")) {

		t.Error(m.KValueToString(mm.Get("message")))
	}

	t.Log("success catch quiz, data:", mm.Get("data"))
}

func TestTakeQuiz(t *testing.T) {

	var err error
	var data *Map

	origin.Path = "/api/v1/users/quiz"
	origin.RawPath = origin.Path
	origin.RawQuery = "id=" + moduleId

	if data, err = Req(http.MethodPost, origin, Map{
		"quizzes": []m.KMapImpl{
			&m.KMap{
				"question": "what creator play for waste time in everyday?",
				"choices": []m.KMapImpl{
					&m.KMap{
						"text":  "mobile bang bang",
						"valid": false,
					},
					&m.KMap{
						"text":  "call of duty mobile",
						"valid": true,
					},
					&m.KMap{
						"text":  "watch anime",
						"valid": true,
					},
				},
			},
		},
	}, Token(token)); err != nil {

		t.Error("failed take quiz")
	}

	mm := Mapping(*data)

	if m.KValueToBool(mm.Get("error")) {

		t.Error(m.KValueToString(mm.Get("message")))
	}

	t.Log("success take quiz, data:", mm.Get("data"))
}

func TestResumeMock(t *testing.T) {

	var err error
	var data *Map

	origin.Path = "/api/v1/users/course/resume"
	origin.RawPath = origin.Path
	origin.RawQuery = "id=" + courseId

	if data, err = Req(http.MethodPost, origin, Map{}, Token(token)); err != nil {

		t.Error("failed take resume")
	}

	mm := Mapping(*data)

	if m.KValueToBool(mm.Get("error")) {

		t.Error(m.KValueToString(mm.Get("message")))
	}

	t.Log("success take resume, data:", mm.Get("data"))
}

func TestCatchReport(t *testing.T) {

	var err error
	var data *Map

	origin.Path = "/api/v1/users/course/report"
	origin.RawPath = origin.Path
	origin.RawQuery = "id=" + courseId

	if data, err = Req(http.MethodGet, origin, Map{}, Token(token)); err != nil {

		t.Error("failed catch course report")
	}

	mm := Mapping(*data)

	if m.KValueToBool(mm.Get("error")) {

		t.Error(m.KValueToString(mm.Get("message")))
	}

	t.Log("success catch course report, data:", mm.Get("data"))
}

func TestCatchCertificate(t *testing.T) {

	var err error
	var data *Map

	origin.Path = "/api/v1/users/course/certificate"
	origin.RawPath = origin.Path
	origin.RawQuery = "id=" + courseId

	if data, err = Req(http.MethodGet, origin, Map{}, Token(token)); err != nil {

		t.Error("failed catch certificate")
	}

	mm := Mapping(*data)

	if m.KValueToBool(mm.Get("error")) {

		t.Error(m.KValueToString(mm.Get("message")))
	}

	t.Log("success catch certificate, data:", mm.Get("data"))
}
