package tree

import (
	"bytes"
	"reflect"
	"testing"
)

var tree = NewTree()

func TestTree_Insert(t *testing.T) {

	type args struct {
		key   string
		value []byte
		term  string
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
				term:  "search sid",
			},
			wantErr: false,
		},
		{
			name: "Test Tree Insert Case 1.5",
			args: args{
				key:   "search sid",
				value: []byte(`data...2`),
				term:  "search sid",
			},
			wantErr: true,
		},
		{
			name: "Test Tree Insert Case 2",
			args: args{
				key:   "search sid ignore",
				value: []byte(`data2...`),
				term:  "search sid ignore",
			},
			wantErr: false,
		},
		{
			name: "Test Tree Insert Case 3",
			args: args{
				key:   "search sid ${0}",
				value: []byte(`data3...`),
				term:  "search sid abcdefg",
			},
			wantErr: false,
		},
		{
			name: "Test Tree Insert Case 4",
			args: args{
				key:   "search sid ${0} ${1}",
				value: []byte(`data4...`),
				term:  "search sid abc efga",
			},
			wantErr: false,
		},
		{
			name: "Test Tree Insert Case 5",
			args: args{
				key:   "help ",
				value: []byte(`data5...`),
				term:  "help",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := tree.Insert(tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Tree.Insert() error = %v, wantErr %v", err, tt.wantErr)
			}

			if data, _ := tree.Search(tt.args.term); !bytes.Equal(data, tt.args.value) && !tt.wantErr {
				t.Errorf("Tree.Insert() valid failed = %v, origin: %v", string(data), string(tt.args.value))
			}
		})
	}
}

func TestTree_Update(t *testing.T) {
	tree.Insert("search sid ${0}", []byte(`abc`))
	type args struct {
		key   string
		value []byte
		term  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test Tree Update Case 1",
			args: args{
				key:   "search sid ${0}",
				value: []byte(`data...`),
				term:  "search sid a",
			},
			wantErr: false,
		},
		{
			name: "Test Tree Update Case 2",
			args: args{
				key:   "search sid ${0}",
				value: []byte(`data.....`),
				term:  "search sid a",
			},
			wantErr: false,
		},
		{
			name: "Test Tree Update Case 3",
			args: args{
				key:   "search sid ${0}",
				value: []byte(`data........`),
				term:  "search sid a",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := tree.Update(tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Tree.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
			if data, _ := tree.Search(tt.args.term); !bytes.Equal(data, tt.args.value) {
				t.Errorf("Tree.Update() valid failed = %v, origin: %v", string(data), string(tt.args.value))
			}
		})
	}
}

func TestTree_Delete(t *testing.T) {

	type args struct {
		key   string
		value []byte
		term  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test Tree Delete Case 1",
			args: args{
				key:   "search sid abc",
				value: []byte(`data1`),
				term:  "search sid abc",
			},
			wantErr: false,
		},
		{
			name: "Test Tree Delete Case 2",
			args: args{
				key:   "search sid ${0}",
				value: []byte(`data2`),
				term:  "search sid data2",
			},
			wantErr: false,
		},
		{
			name: "Test Tree Delete Case 3",
			args: args{
				key:   "search sid def",
				value: []byte(`data3`),
				term:  "search sid def",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree.Insert(tt.args.key, tt.args.value)
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if data, _ := tree.Search(tt.args.term); !bytes.Equal(data, tt.args.value) {
				t.Errorf("Tree.Delete() before valid failed = %v, origin: %v", string(data), string(tt.args.value))
			}
			if err := tree.Delete(tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("Tree.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
			if data, err := tree.Search(tt.args.term); err == nil && bytes.Equal(tt.args.value, data) {
				t.Errorf("Tree.Delete() after valid failed after , find: %v", string(data))
			}
		})
	}
}

func TestTree_Search(t *testing.T) {

	type args struct {
		key   string
		term  string
		value []byte
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		wantEqual bool
	}{
		{
			name: "Test Tree Search Case 1",
			args: args{
				key:   "search sid ${0}",
				term:  "search sid abc",
				value: []byte(`data1`),
			},
			wantErr:   false,
			wantEqual: true,
		},
		{
			name: "Test Tree Search Case 1",
			args: args{
				key:   "search sid ${0}",
				term:  "search sid cdfg",
				value: []byte(`data1`),
			},
			wantErr:   false,
			wantEqual: true,
		},
		{
			name: "Test Tree Search Case 2",
			args: args{
				key:   "search sid ${0} ${1}",
				term:  "search sid cdfg",
				value: []byte(`data2`),
			},
			wantErr:   false,
			wantEqual: false,
		},
		{
			name: "Test Tree Search Case 2",
			args: args{
				key:   "search sid ${0} ${1}",
				term:  "search sid cdfg aa",
				value: []byte(`data2`),
			},
			wantErr:   false,
			wantEqual: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree.Insert(tt.args.key, tt.args.value)
			got, err := tree.Search(tt.args.term)
			if (err != nil) != tt.wantErr {
				t.Errorf("Tree.Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if reflect.DeepEqual(got, tt.args.value) != tt.wantEqual {
				t.Errorf("Tree.Search() = %v, want %v", string(got), string(tt.args.value))
			}
		})
	}
}
