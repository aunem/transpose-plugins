package main_test

import (
	"testing"
)

type A struct {
	Host   string
	Routes []string
	Mux    []Mux
}
type Mux struct {
	Upstream   string
	Downstream string
	Ip         int64
}

var a A

func init() {
	a = A{
		Host:   "my.host",
		Routes: []string{"myroute1", "myroute2"},
		Mux:    []Mux{Mux{Upstream: "my.ip", Downstream: "my.ip", Ip: int64(2903482034)}},
	}

}
func BenchmarkTest10(b *testing.B) {
	for n := 0; n < b.N; n++ {
		test(a)
	}
}

func test(abe interface{}) A {
	return abe.(A)
}
