package adminpage

import (
	"fmt"
	"log"
	"strings"
	"time"

	"smart_house/gosrc/db"
	"smart_house/gosrc/jwt"
	"smart_house/gosrc/structs"

	"golang.org/x/crypto/bcrypt"
)

func returnTime() string {
	loc, _ := time.LoadLocation("Asia/Seoul")
	return time.Now().In(loc).Format("2006-01-02 15:04:05")
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func isSafeColumnName(c map[string]bool, column string) bool {
	// 안전한 컬럼 이름 목록

	_, ok := c[column]
	return ok
}

// 로그인
func Login_select(email string, pwd string) structs.LoginDataForm {
	var result structs.LoginDataForm

	var Email string
	var Name string
	var Password string

	result = structs.LoginDataForm{Success: false, Email: email, TokenPair: structs.TokenPair{}}

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf(" Login_select db err : %s", err)
	}

	// // 암호화
	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pwd), 11) // 10 means difficulty
	// if err != nil {
	// 	log.Printf("failed to GenerateFromPassword: %v", err)
	// 	return result
	// }

	// log.Printf("password : %v", string(hashedPassword))

	if err := conn.QueryRow("SELECT email,name,password FROM adminuser WHERE email = ?", email).Scan(&Email, &Name, &Password); err != nil {
		log.Printf("select user, %v", err)
		return result
	}

	if err := bcrypt.CompareHashAndPassword([]byte(Password), []byte(pwd)); err != nil {
		log.Printf("incorrect password, %v", err)
		return result
	}

	newTokenPair, err := jwt.GenerateTokenPair(email,"admin")
    if err != nil {
        return result
    }

	result = structs.LoginDataForm{Email: Email, Success: true, TokenPair: *newTokenPair}

	defer conn.Close()

	return result
}

// 하우스 설치 정보
func Device_list_select(dconn []string, s_order, e_order string) []structs.DeviceDetail {
	var result []structs.DeviceDetail

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf(" Device_list_select err : %s", err)
	}

	// const sqlString = "SELECT * FROM device ORDER BY id DESC limit ?,?"
	const sqlString = "SELECT A.id, A.comp,A.pump,A.temp, A.humi,A.f_fan,A.e_fan,A.aircon,A.co2_sensor,A.co2_creator,A.ph_ec_sensor,A.led, A.cctv, A.auto, A.b_time FROM device AS A ORDER BY A.id DESC limit ?,?"

	rows, err := conn.Query(sqlString, s_order, e_order)
	if err != nil {
		log.Printf(" Device_list_select err : %s", err)
	}

	for rows.Next() {
		var Id string
		var Comp bool
		var Pump bool
		var Temp float32
		var Humi float32
		var F_Fan bool
		var E_Fan bool
		var Aircon bool
		var Co2_Sensor float32
		var C02_Creator bool
		var Ph_Ec_Sensor float32
		var Led bool
		var Cctv string
		var Auto int
		var B_Time string

		rows.Scan(&Id, &Comp, &Pump, &Temp, &Humi, &F_Fan, &E_Fan, &Aircon, &Co2_Sensor, &C02_Creator, &Ph_Ec_Sensor, &Led, &Cctv, &Auto, &B_Time)

		if !contains(dconn, Id) {
			Auto = 2
		}

		m := structs.DeviceDetail{Id: Id, Comp: Comp, Pump: Pump, Temp: Temp, Humi: Humi, F_Fan: F_Fan, E_Fan: E_Fan, Aircon: Aircon, Co2_Sensor: Co2_Sensor, C02_Creator: C02_Creator, Ph_Ec_Sensor: Ph_Ec_Sensor, Led: Led, Cctv: Cctv, Auto: Auto, B_Time: B_Time}
		result = append(result, m)
	}
	defer rows.Close()
	defer conn.Close()

	return result
}

// 하우스 등록 유저 기기 정보
func Device_user_list_select(s_order, e_order string) []DeviceUserListItem {
	var result []DeviceUserListItem

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf(" Device_user_list_select err : %s", err)
	}

	const sqlString = "SELECT id, user, r_time FROM registration ORDER BY r_time DESC limit ?,?"

	rows, err := conn.Query(sqlString, s_order, e_order)
	if err != nil {
		log.Printf(" Device_user_list_select err : %s", err)
	}

	for rows.Next() {
		var Id string
		var User string
		var R_Time string

		rows.Scan(&Id, &User, &R_Time)

		m := DeviceUserListItem{Id: Id, User: User, R_Time: R_Time}
		result = append(result, m)
	}
	defer rows.Close()
	defer conn.Close()

	return result
}

// 유저 정보
func User_list_select(s_order, e_order string) []UserListItem {
	var result []UserListItem

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf(" User_list_select err : %s", err)
	}

	const sqlString = "SELECT email, name FROM user ORDER BY email DESC limit ?,?"

	rows, err := conn.Query(sqlString, s_order, e_order)
	if err != nil {
		log.Printf(" User_list_select err : %s", err)
	}

	for rows.Next() {
		var Email string
		var Name string

		rows.Scan(&Email, &Name)

		m := UserListItem{Email: Email, Name: Name}
		result = append(result, m)
	}
	defer rows.Close()
	defer conn.Close()

	return result
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// 하우스 기본정보 카운트 반환
func Device_list_count_select() int {
	var result int

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf("Device_list_count_select err : %s", err)
	}

	err = conn.QueryRow("SELECT count(id) FROM device").Scan(&result)
	if err != nil {
		log.Printf("Device_list_count_select : %s", err)
	}

	defer conn.Close()

	return result
}

// 하우스 등록 유저 기기정보 카운트 반환
func Device_user_list_count_select() int {
	var result int

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf("Device_user_list_count_select err : %s", err)
	}

	err = conn.QueryRow("SELECT count(user) FROM registration").Scan(&result)
	if err != nil {
		log.Printf("Device_user_list_count_select : %s", err)
	}

	defer conn.Close()

	return result
}

// 유저 정보 카운트 반환
func User_list_count_select() int {
	var result int

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf("User_list_count_select err : %s", err)
	}

	err = conn.QueryRow("SELECT count(email) FROM user").Scan(&result)
	if err != nil {
		log.Printf("User_list_count_select : %s", err)
	}

	defer conn.Close()

	return result
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// 하우스 기본 정보 검색
func Device_list_select_search(dconn []string, column, data string) []structs.DeviceDetail {
	var result []structs.DeviceDetail

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf(" Device_list_select_search err : %s", err)
	}

	//컬럼이름이 존재하지 않는지 체크	-- sql injection 공격 차단
	c := map[string]bool{
		"id":           true,
		"comp":         true,
		"pump":         true,
		"f_fan":        true,
		"e_fan":        true,
		"temp":         true,
		"humi":         true,
		"aircon":       true,
		"co2_sensor":   true,
		"co2_creator":  true,
		"ph_ec_sensor": true,
		"led":          true,
		"auto":         true,
		"cctv":         true,
	}

	if !isSafeColumnName(c, column) {
		return result
	}

	sqlString := fmt.Sprintf("SELECT A.id, A.comp,A.pump,A.temp, A.humi,A.f_fan,A.e_fan,A.aircon,A.co2_sensor,A.co2_creator,A.ph_ec_sensor,A.led, A.cctv, A.auto , A.b_time FROM device AS A WHERE %s like ? ", column)

	rows, err := conn.Query(sqlString, "%"+strings.Trim(data, " ")+"%")
	if err != nil {
		log.Printf(" Device_list_select_search err : %s", err)
	}

	for rows.Next() {
		var Id string
		var Comp bool
		var Pump bool
		var Temp float32
		var Humi float32
		var F_Fan bool
		var E_Fan bool
		var Aircon bool
		var Co2_Sensor float32
		var C02_Creator bool
		var Ph_Ec_Sensor float32
		var Led bool
		var Cctv string
		var Auto int
		var B_Time string

		rows.Scan(&Id, &Comp, &Pump, &Temp, &Humi, &F_Fan, &E_Fan, &Aircon, &Co2_Sensor, &C02_Creator, &Ph_Ec_Sensor, &Led, &Cctv, &Auto, &B_Time)

		m := structs.DeviceDetail{Id: Id, Comp: Comp, Pump: Pump, Temp: Temp, Humi: Humi, F_Fan: F_Fan, E_Fan: E_Fan, Aircon: Aircon, Co2_Sensor: Co2_Sensor, C02_Creator: C02_Creator, Ph_Ec_Sensor: Ph_Ec_Sensor, Led: Led, Cctv: Cctv, Auto: Auto, B_Time: B_Time}
		result = append(result, m)
	}

	defer rows.Close()
	defer conn.Close()

	return result
}

// 하우스 유저 등록 하우스 간단 리스트
func DeviceSimple_list(dconn []string) []structs.DeviceSimple {
	var result []structs.DeviceSimple

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf("DeviceSimple_list err : %s", err)
	}

	rows, err := conn.Query("SELECT A.id, A.temp, A.humi, A.cctv, A.auto FROM device AS A")
	if err != nil {
		log.Printf("DeviceSimple_list err : %s", err)
	}

	for rows.Next() {
		var Id string
		var Temp float32
		var Humi float32
		var Cctv string
		var Auto int

		rows.Scan(&Id, &Temp, &Humi, &Cctv, &Auto)

		if !contains(dconn, Id) {
			Auto = 2
		}

		m := structs.DeviceSimple{Id: Id, Temp: Temp, Humi: Humi, Cctv: Cctv, Auto: Auto}
		result = append(result, m)
	}
	defer rows.Close()
	defer conn.Close()

	return result
}

// 하우스 등록 유저 기기 정보 검색
func Device_user_list_select_search(column, data string) []DeviceUserListItem {
	var result []DeviceUserListItem

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf("Device_user_list_select_search err : %s", err)
	}

	//컬럼이름이 존재하지 않는지 체크	-- sql injection 공격 차단
	c := map[string]bool{
		"id":     true,
		"user":   true,
		"r_time": true,
	}

	if !isSafeColumnName(c, column) {
		return result
	}

	sqlString := fmt.Sprintf("SELECT id, user,r_time FROM registration AS A WHERE %s like ? ", column)

	rows, err := conn.Query(sqlString, "%"+strings.Trim(data, " ")+"%")
	if err != nil {
		log.Printf("Device_user_list_select_search err : %s", err)
	}

	for rows.Next() {
		var Id string
		var User string
		var R_Time string

		rows.Scan(&Id, &User, &R_Time)

		m := DeviceUserListItem{Id: Id, User: User, R_Time: R_Time}
		result = append(result, m)
	}
	defer rows.Close()
	defer conn.Close()

	return result
}

// 유저 정보 검색
func User_list_select_search(column, data string) []UserListItem {
	var result []UserListItem

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf("User_list_select_search err : %s", err)
	}

	//컬럼이름이 존재하지 않는지 체크	-- sql injection 공격 차단
	c := map[string]bool{
		"email": true,
		"name":  true,
	}

	if !isSafeColumnName(c, column) {
		return result
	}

	sqlString := fmt.Sprintf("SELECT email,name FROM user WHERE %s like ? ", column)

	rows, err := conn.Query(sqlString, "%"+strings.Trim(data, " ")+"%")
	if err != nil {
		log.Printf("User_list_select_search err : %s", err)
	}

	for rows.Next() {
		var Email string
		var Name string

		rows.Scan(&Email, &Name)

		m := UserListItem{Email: Email, Name: Name}
		result = append(result, m)
	}
	defer rows.Close()
	defer conn.Close()

	return result
}

// 하우스 업데이트 관리 리스트 카운트 반환
func Upd_list_count_select() int {
	var result int

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf("Upd_list_count_select err : %s", err)
	}

	err = conn.QueryRow("SELECT count(id) FROM device").Scan(&result)
	if err != nil {
		log.Printf("Upd_list_count_select : %s", err)
	}

	defer conn.Close()

	return result
}

// 하우스 업데이트 관리 리스트 정보
func Upd_list_select(dconn []string, s_order, e_order string) []DeviceUpdListItem {
	var result []DeviceUpdListItem

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf("Upd_list_select err : %s", err)
	}

	rows, err := conn.Query("SELECT id, auto, b_time FROM device ORDER BY id DESC limit ?,?", s_order, e_order)
	if err != nil {
		log.Printf("Upd_list_select err : %s", err)
	}

	for rows.Next() {
		var Id string
		var Auto int
		var B_Time string

		rows.Scan(&Id, &Auto, &B_Time)

		if !contains(dconn, Id) {
			Auto = 2
		}

		m := DeviceUpdListItem{Id: Id, Auto: Auto, B_Time: B_Time}
		result = append(result, m)
	}
	defer rows.Close()
	defer conn.Close()

	return result
}

// 하우스 업데이트 관리 리스트 정보 검색
func Upd_list_select_search(dconn []string, column, data string) []DeviceUpdListItem {
	var result []DeviceUpdListItem

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf("Upd_list_select_search err : %s", err)
	}

	//컬럼이름이 존재하지 않는지 체크	-- sql injection 공격 차단
	c := map[string]bool{
		"id": true,
	}

	if !isSafeColumnName(c, column) {
		return result
	}

	queryStirng := fmt.Sprintf("SELECT id, auto, b_time FROM device where %s like ?  ORDER BY id DESC ", column)

	rows, err := conn.Query(queryStirng, "%"+strings.Trim(data, " ")+"%")
	if err != nil {
		log.Printf("Upd_list_select_search err : %s", err)
	}

	for rows.Next() {
		var Id string
		var Auto int
		var B_Time string

		rows.Scan(&Id, &Auto, &B_Time)

		if !contains(dconn, Id) {
			Auto = 2
		}

		m := DeviceUpdListItem{Id: Id, Auto: Auto, B_Time: B_Time}
		result = append(result, m)
	}
	defer rows.Close()
	defer conn.Close()

	return result
}

// 하우스 로그 리스트 카운트 반환
func Log_list_count_select() int {
	var result int

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf("Log_list_count_select err : %s", err)
	}

	err = conn.QueryRow("SELECT count(id) FROM device").Scan(&result)
	if err != nil {
		log.Printf("Log_list_count_select : %s", err)
	}

	defer conn.Close()

	return result
}

// 하우스 로그 리스트 정보
func Log_list_select(dconn []string, s_order, e_order string) []DeviceLogListItem {
	var result []DeviceLogListItem

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf("Log_list_select err : %s", err)
	}

	rows, err := conn.Query("SELECT id, auto FROM device ORDER BY id DESC limit ?,?", s_order, e_order)
	if err != nil {
		log.Printf("Log_list_select err : %s", err)
	}

	for rows.Next() {
		var Id string
		var Auto int

		rows.Scan(&Id, &Auto)

		if !contains(dconn, Id) {
			Auto = 2
		}

		m := DeviceLogListItem{Id: Id, Auto: Auto}
		result = append(result, m)
	}
	defer rows.Close()
	defer conn.Close()

	return result
}

// 하우스 로그 리스트 정보 검색
func Log_list_select_search(dconn []string, column, data string) []DeviceLogListItem {
	var result []DeviceLogListItem

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf("Log_list_select_search err : %s", err)
	}

	//컬럼이름이 존재하지 않는지 체크	-- sql injection 공격 차단
	c := map[string]bool{
		"id": true,
	}

	if !isSafeColumnName(c, column) {
		return result
	}

	queryStirng := fmt.Sprintf("SELECT id, auto  FROM device where %s like ?  ORDER BY id DESC ", column)

	rows, err := conn.Query(queryStirng, "%"+strings.Trim(data, " ")+"%")
	if err != nil {
		log.Printf("Log_list_select_search err : %s", err)
	}

	for rows.Next() {
		var Id string
		var Auto int

		rows.Scan(&Id, &Auto)

		if !contains(dconn, Id) {
			Auto = 2
		}

		m := DeviceLogListItem{Id: Id, Auto: Auto}
		result = append(result, m)
	}
	defer rows.Close()
	defer conn.Close()

	return result
}

// 일자별 센서값 카운트 반환
func Device_daily_list_count_select() int {
	var result int

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf("Device_daily_list_count_select err : %s", err)
	}

	err = conn.QueryRow("SELECT count(A.id) FROM valhistory AS A").Scan(&result)
	if err != nil {
		log.Printf("Device_daily_list_count_select : %s", err)
	}

	defer conn.Close()

	return result
}

// 일자별 센서값 조회
func Device_daily_list_select(s_order, e_order string) []structs.DeviceDailyListItem {
	var result []structs.DeviceDailyListItem

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf("Device_daily_list_select err : %s", err)
	}

	const sqlString = "SELECT A.id,A.temp,A.humi,A.co2_sensor,A.ph_ec_sensor,A.time FROM valhistory AS A ORDER BY time DESC limit ?,?"

	rows, err := conn.Query(sqlString, s_order, e_order)
	if err != nil {
		log.Printf("Device_daily_list_select err : %s", err)
	}

	for rows.Next() {
		var Id string
		var Temp float32
		var Humi float32
		var Co2_Sensor float32
		var Ph_Ec_Sensor float32
		var Time string

		rows.Scan(&Id, &Temp, &Humi, &Co2_Sensor, &Ph_Ec_Sensor, &Time)

		m := structs.DeviceDailyListItem{Id: Id, Temp: Temp, Humi: Humi, Co2_Sensor: Co2_Sensor, Ph_Ec_Sensor: Ph_Ec_Sensor, Time: Time}
		result = append(result, m)

	}
	defer rows.Close()
	defer conn.Close()

	// log.Printf(" charge list : %s", result)

	return result
}

// 일자별 센서값 검색
func Device_daily_list_select_search(column, data string) []structs.DeviceDailyListItem {
	var result []structs.DeviceDailyListItem

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf("Device_daily_list_select_search err : %s", err)
	}

	//컬럼이름이 존재하지 않는지 체크	-- sql injection 공격 차단
	c := map[string]bool{
		"id":           true,
		"temp":         true,
		"humi":         true,
		"co2_sensor":   true,
		"ph_ec_sensor": true,
		"time":         true,
	}

	if !isSafeColumnName(c, column) {
		return result
	}

	queryStirng := fmt.Sprintf("SELECT A.id,A.temp,A.humi,A.co2_sensor,A.ph_ec_sensor,A.time FROM valhistory AS A WHERE %s LIKE ? ORDER BY time DESC ", column)

	rows, err := conn.Query(queryStirng, "%"+strings.Trim(data, " ")+"%")
	if err != nil {
		log.Printf("Device_daily_list_select_search err : %s", err)
	}

	for rows.Next() {
		var Id string
		var Temp float32
		var Humi float32
		var Co2_Sensor float32
		var Ph_Ec_Sensor float32
		var Time string

		rows.Scan(&Id, &Temp, &Humi, &Co2_Sensor, &Ph_Ec_Sensor, &Time)

		m := structs.DeviceDailyListItem{Id: Id, Temp: Temp, Humi: Humi, Co2_Sensor: Co2_Sensor, Ph_Ec_Sensor: Ph_Ec_Sensor, Time: Time}
		result = append(result, m)

	}
	defer rows.Close()
	defer conn.Close()

	// log.Printf(" charge list : %s", result)

	return result
}

// 센서값 그래프 기기 리스트 조회
func Device_daily_graph_list_select() []structs.DeviceDailyGraphListItem {
	var result []structs.DeviceDailyGraphListItem

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf("Device_daily_graph_list_select err : %s", err)
	}

	rows, err := conn.Query("SELECT A.id FROM device AS A")
	if err != nil {
		log.Printf("DeviceSimple_list err : %s", err)
	}

	for rows.Next() {
		var Id string

		rows.Scan(&Id)

		m := structs.DeviceDailyGraphListItem{Id: Id}

		result = append(result, m)

	}
	defer rows.Close()
	defer conn.Close()

	// log.Printf(" charge list : %s", result)

	return result
}

// 센서값 그래프 조회
func Device_daily_graph_select(starttime, endtime, device_id string) []structs.DeviceDailyGraphItem {
	var result []structs.DeviceDailyGraphItem

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf("Device_daily_graph_select err : %s", err)
	}
	log.Printf("Device_daily_graph_select time : %s , %s", starttime, endtime)
	//검색하는 기기 이름이 전체 or 개별
	if device_id == "전체" {

		const sqlString = "SELECT HOUR(time) AS times , AVG(temp) AS temp, AVG(humi) AS humi, AVG(co2_sensor) AS co2_sensor,AVG(ph_ec_sensor) AS ph_ec_sensor FROM valhistory WHERE time BETWEEN ? AND ? GROUP BY times ORDER BY times"

		rows, err := conn.Query(sqlString, starttime, endtime)
		if err != nil {
			log.Printf("Device_daily_graph_select err : %s", err)
		}

		for rows.Next() {
			var Temp float32
			var Humi float32
			var Co2_Sensor float32
			var Ph_Ec_Sensor float32
			var Times string

			rows.Scan(&Times, &Temp, &Humi, &Co2_Sensor, &Ph_Ec_Sensor)

			m := structs.DeviceDailyGraphItem{Time: Times, Temp: Temp, Humi: Humi, Co2_Sensor: Co2_Sensor, Ph_Ec_Sensor: Ph_Ec_Sensor}

			result = append(result, m)

		}
		defer rows.Close()
		defer conn.Close()

		// log.Printf(" charge list : %s", result)

		return result
	} else {
		rows, err := conn.Query("SELECT HOUR(time) AS times , AVG(temp) AS temp, AVG(humi) AS humi, AVG(co2_sensor) AS co2_sensor,AVG(ph_ec_sensor) AS ph_ec_sensor FROM valhistory WHERE time BETWEEN ? AND ? AND id = ? GROUP BY times ORDER BY times", starttime, endtime, device_id)
		if err != nil {
			log.Printf("Device_daily_graph_select err : %s", err)
		}

		for rows.Next() {
			var Temp float32
			var Humi float32
			var Co2_Sensor float32
			var Ph_Ec_Sensor float32
			var Times string

			rows.Scan(&Times, &Temp, &Humi, &Co2_Sensor, &Ph_Ec_Sensor)

			m := structs.DeviceDailyGraphItem{Time: Times, Temp: Temp, Humi: Humi, Co2_Sensor: Co2_Sensor, Ph_Ec_Sensor: Ph_Ec_Sensor}

			result = append(result, m)

		}
		defer rows.Close()
		defer conn.Close()

		// log.Printf(" charge list : %s", result)

		return result
	}

}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// 하우스 등록 유저 기기 삽입 명령
func Insert_deviceuser(Id, user string) bool {

	conn, err := db.OpenDB()
	if err != nil {
		fmt.Println(err)
		return false
	}

	rows, err := conn.Exec("INSERT INTO registration(id,user,r_time) SELECT ?,?,? FROM DUAL WHERE NOT EXISTS (SELECT * FROM registration WHERE user = ? AND id = ?)", Id, user, returnTime(), user, Id)

	if err != nil {
		fmt.Printf("rows1 : %s", err)
		return false
	}

	log.Printf("insert_deviceuser : %v  ", rows)

	defer conn.Close()

	return true
}

// 기기 삭제
func Delete_deviceuser(Id, user string) bool {

	conn, err := db.OpenDB()
	if err != nil {
		fmt.Println(err)
		return false
	}

	rows, err := conn.Exec("DELETE FROM registration WHERE id =? AND user =?", Id, user)
	if err != nil {
		fmt.Printf("rows1 : %s", err)
		return false
	}

	log.Printf("delete_deviceuser : %v ", rows)

	defer conn.Close()

	return true
}

// 유저 삽입 명령
func Insert_user(Email, Name, pwd string) bool {

	var Password string

	conn, err := db.OpenDB()
	if err != nil {
		fmt.Println(err)
		return false
	}

	if err := bcrypt.CompareHashAndPassword([]byte(Password), []byte(pwd)); err != nil {
		log.Printf("incorrect password, %v", err)
		return false
	}

	rows, err := conn.Exec("INSERT INTO user(email,name,password) SELECT ?,?,? FROM DUAL WHERE NOT EXISTS (SELECT * FROM user WHERE email = ?)", Email, Name, Password, Email)

	if err != nil {
		fmt.Printf("rows1 : %s", err)
		return false
	}

	log.Printf("insert_user : %v  ", rows)

	defer conn.Close()

	return true
}

// 유저 삭제
func Delete_user(Email, Name string) bool {

	conn, err := db.OpenDB()
	if err != nil {
		fmt.Println(err)
		return false
	}

	rows, err := conn.Exec("DELETE FROM user WHERE email =? AND name =?", Email, Name)
	if err != nil {
		fmt.Printf("rows1 : %s", err)
		return false
	}

	log.Printf("delete_user : %v ", rows)

	defer conn.Close()

	return true
}

// 유저 수정 명령
func Update_user(Email, Name, pwd string) bool {

	var Password string

	conn, err := db.OpenDB()
	if err != nil {
		fmt.Println(err)
		return false
	}

	if err := bcrypt.CompareHashAndPassword([]byte(Password), []byte(pwd)); err != nil {
		log.Printf("incorrect password, %v", err)
		return false
	}

	rows, err := conn.Exec("UPDATE user SET name = ? , password = ? WHERE email = ?)", Name, Password, Email)

	if err != nil {
		fmt.Printf("rows1 : %s", err)
		return false
	}

	log.Printf("update_user : %v  ", rows)

	defer conn.Close()

	return true
}

// 유저 로그아웃 명령
func Logout_user(Email, Type string) bool {

	conn, err := db.OpenDB()
	if err != nil {
		fmt.Println(err)
		return false
	}

	if Type == "all" {
		rows, err := conn.Exec("DELETE FROM refresh_tokens_user WHERE email =?", Email)
		if err != nil {
			fmt.Printf("rows1 : %s", err)
			return false
		}

		log.Printf("logout_user : %v ", rows)
	}else{
		rows, err := conn.Exec("DELETE FROM refresh_tokens_user WHERE email =? AND login_type =?", Email, Type)
		if err != nil {
			fmt.Printf("rows1 : %s", err)
			return false
		}

		log.Printf("logout_user : %v ", rows)
	}

	defer conn.Close()

	return true
}

// 유저 메시지 id 반환
func User_AlertId_select(user string) string {
	var result string

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf("User_AlertId_select err : %s", err)
	}

	err = conn.QueryRow("SELECT os_id FROM user WHERE email = ?", user).Scan(&result)
	if err != nil {
		log.Printf("User_AlertId_select err : %s", err)
	}

	defer conn.Close()

	return result
}
