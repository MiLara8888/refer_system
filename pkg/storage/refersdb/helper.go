package storage

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	er "refers_rest/pkg/errors"

	"github.com/golang-jwt/jwt/v4"
)

const (
	saltSize = 10
	letters  = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	iter     = 32000
)

func hashPassword(pw, salt []byte) []byte {
	ret := make([]byte, len(salt))
	copy(ret, salt)
	return append(ret, Key(pw, salt, iter, sha256.Size, sha256.New)...)
}

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	return b, err
}

func GenPassword(pass, salt []byte) (string, string, string, error) {
	var err error

	if salt == nil {
		salt, err = GenerateRandomBytes(saltSize)
		if err != nil {
			return "", "", "", err
		}
	}
	hash := hashPassword([]byte(pass), salt)

	dst := make([]byte, hex.EncodedLen(len(hash)))
	hex.Encode(dst, hash)
	return "pbkdf2", hex.EncodeToString(hash), hex.EncodeToString(salt), nil
}

func PasswordMatched(passw, dbPassw, s string) error {

	salt, err := hex.DecodeString(s)

	if err != nil {
		return er.ErrSaltDecode
	}
	_, pasw, _, err := GenPassword([]byte(passw), salt)
	if err != nil {
		return er.ErrGenPassword
	}
	if dbPassw != pasw {
		return er.ErrIsNotMatch
	}
	return nil

}

// copied from "golang.org/x/crypto/pbkdf2" because it's not available in playground
func Key(password, salt []byte, iter, keyLen int, h func() hash.Hash) []byte {
	prf := hmac.New(h, password)
	hashLen := prf.Size()
	numBlocks := (keyLen + hashLen - 1) / hashLen

	var buf [4]byte
	dk := make([]byte, 0, numBlocks*hashLen)
	U := make([]byte, hashLen)
	for block := 1; block <= numBlocks; block++ {
		// N.B.: || means concatenation, ^ means XOR
		// for each block T_i = U_1 ^ U_2 ^ ... ^ U_iter
		// U_1 = PRF(password, salt || uint(i))
		prf.Reset()
		prf.Write(salt)
		buf[0] = byte(block >> 24)
		buf[1] = byte(block >> 16)
		buf[2] = byte(block >> 8)
		buf[3] = byte(block)
		prf.Write(buf[:4])
		dk = prf.Sum(dk)
		T := dk[len(dk)-hashLen:]
		copy(U, T)

		// U_n = PRF(password, U_(n-1))
		for n := 2; n <= iter; n++ {
			prf.Reset()
			prf.Write(U)
			U = U[:0]
			U = prf.Sum(U)
			for x := range U {
				T[x] ^= U[x]
			}
		}
	}
	return dk[:keyLen]
}

// токен
type JwtToken struct {
	CreateTime           int `json:"crt"`
	ExpDate              int `json:"exp"`
	Email                string
	jwt.RegisteredClaims `json:"-"`
}



// проверка токена
func TokenValid(tokenString, secretKey string) (*JwtToken, bool, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtToken{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, false, err
	}
	if !token.Valid {
		return nil, false, errors.New("")
	} else if errors.Is(err, jwt.ErrTokenMalformed) {
		return nil, false, err
	} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
		return nil, false, err
	}

	claims, ok := token.Claims.(*JwtToken)
	return claims, ok, nil
}
