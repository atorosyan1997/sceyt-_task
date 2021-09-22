package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"my-bank-service/internal/config"
	"my-bank-service/internal/data"
	"my-bank-service/internal/reposytory"
	"my-bank-service/pkg/logging"
	"time"
)

// Authentication interface lists the methods that our authentication service should implement
type Authentication interface {
	Authenticate(reqUser *data.User, user *data.User) bool
	GenerateAccessToken(authD *data.AuthDetails) (string, error)
	GenerateTokens(user *data.User) (string, string, error)
	GenerateCustomKey(userID string, password string) string
	ValidateAccessToken(token string) (string, string, error)
	ValidateRefreshToken(token string) (string, string, string, error)
	InvalidateTokens(authD *data.AuthDetails) error
}

// AccessTokenCustomClaims specifies the claims for access token
type AccessTokenCustomClaims struct {
	UserID   string
	AuthUUID string
	KeyType  string
	jwt.StandardClaims
}

// RefreshTokenCustomClaims specifies the claims for refresh token
type RefreshTokenCustomClaims struct {
	UserID    string
	AuthUUID  string
	CustomKey string
	KeyType   string
	jwt.StandardClaims
}

// AuthService is the implementation of our Authentication
type AuthService struct {
	logger  logging.Logger
	configs *config.Configurations
	repo    reposytory.AuthRepository
}

// NewAuthService returns a new instance of the auth service
func NewAuthService(logger logging.Logger, configs *config.Configurations, repo reposytory.AuthRepository) *AuthService {
	return &AuthService{logger, configs, repo}
}

// Authenticate checks the user credentials in request against the db and authenticates the request
func (a *AuthService) Authenticate(reqUser *data.User, user *data.User) bool {

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqUser.Password)); err != nil {
		a.logger.Debug("password hashes are not same")
		return false
	}
	return true
}

// GenerateTokens generate a new access and refresh tokens for the given user
func (a *AuthService) GenerateTokens(user *data.User) (string, string, error) {
	authD := &data.AuthDetails{}
	authD.AuthUuid = uuid.NewV4().String() //generate a new UUID each time
	authD.UserId = user.ID

	accessToken, err := a.GenerateAccessToken(authD)
	if err != nil {
		a.logger.Error("unable to generate access token", "error", err)
		return "", "", err
	}

	refreshToken, err := a.generateRefreshToken(user, authD)
	if err != nil {
		a.logger.Error("unable to generate refresh token", "error", err)
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

// generateRefreshToken generate a new refresh token for the given user
func (a *AuthService) generateRefreshToken(user *data.User, authD *data.AuthDetails) (string, error) {

	// since after the user logged out, we destroyed that record in the database
	// so that same jwt token can't be used twice. We need to create the token again
	authData, err := a.repo.CreateAuth(authD)
	if err != nil {
		a.logger.Error("unable to creat refresh auth key", "error", err)
		return "", errors.New("could not generate refresh key. please try again later")
	}
	cusKey := a.GenerateCustomKey(user.ID, user.TokenHash)
	tokenType := "refresh"

	claims := RefreshTokenCustomClaims{
		authData.UserID,
		authData.AuthUUID,
		cusKey,
		tokenType,
		jwt.StandardClaims{
			Issuer: "bookite.auth.service",
		},
	}

	signBytes, err := ioutil.ReadFile(a.configs.RefreshTokenPrivateKeyPath)
	if err != nil {
		a.logger.Error("unable to read private key", "error", err)
		return "", errors.New("could not generate refresh token. please try again later")
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		a.logger.Error("unable to parse private key", "error", err)
		return "", errors.New("could not generate refresh token. please try again later")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	return token.SignedString(signKey)
}

// GenerateAccessToken generates a new access token for the given user
func (a *AuthService) GenerateAccessToken(authD *data.AuthDetails) (string, error) {

	// since after the user logged out, we destroyed that record in the database
	// so that same jwt token can't be used twice. We need to create the token again
	authData, err := a.repo.CreateAuth(authD)
	if err != nil {
		a.logger.Error("unable to create access auth key", "error", err)
		return "", errors.New("could not generate refresh key. please try again later")
	}
	tokenType := "access"

	claims := AccessTokenCustomClaims{
		authData.UserID,
		authData.AuthUUID,
		tokenType,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(a.configs.JwtExpiration)).Unix(),
			Issuer:    "bookite.auth.service",
		},
	}

	signBytes, err := ioutil.ReadFile(a.configs.AccessTokenPrivateKeyPath)
	if err != nil {
		a.logger.Error("unable to read private key", "error", err)
		return "", errors.New("could not generate access token. please try again later")
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		a.logger.Error("unable to parse private key", "error", err)
		return "", errors.New("could not generate access token. please try again later")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	return token.SignedString(signKey)
}

// GenerateCustomKey creates a new key for our jwt payload
// the key is a hashed combination of the userID and user tokenhash
func (a *AuthService) GenerateCustomKey(userID string, tokenHash string) string {

	// data := userID + tokenHash
	h := hmac.New(sha256.New, []byte(tokenHash))
	h.Write([]byte(userID))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}

// ValidateAccessToken parses and validates the given access token
// returns the userId present in the token payload
func (a *AuthService) ValidateAccessToken(tokenString string) (string, string, error) {

	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			a.logger.Error("unexpected signing method in auth token")
			return nil, errors.New("unexpected signing method in auth token")
		}
		verifyBytes, err := ioutil.ReadFile(a.configs.AccessTokenPublicKeyPath)
		if err != nil {
			a.logger.Error("unable to read public key", "error", err)
			return nil, err
		}

		verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
		if err != nil {
			a.logger.Error("unable to parse public key", "error", err)
			return nil, err
		}

		return verifyKey, nil
	})

	if err != nil {
		a.logger.Error("unable to parse claims", "error", err)
		return "", "", err
	}

	claims, ok := token.Claims.(*AccessTokenCustomClaims)
	if !ok || !token.Valid || claims.UserID == "" || claims.KeyType != "access" {
		return "", "", errors.New("invalid token: authentication failed")
	}

	authD := &data.AuthDetails{
		UserId:   claims.UserID,
		AuthUuid: claims.AuthUUID,
	}
	_, err = a.repo.FetchAuth(authD)
	if err != nil {
		return "", "", err
	}
	return claims.UserID, claims.AuthUUID, nil
}

// ValidateRefreshToken parses and validates the given refresh token
// returns the userId and customkey present in the token payload
func (a *AuthService) ValidateRefreshToken(tokenString string) (string, string, string, error) {

	token, err := jwt.ParseWithClaims(tokenString, &RefreshTokenCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			a.logger.Error("unexpected signing method in auth token")
			return nil, errors.New("unexpected signing method in auth token")
		}
		verifyBytes, err := ioutil.ReadFile(a.configs.RefreshTokenPublicKeyPath)
		if err != nil {
			a.logger.Error("unable to read public key", "error", err)
			return nil, err
		}

		verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
		if err != nil {
			a.logger.Error("unable to parse public key", "error", err)
			return nil, err
		}

		return verifyKey, nil
	})

	if err != nil {
		a.logger.Error("unable to parse claims", "error", err)
		return "", "", "", err
	}

	claims, ok := token.Claims.(*RefreshTokenCustomClaims)
	a.logger.Debug("ok", ok)
	if !ok || !token.Valid || claims.UserID == "" || claims.KeyType != "refresh" {
		a.logger.Debug("could not extract claims from token")
		return "", "", "", errors.New("invalid token: authentication failed")
	}

	authD := &data.AuthDetails{
		UserId:   claims.UserID,
		AuthUuid: claims.AuthUUID,
	}
	_, err = a.repo.FetchAuth(authD)
	if err != nil {
		return "", "", "", err
	}

	return claims.UserID, claims.AuthUUID, claims.CustomKey, nil
}

// InvalidateTokens invalidates the user's tokens
func (a *AuthService) InvalidateTokens(authD *data.AuthDetails) error {
	err := a.repo.DeleteAuth(authD)
	if err != nil {
		a.logger.Error(err)
		return err
	}
	return nil
}
