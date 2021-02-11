package mail

import "fmt"

var (
	LinkForRegisterTeacher = "https://lk.bashdiktant.ru/registrationExpert/ycAy54mvTWhpsgV7?email=%s"
	LinkForAuthTeacher     = "https://lk.bashdiktant.ru/authorization"

	RecoveryPasswordTmpl = "Восстановление пароля"
	RecoveryPasswordText = "Добрый день!\n\nВы востановили пароль на сайте lk.bashdiktant.ru!\nВаш новый пароль: %s"
	RecoveryPasswordHTML = "<h1 style=\"color: #000\">Добрый день!</h1><br><br><h2 style=\"color: #000\">Вы востановили пароль на сайте lk.bashdiktant.ru!</h2><h2 style=\"color: #000\">Ваш новый пароль: %s</h2><br><br><br><br><br><br>"

	DictantTmplt            = "Международный Башкирский Диктант"
	BodyActivateTeacherText = fmt.Sprintf("Добрый день!\n\nВаша учетная запись в качестве Эксперта подтверждена для проведения Международного Башкирского диктанта!\nДля авторизации - перейдите по ссылке указанной ниже.\nВаша ссылка для авторизации: %s", LinkForAuthTeacher)
	BodyActivateTeacherHTML = fmt.Sprintf("<h1 style=\"color: #000\">Добрый день!</h1><br><br><h2 style=\"color: #000\">Ваша учетная запись в качестве Эксперта подтверждена для проведения Международного Башкирского диктанта!</h2><h2 style=\"color: #000\">Для авторизации - перейдите по ссылке указанной ниже.</h2><br><br><h2 style=\"color: #000\">Ваша ссылка для авторизации: %s</h2><br><br><br><br><br><br>", LinkForAuthTeacher)

	BodyInvateTeacherText        = "Добрый день!\n\nВы приглашены в качестве Эксперта для проведения Международного Башкирского диктанта!\nДля регистрации - перейдите по ссылке указанной ниже.\nВаша ссылка регистрации: %s"
	BodyInvateTeacherHTML        = "<h1 style=\"color: #000\">Добрый день!</h1><br><br><h2 style=\"color: #000\">Вы приглашены в качестве Эксперта для проведения Международного Башкирского диктанта!</h2><h2 style=\"color: #000\">Для регистрации - перейдите по ссылке указанной ниже.</h2><br><br><h2 style=\"color: #000\">Ваша ссылка регистрации: %s</h2><br><br><br><br><br><br>"
	BodyWaitTeacherText          = "Добрый день!\n\nВы зарегистрированны в качестве Эксперта для проведения Международного Башкирского диктанта!\nОднако мы не нашли Вас в списке учителей, для подтверждения Вашего аккаунта - свяжитесь с организатором и передайте ему email указанный при регистрации.\nВаш email: %s"
	BodyWaitTeacherHTML          = "<h1 style=\"color: #000\">Добрый день!</h1><br><br><h2 style=\"color: #000\">Вы зарегистрированны в качестве Эксперта для проведения Международного Башкирского диктанта!</h2><h2 style=\"color: #000\">ДОднако мы не нашли Вас в списке учителей, для подтверждения Вашего аккаунта - свяжитесь с организатором и передайте ему email указанный при регистрации.</h2><br><br><h2 style=\"color: #000\">Ваш email: %s</h2><br><br><br><br><br><br>"
	BodyWaitTeacherForSocialText = "Добрый день!\n\nВы зарегистрированны в качестве Эксперта для проведения Международного Башкирского диктанта!\nОднако мы не нашли Вас в списке учителей, для подтверждения Вашего аккаунта - свяжитесь с организатором и передайте ему email указанный при регистрации.\nВаш email: %s\nВаш пароль для входа: %s"
	BodyWaitTeacherForSocialHTML = "<h1 style=\"color: #000\">Добрый день!</h1><br><br><h2 style=\"color: #000\">Вы зарегистрированны в качестве Эксперта для проведения Международного Башкирского диктанта!</h2><h2 style=\"color: #000\">Однако мы не нашли Вас в списке учителей, для подтверждения Вашего аккаунта - свяжитесь с организатором и передайте ему email указанный при регистрации.</h2><br><br><h2 style=\"color: #000\">Ваш email: %s\n</h2><h2 style=\"color: #000\">Ваш пароль для входа: %s</h2><br><br><br><br><br><br>"

	RegistrationTmpl               = "Регистрация на сайте"
	BodyRegisterSendConfirmURLText = "Добрый день!\n\nСпасибо за регистрацию на сайте lk.bashdiktant.ru!\nДля завершения регистрации - перейдите по ссылке указанной ниже.\nВаша ссылка подтверждения: %s"
	BodyRegisterSendConfirmURLHTML = "<h1 style=\"color: #000\">Добрый день!</h1><br><br><h2 style=\"color: #000\">Спасибо за регистрацию на сайте bashdiktant.ru!</h2><h2 style=\"color: #000\">Для завершения регистрации - перейдите по ссылке указанной ниже.</h2><br><br><h2 style=\"color: #000\">Ваша ссылка подтверждения: %s</h2><br><br><br><br><br><br>"

	BodyRegisterBySocialText = "Добрый день!\n\nСпасибо за регистрацию на сайте lk.bashdiktant.ru!\nВаш пароль: %s"
	BodyRegisterBySocialHTML = "<h1 style=\"color: #000\">Добрый день!</h1><br><br><h2 style=\"color: #000\">Спасибо за регистрацию на сайте bashdiktant.ru!</h2><h2 style=\"color: #000\">Ваш пароль: %s</h2><br><br><br><br><br><br>"

	BodyYouWriteDictantText = "Добрый день!\n\nСпасибо за написание диктанта!\nМы приняли Вашу работу и скоро оповестим Вас о результатах"
	BodyYouWriteDictantHTML = "<h1 style=\"color: #000\">Добрый день!</h1><br><br><h2 style=\"color: #000\">Спасибо за написание диктанта!</h2><h2 style=\"color: #000\">Мы приняли Вашу работу и скоро оповестим Вас о результатах</h2><br><br><br><br><br><br>"

	BodyYourDictationIsCheckedText = "Добрый день!\n\nМы проверили вашу работу!\nВы можете получить сертификат в личном кабинете в разделе \"Мой диктант\""
	BodyYourDictationIsCheckedHTML = "<h1 style=\"color: #000\">Добрый день!</h1><br><br><h2 style=\"color: #000\">Мы проверили вашу работу!</h2><h2 style=\"color: #000\">Вы можете получить сертификат в личном кабинете в разделе \"Мой диктант\"</h2><br><br><br><br><br><br>"

	BodyCertificateInAttachText = "Добрый день!\n\nСпасибо за участие в Международном Башкирском диктанте\nВаш сертификат во вложении!"
	BodyCertificateInAttachHTML = "<h1 style=\"color: #000\">Добрый день!</h1><br><br><h2 style=\"color: #000\">Спасибо за участие в диктанте!</h2><h2 style=\"color: #000\">Ваш сертификат во вложении</h2><br><br><br><br><br><br>"
)
