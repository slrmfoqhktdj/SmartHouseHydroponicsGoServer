package device_handler

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	// _ "github.com/go-sql-driver/mysql"

	"smart_house/gosrc/db"
	"smart_house/gosrc/structs"
)

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// 스마트 하우스 기기 boot update
func Device_Boot_update(id, b_time string, state structs.BootState, appoints []structs.DeviceAppoint) bool {

	conn, err := db.OpenDB()
	if err != nil {
		fmt.Println(err)
		return false
	}

	auto_int := 0

	if state.Auto {
		auto_int = 1
	}

	appoints_marsal, m_err := json.Marshal(appoints)
	if m_err != nil {
		fmt.Printf("marsal err : %s", err)
		return false
	}

	appoints_list := string(appoints_marsal)

	fmt.Println("lenght data : %d", len(appoints_list))

	// rows, err := conn.Exec("UPDATE device SET comp = ? , pump = ? , temp = ? , humi = ? , f_fan = ? , e_fan = ?, aircon = ? , co2_sensor = ? , co2_creator = ? , ph_ec_sensor = ?, led = ?, cctv = ?, auto = ? WHERE id = ?", comp, pump, temp, humi, f_fan, e_fan, aircon, co2_sensor, co2_creator, ph_ec_sensor, led, cctv, auto_int, id)
	rows, err := conn.Exec("INSERT INTO device (id,comp,pump,temp,humi,f_fan,e_fan,aircon,co2_sensor,co2_creator,ph_ec_sensor,led,cctv,auto,b_time) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE comp = ? , pump = ? , temp = ? , humi = ? , f_fan = ? , e_fan = ?, aircon = ? , co2_sensor = ? , co2_creator = ? , ph_ec_sensor = ?, led = ?, cctv = ?, auto = ?,b_time = ?", id, state.Comp, state.Pump, state.Temp, state.Humi, state.F_Fan, state.E_Fan, state.Aircon, state.Co2_Sensor, state.C02_Creator, state.Ph_Ec_Sensor, state.Led, state.Cctv, auto_int, b_time, state.Comp, state.Pump, state.Temp, state.Humi, state.F_Fan, state.E_Fan, state.Aircon, state.Co2_Sensor, state.C02_Creator, state.Ph_Ec_Sensor, state.Led, state.Cctv, auto_int, b_time)

	if err != nil {
		fmt.Printf("rows1 : %s", err)
		return false
	}

	rows2, err2 := conn.Exec("INSERT INTO appoint (id,appoint) VALUES (?,?) ON DUPLICATE KEY UPDATE appoint = ? ", id, appoints_list, appoints_list)
	if err2 != nil {
		fmt.Printf("rows2 : %s", err2)
		return false
	}

	log.Printf("Device_Boot_update : %v,%v", rows, rows2)

	defer conn.Close()

	return true
}

// 스마트 하우스 기기 comp onoff
func Device_Comp_update(id string, value bool) bool {

	conn, err := db.OpenDB()
	if err != nil {
		fmt.Println(err)
		return false
	}

	rows, err := conn.Exec("UPDATE device SET comp = ? WHERE id = ?", value, id)
	if err != nil {
		fmt.Printf("rows1 : %s", err)
		return false
	}

	log.Printf("Device_Comp_update : %v", rows)

	defer conn.Close()

	return true
}

// 스마트 하우스 기기 comp onoff
func Device_Pump_update(id string, value bool) bool {

	conn, err := db.OpenDB()
	if err != nil {
		fmt.Println(err)
		return false
	}

	rows, err := conn.Exec("UPDATE device SET pump = ? WHERE id = ?", value, id)
	if err != nil {
		fmt.Printf("rows1 : %s", err)
		return false
	}

	log.Printf("Device_Pump_update : %v", rows)

	defer conn.Close()
	return true
}

// 스마트 하우스 기기 FFan onoff
func Device_FFan_update(id string, value bool) bool {

	conn, err := db.OpenDB()
	if err != nil {
		fmt.Println(err)
		return false
	}

	rows, err := conn.Exec("UPDATE device SET f_fan = ? WHERE id = ?", value, id)
	if err != nil {
		fmt.Printf("rows1 : %s", err)
		return false
	}

	log.Printf("Device_FFan_update : %v", rows)

	defer conn.Close()
	return true
}

// 스마트 하우스 기기 EFan onoff
func Device_EFan_update(id string, value bool) bool {

	conn, err := db.OpenDB()
	if err != nil {
		fmt.Println(err)
		return false
	}

	rows, err := conn.Exec("UPDATE device SET e_fan = ? WHERE id = ?", value, id)
	if err != nil {
		fmt.Printf("rows1 : %s", err)
		return false
	}

	log.Printf("Device_EFan_update : %v", rows)

	defer conn.Close()
	return true
}

// 스마트 하우스 기기 Aircon onoff
func Device_Aircon_update(id string, value bool) bool {

	conn, err := db.OpenDB()
	if err != nil {
		fmt.Println(err)
		return false
	}

	rows, err := conn.Exec("UPDATE device SET aircon = ? WHERE id = ?", value, id)
	if err != nil {
		fmt.Printf("rows1 : %s", err)
		return false
	}

	log.Printf("Device_Aircon_update : %v", rows)

	defer conn.Close()
	return true
}

// 스마트 하우스 기기 Co2Creator onoff
func Device_Co2Creator_update(id string, value bool) bool {

	conn, err := db.OpenDB()
	if err != nil {
		fmt.Println(err)
		return false
	}

	rows, err := conn.Exec("UPDATE device SET co2_creator = ? WHERE id = ?", value, id)
	if err != nil {
		fmt.Printf("rows1 : %s", err)
		return false
	}

	log.Printf("Device_Co2Creator_update : %v", rows)

	defer conn.Close()
	return true
}

// 스마트 하우스 기기 LED onoff
func Device_LED_update(id string, value bool) bool {

	conn, err := db.OpenDB()
	if err != nil {
		fmt.Println(err)
		return false
	}

	rows, err := conn.Exec("UPDATE device SET led = ? WHERE id = ?", value, id)
	if err != nil {
		fmt.Printf("rows1 : %s", err)
		return false
	}

	log.Printf("Device_LED_update : %v", rows)

	defer conn.Close()
	return true
}

// 스마트 하우스 기기 Temp update
func Device_Temp_update(id string, value float32) bool {

	conn, err := db.OpenDB()
	if err != nil {
		fmt.Println(err)
		return false
	}

	rows, err := conn.Exec("UPDATE device SET temp = ? WHERE id = ?", value, id)
	if err != nil {
		fmt.Printf("rows1 : %s", err)
		return false
	}

	log.Printf("Device_Temp_update : %v", rows)

	defer conn.Close()
	return true
}

// 스마트 하우스 기기 Humi
func Device_Humi_update(id string, value float32) bool {

	conn, err := db.OpenDB()
	if err != nil {
		fmt.Println(err)
		return false
	}

	rows, err := conn.Exec("UPDATE device SET humi = ? WHERE id = ?", value, id)
	if err != nil {
		fmt.Printf("rows1 : %s", err)
		return false
	}

	log.Printf("Device_Humi_update : %v", rows)

	defer conn.Close()
	return true
}

// 스마트 하우스 기기 LED onoff
func Device_Co2Sensor_update(id string, value float32) bool {

	conn, err := db.OpenDB()
	if err != nil {
		fmt.Println(err)
		return false
	}

	rows, err := conn.Exec("UPDATE device SET co2_sensor = ? WHERE id = ?", value, id)
	if err != nil {
		fmt.Printf("rows1 : %s", err)
		return false
	}

	log.Printf("Device_Co2Sensor_update : %v", rows)

	defer conn.Close()
	return true
}

// 스마트 하우스 기기 LED onoff
func Device_PHECSensor_update(id string, value float32) bool {

	conn, err := db.OpenDB()
	if err != nil {
		fmt.Println(err)
		return false
	}

	rows, err := conn.Exec("UPDATE device SET ph_ec_sensor = ? WHERE id = ?", value, id)
	if err != nil {
		fmt.Printf("rows1 : %s", err)
		return false
	}

	log.Printf("Device_PHECSensor_update : %v", rows)

	defer conn.Close()
	return true
}

// 스마트 하우스 기기 자동수동 onoff
func Device_state_update(id string, value int) bool {

	conn, err := db.OpenDB()
	if err != nil {
		fmt.Println(err)
		return false
	}

	rows, err := conn.Exec("UPDATE device SET auto = ? WHERE id = ?", value, id)
	if err != nil {
		fmt.Printf("rows1 : %s", err)
		return false
	}

	log.Printf("Device_state_update : %v", rows)

	defer conn.Close()
	return true
}

// 스마트 하우스 기기 정기 보고
func Device_state_report_insert(id string, temp float32, humi float32, co2_sensor float32, ph_ec_sensor float32, time string) bool {

	conn, err := db.OpenDB()
	if err != nil {
		fmt.Println(err)
		return false
	}

	rows, err := conn.Exec("INSERT INTO valhistory(id,temp,humi,co2_sensor,ph_ec_sensor,time) VALUES(?,?,?,?,?,?)", id, temp, humi, co2_sensor, ph_ec_sensor, time)
	if err != nil {
		fmt.Printf("rows1 : %s", err)
		return false
	}

	log.Printf("Device_state_report_insert : %v", rows)

	defer conn.Close()

	return true
}

// 스마트 하우스 기기 알람 , 등록 기기 반환
func Device_alarm_user_select(id string) ([]string, []string) {

	var os_id_result []string
	var user_result []string

	conn, err := db.OpenDB()
	if err != nil {
		fmt.Println(err)
		return []string{}, []string{}
	}

	rows, err := conn.Query("SELECT A.os_id,A.email FROM user AS A INNER JOIN registration AS B ON A.email = B.user WHERE B.id = ?", id)
	if err != nil {
		log.Printf("Device_alarm_user_select err : %s", err)
	}

	for rows.Next() {
		var User string
		var Email string

		rows.Scan(&User, &Email)

		if User == "" {
			continue
		}

		os_id_result = append(os_id_result, User)
		user_result = append(user_result, Email)
	}

	defer rows.Close()
	defer conn.Close()

	// log.Printf("Device_alarm_user_select result : %v", result)

	return os_id_result, user_result
}

// 스마트 하우스 기기 알람 저장
func Device_alarm_insert(id, title, content, time string, user []string) bool {

	conn, err := db.OpenDB()
	if err != nil {
		fmt.Println(err)
		return false
	}

	updateString, err := json.Marshal(&structs.UserAlarmList{Id: id, Title: title, Content: content, Time: time, Read: false})
	if err != nil {
		fmt.Println(err)
		return false
	}

	whereString := strings.Join(user, ",")

	log.Printf("Device_state_report_insert : %s, update string : %s", whereString, updateString)

	rows, err := conn.Exec("UPDATE user SET alarm = IF(alarm IS NULL OR alarm = '',?, CONCAT(alarm,'&&',?)) WHERE email IN (?)", updateString, updateString, whereString)
	if err != nil {
		fmt.Printf("rows1 : %s", err)
		return false
	}

	i, err := rows.RowsAffected()

	log.Printf("Device_state_report_insert : %v", i)

	defer conn.Close()

	return true
}
