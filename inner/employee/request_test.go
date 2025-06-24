package employee

import (
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestCreateRequestValidation(t *testing.T) {
	roleId := int64(1)
	tests := []struct {
		name     string
		request  CreateRequest
		wantErr  bool
		errField string
	}{
		{
			name: "valid request",
			request: CreateRequest{
				Name:   "Valid name",
				RoleId: &roleId,
			},
			wantErr: false,
		},
		{
			name: "empty name",
			request: CreateRequest{
				Name:   "",
				RoleId: &roleId,
			},
			wantErr:  true,
			errField: "Name",
		},
		{
			name: "short name",
			request: CreateRequest{
				Name:   "J",
				RoleId: &roleId,
			},
			wantErr:  true,
			errField: "Name",
		},
		{
			name: "invalid roleId",
			request: CreateRequest{
				Name:   "John Doe",
				RoleId: nil,
			},
			wantErr:  true,
			errField: "RoleId",
		},
		{
			name: "roleId too long",
			request: CreateRequest{
				Name:   strings.Repeat("a", 156),
				RoleId: &roleId,
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
