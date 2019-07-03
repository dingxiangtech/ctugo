package ctugo

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"reflect"
	"sort"
	"sync"
	"time"

	json "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
)

var (
	// sign pool, for memory reuse
	signBufPool = sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}

	client = &fasthttp.Client{
		MaxConnsPerHost: 16384, // MaxConnsPerHost  default is 512, increase to 16384
		ReadTimeout:     5 * time.Second,
		WriteTimeout:    5 * time.Second,
	}
)

type EngineResponse struct {
	UUID   string       `json:"uuid"`
	Status string       `json:"status"`
	Result engineResult `json:"result"`
}

type engineResult struct {
	RiskLevel       string                   `json:"riskLevel"`
	RiskType        string                   `json:"riskType"`
	Suggestion      []map[string]string      `json:"suggestion"`
	HitPolicyCode   string                   `json:"hitPolicyCode"`
	HitPolicyName   string                   `json:"hitPolicyName"`
	HitRules        []map[string]interface{} `json:"hitRules"`
	SuggestPolicies []map[string]string      `json:"suggestPolicies"`
	Flag            string                   `json:"flag"`
	ExtraInfo       map[string]interface{}   `json:"extraInfo"`
	NameListJSON    map[string]string        `json:"nameListJson"`
}

// EngineConnection object of engine pet tester
type EngineConnection struct {
	URLWithoutSign string

	Host      string
	AppKey    string
	AppSecret string
}

// NewEngineConnection generate a new EngineConnection object
func NewEngineConnection(host, appKey, appSecret string) *EngineConnection {
	return &EngineConnection{
		Host:      host,
		AppKey:    appKey,
		AppSecret: appSecret,

		URLWithoutSign: fmt.Sprintf("%s?appKey=%s&version=1&sign=", host, appKey),
	}
}

// CallRiskEngine is used to send request to ctu engine, and got EngineResponse and error
func (e *EngineConnection) CallRiskEngine(eventCode string, fields map[string]interface{}) (*EngineResponse, error) {
	data, err := e.getData(eventCode, fields)
	if err != nil {
		return nil, err
	}

	req := fasthttp.AcquireRequest()
	req.SetRequestURI(e.URLWithoutSign + e.getSign(eventCode, fields))
	req.Header.SetMethod("POST")
	req.Header.SetContentType("text/plain")
	req.SetBody(data)

	resp := fasthttp.AcquireResponse()

	defer fasthttp.ReleaseResponse(resp)
	defer fasthttp.ReleaseRequest(req)

	if err := client.Do(req, resp); err != nil {
		return nil, err
	}

	if resp.StatusCode() != fasthttp.StatusOK {
		return nil, fmt.Errorf("%d", resp.StatusCode())
	}

	result := &EngineResponse{}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (e *EngineConnection) getSign(eventCode string, fields map[string]interface{}) string {
	data := signBufPool.Get().(*bytes.Buffer)
	data.Reset()
	data.WriteString(e.AppSecret)
	data.WriteString("eventCode")
	data.WriteString(eventCode)
	data.WriteString("flag")
	data.WriteString(eventCode)

	keys := make([]string, len(fields))
	idx := 0
	for key := range fields {
		keys[idx] = key
		idx++
	}
	sort.Strings(keys)

	for _, key := range keys {
		data.WriteString(key)
		value := fields[key]
		if value != nil && reflect.TypeOf(value).Name() == "string" {
			data.WriteString(value.(string))
		} else {
			valueStr, _ := json.Marshal(value)
			data.Write(valueStr)
		}
	}

	data.WriteString(e.AppSecret)

	sum := md5.Sum(data.Bytes())
	signBufPool.Put(data)

	return hex.EncodeToString(sum[:])
}

func (e *EngineConnection) getData(eventCode string, fields map[string]interface{}) ([]byte, error) {
	var data []byte
	var err error

	data, err = json.Marshal(map[string]interface{}{
		"flag":      eventCode,
		"data":      fields,
		"eventCode": eventCode,
	})
	if err != nil {
		return nil, err
	}

	dst := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(dst, data)

	return dst, nil
}
