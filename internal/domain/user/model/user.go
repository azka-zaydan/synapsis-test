package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
)

type User struct {
	ID            uuid.UUID     `db:"id"`
	Username      string        `db:"username"`
	Password      string        `db:"password"`
	MetaCreatedAt time.Time     `db:"meta_created_at"`
	MetaUpdatedAt time.Time     `db:"meta_updated_at"`
	MetaDeletedAt null.Time     `db:"meta_deleted_at"`
	CreatedBy     uuid.UUID     `db:"created_by"`
	UpdatedBy     uuid.UUID     `db:"updated_by"`
	DeletedBy     uuid.NullUUID `db:"deleted_by"`
}
