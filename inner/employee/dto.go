package employee

import "time"

type Entity struct {
	Id        int64     `db:"id"`
	Name      string    `db:"name"`
	RoleID    *int64    `db:"role_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (e *Entity) toResponse() Response {
	return Response{
		Id:        e.Id,
		Name:      e.Name,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

func toSliceResponse(e []Entity) []Response {
	responses := make([]Response, len(e))
	for i := range e {
		responses[i] = e[i].toResponse()
	}
	return responses
}

type Response struct {
	Id        int64     `*json:"id"`
	Name      string    `*json:"name"`
	RoleId    *int64    `*json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateRequest struct {
	Name   string `json:"name" validate:"required,min=2,max=155"`
	RoleId *int64 `json:"roleId" validate:"required,min=1,max=155"`
}

func (req *CreateRequest) ToEntity() Entity {
	return Entity{Name: req.Name,
		RoleID: req.RoleId}
}
