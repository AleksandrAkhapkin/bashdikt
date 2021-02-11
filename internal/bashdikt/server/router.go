package server

import (
	"github.com/AleksandrAkhapkin/bashdikt/internal/bashdikt/server/handlers"
	"github.com/gorilla/mux"

	"net/http"
)

func NewRouter(h *handlers.Handlers) *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	router.Use(h.RecoveryPanic)
	router.Use(h.RequestLog)

	studentRouter := router.PathPrefix("").Subrouter()
	teacherRouter := router.PathPrefix("").Subrouter()
	organizerRouter := router.PathPrefix("").Subrouter()
	adminRouter := router.PathPrefix("").Subrouter()

	authorizeRouter := router.PathPrefix("").Subrouter()

	studentRouter.Use(h.CheckRoleStudent)
	studentRouter.Use(h.CheckUserInDBUsers)

	teacherRouter.Use(h.CheckRoleTeacher)
	teacherRouter.Use(h.CheckUserInDBUsers)

	organizerRouter.Use(h.CheckRoleOrganizer)
	organizerRouter.Use(h.CheckUserInDBUsers)

	authorizeRouter.Use(h.CheckUserInDBUsers)

	adminRouter.Use(h.CheckUserInDBUsers)
	adminRouter.Use(h.CheckRoleAdmin)

	//todo////////////////////////////////////////
	//todo для трансляции необходимо добавить аккаунт ведущего трансляцию в гугл и передать ссылку!!!!!!!!!!
	//todo///////////////////////////////////////

	router.Methods(http.MethodGet).Path("/ping").HandlerFunc(h.Ping)
	router.Methods(http.MethodGet).Path("/timedictation").HandlerFunc(h.TimeDictation)

	router.Methods(http.MethodGet).Path("/vk/callback").HandlerFunc(h.RegisterOrAuthStudentVK)
	router.Methods(http.MethodGet).Path("/vk/callback/teacher").HandlerFunc(h.RegisterOrAuthTeacherVK)
	router.Methods(http.MethodGet).Path("/vk/callback/auth").HandlerFunc(h.VKAuth)
	router.Methods(http.MethodGet).Path("/vk/callback/app").HandlerFunc(h.RegiserOrAuthStudentVKApp)
	router.Methods(http.MethodGet).Path("/vk/callback/teacher/app").HandlerFunc(h.RegiserOrAuthTeacherVKApp)
	router.Methods(http.MethodGet).Path("/vk/callback/auth/app").HandlerFunc(h.AuthVKApp)

	router.Methods(http.MethodGet).Path("/ok/callback").HandlerFunc(h.RegisterOrAuthStudentOK)
	router.Methods(http.MethodGet).Path("/ok/callback/teacher").HandlerFunc(h.RegisterOrAuthTeacherOK)
	router.Methods(http.MethodGet).Path("/ok/callback/auth").HandlerFunc(h.AuthOK)
	router.Methods(http.MethodGet).Path("/ok/callback/app").HandlerFunc(h.RegisterOrAuthStudentOKApp)
	router.Methods(http.MethodGet).Path("/ok/callback/teacher/app").HandlerFunc(h.RegisterOrAuthTeacherOKApp)
	router.Methods(http.MethodGet).Path("/ok/callback/auth/app").HandlerFunc(h.AuthOKApp)

	router.Methods(http.MethodGet).Path("/fb/callback").HandlerFunc(h.RegisterOrAuthStudentFB)
	router.Methods(http.MethodGet).Path("/fb/callback/teacher").HandlerFunc(h.RegisterOrAuthTeacherFB)
	router.Methods(http.MethodGet).Path("/fb/callback/auth").HandlerFunc(h.AuthFB)
	router.Methods(http.MethodGet).Path("/fb/callback/app").HandlerFunc(h.RegisterOrAuthStudentFBApp)
	router.Methods(http.MethodGet).Path("/fb/callback/teacher/app").HandlerFunc(h.RegisterOrAuthTeacherFBApp)
	router.Methods(http.MethodGet).Path("/fb/callback/auth/app").HandlerFunc(h.AuthFBApp)

	router.Methods(http.MethodPost).Path("/apple/callback").HandlerFunc(h.RegisterOrAuthStudentAppleApp)
	router.Methods(http.MethodGet).Path("/apple/callback/teacher").HandlerFunc(h.RegisterOrAuthTeacherAppleApp)
	router.Methods(http.MethodGet).Path("/apple/callback/auth").HandlerFunc(h.AuthAppleApp)

	router.Methods(http.MethodPost).Path("/register").HandlerFunc(h.RegisterNew)
	//router.Methods(http.MethodPost).Path("/register/confirm").HandlerFunc(h.ConfirmEmail)
	router.Methods(http.MethodGet).Path("/register/link").HandlerFunc(h.ConfirmByLink)
	router.Methods(http.MethodPost).Path("/authorization").HandlerFunc(h.Auth)
	router.Methods(http.MethodPost).Path("/recoverpassword").HandlerFunc(h.RecoverPassword)
	authorizeRouter.Methods(http.MethodGet).Path("/cabinet").HandlerFunc(h.Cabinet)
	authorizeRouter.Methods(http.MethodPut).Path("/cabinet").HandlerFunc(h.PutCabinet)

	organizerRouter.Methods(http.MethodGet).Path("/cabinet/organizer/info").HandlerFunc(h.GetOrgCabinetInfo)

	teacherRouter.Methods(http.MethodGet).Path("/cabinet/teacher/info").HandlerFunc(h.GetTeacherCabinetInfo)

	studentRouter.Methods(http.MethodGet).Path("/cabinet/student/info").HandlerFunc(h.GetMyDictation)

	studentRouter.Methods(http.MethodGet).Path("/cabinet/student/info/file").HandlerFunc(h.GetDictationFileForStudent)
	studentRouter.Methods(http.MethodDelete).Path("/cabinet/student/info/file").HandlerFunc(h.DeleteFileForStudent)

	//возвращает массив названий файлов которые загрузил студент
	teacherRouter.Methods(http.MethodGet).Path("/dictation/{idStudent:[0-9]+}").HandlerFunc(h.GetDictationsNameByStudentID)
	teacherRouter.Methods(http.MethodGet).Path("/dictation/{idStudent:[0-9]+}/file").HandlerFunc(h.GetDictationFile)

	//ответить на диктант ученика(ток для учителя)
	teacherRouter.Methods(http.MethodPost).Path("/dictation/{idStudent:[0-9]+}/reply").HandlerFunc(h.ReplyDictation)

	studentRouter.Methods(http.MethodPost).Path("/dictation/upload").HandlerFunc(h.UploadFiles)

	//отправка поля диктант
	studentRouter.Methods(http.MethodPost).Path("/dictation/write").HandlerFunc(h.WriteDictation)

	authorizeRouter.Methods(http.MethodGet).Path("/cabinet/student/info/certificate").HandlerFunc(h.MakeAndSendMyCertificate)

	adminRouter.Methods(http.MethodPost).Path("/a/addemail").HandlerFunc(h.AddWhiteEmail)
	organizerRouter.Methods(http.MethodPost).Path("/addemail").HandlerFunc(h.AddWhiteEmail)

	router.Methods(http.MethodPost).Path("/fronterror").HandlerFunc(h.SendFrontError)

	router.Methods(http.MethodGet).Path("/chat").HandlerFunc(h.GetChat)
	adminRouter.Methods(http.MethodGet).Path("/a/auth/start").HandlerFunc(h.StartAuthYoutube)
	router.Methods(http.MethodGet).Path("/a/auth/finish").HandlerFunc(h.CallbackYoutube)
	return router
}
