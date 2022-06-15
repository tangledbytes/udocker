package main

import (
	"reflect"
	"testing"
)

func Test_convertStringSetToFlags(t *testing.T) {
	type args struct {
		set []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "When slice is empty",
			args: args{
				set: []string{},
			},
			want: []string{},
		},
		{
			name: "When slice has set strings with len(key) > 1",
			args: args{
				set: []string{"IJK=1", "hello=world", "platform=linux/amd64,linux/arm64", "build_arg=uid=501,gid=20,username=utkarsh"},
			},
			want: []string{"--IJK", "1", "--hello", "world", "--platform", "linux/amd64,linux/arm64", "--build-arg", "uid=501,gid=20,username=utkarsh"},
		},
		{
			name: "When slice has set strings with len(key) = 1",
			args: args{
				set: []string{"i=", "t=", "p=access"},
			},
			want: []string{"-i", "-t", "-p", "access"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertStringSetToFlags(tt.args.set); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertStringSetToFlags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_transformArgsWithSet(t *testing.T) {
	type args struct {
		args   []string
		set    []string
		target string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := transformArgsWithSet(tt.args.args, tt.args.set, tt.args.target); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("transformArgsWithSet() = %v, want %v", got, tt.want)
			}
		})
	}
}
