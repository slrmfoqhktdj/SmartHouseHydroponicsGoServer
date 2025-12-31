package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"smart_house/gosrc/adminpage"
	"smart_house/gosrc/structs"
)

const (
	MessageTypeIdCall       = 2
	MessageTypeIdCallResult = 3
	MessageTypeIdCallError  = 4
)

var errInvalidAction = errors.New("invalid action")

var reqmapv16, resmapv16 map[string]func(json.RawMessage) (Payload, error)

type ocppError struct {
	id    string
	code  string
	cause string
}

func (e *ocppError) Error() string {
	return e.code + ": " + e.cause
}

// Call represents OCPP Call
type Call struct {
	MessageTypeId uint8
	UniqueId      string
	Action        string
	Payload       Payload
}

func (c *Call) getID() string {
	return c.UniqueId
}

// Create CallResult from a received Call
func (call *Call) createCallResult(r Payload) []byte {
	out := [3]interface{}{
		3,
		call.UniqueId,
		r,
	}
	raw, _ := json.Marshal(out)
	return raw
}

// Creates a CallError from a received Call
// TODO: organize error codes
func (call *Call) createCallError(err error) []byte {
	var id, code, cause string
	var ocppErr *ocppError
	if errors.As(err, &ocppErr) {
		id = ocppErr.id
		code = ocppErr.code
		cause = ocppErr.cause
	}
	if id == "" {
		id = "-1"
	}
	callError := &CallError{
		UniqueId:     id,
		ErrorCode:    code,
		ErrorDetails: cause,
	}
	switch code {
	case "ProtocolError":
		callError.ErrorDescription = "Payload for Action is incomplete"
	case "PropertyConstraintViolation":
		callError.ErrorDescription = "Payload is syntactically correct but at least one field contains an invalid value"
	case "NotImplemented":
		callError.ErrorDescription = "Requested Action is not known by receiver"
	case "TypeConstraintViolationError":
		callError.ErrorDescription = "Payload for Action is syntactically correct but at least one of the fields violates data type constraints (e.g. “somestring”: 12)"
	case "PropertyConstraintViolationError":
		callError.ErrorDescription = "Payload is syntactically correct but at least one field contains an invalid value"
	case "MessageTypeNotSupported":
		callError.ErrorDescription = "A message with an Message Type Number received that is not supported by this implementation"
	default:
		callError.ErrorDescription = "Unknown error"
	}
	return callError.marshal()
}

// CallResult represents OCPP CallResult
type CallResult struct {
	MessageTypeId uint8
	UniqueId      string
	Payload       json.RawMessage
}

func (cr *CallResult) getID() string {
	return cr.UniqueId
}

// CallError represents OCPP CallError
type CallError struct {
	MessageTypeId    uint8
	UniqueId         string
	ErrorCode        string
	ErrorDescription string
	ErrorDetails     interface{}
}

func (ce *CallError) marshal() []byte {
	ed := ce.ErrorDetails.(string)
	out := [5]interface{}{
		4,
		ce.UniqueId,
		ce.ErrorCode,
		ce.ErrorDescription,
		`{"cause":` + ed + `}`,
	}
	raw, _ := json.Marshal(out)
	return raw
}

func (ce *CallError) Error() string {
	return fmt.Sprintf("CallError: UniqueId=%s, ErrorCode=%s, ErrorDescription=%s, ErrorDetails=%s",
		ce.UniqueId, ce.ErrorCode, ce.ErrorDescription, ce.ErrorDetails)
}

func (ce *CallError) getID() string {
	return ce.UniqueId
}

type OcppMessage interface {
	getID() string
}

// umpack converts json byte to one of the ocpp messages, if not successful
// returns an error
// unpack expects ocpp messages in the below forms:
//   - [<MessageTypeId>, "<UniqueId>", "<Action>", {<Payload>}] -> Call
//   - [<MessageTypeId>, "<UniqueId>", {<Payload>}]             -> CallResult
//   - [<MessageTypeId>, "<UniqueId>", "<ErrorCode>", "<ErrorDescription>" , {<ErrorDetails>}] -> CallError
//
// if json byte does not conform with these three formats, OcppError is created
// and is returned as error
//
// TODO: improve error construction depending on protocol
func unpack(b []byte, proto string) (OcppMessage, error) {
	var rm []json.RawMessage //  raw message
	var mti uint8            //  MessageTypeId
	var ui string            //  UniqueId
	var a string             //  Action
	var p Payload            //  Payload
	var ocppMsg OcppMessage  //  OcppMessage
	var e *ocppError

	// unmarshal []byte to []json.RawMessage
	err := json.Unmarshal(b, &rm)
	if err != nil {
		e = &ocppError{
			id:    "-1",
			code:  "ProtocolError",
			cause: "Invalid JSON format",
		}
		return nil, e
	}

	// unmarshal [0]json.RawMessage to MessageTypeId
	err1 := json.Unmarshal(rm[0], &mti)
	// unmarshal [1]json.RawMessage to UniqueId
	err2 := json.Unmarshal(rm[1], &ui)
	if err1 != nil {
		if err2 != nil {
			e.id = "-1"
			e.cause = e.cause + "," + fmt.Sprintf("UniqueId: %v is not valid", rm[1])
		} else {
			e.id = ui
		}
		return nil, e
	}
	if len(ui) > 36 {
		e = &ocppError{
			code:  "ProtocolError",
			cause: fmt.Sprintf("UniqueId: %v is too long", ui),
		}
	}
	switch mti {
	case MessageTypeIdCall:
		call := &Call{
			MessageTypeId: mti,
			UniqueId:      ui,
		}
		if e != nil {
			return call, e
		}
		err = json.Unmarshal(rm[2], &a)
		if err != nil {
			e := &ocppError{
				id:    ui,
				code:  "ProtocolError",
				cause: "Message does not contain Action",
			}
			return call, e
		}
		call.Action = a

		p, err = unmarshalRequestPayload(a, rm[3], proto)

		var ocppErr *ocppError
		if err != nil {
			if errors.As(err, &ocppErr) {
				ocppErr.id = ui
			}
			return call, err
		}
		call.Payload = p
		ocppMsg = call
	case MessageTypeIdCallResult:
		p := rm[2]
		ocppMsg = &CallResult{
			MessageTypeId: mti,
			UniqueId:      ui,
			Payload:       p,
		}
	case MessageTypeIdCallError:
		var me [5]interface{}
		_ = json.Unmarshal(b, &me)
		ocppMsg = &CallError{
			MessageTypeId:    mti,
			UniqueId:         ui,
			ErrorCode:        me[2].(string),
			ErrorDescription: me[3].(string),
			ErrorDetails:     me[4],
		}
	default:
		e := &ocppError{
			id:    ui,
			code:  "MessageTypeNotSupported",
			cause: fmt.Sprintf("A message with: %v is not supported by this implementation", mti),
		}
		return nil, e
	}
	return ocppMsg, nil
}

// unmarshalRequestPayload unmarshals raw bytes of request type payload
// to a corresponding struct depending on Action and ocpp protocol
func unmarshalRequestPayload(actionName string, rawPayload json.RawMessage, proto string) (Payload, error) {
	var uf func(json.RawMessage) (Payload, error) // uf unmarshal function for a specific action request
	var ok bool

	switch proto {
	case OcppV16:
		uf, ok = reqmapv16[actionName]
	}
	if !ok {
		e := &ocppError{
			code:  "NotImplemented",
			cause: fmt.Sprintf("Action %v is not implemented", actionName),
		}
		return nil, e
	}
	return uf(rawPayload)
}

// unmarshalRequestPayloadv16 unmarshals raw request type payload to a ***Req type struct of ocppv16
func unmarshalRequestPayloadv16[T any](rawPayload json.RawMessage) (Payload, error) {
	var p T
	var payload Payload

	err := json.Unmarshal(rawPayload, &p)

	// log.Print(rawPayload)

	if err != nil {
		e := &ocppError{
			code:  "TypeConstraintViolationError",
			cause: "Call Payload is not valid",
		}

		return nil, e
	}

	// log.Print(p)

	err = ValidateV16.Struct(p)
	if err != nil {
		// TODO: construct more detailed error
		e := &ocppError{
			code:  "PropertyConstraintViolationError",
			cause: "Call Payload is not valid",
		}

		return nil, e
	}
	payload = &p

	return payload, nil
}

// unmarshalResponsePv16 unmarshals raw bytes of request type payload
// to a corresponding struct depending on actionName
func unmarshalResponsePv16(actionName string, rawPayload json.RawMessage) (Payload, error) {
	uf, ok := resmapv16[actionName] // uf unmarshal function for a specific action request
	if !ok {
		return nil, errInvalidAction
	}
	return uf(rawPayload)
}

// unmarshalResponsePayloadv16 unmarshals raw response type payload to a ***Conf type struct of ocppv16
func unmarshalResponsePayloadv16[T any](rawPayload json.RawMessage) (Payload, error) {
	var p T
	var payload Payload
	err := json.Unmarshal(rawPayload, &p)
	if err != nil {
		return nil, err
	}
	err = ValidateV16.Struct(p)
	if err != nil {
		return nil, err
	}
	payload = &p
	return payload, nil
}

func init() {

	reqmapv16 = map[string]func(json.RawMessage) (Payload, error){
		"Heartbeat": unmarshalRequestPayloadv16[structs.Heartbeat],

		"BOOTDEVICE":       unmarshalRequestPayloadv16[structs.Boot],
		"CMDCOMP":          unmarshalRequestPayloadv16[structs.CMDComp],
		"CMDPUMP":          unmarshalRequestPayloadv16[structs.CMDPump],
		"CMDFFAN":          unmarshalRequestPayloadv16[structs.CMDFFan],
		"CMDEFAN":          unmarshalRequestPayloadv16[structs.CMDEFan],
		"CMDAIRCON":        unmarshalRequestPayloadv16[structs.CMDAircon],
		"CMDCO2CREATOR":    unmarshalRequestPayloadv16[structs.CMDCo2Creator],
		"CMDLED":           unmarshalRequestPayloadv16[structs.CMDLED],
		"CMDALLONOFF":      unmarshalRequestPayloadv16[structs.CMDAllOnoff],
		"CMDREBOOTSOFT":    unmarshalRequestPayloadv16[structs.CMDRebootSoftware],
		"CMDREBOOTHARD":    unmarshalRequestPayloadv16[structs.CMDRebootHardware],
		"CMDLOGUPLOAD":     unmarshalRequestPayloadv16[adminpage.CMDLogUpload],
		"CMDPRGUPLOAD":     unmarshalRequestPayloadv16[adminpage.CMDPrgUpdate],
		"REPORTCOMP":       unmarshalRequestPayloadv16[structs.REPORTComp],
		"REPORTPUMP":       unmarshalRequestPayloadv16[structs.REPORTPump],
		"REPORTFFAN":       unmarshalRequestPayloadv16[structs.REPORTFFan],
		"REPORTEFAN":       unmarshalRequestPayloadv16[structs.REPORTEFan],
		"REPORTAIRCON":     unmarshalRequestPayloadv16[structs.REPORTAircon],
		"REPORTCO2CREATOR": unmarshalRequestPayloadv16[structs.REPORTCo2Creator],
		"REPORTLED":        unmarshalRequestPayloadv16[structs.REPORTLED],
		"REPORTTEMP":       unmarshalRequestPayloadv16[structs.REPORTTemp],
		"REPORTHUMI":       unmarshalRequestPayloadv16[structs.REPORTHumi],
		"REPORTCO2SENSOR":  unmarshalRequestPayloadv16[structs.REPORTCo2Sensor],
		"REPORTPHECSENSOR": unmarshalRequestPayloadv16[structs.REPORTPHECSensor],
		"REPORTAuto":       unmarshalRequestPayloadv16[structs.REPORTAuto],
		"REPORTALLSENSOR":  unmarshalRequestPayloadv16[structs.REPORTAllSensor],
		"ALARM":            unmarshalRequestPayloadv16[structs.Alarm],
	}

	resmapv16 = map[string]func(json.RawMessage) (Payload, error){
		"Heartbeat": unmarshalResponsePayloadv16[structs.HeartbeatConf],

		"BOOTDEVICE":       unmarshalResponsePayloadv16[structs.BootConf],
		"CMDCOMP":          unmarshalResponsePayloadv16[structs.CMDCompConf],
		"CMDPUMP":          unmarshalResponsePayloadv16[structs.CMDPumpConf],
		"CMDFFAN":          unmarshalResponsePayloadv16[structs.CMDFFanConf],
		"CMDEFAN":          unmarshalResponsePayloadv16[structs.CMDEFanConf],
		"CMDAIRCON":        unmarshalResponsePayloadv16[structs.CMDAirconConf],
		"CMDCO2CREATOR":    unmarshalResponsePayloadv16[structs.CMDCo2CreatorConf],
		"CMDLED":           unmarshalResponsePayloadv16[structs.CMDLEDConf],
		"CMDALLONOFF":      unmarshalResponsePayloadv16[structs.CMDAllOnoffConf],
		"CMDREBOOTSOFT":    unmarshalResponsePayloadv16[structs.CMDRebootSoftwareConf],
		"CMDREBOOTHARD":    unmarshalResponsePayloadv16[structs.CMDRebootHardwareConf],
		"CMDLOGUPLOAD":     unmarshalResponsePayloadv16[adminpage.CMDLogUploadConf],
		"CMDPRGUPLOAD":     unmarshalResponsePayloadv16[adminpage.CMDPrgUpdateConf],
		"REPORTCOMP":       unmarshalResponsePayloadv16[structs.REPORTCompConf],
		"REPORTPUMP":       unmarshalResponsePayloadv16[structs.REPORTPumpConf],
		"REPORTFFAN":       unmarshalResponsePayloadv16[structs.REPORTFFanConf],
		"REPORTEFAN":       unmarshalResponsePayloadv16[structs.REPORTEFanConf],
		"REPORTAIRCON":     unmarshalResponsePayloadv16[structs.REPORTAirconConf],
		"REPORTCO2CREATOR": unmarshalResponsePayloadv16[structs.REPORTCo2CreatorConf],
		"REPORTLED":        unmarshalResponsePayloadv16[structs.REPORTLEDConf],
		"REPORTTEMP":       unmarshalResponsePayloadv16[structs.REPORTTempConf],
		"REPORTHUMI":       unmarshalResponsePayloadv16[structs.REPORTHumiConf],
		"REPORTCO2SENSOR":  unmarshalResponsePayloadv16[structs.REPORTCo2SensorConf],
		"REPORTPHECSENSOR": unmarshalResponsePayloadv16[structs.REPORTPHECSensorConf],
		"REPORTAuto":       unmarshalResponsePayloadv16[structs.REPORTAutoConf],
		"ALARM":            unmarshalResponsePayloadv16[structs.AlarmConf],

		"ACMDCOMPUMP":    unmarshalResponsePayloadv16[structs.CMDCompConf],
		"ACMDPUMP":       unmarshalResponsePayloadv16[structs.CMDPumpConf],
		"ACMDFFAN":       unmarshalResponsePayloadv16[structs.CMDFFanConf],
		"ACMDEFAN":       unmarshalResponsePayloadv16[structs.CMDEFanConf],
		"ACMDAIRCON":     unmarshalResponsePayloadv16[structs.CMDAirconConf],
		"ACMDCO2CREATOR": unmarshalResponsePayloadv16[structs.CMDCo2CreatorConf],
		"ACMDLED":        unmarshalResponsePayloadv16[structs.CMDLEDConf],

		"AppointList": unmarshalResponsePayloadv16[structs.DeviceAppointList],

		"REPORTALLSENSOR": unmarshalResponsePayloadv16[structs.REPORTAllSensorConf],
	}

}
