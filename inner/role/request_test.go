package role

import (
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateRequestValidation(t *testing.T) {
	tests := []struct {
		name     string
		request  CreateRequest
		wantErr  bool
		errField string
	}{
		{
			name: "valid request",
			request: CreateRequest{
				Name: "Valid name",
			},
			wantErr: false,
		},
		{
			name: "empty name",
			request: CreateRequest{
				Name: "",
			},
			wantErr:  true,
			errField: "Name",
		},
		{
			name: "short name",
			request: CreateRequest{
				Name: "J",
			},
			wantErr:  true,
			errField: "Name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.New().Struct(tt.request)
			if err != nil {
				return
			}

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errField)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
