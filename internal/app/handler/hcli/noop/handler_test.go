package noop

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/oprekable/bank-reconcile/internal/app/component"
	"github.com/oprekable/bank-reconcile/internal/app/repository"
	"github.com/oprekable/bank-reconcile/internal/app/service"
)

func TestHandlerExec(t *testing.T) {
	var bf bytes.Buffer
	type fields struct {
		comp *component.Components
		svc  *service.Services
		repo *repository.Repositories
	}

	tests := []struct {
		fields  fields
		name    string
		wantErr bool
	}{
		{
			name: "Ok",
			fields: fields{
				comp: &component.Components{},
				svc:  &service.Services{},
				repo: &repository.Repositories{},
			},
			wantErr: false,
		},
		{
			name: "Ok",
			fields: fields{
				comp: nil,
				svc:  nil,
				repo: nil,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				comp:   tt.fields.comp,
				svc:    tt.fields.svc,
				repo:   tt.fields.repo,
				writer: &bf,
			}

			if err := h.Exec(); (err != nil) != tt.wantErr {
				t.Errorf("Exec() error = %v, wantErr %v", err, tt.wantErr)
			}

			bf.Reset()
		})
	}
}

func TestHandlerName(t *testing.T) {
	var bf bytes.Buffer
	type fields struct {
		comp *component.Components
		svc  *service.Services
		repo *repository.Repositories
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Ok",
			fields: fields{
				comp: nil,
				svc:  nil,
				repo: nil,
			},
			want: "noop",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				comp:   tt.fields.comp,
				svc:    tt.fields.svc,
				repo:   tt.fields.repo,
				writer: &bf,
			}

			if got := h.Name(); got != tt.want {
				t.Errorf("Name() = %v, want %v", got, tt.want)
			}
			bf.Reset()
		})
	}
}

func TestHandlerSetComponents(t *testing.T) {
	var bf bytes.Buffer
	type fields struct {
		comp *component.Components
		svc  *service.Services
		repo *repository.Repositories
	}

	type args struct {
		c *component.Components
	}

	tests := []struct {
		fields fields
		args   args
		name   string
	}{
		{
			name: "Ok",
			fields: fields{
				comp: nil,
				svc:  nil,
				repo: nil,
			},
			args: args{
				c: &component.Components{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				comp:   tt.fields.comp,
				svc:    tt.fields.svc,
				repo:   tt.fields.repo,
				writer: &bf,
			}

			h.SetComponents(tt.args.c)
			bf.Reset()
		})
	}
}

func TestHandlerSetRepositories(t *testing.T) {
	var bf bytes.Buffer
	type fields struct {
		comp *component.Components
		svc  *service.Services
		repo *repository.Repositories
	}

	type args struct {
		r *repository.Repositories
	}

	tests := []struct {
		fields fields
		args   args
		name   string
	}{
		{
			name: "Ok",
			fields: fields{
				comp: nil,
				svc:  nil,
				repo: nil,
			},
			args: args{
				r: &repository.Repositories{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				comp:   tt.fields.comp,
				svc:    tt.fields.svc,
				repo:   tt.fields.repo,
				writer: &bf,
			}

			h.SetRepositories(tt.args.r)
			bf.Reset()
		})
	}
}

func TestHandlerSetServices(t *testing.T) {
	var bf bytes.Buffer
	type fields struct {
		comp *component.Components
		svc  *service.Services
		repo *repository.Repositories
	}

	type args struct {
		s *service.Services
	}

	tests := []struct {
		fields fields
		args   args
		name   string
	}{
		{
			name: "Ok",
			fields: fields{
				comp: nil,
				svc:  nil,
				repo: nil,
			},
			args: args{
				s: &service.Services{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				comp:   tt.fields.comp,
				svc:    tt.fields.svc,
				repo:   tt.fields.repo,
				writer: &bf,
			}
			h.SetServices(tt.args.s)
			bf.Reset()
		})
	}
}

func TestNewHandler(t *testing.T) {
	var bf bytes.Buffer
	tests := []struct {
		want *Handler
		name string
	}{
		{
			name: "Ok",
			want: &Handler{
				writer: &bf,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHandler(&bf); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHandler() = %v, want %v", got, tt.want)
			}

			bf.Reset()
		})
	}
}
