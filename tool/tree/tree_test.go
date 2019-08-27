package tree

import (
	"bytes"
	"testing"
)

var tree = NewTree()

func TestTree_Insert(t *testing.T) {

	type args struct {
		key   string
		value []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test Tree Insert Case 1",
			args: args{
				key:   "search sid",
				value: []byte(`data...`),
			},
			wantErr: false,
		},
		{
			name: "Test Tree Insert Case 2",
			args: args{
				key:   "search sid ignore",
				value: []byte(`data2...`),
			},
			wantErr: false,
		},
		{
			name: "Test Tree Insert Case 3",
			args: args{
				key:   "help ",
				value: []byte(`data3...`),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := tree.Insert(tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Tree.Insert() error = %v, wantErr %v", err, tt.wantErr)
			}
			if data, _ := tree.Search(tt.args.key); !bytes.Equal(data, tt.args.value) {
				t.Errorf("Tree.Insert() valid failed = %v, origin: %v", string(data), string(tt.args.value))
			}
		})
	}
}

func TestTree_Update(t *testing.T) {
	tree.Insert("search sid", []byte(`abc`))
	type args struct {
		key   string
		value []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test Tree Update Case 1",
			args: args{
				key:   "search sid",
				value: []byte(`data...`),
			},
			wantErr: false,
		},
		{
			name: "Test Tree Update Case 2",
			args: args{
				key:   "search sid",
				value: []byte(`data.....`),
			},
			wantErr: false,
		},
		{
			name: "Test Tree Update Case 3",
			args: args{
				key:   "search sid",
				value: []byte(`data........`),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := tree.Update(tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Tree.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
			if data, _ := tree.Search(tt.args.key); !bytes.Equal(data, tt.args.value) {
				t.Errorf("Tree.Update() valid failed = %v, origin: %v", string(data), string(tt.args.value))
			}
		})
	}
}

func TestTree_Delete(t *testing.T) {
	tree.Insert("search sid", []byte(`abc`))
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test Tree Delete Case 1",
			args: args{
				key: "search sid",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tree.Delete(tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("Tree.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
			if _, err := tree.Search(tt.args.key); err.Error() != "wrong key" {
				t.Errorf("Tree.Delete() valid failed , error = %v", err)
			}
		})
	}
}
