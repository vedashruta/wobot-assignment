package files

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type File struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID     primitive.ObjectID `json:"user_id" bson:"user_id"`
	FileName   string             `json:"file_name" bson:"file_name"`
	Extension  string             `json:"extension" bson:"extension"`
	UUID       string             `json:"uuid" bson:"uuid"`
	Size       int                `json:"size" bson:"size"`
	UploadedAt time.Time          `json:"uploaded_at" bson:"uploaded_at"`
}
