go1.23.7 windows/amd64

[개요] 2025.11~2026.02 /  컨테이너형 식물공장 베이스 시스템의 중앙 서버. 
(충전기 OCPP 서버의 WebSocket/Call 구조를 재사용·확장)

[사용기술]
 - Go 1.23, gorilla/websocket, gorilla/mux, jwt-go, bcrypt(x/crypto), zap, godotenv,
   go-sql-driver/mysql, validator
 - MySQL, JWT(Access 15분 / Refresh 30일), OneSignal 푸시 연동, CFSSL 인증서 연동

[주요기능]
 - 장치(하우스) WebSocket 서버(서브프로토콜 "sh1.0")
   · 기기 → 서버 보고: BOOTDEVICE(부팅·예약정보 등록), REPORTCOMP/PUMP/FFAN/EFAN/
     AIRCON/CO2CREATOR/LED(장치 ON/OFF), REPORTTEMP/HUMI/CO2SENSOR/PHECSENSOR(센서값),
     REPORTAuto(자동/수동), REPORTALLSENSOR(정기 센서 일괄 보고), ALARM(알림)
   · 서버 → 기기 명령: 일괄 ON/OFF, 소프트웨어/하드웨어 재부팅, 컴프·펌프·유동팬·배기팬·
     냉난방기·CO2 발생기·LED 개별 제어, 장치별 예약(스케줄) 설정, 로그 업로드, 프로그램 업데이트
 - 3종 클라이언트용 REST API: 사용자 웹(/login), 앱(/app/*), 관리자(/if/*)
 - 관리자: 사용자/기기 소유자 관리, 강제 로그아웃, 전체/개별 푸시 알림 발송
 - 일자별 센서 통계/그래프, 알림 목록, 장치 예약 목록

[핵심 시스템 구조와 구현 방식]
 - 장치 통신: OCPP 스타일 Call/CallResult/CallError 프레임 위에 도메인 액션
   (식물공장 장치 제어)을 올린 WebSocket 서버. 장치 ID로 연결을 Store/Load 하여
   REST 요청을 해당 장치로 Call하고 결과를 응답.
 - 인증 3계층: 웹/앱/관리자를 분리한 로그인·리프레시 토큰 엔드포인트,
   Refresh Token은 DB 저장 + 재발급 시 갱신(RTR), 일반/관리자 JWT 미들웨어 구분.
 - 패키지: jwt, clientpage/clientpagedbconn(사용자·앱), adminpage(관리자),
   device_handler(장치 상태→DB), validation(입력 검증), shlogger, structs.
 - 센서값/장치 상태는 DB에 즉시 반영, 정기 보고(REPORTALLSENSOR)는 이력 테이블에 INSERT.
 - 푸시: 장치 ALARM 수신 시 해당 기기 사용자 목록을 조회해 OneSignal로 발송하고 알림 이력 저장.
 - 장치용 인증서 재발급 엔드포인트(CFSSL 기반 CA/클라이언트 인증서 제공).

[데이터 흐름]
 센서/장치(컨테이너 컨트롤러, Python 더미 기기 예제 있음) --WebSocket--> Go
   → MySQL(현재 상태 + 이력) → REST로 웹/앱에 조회
 사용자(웹/앱) --REST(JWT)--> Go --Call--> 장치 제어 → 결과 응답
 장치 알림 → Go → OneSignal → 사용자 스마트폰 푸시
