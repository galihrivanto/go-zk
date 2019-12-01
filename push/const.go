package push

// client / device response (Appendix 1)
const (
	// Successful
	Success = 1

	// Enrolling User Print, the finger print
	// of corresponding user already exists
	FingerPrintExists = 2

	// Enrolling User Print, enrollment failure
	// which is usually due to poor quality
	// or not them same fingerprint enrolled
	// for three times
	FingerPrintEnrollmentFailed = 4

	// Enrolling user Print, the enrolled fingerptint
	// already exists in the database
	FingerPrintExistsInDatabase = 5

	// Enrolling User Print, enrollment is cancelled
	FingerPrintEnrollmentCanceled = 6

	// Enrollment User Print, the equipment is busy and
	// enrollment cannot be conducted
	FingerPrintEnrollmentFailedDeviceIsBusy = 7

	// The parameter is incorrect
	ParameterIncorrect = -1

	// The transmitted user photo data doesn't match
	// the given size
	MissmatchUserPhotoDataSize = -2

	// Reading or writing is incorrect
	ReadingWritingIncorrect = -3

	// Transmitted template data does not match the given size
	MissmatchTemplateDataSize = -9

	// User specified by PIN does not exists in device
	UserPINNotExists = -10

	// The fingerprint format is illegal
	IllegalFingerPrintFormat = -11

	// The fingerprint template is illegal
	IllegalFingerPrintTemplate = -12

	// Limited capacity
	LimitedCapacity = -1001

	// Not supported by equipment
	NotSupported = -1002

	// Command Execution Timeout
	CommandTimeout = -1003

	// The data and equipment configuration are incosistent
	InconsistentConfiguration = -1004

	// The equipment is busy
	EquipmentIsBusy = -1005

	// Data is too long
	DataTooLong = -1006

	// MemoryError
	MemoryError = -1007
)

// command / operation (appendix 3)
const (
	// Startup
	Startup = 0

	// Shutdown
	Shutdown = 1

	// Authentication Fails
	AuthenticationFails = 2

	// Alarm
	Alarm = 3

	// Access Menu
	AccessMenu = 4

	// Change Setting
	ChangeSettings = 5

	// Enroll fingerprint
	EnrollFingerPrint = 6

	// Enroll password
	EnrollPassword = 7

	// Enroll HID Card
	EnrollHIDCard = 8

	// Delete User
	DeleteUser = 9

	// Delete fingerprint
	DeleteFingerPrint = 10

	// Delete Password
	DeletePassword = 11

	// Delete RF Card
	DeleteRFCard = 12

	// Clear Data
	ClearData = 13

	// Create MF Data
	CreateMFCard = 14

	// Enroll MF Card
	EnrollMFCard = 15

	// Register MF Card
	RegisterMFCard = 16

	// Delete MF Card
	DeleteMFCard = 17

	// Clear MF Card content
	ClearMFCardContent = 18

	// Move enrollment data into card
	MoveEnrolledDataIntoCard = 19

	// Copy data in the card to the machine
	CopyDataCardToMachine = 20

	// Set time
	SetTime = 21

	// Delivery Configuration
	DeliveryConfiguration = 22

	// Delete entry and exit records
	DeleteEntryAndExitRecords = 23

	// Clear administrator priviledge
	ClearAdministratorPriviledge = 24

	// Modify access group setting
	ModifyAccessGroupSetting = 25

	// Modify User access setting
	ModifyUserAccessSetting = 26

	// Modify access time period
	ModifyAccessTimePeriod = 27

	// Modify Unlocking Combination
	ModifyUnlockingCombination = 28

	// Unlock
	Unlock = 29

	// Enroll new user
	EnrollNewUser = 30

	// Change fingerprint attribute
	ChangeFingerPrintAttribute = 31

	// Duress Alarm
	DuressAlarm = 32
)

// Alarm Reason
const (
	// Door Close Detected
	DoorCloseDetected = 50

	// Door Open Detected
	DoorOpenDetected = 51

	// Out Door Button
	OutDoorButton = 53

	// Door Broken Accidentally
	DoorBrokenAccidentally = 54

	// Machine been broken
	MachineBeenBroken = 55

	// Try Invalid Verification
	TryInvalidVerfication = 58

	// Alarm Cancelled
	AlarmCancelled = 65535
)

// languages
const (
	// Simplied Chinese
	LangCN = 83

	// English
	LangEN = 69

	// Spain
	LangES = 97

	// Frech
	LangFR = 70

	// Arabic
	LangAR = 66

	// Portuguese
	LangPT = 80

	// Russia
	LangRU = 82

	// German
	LangDE = 71

	// Persian
	LangFA = 65

	// Thai
	LangTH = 76

	// Indonesian
	LangID = 73

	// Japanese
	LangJA = 74

	// Korean
	LangKO = 75

	// Vietnamese
	LangVI = 86

	// Turkish
	LangTK = 116

	// Hebrew
	LangHE = 72

	// Czech
	LangCS = 90

	// Dutch
	LangNL = 68

	// Italian
	LangIT = 105

	// Slovak
	LangSK = 89

	// Greek
	LangEL = 103

	// Polish
	LangPL = 112

	// Traditional Chinese
	LangTW = 84
)

// LangText return language code name
func LangText(code int) string {
	switch code {
	case LangAR:
		return "Arabic"
	case LangCS:
		return "Czech"
	case LangDE:
		return "German"
	case LangEL:
		return "Greek"
	case LangEN:
		return "English"
	case LangES:
		return "Spanish"
	case LangFR:
		return "French"
	case LangHE:
		return "Hebrew"
	case LangID:
		return "Indonesian"
	case LangIT:
		return "Italian"
	case LangJA:
		return "Japanese"
	case LangKO:
		return "Korean"
	case LangNL:
		return "Dutch"
	case LangPL:
		return "Polish"
	case LangPT:
		return "Portuguese"
	case LangRU:
		return "Russian"
	case LangSK:
		return "Slovak"
	case LangTH:
		return "Thai"
	case LangTK:
		return "Turkish"
	case LangTW:
		return "Traditional Chinese"
	case LangVI:
		return "Vietnamese"
	case LangCN:
		return "Simplified Chinese"
	default:
		return ""
	}
}
