package core

import "testing"

func TestErrorTypeError(t *testing.T) {
	tests := []struct {
		name    string
		wantErr string
		e       ErrorType
	}{
		{
			name:    "Ok",
			e:       CErrDBConn,
			wantErr: "100001",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.e.Error()
			if err.Error() != tt.wantErr {
				t.Errorf("Error() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
