package v1

import (
	"bytes"
	"strings"
	"testing"
)

var (
	map1 = map[string]interface{}{
		"test": "me",
		"check": []interface{}{
			map[string]interface{}{
				"this": "test",
				"out":  nil,
			},
			nil,
			true,
		},
		"map": map[string]interface{}{
			"sth": nil,
			"test": []interface{}{
				123,
				nil,
				"abc",
				[]interface{}{
					nil,
				},
				[]interface{}{
					"test",
					nil,
					"me",
					"plz",
				},
			},
		},
	}
)

func Test_defaultJSONEncoder(t *testing.T) {
	type args struct {
		i interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
		{
			name: "OK",
			args: args{
				i: &map1,
			},
			wantW:   `{"check":[{"this":"test"},true],"map":{"test":[123,"abc",[],["test","me","plz"]]},"test":"me"}`,
			wantErr: false,
		},
		{
			name: "Err",
			args: args{
				i: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if err := defaultJSONEncoder(tt.args.i, w); (err != nil) != tt.wantErr {
				t.Errorf("defaultJSONEncoder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); strings.TrimRight(gotW, "\n") != tt.wantW {
				t.Errorf("defaultJSONEncoder() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
