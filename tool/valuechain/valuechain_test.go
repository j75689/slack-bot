package valuechain

import "testing"

func TestExecute(t *testing.T) {
	type args struct {
		tag   string
		value string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test Basic",
			args: args{
				tag:   "",
				value: "a!!@#aa",
			},
			want: "a!!@#aa",
		},
		{
			name: "Test UrlEncode",
			args: args{
				tag:   "URLEncode",
				value: "a!!@#aa",
			},
			want: "a%21%21%40%23aa",
		},
		{
			name: "Test UrlDecode",
			args: args{
				tag:   "URLDecode",
				value: "a%21%21%40%23aa",
			},
			want: "a!!@#aa",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Execute(tt.args.tag, tt.args.value); got != tt.want {
				t.Errorf("Execute() = %v, want %v", got, tt.want)
			}
		})
	}
}
