package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	m "skfw/papaya/koala/mapping"
	"strconv"
	"testing"
	"time"
)

type Map map[string]any
type Handler func(req *http.Request) error

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

type PingPong struct {
	URL     *url.URL
	Handler Handler
}

func (p *PingPong) Ping() error {

	var err error
	var data *Map

	p.URL.Path = "/api/v1/test/ping"
	p.URL.RawPath = p.URL.Path

	if data, err = Req(http.MethodGet, origin, nil, p.Handler); err != nil {

		return err
	}

	if m.KValueToBool(Mapping(*data).Get("error")) {

		return errors.New(m.KValueToString(Mapping(*data).Get("message")))
	}

	return nil
}

var origin = url.URL{
	Scheme: "https",
	Host:   "skfw.net",
}

var token string

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

		t.Error("Connection refuse:", origin.String())
	}

	mm := Mapping(*data)
	token = m.KValueToString(mm.Get("data.token"))

	if m.KValueToBool(mm.Get("error")) {

		t.Error(m.KValueToString(mm.Get("message")))
	}

	t.Log("success login, data:", data)
}

func TestPingPong(t *testing.T) {

	var err error
	var data *Map

	origin.Path = "/api/v1/test/ping"
	origin.RawPath = origin.Path

	if data, err = Req(http.MethodGet, origin, nil, nil); err != nil {

		t.Error("Connection refuse:", origin.String())
	}

	mm := Mapping(*data)

	if m.KValueToBool(mm.Get("error")) {

		t.Error(m.KValueToString(mm.Get("message")))
	}

	t.Log("success ping pong, data:", data)
}

func TestAuthPingPong(t *testing.T) {

	var err error
	var data *Map

	origin.Path = "/api/v1/test/auth/ping"
	origin.RawPath = origin.Path

	if data, err = Req(http.MethodGet, origin, nil, Token(token)); err != nil {

		t.Error("Connection refuse:", origin.String())
	}

	mm := Mapping(*data)

	if m.KValueToBool(mm.Get("error")) {

		t.Error(m.KValueToString(mm.Get("message")))
	}

	t.Log("success ping pong, data:", data)
}

func TestAuthDevicePingPong(t *testing.T) {

	var err error
	var data *Map

	origin.Path = "/api/v1/test/auth/strict/ping"
	origin.RawPath = origin.Path

	if data, err = Req(http.MethodGet, origin, nil, Token(token)); err != nil {

		t.Error("Connection refuse:", origin.String())
	}

	mm := Mapping(*data)

	if m.KValueToBool(mm.Get("error")) {

		t.Error(m.KValueToString(mm.Get("message")))
	}

	t.Log("success ping pong, data:", data)
}

func TestPerformPingPong(t *testing.T) {

	var err error

	origin.Path = "/api/v1/test/ping"
	origin.RawPath = origin.Path

	p := &PingPong{
		URL:     &origin,
		Handler: nil,
	}

	var rtt time.Duration
	var after, before time.Time
	var Min, Avg, Max time.Duration
	Iteration := 20

	for i := 0; i < Iteration; i++ {

		before = time.Now().UTC()

		if err = p.Ping(); err != nil {

			t.Error(err)
			break
		}

		after = time.Now().UTC()
		rtt = after.Sub(before)

		Avg += rtt
		if i > 0 {

			Avg /= 2

			if rtt < Min {

				Min = rtt
			}

			if rtt > Max {

				Max = rtt
			}

		} else {

			Min = rtt
			Max = rtt
		}

		t.Logf("RTT #%d %s", i, rtt)
	}

	t.Log("Min", Min)
	t.Log("Avg", Avg)
	t.Log("Max", Max)
}

func TestPerformAuthPingPong(t *testing.T) {

	var err error

	origin.Path = "/api/v1/test/auth/ping"
	origin.RawPath = origin.Path

	p := &PingPong{
		URL:     &origin,
		Handler: Token(token),
	}

	var rtt time.Duration
	var after, before time.Time
	var Min, Avg, Max time.Duration
	Iteration := 20

	for i := 0; i < Iteration; i++ {

		before = time.Now().UTC()

		if err = p.Ping(); err != nil {

			t.Error(err)
			break
		}

		after = time.Now().UTC()
		rtt = after.Sub(before)

		Avg += rtt
		if i > 0 {

			Avg /= 2

			if rtt < Min {

				Min = rtt
			}

			if rtt > Max {

				Max = rtt
			}

		} else {

			Min = rtt
			Max = rtt
		}

		t.Logf("RTT #%d %s", i, rtt)
	}

	t.Log("Min", Min)
	t.Log("Avg", Avg)
	t.Log("Max", Max)
}

func TestPerformAuthDevicePingPong(t *testing.T) {

	var err error

	origin.Path = "/api/v1/test/auth/strict/ping"
	origin.RawPath = origin.Path

	p := &PingPong{
		URL:     &origin,
		Handler: Token(token),
	}

	var rtt time.Duration
	var after, before time.Time
	var Min, Avg, Max time.Duration
	Iteration := 20

	for i := 0; i < Iteration; i++ {

		before = time.Now().UTC()

		if err = p.Ping(); err != nil {

			t.Error(err)
			break
		}

		after = time.Now().UTC()
		rtt = after.Sub(before)

		Avg += rtt
		if i > 0 {

			Avg /= 2

			if rtt < Min {

				Min = rtt
			}

			if rtt > Max {

				Max = rtt
			}

		} else {

			Min = rtt
			Max = rtt
		}

		t.Logf("RTT #%d %s", i, rtt)
	}

	t.Log("Min", Min)
	t.Log("Avg", Avg)
	t.Log("Max", Max)
}

func TestLogout(t *testing.T) {

	var err error
	var data *Map

	origin.Path = "/api/v1/users/sessions"
	origin.RawPath = origin.Path

	if data, err = Req(http.MethodDelete, origin, nil, Token(token)); err != nil {

		t.Error("Connection refuse:", origin.String())
	}

	mm := Mapping(*data)

	if m.KValueToBool(mm.Get("error")) {

		t.Error(m.KValueToString(mm.Get("message")))
	}

	t.Log("success logout, data:", data)
}
