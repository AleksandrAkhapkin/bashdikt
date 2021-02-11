package types

import (
	"io"
	"mime/multipart"
)

const (
	RoleOrganizer string = "organizer"
	RoleTeacher   string = "teacher"
	RoleStudent   string = "student"
	RoleAdmin     string = "admin"

	LevelDicStart    string = "start"
	LevelDicAdvanced string = "advanced"
	LevelDicDialect  string = "dialect"

	FormatDicFile   = "File"
	FormatDicOnline = "Online"

	StatusStudentNotWrite   = "Не написан"
	StatusStudentChecked    = "Проверен"
	StatusStudentNotChecked = "Проверяется"
	StatusStudentNotValid   = "Отклонен"

	EmailAPIPaymentCredit     = "credit"
	EmailAPIPaymentSubscribes = "subscriber"
	EmailAPIURL               = "https://api.mailopost.ru/v1/email/messages"

	RedirectForLocal = "https://lk.bashdiktant.ru/registration"
)

type Token struct {
	Token string `json:"token"`
}

type Auth struct {
	ID    int    `json:"id"`
	Pass  string `json:"pass"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type User struct {
	ID           int    `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	MiddleName   string `json:"middle_name"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	ConfirmEmail bool   `json:"confirm_email"`
}

type ConfirmationEmail struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type RecoverPass struct {
	GeneratePass string `json:"generate_pass"`
	Email        string `json:"email"`
}

type NameFiles struct {
	Names []string `json:"names,omitempty"`
}

////////////    РЕГИСТРАЦИЯ     ////////////
type Register struct {
	ID         int    `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
	Email      string `json:"email"`
	Pass       string `json:"pass"`
	RepeatPass string `json:"repeat_pass"`
	Address    string `json:"address"`
	Role       string `json:"role"`
}

type StudentRegister struct {
	Register
	Level string `json:"level"`
}
type TeacherRegister struct {
	Register
	Info string `json:"info"`
}

type OrganizerRegister struct {
	Register
	Phone           string `json:"phone"`
	SocURL          string `json:"soc_url"`
	CountStudent    int    `json:"count_students"`
	FormatDictation string `json:"format_dictation"`
}

////////////    Изменение кабинета     ////////////
type StandardProfileForPUT struct {
	ID         int    `json:"id"`
	Email      string `json:"email"`
	LastName   string `json:"last_name"`
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	Address    string `json:"address"`
	OldPass    string `json:"old_pass"`
	Pass       string `json:"pass"`
	Role       string `json:"role"`
}

type StudentProfileForPUT struct {
	StandardProfileForPUT
	Level string `json:"level"`
}

type TeacherProfileForPUT struct {
	StandardProfileForPUT
	Info string `json:"info"`
}

type OrganizerProfileForPUT struct {
	StandardProfileForPUT
	Phone           string   `json:"phone"`
	AddPhone        []string `json:"add_phone"`
	AddEmail        []string `json:"add_email"`
	SocURL          string   `json:"soc_url"`
	CountStudent    int      `json:"count_students"`
	FormatDictation string   `json:"format_dictation"`
}

////////////    ГЕТ КАБИНЕТ     ////////////
type StandardProfile struct {
	ID         int    `json:"id"`
	Role       string `json:"role"`
	Email      string `json:"email"`
	LastName   string `json:"last_name"`
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	Address    string `json:"address"`
}

type StudentProfile struct {
	StandardProfile
	Level string `json:"level"`
}

type TeacherProfile struct {
	StandardProfile
	Info string `json:"info"`
}

type OrganizerProfile struct {
	StandardProfile
	Phone           string   `json:"phone"`
	AddPhone        []string `json:"add_phone"`
	AddEmail        []string `json:"add_email"`
	SocURL          string   `json:"soc_url"`
	CountStudent    int      `json:"count_students"`
	FormatDictation string   `json:"format_dictation"`
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type GetInfoForCabinet struct {
	Role   string
	Limit  int
	Offset int
}

type AllStudentsForCabinet struct {
	Total       int           `json:"total"`
	AllStudents []AllStudents `json:"all_people"`
}
type AllTeachersForCabinet struct {
	Total       int           `json:"total"`
	AllTeachers []AllTeachers `json:"all_people"`
}
type AllPinStudentsForCabinet struct {
	Total          int              `json:"total"`
	AllPinStudents []AllPinStudents `json:"all_people"`
}

type AllPinStudents struct {
	StandardProfileNotPass
	Status string `json:"status"`
}

type AllStudents struct {
	StandardProfileNotPass
	Status string `json:"status"`
	Level  string `json:"level"`
}

type AllTeachers struct {
	StandardProfileNotPass
	InfoAboutWork string `json:"info"`
}

type StandardProfileNotPass struct {
	ID         int    `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
	Email      string `json:"email"`
	Address    string `json:"address"`
	Role       string `json:"role"`
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

//uploadFile
type UploadFile struct {
	UserID int
	Body   io.Reader
	Head   string
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type Chat struct {
	Messages []MessageForChat `json:"messages"`
}

type MessageForChat struct {
	Message   string `json:"message"`
	MessageID int    `json:"message_id"`
	Name      string `json:"name"`
	Time      string `json:"time"`
}

type Dictation struct {
	Text   string `json:"text"`
	UserID int    `json:"user_id"`
}

type ReplyDictation struct {
	UserID    int       `json:"user_id"`
	TeacherID int       `json:"teacher_id"`
	Rating    int       `json:"rating"`
	Markers   []Markers `json:"markers"`
}

type Markers struct {
	Text     string `json:"text"`
	Position int    `json:"position"`
}

type MyDictation struct {
	NameFiles
	Text     string    `json:"text,omitempty"`
	UserID   int       `json:"user_id"`
	Status   string    `json:"status"`
	Rating   int       `json:"rating"`
	SendCert bool      `json:"send_cert"`
	Markers  []Markers `json:"markers,omitempty"`
}

type MakeCertificate struct {
	LastName     string
	FirstName    string
	MiddleName   string
	EmailTo      string
	PathForCert  string
	PathForFonts string
	PathForSave  string
	UserID       int
}

type RegisterOrAuthSocial struct {
	Email     string
	Firstname string
	Lastname  string
}

type AccessOK struct {
	AccessToken      string `json:"access_token"`
	TokenType        string `json:"token_type"`
	RefreshToken     string `json:"refresh_token"`
	ExpiresIn        string `json:"expires_in"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type AccessFB struct {
	AccessToken      string `json:"access_token"`
	TokenType        string `json:"token_type"`
	ExpiresIn        int    `json:"expires_in"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type RegOrAuthOK struct {
	Email     string `json:"email"`
	Firstname string `json:"first_name"`
	LastName  string `json:"last_name"`
	Error     string `json:"error_msg"`
}

type EmailAPI struct {
	FromEmail string `json:"from_email"`
	FromName  string `json:"from_name"`
	To        string `json:"to"`
	Subject   string `json:"subject"`
	Text      string `json:"text"`
	HTML      string `json:"html"`
	Payment   string `json:"payment"`
}
type EmailAPIAttach struct {
	EmailAPI
	Attachments multipart.Form `json:"attachments"`
}

//front error send to telega
type FrontError struct {
	Route   string   `json:"route"`
	Error   *Error   `json:"error"`
	Device  *Device  `json:"device"`
	Browser *Browser `json:"browser"`
}

type Error struct {
	Line      int    `json:"line"`
	Column    int    `json:"column"`
	SourceURL string `json:"sourceURL"`
}

type Device struct {
	IsBrowser           bool   `json:"isBrowser"`
	IsMobile            bool   `json:"isMobile"`
	Vendor              string `json:"vendor"`
	Model               string `json:"model"`
	OS                  string `json:"os"`
	UA                  string `json:"ua"`
	BrowserMajorVersion string `json:"browserMajorVersion"`
	BrowserFullVersion  string `json:"browserFullVersion"`
	BrowserName         string `json:"browserName"`
	EngineName          string `json:"engineName"`
	EngineVersion       string `json:"engineVersion"`
	OsName              string `json:"osName"`
	OsVersion           string `json:"osVersion"`
	UserAgent           string `json:"userAgent"`
}

type Browser struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type AddEmail struct {
	Email string `json:"email"`
}
