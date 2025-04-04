package atexit

import (
	"bytes"
	"log"
	"os"
	"testing"
)

func TestAdd(t *testing.T) {
	type args struct {
		y []func()
	}

	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "ok",
			args: args{
				y: []func(){
					func() {
						log.Println(1)
					},
					func() {
						log.Println(2)
					},
				},
			},
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Add(tt.args.y...)
			if got := len(functions); got != tt.want {
				t.Errorf("len(functions) = %v, want %v", got, tt.want)
			}
			functions = functions[:0]
		})
	}
}

func TestAtExit(t *testing.T) {
	type args struct {
		y []func()
	}

	tests := []struct {
		name string
		want string
		args args
	}{
		{
			name: "ok",
			args: args{
				y: []func(){
					func() {
						log.Println(1)
					},
					func() {
						log.Println(2)
					},
				},
			},
			want: "1\n2\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bf bytes.Buffer
			log.SetOutput(&bf)
			log.SetFlags(0)
			t.Cleanup(func() {
				log.SetOutput(os.Stdout)
			})
			Add(tt.args.y...)
			AtExit()

			if got := bf.String(); got != tt.want {
				t.Errorf("AtExit() output = %v, want %v", got, tt.want)
			}

			functions = functions[:0]
			bf.Reset()
		})
	}
}
