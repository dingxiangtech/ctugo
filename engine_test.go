package ctugo

import (
	"encoding/base64"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

var (
	engineConn = NewEngineConnection("http://127.0.0.1:7090/ctu/event.do", "testappid", "testappkey")
)

func TestNewConn(t *testing.T) {
	Convey("new conn url should be OK", t, func() {
		So(engineConn.URLWithoutSign, ShouldEqual, "http://127.0.0.1:7090/ctu/event.do?appKey=testappid&version=1&sign=")
	})
}

func TestGetSign(t *testing.T) {
	Convey("get sign of one field", t, func() {
		sign := engineConn.getSign("test_event", "test_event", map[string]interface{}{
			"ip": "1.2.3.4",
		})
		So(sign, ShouldEqual, "5db1434d41ba1d8029de5abdde40060b")
	})

	Convey("get sign of multi fields", t, func() {
		sign := engineConn.getSign("test_event", "test_event", map[string]interface{}{
			"ip":      "1.2.3.4",
			"email":   "a@b.c",
			"user_id": "hello_world",
		})
		So(sign, ShouldEqual, "2d4d5c7245ece671c359f454e565e223")
	})

	Convey("get sign of non-string fields", t, func() {
		sign := engineConn.getSign("test_event", "test_event", map[string]interface{}{
			"ip":          "1.2.3.4",
			"androidRoot": false,
			"ratio":       1,
		})
		So(sign, ShouldEqual, "996aa5fdffcbb8e5e0e832cbcbd19c44")
	})

	Convey("get sign of null fields", t, func() {
		sign := engineConn.getSign("test_event", "test_event", map[string]interface{}{
			"ip":          nil,
			"androidRoot": false,
			"ratio":       1,
		})
		So(sign, ShouldEqual, "bc11ed1d9dff21f03642ba9ee1684e52")
	})
}

func TestGetData(t *testing.T) {
	Convey("get data of single field", t, func() {
		data, err := engineConn.getData("test_event", "test_event", map[string]interface{}{
			"ip": "1.2.3.4",
		})
		So(err, ShouldEqual, nil)
		dst := make([]byte, base64.StdEncoding.DecodedLen(len(data)))
		size, err := base64.StdEncoding.Decode(dst, data)
		So(err, ShouldEqual, nil)
		So(size, ShouldEqual, 70)
	})

	Convey("get data of multi field", t, func() {
		data, err := engineConn.getData("test_event", "test_event", map[string]interface{}{
			"ip":          "1.2.3.4",
			"androidRoot": true,
		})
		So(err, ShouldEqual, nil)
		dst := make([]byte, base64.StdEncoding.DecodedLen(len(data)))
		size, err := base64.StdEncoding.Decode(dst, data)
		So(err, ShouldEqual, nil)
		So(size, ShouldEqual, 89)
	})
}

func BenchmarkGetData(b *testing.B) {
	fields := map[string]interface{}{
		"ip":         "1.2.3.4",
		"email":      "a@b.c",
		"ext_amount": 123,
	}
	for i := 0; i < b.N; i++ {
		engineConn.getData("test_event", "test_event", fields)
	}
}

func BenchmarkGetSign(b *testing.B) {
	fields := map[string]interface{}{
		"ip":         "1.2.3.4",
		"email":      "a@b.c",
		"ext_amount": 123,
	}
	for i := 0; i < b.N; i++ {
		engineConn.getSign("test_event", "test_event", fields)
	}
}
