package config

type Config struct {
	PostgresDsn             string               `yaml:"postgres_dsn"`
	ServerPort              string               `yaml:"server_port"`
	SecretKeyJWT            string               `yaml:"secret_key_jwt"`
	Soc                     *SocAuth             `yaml:"soc_auth"`
	Telegram                *Telegram            `yaml:"telegram"`
	Email                   *ForSendEmail        `yaml:"server_email"`
	PathForDictation        string               `yaml:"path_for_dictation"`
	GenerateCertificate     *GenerateCertificate `yaml:"generate_certificates"`
	PathSECRETForYouTubeApi string               `yaml:"path_for_credentional_youtube_api"`
	PathForClientSecret     string               `yaml:"path_for_save_user_credentional_youtube"`
}

type SocAuth struct {
	VKAppID         string `yaml:"vk_app_id"`
	VKCallbackURL   string `yaml:"vk_callback"`
	VKSecretKey     string `yaml:"vk_secret_key"`
	OKAppID         string `yaml:"ok_app_id"`
	OKCallbackURL   string `yaml:"ok_callback"`
	OKSecretKey     string `yaml:"ok_secret_key"`
	OKPublicKey     string `yaml:"ok_public_key"`
	FBAppID         string `yaml:"fb_app_id"`
	FBCallbackURL   string `yaml:"fb_callback"`
	FBSecretKey     string `yaml:"fb_secret_key"`
	AppleTeamID     string `yaml:"apple_team_id"`
	AppleClientID   string `yaml:"apple_client_id"`
	AppleCallback   string `yaml:"apple_callback"`
	AppleKeyID      string `yaml:"apple_key_id"`
	AppleSecretFile string `yaml:"apple_secret_file_auth_key"`
}

type Telegram struct {
	TelegramToken string `yaml:"telegram_token"`
	ChatID        string `yaml:"chat_id"`
}

type ForSendEmail struct {
	EmailHost        string `yaml:"host"`
	EmailPort        string `yaml:"port"`
	EmailLogin       string `yaml:"login"`
	EmailPass        string `yaml:"pass"`
	EmailUnsubscribe string `yaml:"email_unsubscribe"`
	NameSender       string `yaml:"email_name_sender"`
}

type GenerateCertificate struct {
	PathForFonts               string `yaml:"path_for_fonts"`
	PathForStudentTemplateBash string `yaml:"path_for_student_certificate_bash"`
	PathForStudentTemplateRus  string `yaml:"path_for_student_certificate_rus"`
	PathForTeacherTemplateRus  string `yaml:"path_for_teacher_certificate_rus"`
	PathForTeacherTemplateBash string `yaml:"path_for_teacher_certificate_bash"`
	PathForSaveCert            string `yaml:"path_for_save_certificate"`
}
