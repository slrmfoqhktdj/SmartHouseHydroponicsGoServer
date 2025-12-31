package jwt

import (
    "encoding/json"
	"fmt"
	"time"
    "log"
	"os"
    "net/http"
    "strings"
    "context"

	"github.com/dgrijalva/jwt-go"
	"smart_house/gosrc/structs"
)

// KeyStore는 비밀키를 관리합니다.
// type KeyStore struct {
// 	keys map[int64]string // 시간 기반 키 저장
// }

// // NewKeyStore는 새로운 KeyStore를 생성합니다.
// func NewKeyStore() *KeyStore {
// 	return &KeyStore{
// 		keys: make(map[int64]string),
// 	}
// }

// // SetKey는 주어진 시간에 대한 비밀키를 설정합니다.
// func (ks *KeyStore) SetKey(t time.Time, key string) {
// 	ks.keys[t.Unix()] = key
// }

// // GetKey는 주어진 시간에 가장 가까운 비밀키를 반환합니다.
// func (ks *KeyStore) GetKey(t time.Time) string {
// 	// 여기서는 단순화를 위해 직접 구현하지 않고 개념만 설명합니다.
// 	// 실제로는 주어진 시간에 가장 가까운 키를 찾아 반환해야 합니다.
// 	return ks.keys[t.Unix()]
// }

var mySigningKey = []byte(os.Getenv("JWT_SIGNKEY"))

// 토큰 생성
// func GenerateToken(userId string) (string, error) {
// 	token := jwt.New(jwt.SigningMethodHS256)
// 	claims := token.Claims.(jwt.MapClaims)

// 	claims["authorized"] = true
// 	claims["user"] = userId
// 	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

// 	tokenString, err := token.SignedString(mySigningKey)

// 	if err != nil {
// 		return "", fmt.Errorf("토큰생성실패: %s", err.Error())
// 	}

// 	return tokenString, nil
// }


func JwtMiddleware(next http.Handler) http.Handler {

	// log.Print("jwtMiddleware")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")

		if tokenString == "" {
			log.Print("jwtMiddleware err 1")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

		// log.Print("JWT toke is " + tokenString)

		token, err := ParseToken(tokenString)
		if err != nil {
			log.Print("Error parsing token")
			//유효기간 시간 초과시
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

        // 토큰에서 사용자 정보 추출
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Print("jwtMiddleware err 3")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

        //토큰 타입 확인
        tokenType, ok := claims["type"].(string)
        if !ok || tokenType != "access" {
            log.Print("jwtMiddleware err 4")
            w.WriteHeader(http.StatusUnauthorized)
            return
        }

        // 토큰 만료 시간 확인
        exp, ok := claims["exp"].(float64)
        if !ok {
            log.Print("jwtMiddleware err 5")
            w.WriteHeader(http.StatusUnauthorized)
            return
        }

        if int64(exp) < time.Now().Unix() {
            log.Print("jwtMiddleware err 6")
            w.WriteHeader(http.StatusUnauthorized)
            return
        }

        // 사용자 ID 추출
		userID, ok := claims["user"].(string)
		if !ok {
			log.Print("jwtMiddleware err 6")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), structs.ContextStruct{}, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// 토큰 검증
func ParseToken(tokenStr string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return mySigningKey, nil
	})

    // 토큰 파싱 중 에러 발생 시
    if err != nil {
        return nil, err
    }

    // 토큰 유효성 확인
    if !token.Valid {
        return nil, fmt.Errorf("invalid token")
    }
    
	return token, err
}


// 토큰 페어 생성 함수
func GenerateTokenPair(userEmail, login_type string) (*structs.TokenPair, error) {

    // 규격화된 타입 아니면 에러 반환
	if login_type != "device" && login_type != "app" && login_type != "web"	&& login_type != "admin" {
		return nil, fmt.Errorf("invalid login type: %s", login_type)
	}

    // Access Token (15분)
    accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user": userEmail,
        "exp":  time.Now().Add(30 * time.Minute).Unix(),
        "iat":  time.Now().Unix(),
        "type": "access",
    })

    accessTokenString, err := accessToken.SignedString(mySigningKey)
    if err != nil {
        return nil, err
    }

    // Refresh Token (30일)
    refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user": userEmail,
        "exp":  time.Now().Add(30 * 24 * time.Hour).Unix(),
        "iat":  time.Now().Unix(),
        "type": "refresh",
    })

    refreshTokenString, err := refreshToken.SignedString(mySigningKey)
    if err != nil {
        return nil, err
    }

    if login_type == "device" {
        return &structs.TokenPair{
            Access_Token:  accessTokenString,
            Refresh_Token: "",
            ExpiresIn:    1800, // 30분 (초 단위)
        }, nil
    }

    saveRefreshTokenHandler(userEmail, refreshTokenString, login_type)

    return &structs.TokenPair{
        Access_Token:  accessTokenString,
        Refresh_Token: refreshTokenString,
        ExpiresIn:    1800, // 30분 (초 단위)
    }, nil
}

//웹 Refresh Token으로 새 Access Token 발급
func RefreshToken_Handler(w http.ResponseWriter, r *http.Request) {

    var receiveresult structs.RefreshTokenRequest

    err := json.NewDecoder(r.Body).Decode(&receiveresult)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    defer r.Body.Close()

    // Refresh Token 검증
    token, err := jwt.Parse(receiveresult.RefreshToken, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method")
        }
        return mySigningKey, nil
    })

    if err != nil || !token.Valid {
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(map[string]string{
            "error": "Invalid refresh token",
        })
        return
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        w.WriteHeader(http.StatusUnauthorized)
        return
    }

    // 타입 확인
    if claims["type"] != "refresh" {
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(map[string]string{
            "error": "Invalid token type",
        })
        return
    }

    userEmail := claims["user"].(string)

    // DB에 저장된 Refresh Token과 비교 (보안 강화) = 토큰을 탈취 당했을 경우에 db에 저장된 토큰을 지움으로써 무효화 시킬 수 있다.
    if !validateRefreshToken(userEmail, receiveresult.RefreshToken) {
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(map[string]string{
            "error": "Refresh token not found or expired",
        })
        return
    }

    // 새 토큰 페어 생성
    newTokenPair, err := GenerateTokenPair(userEmail,"web")
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    // // 새 Refresh Token을 DB에 저장하고 기존 것 삭제
    // updateRefreshToken(userEmail, receiveresult.RefreshToken, newTokenPair.Refresh_Token)


    json, _ := json.Marshal(&structs.RenewedAccessToken{Success: true, Email: userEmail, TokenPair: *newTokenPair})
    w.WriteHeader(http.StatusOK)
    w.Header().Set("Content-Type", "application/json")
    w.Write(json)
}

//어플 Refresh Token으로 새 Access Token 발급
func RefreshToken_App_Handler(w http.ResponseWriter, r *http.Request) {

    var receiveresult structs.RefreshTokenAppRequest

    err := json.NewDecoder(r.Body).Decode(&receiveresult)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    defer r.Body.Close()

    // Refresh Token 검증
    token, err := jwt.Parse(receiveresult.RefreshToken, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method")
        }
        return mySigningKey, nil
    })

    if err != nil || !token.Valid {
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(map[string]string{
            "error": "Invalid refresh token",
        })
        return
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        w.WriteHeader(http.StatusUnauthorized)
        return
    }

    // 타입 확인
    if claims["type"] != "refresh" {
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(map[string]string{
            "error": "Invalid token type",
        })
        return
    }

    userEmail := claims["user"].(string)

    // DB에 저장된 Refresh Token과 비교 (보안 강화) = 토큰을 탈취 당했을 경우에 db에 저장된 토큰을 지움으로써 무효화 시킬 수 있다.
    if !validateRefreshToken(userEmail, receiveresult.RefreshToken) {
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(map[string]string{
            "error": "Refresh token not found or expired",
        })
        return
    }

    // 새 토큰 페어 생성
    newTokenPair, err := GenerateTokenPair(userEmail,"app")
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    // // 새 Refresh Token을 DB에 저장하고 기존 것 삭제
    // updateRefreshToken(userEmail, receiveresult.RefreshToken, newTokenPair.Refresh_Token)


    json, _ := json.Marshal(&structs.RenewedAccessToken{Success: true, Email: userEmail, TokenPair: *newTokenPair})
    w.WriteHeader(http.StatusOK)
    w.Header().Set("Content-Type", "application/json")
    w.Write(json)
}

// Refresh Token DB 관리 함수들
func saveRefreshTokenHandler(userEmail, refreshToken, login_type string) error {
    saveRefreshToken(userEmail, refreshToken, login_type)
    return nil
}


// 로그아웃 - Refresh Token 삭제
func Logout_Handler(w http.ResponseWriter, r *http.Request) {
    // var logoutReq struct {
    //     UserEmail    string `json:"user_email"`
    //     RefreshToken string `json:"refresh_token"`
    // }

    var receiveresult structs.LogoutDataForm

    err := json.NewDecoder(r.Body).Decode(&receiveresult)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    // DB에서 Refresh Token 삭제
    if err := deleteRefreshToken(receiveresult.Email, receiveresult.Refresh_Token); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    
    json, _ := json.Marshal(&structs.LogoutResult{Success: true})
    w.WriteHeader(http.StatusOK)
    w.Header().Set("Content-Type", "application/json")
    w.Write(json)
}