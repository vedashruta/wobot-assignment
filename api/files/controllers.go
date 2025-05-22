package files

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"server/api/users"
	"server/env"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func upload(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "error": err.Error()})
	}
	userID := c.FormValue("user_id")
	if userID == "" {
		err = fmt.Errorf("user_id key not found")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "error": err.Error()})
	}
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	filter := bson.D{
		{
			Key:   "_id",
			Value: userID,
		},
	}
	var user users.User
	err = env.UsersCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "failure", "error": "unauthorized"})
	}
	remainingStorage := user.StorageLimit - user.UsedStorage
	fileSize := file.Size
	if fileSize > int64(remainingStorage) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "error": "file size exceeds your remaining storage quota"})
	}
	uuid := uuid.NewString()
	uploadPath := fmt.Sprintf("%[1]s/%[2]s", user.Username, uuid)
	err = os.MkdirAll(uploadPath, os.ModePerm)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "failure", "error": err.Error()})
	}
	localPath := fmt.Sprintf("%[1]s/%[2]s", uploadPath, uuid+filepath.Ext(file.Filename))
	err = c.SaveFile(file, localPath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "failure", "error": err.Error()})
	}
	meta := File{
		UserID:     user.ID,
		UUID:       uuid,
		FileName:   file.Filename,
		Extension:  filepath.Ext(file.Filename),
		Size:       int(file.Size),
		UploadedAt: time.Now(),
	}
	_, err = env.FilesCollection.InsertOne(ctx, meta)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "failure", "error": err.Error()})
	}
	filter = bson.D{
		{
			Key:   "_id",
			Value: user.ID,
		},
	}
	inc := bson.D{
		{
			Key: "$inc",
			Value: bson.D{
				{
					Key:   "used_storage",
					Value: file.Size,
				},
			},
		},
	}
	_, err = env.UsersCollection.UpdateOne(ctx, filter, inc)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "upload successful"})
}

func remaining(c *fiber.Ctx) error {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "error": "user_id required"})
	}
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "error": err.Error()})
	}
	var user users.User
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	filter := bson.D{
		{
			Key:   "_id",
			Value: userID,
		},
	}
	err = env.UsersCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	remaining := user.StorageLimit - user.UsedStorage
	response := fiber.Map{
		"total_storage":     user.StorageLimit,
		"used_storage":      user.UsedStorage,
		"remaining_storage": remaining,
	}
	return c.JSON(response)
}

func files(c *fiber.Ctx) error {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "error": "user_id required"})
	}
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "error": "invalid user_id"})
	}
	// Pagination params
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}
	if limit > 30 {
		limit = 20
	}
	skip := (page - 1) * limit
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	filter := bson.D{
		{
			Key:   "user_id",
			Value: userID,
		},
	}
	sort := bson.D{
		{
			Key:   "_id",
			Value: -1,
		},
	}
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))
	findOptions.SetSort(sort)

	cursor, err := env.FilesCollection.Find(ctx, filter, findOptions)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "failure", "error": "failed to fetch files"})
	}
	defer cursor.Close(ctx)
	var files []File
	err = cursor.All(ctx, &files)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "failure", "error": "failed to parse files"})
	}
	resp := make([]fiber.Map, len(files))
	for i, f := range files {
		resp[i] = fiber.Map{
			"uuid":        f.UUID,
			"filename":    f.FileName,
			"size":        f.Size,
			"uploaded_at": f.UploadedAt,
		}
	}
	return c.JSON(fiber.Map{
		"page":  page,
		"limit": limit,
		"files": resp,
	})
}
