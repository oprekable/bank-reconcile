package shutdown

import (
	"bytes"
	"context"
	"log"
	"os"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestSignalTrapWait(t *testing.T) {
	type args struct {
		ctx context.Context
		f   func()
	}

	tests := []struct {
		name    string
		t       SignalTrap
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Context canceled",
			t:    TermSignalTrap(),
			args: args{
				ctx: func() context.Context {
					ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
					defer cancel()
					return ctx
				}(),
				f: func() {
					log.Println("foo")
				},
			},
			wantErr: true,
			want:    "foo",
		},
		{
			name: "Syscall signal",
			t: func() (r SignalTrap) {
				trap := SignalTrap(make(chan os.Signal, 1))
				trap <- syscall.SIGTERM
				return trap
			}(),
			args: args{
				ctx: context.Background(),
				f: func() {
					log.Println("bar")
				},
			},
			wantErr: true,
			want:    "bar",
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

			err := tt.t.Wait(tt.args.ctx, tt.args.f)
			if (err != nil) != tt.wantErr {
				t.Errorf("Wait() error = %v, wantErr %v", err, tt.wantErr)
			}

			if got := bf.String(); strings.TrimRight(got, "\n") != tt.want {
				t.Errorf("Wait() output = %v, want %v", got, tt.want)
			}

			bf.Reset()
		})
	}
}
