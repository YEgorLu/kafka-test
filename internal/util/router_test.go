package util

import (
	"net/http"
	"testing"
)

func TestRpm(t *testing.T) {
	basePath := "basePath"
	f := Rpm(basePath)
	type args struct {
		method string
		path   []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"parses path if path not empty", args{http.MethodGet, []string{"path1"}}, http.MethodGet + " /" + basePath + "/path1"},
		{"skips if path empty", args{http.MethodGet, []string{}}, http.MethodGet + " /" + basePath},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := f(tt.args.method, tt.args.path...); got != tt.want {
				t.Errorf("Rpm() = %v, want %v", got, tt.want)
			}
		})
	}
}
