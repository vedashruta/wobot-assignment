package jwt

import (
	"os"
	"server/services/json"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Init(fileName string) (err error) {
	err = godotenv.Load(fileName)
	if err != nil {
		return
	}
	m = &model{}
	secret := os.Getenv("JWT_SECRET")
	m.secret = []byte(secret)
	m.subject = os.Getenv("JWT_SUBJECT")
	m.issuer = os.Getenv("JWT_ISSUER")
	audience := os.Getenv("JWT_AUDIENCE")
	m.audience = append(m.audience, audience)
	signinMethod := os.Getenv("SIGNING_METHOD")
	switch signinMethod {
	default:
		m.signingMethod = jwt.SigningMethodHS256
	}
	expiry := os.Getenv("JWT_EXPIRY")
	duration, err := time.ParseDuration(expiry)
	if err != nil {
		return
	}
	if duration.Seconds() <= 100 {
		duration = 1800 * time.Second
	}
	m.expiry = duration
	pathsStr := os.Getenv("JWT_SKIP")
	paths := strings.Split(pathsStr, ",")
	m.skip = make(map[string]struct{}, len(paths))
	for _, path := range paths {
		m.skip[path] = struct{}{}
	}
	return
}

func Auth(c *fiber.Ctx) (err error) {
	path := c.Path()
	_, ok := m.skip[path]
	if !ok {
		method := c.Method()
		var userID primitive.ObjectID
		switch method {
		case fiber.MethodPost:
			contentType := string(c.Request().Header.ContentType())
			if strings.Contains(contentType, fiber.MIMEMultipartForm) {
				contentType = fiber.MIMEMultipartForm
			}
			switch contentType {
			case fiber.MIMEMultipartForm:
				userID, err = primitive.ObjectIDFromHex(c.FormValue("user_id"))
				if err != nil {
					return
				}
			case fiber.MIMEApplicationJSON, fiber.MIMETextPlain:
				reqBody := map[string]any{}
				body := c.Request().Body()
				err = json.Decode(body, &reqBody)
				if err != nil {
					return
				}
				userIDStr, ok := reqBody["user_id"].(string)
				if !ok {
					err = fiber.ErrBadRequest
					return
				}
				userID, err = primitive.ObjectIDFromHex(userIDStr)
				if err != nil {
					return
				}
			default:
				err = fiber.ErrBadRequest
				return
			}
		case fiber.MethodGet:
			userID, err = primitive.ObjectIDFromHex(c.Query("user_id"))
			if err != nil {
				return
			}
		default:
			err = fiber.ErrBadRequest
			return
		}
		id := userID.Hex()
		signature := c.Cookies("token")
		ok := m.verify(signature, id)
		if !ok {
			err = fiber.ErrUnauthorized
			return
		}
	}
	return c.Next()
	// }
	// var authToken string
	// authToken = c.Get("X-Authorization-Bearer")
	// if authToken == "" {
	// 	authToken = c.Cookies("token")
	// }
	// if authToken == "" {
	// 	err = fiber.ErrUnauthorized
	// 	return
	// }
	// ok := m.verify(authToken)
	// if !ok {
	// 	err = fiber.ErrUnauthorized
	// 	return
	// }
	// return c.Next()
}

func (m model) verify(signedToken string, id string) (ok bool) {
	claims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(signedToken, &claims, func(token *jwt.Token) (interface{}, error) {
		return m.secret, nil
	})
	if err != nil {
		return
	}
	if !token.Valid {
		return
	}
	if claims.Issuer != m.issuer {
		return
	}
	if claims.Subject != m.subject {
		return
	}
	if len(m.audience) != len(claims.Audience) {
		return
	}
	if claims.ID != id {
		return
	}
	for i := 0; i < len(m.audience); i++ {
		if m.audience[i] != claims.Audience[i] {
			return
		}
	}
	ok = true
	return
}

func NewClaims(userID primitive.ObjectID) (signedToken string, err error) {
	claims := jwt.RegisteredClaims{
		Issuer:    m.issuer,
		Subject:   m.subject,
		Audience:  m.audience,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.expiry)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        userID.Hex(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err = token.SignedString([]byte(m.secret))
	if err != nil {
		return
	}
	return
}
