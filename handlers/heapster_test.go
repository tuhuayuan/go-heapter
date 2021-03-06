package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"testing"
	"zonst/qipai-golang-libs/httputil"
	"zonst/qipai/gamehealthysrv/middlewares"
	"zonst/qipai/gamehealthysrv/models"

	"net/url"

	"github.com/stretchr/testify/assert"
)

var heapsterTestID string

func TestCreateHeapster(t *testing.T) {
	ctx := httputil.WithHTTPContext(nil)
	httputil.Use(ctx, middlewares.RedisConnHandler("0.0.0.0:6379", "", 9))
	handler := httputil.HandleFunc(ctx,
		middlewares.BindBody(&CreateHeapsterReq{}),
		CreateHeapsterHandler)
	data := []byte(`
    {
        "name": "sample_heapster",
        "type": "http",
        "port": 5050,
        "accept_code": [
            400
        ],
        "host": "zonst.local2",
        "location": "/healthz",
        "timeout": 3,
        "interval": 5,
        "healthy_threshold": 3,
        "unhealthy_threshold": 3,
        "groups": [
            "test_group"
        ],
        "notifiers": [
            "test_notifiers"
        ]
    }
    `)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(data))
	req.Header.Add("Content-Type", "json")
	resp := httptest.NewRecorder()
	handler(resp, req)
	assert.Equal(t, 200, resp.Code)
	body, _ := ioutil.ReadAll(resp.Body)
	hp := models.Heapster{}
	assert.NoError(t, json.Unmarshal(body, &hp))
	heapsterTestID = string(hp.ID)
	fmt.Println(heapsterTestID)
}

func TestUpdateHeapter(t *testing.T) {
	ctx := httputil.WithHTTPContext(nil)
	httputil.Use(ctx, middlewares.RedisConnHandler("0.0.0.0:6379", "", 9))
	handler := httputil.HandleFunc(ctx,
		middlewares.BindBody(&UpdateHeapsterReq{}),
		UpdateHeapsterHandler)
	data := []byte(`
    {
        "id": "` + heapsterTestID + `",
        "name": "sample_heapster",
        "type": "http",
        "port": 5050,
        "accept_code": [
            400, 200
        ],
        "host": "zonst.local2",
        "location": "/healthz",
        "timeout": 3,
        "interval": 5,
        "healthy_threshold": 3,
        "unhealthy_threshold": 3,
        "groups": [
            "test_group"
        ],
        "notifiers": [
            "test_notifiers"
        ]
    }
    `)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(data))
	req.Header.Add("Content-Type", "json")
	resp := httptest.NewRecorder()
	handler(resp, req)
	assert.Equal(t, 200, resp.Code)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func TestFetchHeapster(t *testing.T) {
	ctx := httputil.WithHTTPContext(nil)
	httputil.Use(ctx, middlewares.RedisConnHandler("0.0.0.0:6379", "", 9))
	handler := httputil.HandleFunc(ctx,
		middlewares.BindBody(&FetchHeapsterReq{}),
		FetchHeapsterHandler)

	req := httptest.NewRequest("GET", "/", nil)
	resp := httptest.NewRecorder()
	handler(resp, req)
	assert.Equal(t, 200, resp.Code)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func TestMuteHeapster(t *testing.T) {
	ctx := httputil.WithHTTPContext(nil)
	httputil.Use(ctx, middlewares.RedisConnHandler("0.0.0.0:6379", "", 9))
	handler := httputil.HandleFunc(ctx,
		middlewares.BindBody(&MuteHeapsterReq{}),
		MuteHeapsterHandler)

	data := make(url.Values)
	data.Add("id", heapsterTestID)
	data.Add("mute", "1")
	fmt.Println(data.Encode())
	req := httptest.NewRequest("POST", "/", bytes.NewReader([]byte(data.Encode())))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp := httptest.NewRecorder()
	handler(resp, req)
	assert.Equal(t, 200, resp.Code)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func TestDeleteHeapster(t *testing.T) {
	ctx := httputil.WithHTTPContext(nil)
	httputil.Use(ctx, middlewares.RedisConnHandler("0.0.0.0:6379", "", 9))
	handler := httputil.HandleFunc(ctx,
		middlewares.BindBody(&DeleteHeapsterReq{}),
		DeleteHeapsterHandler)
	data := []byte(`
    {
        "id": "` + heapsterTestID + `"
    }
    `)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(data))
	req.Header.Add("Content-Type", "json")
	resp := httptest.NewRecorder()
	handler(resp, req)
	assert.Equal(t, 200, resp.Code)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
