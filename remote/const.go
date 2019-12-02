package remote

// fixed arrays
var (
	StartTag  = []byte{0x50, 0x50, 0x82, 0x7D}
	ShortZero = []byte{0x00, 0x00}
)

// command, reply and realtime codes
const (
	// command
	CmdConnect        uint16 = 0x03e8
	CmdExit           uint16 = 0x03e9
	CmdEnabledevice   uint16 = 0x03ea
	CmdDisabledevice  uint16 = 0x03eb
	CmdRestart        uint16 = 0x03ec
	CmdPoweroff       uint16 = 0x03ed
	CmdSleep          uint16 = 0x03ee
	CmdResume         uint16 = 0x03ef
	CmdCapturefinger  uint16 = 0x03f1
	CmdTestTemp       uint16 = 0x03f3
	CmdCaptureimage   uint16 = 0x03f4
	CmdRefreshdata    uint16 = 0x03f5
	CmdRefreshoption  uint16 = 0x03f6
	CmdTestvoice      uint16 = 0x03f9
	CmdGetVersion     uint16 = 0x044c
	CmdChangeSpeed    uint16 = 0x044d
	CmdAuth           uint16 = 0x044e
	CmdPrepareData    uint16 = 0x05dc
	CmdData           uint16 = 0x05dd
	CmdFreeData       uint16 = 0x05de
	CmdDataWrrq       uint16 = 0x05df
	CmdDataRdy        uint16 = 0x05e0
	CmdDbRrq          uint16 = 0x0007
	CmdUserWrq        uint16 = 0x0008
	CmdUsertempRrq    uint16 = 0x0009
	CmdUsertempWrq    uint16 = 0x000a
	CmdOptionsRrq     uint16 = 0x000b
	CmdOptionsWrq     uint16 = 0x000c
	CmdAttlogRrq      uint16 = 0x000d
	CmdClearData      uint16 = 0x000e
	CmdClearAttlog    uint16 = 0x000f
	CmdDeleteUser     uint16 = 0x0012
	CmdDeleteUsertemp uint16 = 0x0013
	CmdClearAdmin     uint16 = 0x0014
	CmdUsergrpRrq     uint16 = 0x0015
	CmdUsergrpWrq     uint16 = 0x0016
	CmdUsertzRrq      uint16 = 0x0017
	CmdUsertzWrq      uint16 = 0x0018
	CmdGrptzRrq       uint16 = 0x0019
	CmdGrptzWrq       uint16 = 0x001a
	CmdTzRrq          uint16 = 0x001b
	CmdTzWrq          uint16 = 0x001c
	CmdUlgRrq         uint16 = 0x001d
	CmdUlgWrq         uint16 = 0x001e
	CmdUnlock         uint16 = 0x001f
	CmdClearAcc       uint16 = 0x0020
	CmdClearOplog     uint16 = 0x0021
	CmdOplogRrq       uint16 = 0x0022
	CmdGetFreeSizes   uint16 = 0x0032
	CmdEnableClock    uint16 = 0x0039
	CmdStartverify    uint16 = 0x003c
	CmdStartenroll    uint16 = 0x003d
	CmdCancelcapture  uint16 = 0x003e
	CmdStateRrq       uint16 = 0x0040
	CmdWriteLcd       uint16 = 0x0042
	CmdClearLcd       uint16 = 0x0043
	CmdGetPinwidth    uint16 = 0x0045
	CmdSmsWrq         uint16 = 0x0046
	CmdSmsRrq         uint16 = 0x0047
	CmdDeleteSms      uint16 = 0x0048
	CmdUdataWrq       uint16 = 0x0049
	CmdDeleteUdata    uint16 = 0x004a
	CmdDoorstateRrq   uint16 = 0x004b
	CmdWriteMifare    uint16 = 0x004c
	CmdEmptyMifare    uint16 = 0x004e
	CmdVerifyWrq      uint16 = 0x004f
	CmdVerifyRrq      uint16 = 0x0050
	CmdTmpWrite       uint16 = 0x0057
	CmdChecksumBuffer uint16 = 0x0077
	CmdDelFptmp       uint16 = 0x0086
	CmdGetTime        uint16 = 0x00c9
	CmdSetTime        uint16 = 0x00ca
	CmdRegEvent       uint16 = 0x01f4

	// reply
	CmdAckOk        uint16 = 0x07d0
	CmdAckError     uint16 = 0x07d1
	CmdAckData      uint16 = 0x07d2
	CmdAckRetry     uint16 = 0x07d3
	CmdAckRepeat    uint16 = 0x07d4
	CmdAckUnauth    uint16 = 0x07d5
	CmdAckUnknown   uint16 = 0xffff
	CmdAckErrorCmd  uint16 = 0xfffd
	CmdAckErrorInit uint16 = 0xfffc
	CmdAckErrorData uint16 = 0xfffb

	// realtime
	EfAttlog       = 0x1
	EfFinger       = 0x2
	EfEnrolluser   = 0x4
	EfEnrollfinger = 0x8
	EfButton       = 0x10
	EfUnlock       = 0x20
	EfVerify       = 0x80
	EfFpftr        = 0x100
	EfAlarm        = 0x200
)

// status position
var (
	Status = map[string]int{
		"admin_count":      48,
		"user_count":       16,
		"fp_count":         24,
		"pwd_count":        52,
		"oplog_count":      40,
		"attlog_count":     32,
		"fp_capacity":      56,
		"user_capacity":    60,
		"attlog_capacity":  64,
		"remaining_fp":     68,
		"remaining_user":   72,
		"remaining_attlog": 76,
		"face_count":       80,
		"face_capacity":    88,
	}
)

// verification
const (
	GroupVerify   uint16 = 0
	FPorPWorRF    uint16 = 0x80
	FP            uint16 = 0x81
	PIN           uint16 = 0x82
	PW            uint16 = 0x83
	RF            uint16 = 0x84
	FPorPW        uint16 = 0x85
	FPorRF        uint16 = 0x86
	PWorRF        uint16 = 0x87
	PINandFP      uint16 = 0x88
	FPandPW       uint16 = 0x89
	FPandRF       uint16 = 0x8a
	PWandRF       uint16 = 0x8b
	FPandPWandRF  uint16 = 0x8c
	PINandFPandPW uint16 = 0x8d
	FPandRForPIN  uint16 = 0x8e
)
