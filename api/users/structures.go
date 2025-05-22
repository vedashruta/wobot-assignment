package users

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Request struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Username     string             `json:"username" bson:"username"`
	Password     string             `json:"password" bson:"password"`
	StorageLimit int                `json:"storage_limit" bson:"storage_limit"`
	UsedStorage  int                `json:"used_storage" bson:"used_storage"`
}
