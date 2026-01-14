package jwt

import (
	"log"
	"smart_house/gosrc/db"
)

// Refresh Token DB 관리 함수들
func saveRefreshToken(userEmail, refreshToken, login_type string) error {
    // DB에 저장 (예시)
    // INSERT INTO refresh_tokens (user_email, token, created_at, expires_at) 
    // VALUES (?, ?, NOW(), DATE_ADD(NOW(), INTERVAL 30 DAY))

	log.Printf("saveRefreshToken refreshToken : %s", refreshToken)

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf("saveRefreshToken db err : %s", err)
		defer conn.Close()
		return err
	}
							
	result, err := conn.Exec("INSERT INTO refresh_tokens_user (email,login_type, token, created_at, expires_at) VALUES (?,?,?,NOW(),DATE_ADD(NOW(), INTERVAL 30 DAY)) ON DUPLICATE KEY UPDATE token = VALUES(token), created_at = VALUES(created_at), expires_at = VALUES(expires_at)", userEmail, login_type , refreshToken); 
	if err != nil {
		log.Printf("select user, %v", err)
		defer conn.Close()
		return err
	}

	n, err := result.RowsAffected()
	if err != nil {
		log.Printf("select user, %v", err)
		defer conn.Close()
		return err
	}

	log.Printf("saveRefreshToken RowsAffected : %d", n)


	defer conn.Close()
    return nil
}

// 관리자 Refresh Token DB 관리 함수들
func saveAdminRefreshToken(userEmail, refreshToken, login_type string) error {
    // DB에 저장 (예시)
    // INSERT INTO refresh_tokens (user_email, token, created_at, expires_at) 
    // VALUES (?, ?, NOW(), DATE_ADD(NOW(), INTERVAL 30 DAY))

	log.Printf("saveAdminRefreshToken refreshToken : %s", refreshToken)

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf("saveAdminRefreshToken db err : %s", err)
		defer conn.Close()
		return err
	}
							
	result, err := conn.Exec("INSERT INTO refresh_tokens_admin (email,login_type, token, created_at, expires_at) VALUES (?,?,?,NOW(),DATE_ADD(NOW(), INTERVAL 30 DAY)) ON DUPLICATE KEY UPDATE token = VALUES(token), created_at = VALUES(created_at), expires_at = VALUES(expires_at)", userEmail, login_type , refreshToken); 
	if err != nil {
		log.Printf("select user, %v", err)
		defer conn.Close()
		return err
	}

	n, err := result.RowsAffected()
	if err != nil {
		log.Printf("select user, %v", err)
		defer conn.Close()
		return err
	}

	log.Printf("saveAdminRefreshToken RowsAffected : %d", n)


	defer conn.Close()
    return nil
}

func validateRefreshToken(userEmail, refreshToken string) bool {
    // DB에서 확인 (예시)
    // SELECT COUNT(*) FROM refresh_tokens 
    // WHERE user_email = ? AND token = ? AND expires_at > NOW()

	var Count int

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf("validateRefreshToken db err : %s", err)
		defer conn.Close()
		return false
	}

	if err := conn.QueryRow("SELECT COUNT(*) FROM refresh_tokens_user WHERE email = ? AND token = ? AND expires_at > NOW()", userEmail , refreshToken).Scan(&Count); err != nil {
		log.Printf("select user, %v", err)
		defer conn.Close()
		return false
	}

	if Count == 0 {
		defer conn.Close()
		return false
	}

	defer conn.Close()
    return true
}

func validateAdminRefreshToken(userEmail, refreshToken string) bool {
    // DB에서 확인 (예시)
    // SELECT COUNT(*) FROM refresh_tokens 
    // WHERE user_email = ? AND token = ? AND expires_at > NOW()

	var Count int

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf("validateAdminRefreshToken db err : %s", err)
		defer conn.Close()
		return false
	}

	if err := conn.QueryRow("SELECT COUNT(*) FROM refresh_tokens_admin WHERE email = ? AND token = ? AND expires_at > NOW()", userEmail , refreshToken).Scan(&Count); err != nil {
		log.Printf("select user, %v", err)
		defer conn.Close()
		return false
	}

	if Count == 0 {
		defer conn.Close()
		return false
	}

	defer conn.Close()
    return true
}

// func updateRefreshToken(userEmail, oldToken, newToken, login_type string) error {
//     // 기존 토큰 삭제하고 새 토큰 저장 (예시)
//     // DELETE FROM refresh_tokens WHERE user_email = ? AND token = ?
//     // INSERT INTO refresh_tokens (user_email, token, created_at, expires_at) 
//     // VALUES (?, ?, NOW(), DATE_ADD(NOW(), INTERVAL 30 DAY))

// 	conn, err := db.OpenDB()
// 	if err != nil {
// 		log.Printf("updateRefreshToken db err : %s", err)
// 		defer conn.Close()
// 		return err
// 	}

// 	result, err := conn.Exec("INSERT INTO refresh_tokens_user (email,login_type, token, created_at, expires_at) VALUES (?,?,?,NOW(),DATE_ADD(NOW(), INTERVAL 30 DAY)) ON DUPLICATE KEY UPDATE (email = ?, login_type = ? ,token = ?, created_at = NOW(), expires_at = DATE_ADD(NOW(), INTERVAL 30 DAY))", userEmail, login_type , newToken, userEmail, login_type,newToken);
// 	if err != nil {
// 		log.Printf("select user, %v", err)
// 		defer conn.Close()
// 		return err
// 	}

// 	n, err := result.RowsAffected()
// 	if err != nil {
// 		log.Printf("select user, %v", err)
// 		defer conn.Close()
// 		return err
// 	}

// 	log.Printf("updateRefreshToken RowsAffected : %d", n)

// 	defer conn.Close()
//     return nil
// }

// Refresh Token 삭제
func deleteRefreshToken(userEmail,token string) error {

	log.Printf("deleteRefreshToken refreshToken : %s", token)

    conn, err := db.OpenDB()
	if err != nil {
		log.Printf("deleteRefreshToken db err : %s", err)
		defer conn.Close()
		return err
	}

	result, err := conn.Exec("DELETE FROM refresh_tokens_user WHERE email = ? AND token = ?", userEmail,token); 
	if err != nil {
		log.Printf("select user, %v", err)
		defer conn.Close()
		return err
	}

	n, err := result.RowsAffected()
	if err != nil {
		log.Printf("select user, %v", err)
		defer conn.Close()
		return err
	}

	log.Printf("deleteRefreshToken RowsAffected : %d", n)
		

	defer conn.Close()
    return nil
}

// 관리자 Refresh Token 삭제
func deleteAdminRefreshToken(userEmail,token string) error {

	log.Printf("deleteAdminRefreshToken refreshToken : %s", token)

    conn, err := db.OpenDB()
	if err != nil {
		log.Printf("deleteAdminRefreshToken db err : %s", err)
		defer conn.Close()
		return err
	}

	result, err := conn.Exec("DELETE FROM refresh_tokens_admin WHERE email = ? AND token = ?", userEmail,token); 
	if err != nil {
		log.Printf("select user, %v", err)
		defer conn.Close()
		return err
	}

	n, err := result.RowsAffected()
	if err != nil {
		log.Printf("select user, %v", err)
		defer conn.Close()
		return err
	}

	log.Printf("deleteAdminRefreshToken RowsAffected : %d", n)
		

	defer conn.Close()
    return nil
}

