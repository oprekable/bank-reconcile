package process

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/oprekable/bank-reconcile/internal/app/component"
	"github.com/oprekable/bank-reconcile/internal/app/component/csqlite"
)

func TestProviderDB(t *testing.T) {
	type args struct {
		comp *component.Components
	}

	tests := []struct {
		args    args
		want    *DB
		name    string
		wantErr bool
	}{
		{
			name: "Ok",
			args: args{
				comp: &component.Components{
					DBSqlite: &csqlite.DBSqlite{
						DBWrite: &sql.DB{},
					},
				},
			},
			want: &DB{
				db:      &sql.DB{},
				stmtMap: make(map[string]*sql.Stmt),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ProviderDB(tt.args.comp)

			if (err != nil) != tt.wantErr {
				t.Errorf("ProviderDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProviderDB() got = %v, want %v", got, tt.want)
			}
		})
	}
}
