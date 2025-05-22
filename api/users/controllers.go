package users

import (
	"context"
	"fmt"
	"server/env"
	"server/middlewares/jwt"
	"server/services/bcrypt"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func register(c *fiber.Ctx) (err error) {
	req := Request{}
	err = c.BodyParser(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "failure",
			"error":  err.Error(),
		})
	}
	if req.Username == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "failure",
			"error":  "username and password are required",
		})
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{
		{
			Key:   "username",
			Value: req.Username,
		},
	}
	var res User
	var ok bool
	err = env.UsersCollection.FindOne(ctx, filter).Decode(&res)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "failure",
				"error":  err.Error(),
			})
		} else {
			ok = true
		}
	}
	if !ok {
		err = fmt.Errorf("user already registered")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "failure",
			"error":  err.Error(),
		})
	}
	hashedPwd, err := bcrypt.Hash(req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "failure",
			"error":  err.Error(),
		})
	}
	user := User{
		Username:     req.Username,
		Password:     hashedPwd,
		UsedStorage:  0,
		StorageLimit: env.DefaultStorageQuota,
	}
	_, err = env.UsersCollection.InsertOne(ctx, user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "failure",
			"error":  err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "user registered successfully",
	})
}

func login(c *fiber.Ctx) (err error) {
	req := Request{}
	err = c.BodyParser(&req)
	if err != nil {
		return c.JSON(fiber.Map{
			"status": "failure",
			"error":  err.Error(),
		})
	}
	var user User
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	filter := bson.D{
		{
			Key:   "username",
			Value: req.Username,
		},
	}
	err = env.UsersCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return c.JSON(fiber.Map{
			"status": "failure",
			"error":  err.Error(),
		})
	}
	ok, err := bcrypt.Verify(req.Password, user.Password)
	if err != nil {
		return c.JSON(fiber.Map{
			"status": "failure",
			"error":  err.Error(),
		})
	}
	if !ok {
		err = fmt.Errorf("invalid credentials")
		return c.JSON(fiber.Map{
			"status": "failure",
			"error":  err.Error(),
		})
	}
	signedToken, err := jwt.NewClaims(user.ID)
	if err != nil {
		return c.JSON(fiber.Map{
			"status": "failure",
			"error":  err.Error(),
		})
	}
	expiry := 30 * time.Minute
	c.Cookie(&fiber.Cookie{
		Name:     "jwt_token",
		Value:    signedToken,
		Expires:  time.Now().Add(expiry),
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Strict",
		Path:     "/",
	})
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "login successful",
		"token":   signedToken,
	})
}
