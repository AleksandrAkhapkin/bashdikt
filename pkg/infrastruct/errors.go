package infrastruct

import "net/http"

type CustomError struct {
	msg  string
	Code int
}

func NewError(msg string, code int) *CustomError {
	return &CustomError{
		msg:  msg,
		Code: code,
	}
}

func (c *CustomError) Error() string {
	return c.msg
}

var (
	ErrorEmailIsExist               = NewError("email уже зарегистрирован", http.StatusBadRequest)
	ErrorPhoneIsExist               = NewError("Телефон уже зарегистрирован", http.StatusBadRequest)
	ErrorInternalServerError        = NewError("внутренняя ошибка сервера", http.StatusInternalServerError)
	ErrorBadRequest                 = NewError("плохие входные данные запроса", http.StatusBadRequest)
	ErrorJWTIsBroken                = NewError("jwt испорчен", http.StatusForbidden)
	ErrorPermissionDenied           = NewError("у вас недостаточно прав", http.StatusForbidden)
	ErrorPasswordOrEmailIsIncorrect = NewError("Неверный пароль или логин", http.StatusForbidden)
	ErrorPasswordsDoNotMatch        = NewError("Пароли не совпадают", http.StatusBadRequest)
	ErrorCodeIsIncorrect            = NewError("Неверный код", http.StatusForbidden)
	ErrorNOTConfirmEmail            = NewError("Ваш email не подтвержден", http.StatusForbidden)
	ErrorOldPassDontMatch           = NewError("Вы указали неверный старый пароль.", http.StatusForbidden)
	ErrorNotHaveDictant             = NewError("Диктант еще не написан", http.StatusBadRequest)
	ErrorNotHaveUser                = NewError("Пользователь не зарегистрирован", http.StatusForbidden)
	ErrorNotEmailByAuth             = NewError("Email не получен!", http.StatusBadRequest)
	ErrorNotWhiteEmail              = NewError("Не найден учитель с таким Email", http.StatusForbidden)
	ErrorAccountWaitConfirm         = NewError("Аккаунт ожидает подтверждения организатором", http.StatusForbidden)
)
