package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"github.com/DuvanAlbarracin/movies_auth/pkg/models"
	"github.com/golang-jwt/jwt/v5"
)

type JwtWrapper struct {
	SecretKey       string
	Issuer          string
	ExpirationHours int64
}

type jwtClaims struct {
	jwt.RegisteredClaims
	Id    int64
	Email string
}

func (w *JwtWrapper) GenerateToken(u models.User) (sT string, err error) {
	claims := &jwtClaims{
		Id:    u.Id,
		Email: u.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Local().Add(
					time.Hour * time.Duration(w.ExpirationHours)),
			),
			Issuer: w.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	privateKey, err := generatePrivateKey(w.SecretKey)
	if err != nil {
		return
	}

	sT, err = token.SignedString(privateKey)
	if err != nil {
		return
	}

	return sT, nil
}

func (w *JwtWrapper) ValidateToken(sT string) (claims *jwtClaims, err error) {
	token, err := jwt.ParseWithClaims(
		sT,
		&jwtClaims{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(w.SecretKey), nil
		},
	)

	if err != nil {
		return
	}

	claims, ok := token.Claims.(*jwtClaims)
	if !ok {
		return nil, errors.New("Could not parse claims")
	}

	if claims.ExpiresAt.Time.Unix() < time.Now().Local().Unix() {
		return nil, errors.New("JWT is expired")
	}

	return claims, nil
}

func generatePrivateKey(secret string) (*ecdsa.PrivateKey, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	hash := sha256.Sum256([]byte(secret))

	// hash is an [32]byte array; hash[:] => returns an slice of the whole hash array
	// hash[1:3] => left most inclusive, right most exclusive
	sig, err := ecdsa.SignASN1(rand.Reader, privateKey, hash[:])
	if err != nil {
		return nil, err
	}

	fmt.Printf("signature: %x\n", sig)

	valid := ecdsa.VerifyASN1(&privateKey.PublicKey, hash[:], sig)

	fmt.Println("signature verified:", valid)

	return privateKey, err
}
