package test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
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

	pp.Void(users, sessions)

	var prepared *sql.DB

	prepared, _ = sessions.DB()
	prepared.Query("DELETE FROM sessions")
}

var origin = url.URL{
	Scheme: "http",
	Host:   "localhost:8000",
}

var token, adminToken string

func TestLogin(t *testing.T) {

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

	t.Log("success login, sessions:", mm.Get("data"))
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

	t.Log("success fetch, statistics:", mm.Get("data"))
}
