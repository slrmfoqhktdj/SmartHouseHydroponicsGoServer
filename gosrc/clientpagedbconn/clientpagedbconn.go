package clientpagedbconn

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"smart_house/gosrc/db"
	"smart_house/gosrc/jwt"
	"smart_house/gosrc/structs"

	"golang.org/x/crypto/bcrypt"
)

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

// 로그인
func Login_select(email string, pwd string) structs.LoginDataForm {
	var result structs.LoginDataForm
	// dbname := os.Getenv("DBNAME")
	// dbdata := os.Getenv("DBDATA")

	// conn, err := sql.Open(dbname, dbdata)

	var Email string
	var Name string
	var Password string

	result = structs.LoginDataForm{Success: false, Email: email, TokenPair: structs.TokenPair{}}

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf(" Login_select db err : %s", err)
		return result
	}

	// // 암호화
	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pwd), 11) // 10 means difficulty
	// if err != nil {
	// 	log.Printf("failed to GenerateFromPassword: %v", err)
	// 	return result
	// }

	// log.Printf("password : %v", string(hashedPassword))

	if err := conn.QueryRow("SELECT email,name,password FROM user WHERE email = ?", email).Scan(&Email, &Name, &Password); err != nil {
		log.Printf("select user, %v", err)
		return result
	}

	if err := bcrypt.CompareHashAndPassword([]byte(Password), []byte(pwd)); err != nil {
		log.Printf("incorrect password, %v", err)
		return result
	}

	newTokenPair, err := jwt.GenerateTokenPair(email,"web")
    if err != nil {
		log.Printf("generateTokenPair err, %v", err)
        return result
    }

	result = structs.LoginDataForm{Email: Email, Success: true, TokenPair: *newTokenPair}

	defer conn.Close()

	return result
}

// 하우스 앱 로그인
func LoginApp_select(email string, pwd string) structs.LoginAppDataForm {
	var result structs.LoginAppDataForm
	// dbname := os.Getenv("DBNAME")
	// dbdata := os.Getenv("DBDATA")

	// conn, err := sql.Open(dbname, dbdata)

	var Email string
	var Name string
	var Password string

	result = structs.LoginAppDataForm{Success: false, Email: email, TokenPair: structs.TokenPair{}}

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf(" LoginApp_select db err : %s", err)
		return result
	}

	// // 암호화
	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pwd), 11) // 10 means difficulty
	// if err != nil {
	// 	log.Printf("failed to GenerateFromPassword: %v", err)
	// 	return result
	// }

	// log.Printf("password : %v", string(hashedPassword))

	if err := conn.QueryRow("SELECT email,name,password FROM user WHERE email = ?", email).Scan(&Email, &Name, &Password); err != nil {
		log.Printf("select user, %v", err)
		return result
	}

	if err := bcrypt.CompareHashAndPassword([]byte(Password), []byte(pwd)); err != nil {
		log.Printf("incorrect password, %v", err)
		return result
	}

	newTokenPair, err := jwt.GenerateTokenPair(email,"app")
    if err != nil {
		log.Printf("generateTokenPair err, %v", err)
        return result
    }

	result = structs.LoginAppDataForm{Email: Email, Success: true, TokenPair: *newTokenPair}

	defer conn.Close()

	return result
}


// 하우스 기기 로그인
func LoginDevice_select(id string) structs.LoginDeviceDataForm {
	var result structs.LoginDeviceDataForm
	// dbname := os.Getenv("DBNAME")
	// dbdata := os.Getenv("DBDATA")

	// conn, err := sql.Open(dbname, dbdata)

	var Id string

	result = structs.LoginDeviceDataForm{Success: false, DeviceId : id , TokenPair: structs.TokenPair{}}

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf(" LoginDevice_select db err : %s", err)
		return result
	}

	// // 암호화
	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pwd), 11) // 10 means difficulty
	// if err != nil {
	// 	log.Printf("failed to GenerateFromPassword: %v", err)
	// 	return result
	// }

	// log.Printf("password : %v", string(hashedPassword))

	if err := conn.QueryRow("SELECT id FROM device WHERE id = ?", id).Scan(&Id); err != nil {
		log.Printf("select user, %v", err)
		return result
	}

	if Id != id {
		return result
	}

	newTokenPair, err := jwt.GenerateTokenPair(id,"device")
    if err != nil {
		log.Printf("generateTokenPair err, %v", err)
        return result
    }

	result = structs.LoginDeviceDataForm{ DeviceId : id , Success: true, TokenPair: *newTokenPair}

	defer conn.Close()

	return result
}

func AppUser_OneSignalID_update(email, os_id string) bool {

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf("AppUser_OneSignalID_update db err : %s", err)
		return false
	}

	rows, err := conn.Exec("UPDATE user SET os_id = ? WHERE email = ?", os_id, email)

	if err != nil {
		log.Printf("rows err : %s", err)
		return false
	}

	r, err_i := rows.RowsAffected()

	if err_i != nil {
		log.Printf("rows_i err : %s", err)
		return false
	}

	log.Printf("AppUser_OneSignalID_update rows : %d", r)

	defer conn.Close()

	return true
}

func UserAlarmList_select(email string) []structs.UserAlarmList {

	var resultString string
	var result []structs.UserAlarmList

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf("UserAlarmList_select err : %s", err)
		return result
	}

	if err := conn.QueryRow("SELECT alarm FROM user WHERE email = ?", email).Scan(&resultString); err != nil {
		log.Printf("UserAlarmList_select err : %s", err)
		return result
	}

	defer conn.Close()

	if resultString == "" {
		return result
	}

	resultSplit := strings.Split(resultString, "&&")

	for _, resultSplited := range resultSplit {
		var resultItem structs.UserAlarmList

		err := json.Unmarshal([]byte(resultSplited), &resultItem)
		if err != nil {
			log.Printf("UserAlarmList_select err : %s", err)
		}

		result = append(result, resultItem)
	}

	return result
}

// 하우스 유저 등록 하우스 간단 리스트
func DeviceSimple_list(dconn []string, user string) []structs.DeviceSimple {
	var result []structs.DeviceSimple

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf("DeviceSimple_list err : %s", err)
	}

	rows, err := conn.Query("SELECT A.id, A.temp, A.humi, A.cctv, A.auto FROM device AS A INNER JOIN registration AS B ON A.id = B.id WHERE B.user = ?", user)
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

// 하우스 현재 상태 보고
func DeviceDetail_select(dconn []string, id string) structs.DeviceDetail {
	var result structs.DeviceDetail

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

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf(" Device_list_select err : %s", err)
	}

	const sqlString = "SELECT id, comp,pump,temp,humi,f_fan,e_fan,aircon,co2_sensor,co2_creator,ph_ec_sensor,led,cctv, auto,b_time FROM device WHERE id = ?"
	log.Printf(" DeviceInfo_select id : %s", id)
	err = conn.QueryRow(sqlString, id).Scan(&Id, &Comp, &Pump, &Temp, &Humi, &F_Fan, &E_Fan, &Aircon, &Co2_Sensor, &C02_Creator, &Ph_Ec_Sensor, &Led, &Cctv, &Auto, &B_Time)
	if err != nil {
		log.Printf(" Device_list_select err : %s", err)
	}

	if !contains(dconn, Id) {
		Auto = 2
	}

	result = structs.DeviceDetail{Id: Id, Comp: Comp, Pump: Pump, Temp: Temp, Humi: Humi, F_Fan: F_Fan, E_Fan: E_Fan, Aircon: Aircon, Co2_Sensor: Co2_Sensor, C02_Creator: C02_Creator, Ph_Ec_Sensor: Ph_Ec_Sensor, Led: Led, Cctv: Cctv, Auto: Auto, B_Time: B_Time}
	// log.Printf(" DeviceInfo_select : id - %s, comp - %v, pump - %v , temp - %f, humi - %f, f_fan - %v, e_fan - %v, aircon - %v, co2_sensor - %f, co2_creator  - %v, ph_ec_sensor - %f, led - %v, cctv - %s", Id, Comp, Pump, Temp, Humi, F_Fan, E_Fan, Aircon, Co2_Sensor, C02_Creator, Ph_Ec_Sensor, Led, Cctv)

	defer conn.Close()

	return result
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// 하우스 설치 정보
func Device_list_select(dconn []string, user, s_order, e_order string) []structs.DeviceDetail {
	var result []structs.DeviceDetail
	// dbname := os.Getenv("DBNAME")
	// dbdata := os.Getenv("DBDATA")

	// conn, err := sql.Open(dbname, dbdata)

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf(" Device_list_select err : %s", err)
	}

	// const sqlString = "SELECT * FROM device ORDER BY id DESC limit ?,?"
	const sqlString = "SELECT A.id, A.comp,A.pump,A.temp, A.humi,A.f_fan,A.e_fan,A.aircon,A.co2_sensor,A.co2_creator,A.ph_ec_sensor,A.led, A.cctv, A.auto ,A.b_time FROM device AS A INNER JOIN registration AS B ON A.id = B.id WHERE B.user = ? ORDER BY A.id DESC limit ?,?"

	rows, err := conn.Query(sqlString, user, s_order, e_order)
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

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// 하우스 기본정보 카운트 반환
func Device_list_count_select(user string) int {
	var result int

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf("Device_list_count_select err : %s", err)
	}

	err = conn.QueryRow("SELECT count(A.id) FROM device AS A INNER JOIN registration AS B ON A.id = B.id WHERE B.user = ?", user).Scan(&result)
	if err != nil {
		log.Printf("Device_list_count_select : %s", err)
	}

	defer conn.Close()

	return result
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// 하우스 기본 정보 검색
func Device_list_select_search(dconn []string, user, column, data string) []structs.DeviceDetail {
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

	sqlString := fmt.Sprintf("SELECT A.id, A.comp,A.pump,A.temp, A.humi,A.f_fan,A.e_fan,A.aircon,A.co2_sensor,A.co2_creator,A.ph_ec_sensor,A.led, A.cctv, A.auto , A.b_time FROM device AS A INNER JOIN registration AS B ON A.id = B.id WHERE B.user = ? AND %s like ? ", column)

	rows, err := conn.Query(sqlString, user, "%"+strings.Trim(data, " ")+"%")
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

// 기기 예약 정보 반환
func Device_list_appoint_select(id string) []structs.DeviceAppoint {
	var result string
	var appoints []structs.DeviceAppoint

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf("Device_list_appoint_select err : %s", err)
	}

	err = conn.QueryRow("SELECT appoint FROM appoint WHERE id = ?", id).Scan(&result)
	if err != nil {
		log.Printf("Device_list_appoint_select : %s", err)
	}

	m_err := json.Unmarshal([]byte(result), &appoints)
	if m_err != nil {
		log.Printf("marsal err : %s", err)
	}

	defer conn.Close()

	return appoints
}

// 기기 예약 정보 수정
func Device_list_appoint_update(id, name, starttime, stoptime string, onoff bool) {
	var result string
	var appoints []structs.DeviceAppoint

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf("Device_list_appoint_update err : %s", err)
	}

	err = conn.QueryRow("SELECT appoint FROM appoint WHERE id = ?", id).Scan(&result)
	if err != nil {
		log.Printf("Device_list_appoint_update : %s", err)
	}

	m_err := json.Unmarshal([]byte(result), &appoints)
	if m_err != nil {
		fmt.Printf("marsal err : %s", err)
	}

	for n, Object := range appoints {
		if Object.Name == name {
			appoints[n].Start = starttime
			appoints[n].Stop = stoptime
			appoints[n].Onoff = onoff
			break
		}
	}

	appoints_marsal, m_err := json.Marshal(appoints)
	if m_err != nil {
		fmt.Printf("marsal err : %s", err)
	}

	appoints_list := string(appoints_marsal)

	rows, err := conn.Exec("UPDATE appoint SET appoint = ? WHERE id = ?", appoints_list, id)

	if err != nil {
		log.Printf("rows err : %s", err)
	}

	r, err_i := rows.RowsAffected()

	if err_i != nil {
		log.Printf("rows_i err : %s", err)
	}

	log.Printf("Device_list_appoint_update rows : %d", r)

	defer conn.Close()
}

func isSafeColumnName(c map[string]bool, column string) bool {
	// 안전한 컬럼 이름 목록

	_, ok := c[column]
	return ok
}

// 일자별 센서값 카운트 반환
func Device_daily_list_count_select(user string) int {
	var result int

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf("Device_daily_list_count_select err : %s", err)
	}

	err = conn.QueryRow("SELECT count(A.id) FROM valhistory AS A INNER JOIN registration AS B ON A.id = B.id WHERE B.user = ?", user).Scan(&result)
	if err != nil {
		log.Printf("Device_daily_list_count_select : %s", err)
	}

	defer conn.Close()

	return result
}

// 일자별 센서값 조회
func Device_daily_list_select(user, s_order, e_order string) []structs.DeviceDailyListItem {
	var result []structs.DeviceDailyListItem
	// dbname := os.Getenv("DBNAME")
	// dbdata := os.Getenv("DBDATA")

	// conn, err := sql.Open(dbname, dbdata)

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf("Device_daily_list_select err : %s", err)
	}

	const sqlString = "SELECT A.id,A.temp,A.humi,A.co2_sensor,A.ph_ec_sensor,A.time FROM valhistory AS A INNER JOIN registration AS B ON A.id = B.id WHERE B.user = ? ORDER BY time DESC limit ?,?"

	rows, err := conn.Query(sqlString, user, s_order, e_order)
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
func Device_daily_list_select_search(user, column, data string) []structs.DeviceDailyListItem {
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

	queryStirng := fmt.Sprintf("SELECT A.id,A.temp,A.humi,A.co2_sensor,A.ph_ec_sensor,A.time FROM valhistory AS A INNER JOIN registration AS B ON A.id = B.id WHERE B.user = ? AND %s LIKE ? ORDER BY time DESC ", column)

	rows, err := conn.Query(queryStirng, user, "%"+strings.Trim(data, " ")+"%")
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
func Device_daily_graph_list_select(user string) []structs.DeviceDailyGraphListItem {
	var result []structs.DeviceDailyGraphListItem
	// dbname := os.Getenv("DBNAME")
	// dbdata := os.Getenv("DBDATA")

	// conn, err := sql.Open(dbname, dbdata)

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf("Device_daily_graph_list_select err : %s", err)
	}

	rows, err := conn.Query("SELECT A.id FROM device AS A INNER JOIN registration AS B ON A.id = B.id WHERE B.user = ?", user)
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
	// dbname := os.Getenv("DBNAME")
	// dbdata := os.Getenv("DBDATA")

	// conn, err := sql.Open(dbname, dbdata)

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
