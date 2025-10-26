package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type TypeModel interface {
	IModel
	any
}

type IModel interface {
	GetId() uuid.UUID
	SetId(id uuid.UUID)
	SetCreatedAt(createdAt int64)
	SetUpdatedAt(updatedAt int64)
}

type Model struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt int64     `json:"created_at"`
	UpdatedAt int64     `json:"updated_at"`
	DeletedAt *int64    `json:"deleted_at"`
}

func (m Model) GetId() uuid.UUID {
	return m.Id
}

func (m *Model) SetId(id uuid.UUID) {
	m.Id = id
}

func (m *Model) SetCreatedAt(createdAt int64) {
	m.CreatedAt = createdAt
}

func (m *Model) SetUpdatedAt(updatedAt int64) {
	m.UpdatedAt = updatedAt
}

type Jsonb map[string]any

func (j Jsonb) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *Jsonb) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSON value: %v", value)
	}
	data := map[string]any{}
	err := json.Unmarshal(bytes, &data)
	*j = Jsonb(data)
	return err
}

type JsonbArray []string

func (j JsonbArray) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	bytes, err := json.Marshal(j)
	return string(bytes), err
}

func (j *JsonbArray) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSON value: %v", value)
	}
	data := []string{}
	if err := json.Unmarshal(bytes, &data); err != nil {
		return fmt.Errorf("failed to unmarshal JSON value: %v [%v]", value, err)
	}
	*j = JsonbArray(data)
	return nil
}
