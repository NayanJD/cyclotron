package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/mail"
	"os"
	"strings"
	"time"
	"user/pkg/domains/user"
	customErrors "user/pkg/errors"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	jwt "github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	InvalidPasswordErr     = customErrors.ConstError("The username and password provided is invalid")
	InternalServerErr      = customErrors.ConstError("Internal Server Error")
	InvalidRefreshTokenErr = customErrors.ConstError("Invalid Refresh Token")
)

// UserService describes the service.
type UserService interface {
	// Add your methods here
	Login(ctx context.Context, username, password string) (token user.AuthToken, err error)

	Register(ctx context.Context, user user.User) (newUser user.User, err error)

	GetUserFromToken(ctx context.Context, token string) (user user.User, err error)

	RefreshAccessToken(ctx context.Context, refreshToken string) (token user.AuthToken, err error)
}

type basicUserService struct {
	userRepository  user.UserRepository
	tokenRepository user.AuthTokenRepository
	jwtSecret       string
	logger          log.Logger
}

type JwtClaims struct {
	UserId  string `json:"user_id"`
	TokenId string `json:"token_id"`
	jwt.RegisteredClaims
}

func (b *basicUserService) Login(ctx context.Context, username string, password string) (token user.AuthToken, err error) {

	var rUser *user.User

	if rUser, err = b.userRepository.FindByUsername(ctx, username); err != nil {
		level.Error(b.logger).Log("msg", err)
		return token, InvalidPasswordErr
	}

	if err = bcrypt.CompareHashAndPassword([]byte(rUser.HashedPassword), []byte(password)); err != nil {
		level.Error(b.logger).Log("msg", err)

		return token, InvalidPasswordErr
	}

	bits := make([]byte, 12)

	_, err = rand.Read(bits)
	if err != nil {
		level.Error(b.logger).Log("token", "accessToken", "msg", err)
		return token, InternalServerErr
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, JwtClaims{
		UserId:  rUser.ID.String(),
		TokenId: base64.StdEncoding.EncodeToString(bits),
	})

	_, err = rand.Read(bits)
	if err != nil {
		level.Error(b.logger).Log("token", "refreshToken", "msg", err)
		return token, InternalServerErr
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, JwtClaims{
		UserId:  rUser.ID.String(),
		TokenId: base64.StdEncoding.EncodeToString(bits),
	})

	// Sign and get the complete encoded token as a string using the secret
	accessTokenString, err := accessToken.SignedString([]byte(b.jwtSecret))

	if err != nil {
		level.Error(b.logger).Log("token", "accessToken", "msg", err)
		return token, InternalServerErr
	}

	refreshTokenString, err := refreshToken.SignedString([]byte(b.jwtSecret))

	if err != nil {
		level.Error(b.logger).Log("token", "refreshToken", "msg", err)
		return token, InternalServerErr
	}

	newToken := &user.AuthToken{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		UserId:       rUser.ID.String(),
	}

	if newToken, err = b.tokenRepository.CreateToken(ctx, newToken, 900); err != nil {
		return user.AuthToken{}, err
	}

	return *newToken, err
}

// NewBasicUserService returns a naive, stateless implementation of UserService.
func NewBasicUserService(userRepository user.UserRepository, tokenRepository user.AuthTokenRepository, logger log.Logger) (UserService, error) {

	jwtSecret, ok := os.LookupEnv("JWT_SECRET")

	if !ok {
		level.Error(logger).Log("msg", "JWT_SECRET env var not provided. Unable to start service")
		return nil, fmt.Errorf("JWT_SECRET env var not provided. Unable to start service")
	}
	return &basicUserService{
		userRepository:  userRepository,
		tokenRepository: tokenRepository,
		jwtSecret:       jwtSecret,
		logger:          logger,
	}, nil
}

// New returns a UserService with all of the expected middleware wired in.
func New(ur user.UserRepository, tr user.AuthTokenRepository, logger log.Logger, middleware []Middleware) (UserService, error) {
	svc, err := NewBasicUserService(ur, tr, logger)

	if err != nil {
		return svc, err
	}

	for _, m := range middleware {
		svc = m(svc)
	}
	return svc, err
}

func (b *basicUserService) Register(ctx context.Context, user user.User) (newUser user.User, err error) {

	var hasValidationErrors bool
	var validationErrors []string

	if len(user.FirstName) < 2 {
		hasValidationErrors = true
		validationErrors = append(validationErrors, "First name should atleast be 2 chars long")
	}

	if len(user.LastName) < 2 {
		hasValidationErrors = true
		validationErrors = append(validationErrors, "Last name should atleast be 2 chars long")
	}

	if _, err := mail.ParseAddress(user.Username); err != nil {
		hasValidationErrors = true
		validationErrors = append(validationErrors, "Username should be a valid email address")
	}

	if len(user.Password) < 8 {
		hasValidationErrors = true
		validationErrors = append(validationErrors, "Password should be atleast 8 characters long")
	}

	if time.Since(user.Dob) > time.Duration(200*365*24*time.Hour) {
		hasValidationErrors = true
		validationErrors = append(validationErrors, "Impossible DOB provided")
	}

	if hasValidationErrors {
		return newUser, fmt.Errorf("Validation Errors: %s", strings.Join(validationErrors, "\n"))
	}

	var hashedPasswordBytes []byte

	if hashedPasswordBytes, err = bcrypt.GenerateFromPassword([]byte(user.Password), 10); err != nil {
		return newUser, fmt.Errorf("Password could not be hashed!")
	}

	user.HashedPassword = string(hashedPasswordBytes)

	userPtr, err := b.userRepository.CreateUser(ctx, &user)

	if err != nil {
		fmt.Printf(err.Error(), "\n")
		return
	}
	newUser = *userPtr
	return
}

func (b *basicUserService) GetUserFromToken(ctx context.Context, token string) (u user.User, err error) {
	parsedToken, err := jwt.ParseWithClaims(token, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(b.jwtSecret), nil
	})

	if err != nil {
		return u, err
	}

	var claims *JwtClaims

	claims, ok := parsedToken.Claims.(*JwtClaims)

	if ok && parsedToken.Valid {
		level.Debug(b.logger).Log("issuer", claims.RegisteredClaims.Issuer)
	} else {
		return u, err
	}

	var dUser *user.User
	if dUser, err = b.userRepository.FindByID(ctx, claims.UserId); err != nil {
		return u, err
	}

	return *dUser, err
}

func (b *basicUserService) RefreshAccessToken(ctx context.Context, refreshToken string) (token user.AuthToken, err error) {
	parsedToken, err := jwt.ParseWithClaims(refreshToken, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(b.jwtSecret), nil
	})

	if err != nil {
		return token, err
	}

	var claims *JwtClaims

	claims, ok := parsedToken.Claims.(*JwtClaims)

	if ok && parsedToken.Valid {
		level.Debug(b.logger).Log("issuer", claims.RegisteredClaims.Issuer)
	} else {
		return token, err
	}

	if _, err = b.userRepository.FindByID(ctx, claims.UserId); err != nil {
		return token, InvalidRefreshTokenErr
	}

	bits := make([]byte, 12)

	_, err = rand.Read(bits)
	if err != nil {
		level.Error(b.logger).Log("token", "accessToken", "msg", err)
		return token, InternalServerErr
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, JwtClaims{
		UserId:  claims.UserId,
		TokenId: base64.StdEncoding.EncodeToString(bits),
	})

	_, err = rand.Read(bits)
	if err != nil {
		level.Error(b.logger).Log("token", "refreshToken", "msg", err)
		return token, InternalServerErr
	}

	newRefreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, JwtClaims{
		UserId:  claims.UserId,
		TokenId: base64.StdEncoding.EncodeToString(bits),
	})

	// Sign and get the complete encoded token as a string using the secret
	accessTokenString, err := accessToken.SignedString([]byte(b.jwtSecret))

	if err != nil {
		level.Error(b.logger).Log("token", "accessToken", "msg", err)
		return token, InternalServerErr
	}

	refreshTokenString, err := newRefreshToken.SignedString([]byte(b.jwtSecret))

	if err != nil {
		level.Error(b.logger).Log("token", "refreshToken", "msg", err)
		return token, InternalServerErr
	}

	newToken := &user.AuthToken{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		UserId:       claims.UserId,
	}

	if newToken, err = b.tokenRepository.CreateToken(ctx, newToken, 900); err != nil {
		return user.AuthToken{}, err
	}

	return *newToken, err
}
