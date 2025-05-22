package jwt

import (
	"os"
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
	m.subject = os.Getenv("SUBJECT")
	m.issuer = os.Getenv("ISSUER")
	audience := os.Getenv("AUDIENCE")
	m.audience = append(m.audience, audience)
	signinMethod := os.Getenv("SIGNING_METHOD")
	switch signinMethod {
	default:
		m.signingMethod = jwt.SigningMethodHS256
	}
	expiry := os.Getenv("EXPIRY")
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
	_, skip := m.skip[path]
	if skip {
		return c.Next()
	}
	method := c.Method()
	var authToken string
	switch method {
	case fiber.MethodPost:
		authToken = c.Get("X-Authorization-Bearer")
		if authToken == "" {
			err = fiber.ErrUnauthorized
			return
		}
	case fiber.MethodGet:
		authToken = c.Get("user_id")
		if authToken == "" {
			err = fiber.ErrUnauthorized
			return
		}
	default:
		err = fiber.ErrBadRequest
		return
	}
	ok := m.verify(authToken)
	if !ok {
		err = fiber.ErrUnauthorized
		return
	}
	return c.Next()
}

func (m model) verify(signedToken string) (ok bool) {
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
	// if claims.ID != id {
	// 	return
	// }
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
