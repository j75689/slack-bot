package tool

import (
	"reflect"
	"testing"
)

func TestReplaceVariables(t *testing.T) {
	type args struct {
		content   []byte
		variables map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "Test case 1",
			args: args{
				content: []byte(`{"name":"${name}","size":"${body.size}","content":"${##URLencode:body.content##}","length":"${body.length}"}`),
				variables: map[string]interface{}{
					"name": "hello",
					"body": map[string]interface{}{
						"size":    5,
						"length":  1024,
						"content": "world!#@",
					},
				},
			},
			want: []byte(`{"name":"hello","size":"5","content":"world%21%23%40","length":"1024"}`),
		},
		{
			name: "Test case 2",
			args: args{
				content: []byte(`{"name":"${name}","size":"${body.size}","content":"${##URLdecode:body.decode##}","length":"${body.length}"}`),
				variables: map[string]interface{}{
					"name": "hello",
					"body": map[string]interface{}{
						"size":    5,
						"length":  1024,
						"content": "world!#@",
						"decode":  "world%21%23%40",
					},
				},
			},
			want: []byte(`{"name":"hello","size":"5","content":"world!#@","length":"1024"}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReplaceVariables(tt.args.content, tt.args.variables); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReplaceVariables() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}
