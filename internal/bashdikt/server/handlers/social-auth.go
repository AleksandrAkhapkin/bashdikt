package handlers

import (
	"crypto/md5"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/AleksandrAkhapkin/bashdikt/internal/bashdikt/types"
	"github.com/AleksandrAkhapkin/bashdikt/pkg/logger"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func (h *Handlers) RegisterOrAuthStudentVK(w http.ResponseWriter, r *http.Request) {

	//https://oauth.vk.com/authorize?client_id=7687118&display=page&redirect_uri=https://lk.bashdiktant.ru/api/vk/callback&scope=email&response_type=code&v=5.126
	errVK := r.FormValue("error")
	errVKDesc := r.FormValue("error_description")
	if errVK != "" {
		logger.LogError(fmt.Errorf("err with auth vk, err:%s\n%s", errVK, errVKDesc))
	}

	code := r.FormValue("code")
	if code == "" {
		//apiErrorEncode(w, infrastruct.ErrorBadRequest)
		http.Redirect(w, r, fmt.Sprintf("%s?err=%v", types.RedirectForLocal, "ErrorBadRequest"), 301)
		return
	}
	accessURL := fmt.Sprintf("https://oauth.vk.com/access_token?client_id=%s&client_secret=%s&redirect_uri=%s&code=%s", h.soc.VKAppID, h.soc.VKSecretKey, h.soc.VKCallbackURL, code)
	res, err := http.Get(accessURL)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with http.Get(accessURL) in RegisterOrAuthStudentVK"))
		//apiErrorEncode(w, infrastruct.ErrorInternalServerError)
		http.Redirect(w, r, fmt.Sprintf("%s?err=%v", types.RedirectForLocal, "ErrorInternalServerError"), 301)
		return
	}

	accessToken := struct {
		AccessToken string `json:"access_token"`
		Email       string `json:"email"`
	}{}

	if err = json.NewDecoder(res.Body).Decode(&accessToken); err != nil {
		logger.LogError(errors.Wrap(err, "err with Decode access token in RegisterOrAuthStudentVK"))
		//apiErrorEncode(w, infrastruct.ErrorInternalServerError)
		http.Redirect(w, r, fmt.Sprintf("%s?err=%v", types.RedirectForLocal, "ErrorInternalServerError"), 301)
		return
	}

	getUserURL := fmt.Sprintf("https://api.vk.com/method/users.get?fields=bdate&access_token=%s&v=5.124", accessToken.AccessToken)

	res, err = http.Get(getUserURL)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with http.http.Get(getUserURL) in RegisterOrAuthStudentVK"))
		//apiErrorEncode(w, infrastruct.ErrorInternalServerError)
		http.Redirect(w, r, fmt.Sprintf("%s?err=%v", types.RedirectForLocal, "ErrorInternalServerError"), 301)
		return
	}

	firstName := struct {
		Response []struct {
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
		}
	}{}

	if err = json.NewDecoder(res.Body).Decode(&firstName); err != nil {
		logger.LogError(errors.Wrap(err, "err with Decode first name in RegisterOrAuthStudentVK"))
		//apiErrorEncode(w, infrastruct.ErrorInternalServerError)
		http.Redirect(w, r, fmt.Sprintf("%s?err=%v", types.RedirectForLocal, "ErrorInternalServerError"), 301)
		return
	}

	auth := &types.RegisterOrAuthSocial{Email: accessToken.Email}
	if len(firstName.Response) > 0 {
		auth.Firstname = firstName.Response[0].FirstName
		auth.Lastname = firstName.Response[0].LastName
	}

	token, err := h.srv.RegisterOrAuthStudentBySocial(auth)
	if err != nil {
		//apiErrorEncode(w, err)
		http.Redirect(w, r, fmt.Sprintf("%s?err=%v", types.RedirectForLocal, err), 301)
		return
	}

	//apiResponseEncoder(w, token)
	http.Redirect(w, r, fmt.Sprintf("%s?jwt=%v", types.RedirectForLocal, token.Token), 301)
}

func (h *Handlers) RegisterOrAuthTeacherVK(w http.ResponseWriter, r *http.Request) {

	//https://oauth.vk.com/authorize?client_id=7687118&display=page&redirect_uri=https://lk.bashdiktant.ru/api/vk/callback/teacher&scope=email&response_type=code&v=5.126
	errVK := r.FormValue("error")
	errVKDesc := r.FormValue("error_description")
	if errVK != "" {
		logger.LogError(fmt.Errorf("err with auth vk, err:%s\n%s", errVK, errVKDesc))
	}

	code := r.FormValue("code")
	if code == "" {
		//apiErrorEncode(w, infrastruct.ErrorBadRequest)
		http.Redirect(w, r, fmt.Sprintf("%s?err=%v", types.RedirectForLocal, "ErrorBadRequest"), 301)
		return
	}
	accessURL := fmt.Sprintf("https://oauth.vk.com/access_token?client_id=%s&client_secret=%s&redirect_uri=%s/teacher&code=%s", h.soc.VKAppID, h.soc.VKSecretKey, h.soc.VKCallbackURL, code)
	res, err := http.Get(accessURL)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with http.Get(accessURL) in RegisterOrAuthTeacherVK"))
		//apiErrorEncode(w, infrastruct.ErrorInternalServerError)
		http.Redirect(w, r, fmt.Sprintf("%s?err=%v", types.RedirectForLocal, "ErrorInternalServerError"), 301)
		return
	}

	accessToken := struct {
		AccessToken string `json:"access_token"`
		Email       string `json:"email"`
	}{}

	if err = json.NewDecoder(res.Body).Decode(&accessToken); err != nil {
		logger.LogError(errors.Wrap(err, "err with Decode access token in RegisterOrAuthTeacherVK"))
		//apiErrorEncode(w, infrastruct.ErrorInternalServerError)
		http.Redirect(w, r, fmt.Sprintf("%s?err=%v", types.RedirectForLocal, "ErrorInternalServerError"), 301)
		return
	}

	getUserURL := fmt.Sprintf("https://api.vk.com/method/users.get?fields=bdate&access_token=%s&v=5.124", accessToken.AccessToken)

	res, err = http.Get(getUserURL)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with http.Get(getUserURL) in RegisterOrAuthTeacherVK"))
		//apiErrorEncode(w, infrastruct.ErrorInternalServerError)
		http.Redirect(w, r, fmt.Sprintf("%s?err=%v", types.RedirectForLocal, "ErrorInternalServerError"), 301)
		return
	}

	firstName := struct {
		Response []struct {
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
		}
	}{}

	if err = json.NewDecoder(res.Body).Decode(&firstName); err != nil {
		logger.LogError(errors.Wrap(err, "err with Decode in RegisterOrAuthTeacherVK"))
		//apiErrorEncode(w, infrastruct.ErrorInternalServerError)
		http.Redirect(w, r, fmt.Sprintf("%s?err=%v", types.RedirectForLocal, "ErrorInternalServerError"), 301)
		return
	}

	auth := &types.RegisterOrAuthSocial{Email: accessToken.Email}
	if len(firstName.Response) > 0 {
		auth.Firstname = firstName.Response[0].FirstName
		auth.Lastname = firstName.Response[0].LastName
	}

	token, err := h.srv.RegisterOrAuthTeacherBySocial(auth)
	if err != nil {
		//apiErrorEncode(w, err)
		http.Redirect(w, r, fmt.Sprintf("%s?err=%v", types.RedirectForLocal, err), 301)
		return
	}

	//apiResponseEncoder(w, token)
	http.Redirect(w, r, fmt.Sprintf("%s?jwt=%v", types.RedirectForLocal, token.Token), 301)
}

func (h *Handlers) VKAuth(w http.ResponseWriter, r *http.Request) {

	//https://oauth.vk.com/authorize?client_id=7687118&display=page&redirect_uri=https://lk.bashdiktant.ru/api/vk/callback/auth&scope=email&response_type=code&v=5.126
	errVK := r.FormValue("error")
	errVKDesc := r.FormValue("error_description")
	if errVK != "" {
		logger.LogError(fmt.Errorf("err with auth vk, err:%s\n%s", errVK, errVKDesc))
	}

	code := r.FormValue("code")
	if code == "" {
		//apiErrorEncode(w, infrastruct.ErrorBadRequest)
		http.Redirect(w, r, fmt.Sprintf("%s?err=%v", types.RedirectForLocal, "ErrorBadRequest"), 301)
		return
	}
	accessURL := fmt.Sprintf("https://oauth.vk.com/access_token?client_id=%s&client_secret=%s&redirect_uri=%s/auth&code=%s", h.soc.VKAppID, h.soc.VKSecretKey, h.soc.VKCallbackURL, code)
	res, err := http.Get(accessURL)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with http.Get(accessURL) accessURL in VKAuth"))
		//apiErrorEncode(w, infrastruct.ErrorInternalServerError)
		http.Redirect(w, r, fmt.Sprintf("%s?err=%v", types.RedirectForLocal, "ErrorInternalServerError"), 301)
		return
	}

	accessToken := struct {
		AccessToken string `json:"access_token"`
		Email       string `json:"email"`
	}{}

	if err = json.NewDecoder(res.Body).Decode(&accessToken); err != nil {
		logger.LogError(errors.Wrap(err, "err with Decode access token in VKAuth"))
		//apiErrorEncode(w, infrastruct.ErrorInternalServerError)
		http.Redirect(w, r, fmt.Sprintf("%s?err=%v", types.RedirectForLocal, "ErrorInternalServerError"), 301)
		return
	}

	getUserURL := fmt.Sprintf("https://api.vk.com/method/users.get?fields=bdate&access_token=%s&v=5.124", accessToken.AccessToken)

	res, err = http.Get(getUserURL)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with  http.Get(getUserURL) in VKAuth"))
		//apiErrorEncode(w, infrastruct.ErrorInternalServerError)
		http.Redirect(w, r, fmt.Sprintf("%s?err=%v", types.RedirectForLocal, "ErrorInternalServerError"), 301)
		return
	}

	firstName := struct {
		Response []struct {
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
		}
	}{}

	if err = json.NewDecoder(res.Body).Decode(&firstName); err != nil {
		logger.LogError(errors.Wrap(err, "err with Decode in VKAuth"))
		//apiErrorEncode(w, infrastruct.ErrorInternalServerError)
		http.Redirect(w, r, fmt.Sprintf("%s?err=%v", types.RedirectForLocal, "ErrorInternalServerError"), 301)
		return
	}

	auth := &types.RegisterOrAuthSocial{Email: accessToken.Email}
	if len(firstName.Response) > 0 {
		auth.Firstname = firstName.Response[0].FirstName
		auth.Lastname = firstName.Response[0].LastName
	}

	token, err := h.srv.AuthSocial(auth)
	if err != nil {
		//apiErrorEncode(w, err)
		http.Redirect(w, r, fmt.Sprintf("%s?err=%v", types.RedirectForLocal, err), 301)
		return
	}

	//apiResponseEncoder(w, token)
	http.Redirect(w, r, fmt.Sprintf("%s?jwt=%v", types.RedirectForLocal, token.Token), 301)
}

func (h *Handlers) RegiserOrAuthStudentVKApp(w http.ResponseWriter, r *http.Request) {

	//https://oauth.vk.com/authorize?client_id=7687118&display=page&redirect_uri=https://lk.bashdiktant.ru/api/vk/callback/app&scope=email&response_type=code&v=5.126
	errVK := r.FormValue("error")
	errVKDesc := r.FormValue("error_description")
	if errVK != "" {
		logger.LogError(fmt.Errorf("err with auth vk, err:%s\n%s", errVK, errVKDesc))
	}

	code := r.FormValue("code")
	if code == "" {
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorBadRequest"), 301)
		return
	}
	accessURL := fmt.Sprintf("https://oauth.vk.com/access_token?client_id=%s&client_secret=%s&redirect_uri=%s/app&code=%s", h.soc.VKAppID, h.soc.VKSecretKey, h.soc.VKCallbackURL, code)
	res, err := http.Get(accessURL)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with http.Get(accessURL) in RegiserOrAuthStudentVKApp"))
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorInternalServerError"), 301)
		return
	}

	accessToken := struct {
		AccessToken string `json:"access_token"`
		Email       string `json:"email"`
	}{}

	if err = json.NewDecoder(res.Body).Decode(&accessToken); err != nil {
		logger.LogError(errors.Wrap(err, "err with Decode access token in RegiserOrAuthStudentVKApp"))
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorInternalServerError"), 301)
		return
	}

	getUserURL := fmt.Sprintf("https://api.vk.com/method/users.get?fields=bdate&access_token=%s&v=5.124", accessToken.AccessToken)

	res, err = http.Get(getUserURL)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with Get getUserURL in RegiserOrAuthStudentVKApp"))
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorInternalServerError"), 301)
		return
	}

	firstName := struct {
		Response []struct {
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
		}
	}{}

	if err = json.NewDecoder(res.Body).Decode(&firstName); err != nil {
		logger.LogError(errors.Wrap(err, "err with Decode first name in RegiserOrAuthStudentVKApp"))
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorInternalServerError"), 301)
		return
	}

	auth := &types.RegisterOrAuthSocial{Email: accessToken.Email}
	if len(firstName.Response) > 0 {
		auth.Firstname = firstName.Response[0].FirstName
		auth.Lastname = firstName.Response[0].LastName
	}

	token, err := h.srv.RegisterOrAuthStudentBySocial(auth)
	if err != nil {
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", err), 301)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?jwt=%v", token.Token), 301)
}

func (h *Handlers) RegiserOrAuthTeacherVKApp(w http.ResponseWriter, r *http.Request) {

	//https://oauth.vk.com/authorize?client_id=7687118&display=page&redirect_uri=https://lk.bashdiktant.ru/api/vk/callback/teacher/app&scope=email&response_type=code&v=5.126
	errVK := r.FormValue("error")
	errVKDesc := r.FormValue("error_description")
	if errVK != "" {
		logger.LogError(fmt.Errorf("err with auth vk, err:%s\n%s", errVK, errVKDesc))
	}

	code := r.FormValue("code")
	if code == "" {
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorBadRequest"), 301)
		return
	}
	accessURL := fmt.Sprintf("https://oauth.vk.com/access_token?client_id=%s&client_secret=%s&redirect_uri=%s/teacher/app&code=%s", h.soc.VKAppID, h.soc.VKSecretKey, h.soc.VKCallbackURL, code)
	res, err := http.Get(accessURL)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with http.Get accessURL in RegiserOrAuthTeacherVKApp"))
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorInternalServerError"), 301)
		return
	}

	accessToken := struct {
		AccessToken string `json:"access_token"`
		Email       string `json:"email"`
	}{}

	if err = json.NewDecoder(res.Body).Decode(&accessToken); err != nil {
		logger.LogError(errors.Wrap(err, "err with Decode access token in RegiserOrAuthTeacherVKApp"))
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorInternalServerError"), 301)
		return
	}

	getUserURL := fmt.Sprintf("https://api.vk.com/method/users.get?fields=bdate&access_token=%s&v=5.124", accessToken.AccessToken)

	res, err = http.Get(getUserURL)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with Get getUserURL in RegiserOrAuthTeacherVKApp"))
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorInternalServerError"), 301)
		return
	}

	firstName := struct {
		Response []struct {
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
		}
	}{}

	if err = json.NewDecoder(res.Body).Decode(&firstName); err != nil {
		logger.LogError(errors.Wrap(err, "err with Decode first name in RegiserOrAuthTeacherVKApp"))
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorInternalServerError"), 301)
		return
	}

	auth := &types.RegisterOrAuthSocial{Email: accessToken.Email}
	if len(firstName.Response) > 0 {
		auth.Firstname = firstName.Response[0].FirstName
		auth.Lastname = firstName.Response[0].LastName
	}

	token, err := h.srv.RegisterOrAuthTeacherBySocial(auth)
	if err != nil {
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", err), 301)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?jwt=%v", token.Token), 301)
}

func (h *Handlers) AuthVKApp(w http.ResponseWriter, r *http.Request) {

	//https://oauth.vk.com/authorize?client_id=7687118&display=page&redirect_uri=https://lk.bashdiktant.ru/api/vk/callback/auth/app&scope=email&response_type=code&v=5.126
	errVK := r.FormValue("error")
	errVKDesc := r.FormValue("error_description")
	if errVK != "" {
		logger.LogError(fmt.Errorf("err with auth vk, err:%s\n%s", errVK, errVKDesc))
	}

	code := r.FormValue("code")
	if code == "" {
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorBadRequest"), 301)
		return
	}
	accessURL := fmt.Sprintf("https://oauth.vk.com/access_token?client_id=%s&client_secret=%s&redirect_uri=%s/auth/app&code=%s", h.soc.VKAppID, h.soc.VKSecretKey, h.soc.VKCallbackURL, code)
	res, err := http.Get(accessURL)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with http.Get accessURL in AuthVKApp"))
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorInternalServerError"), 301)
		return
	}

	accessToken := struct {
		AccessToken string `json:"access_token"`
		Email       string `json:"email"`
	}{}

	if err = json.NewDecoder(res.Body).Decode(&accessToken); err != nil {
		logger.LogError(errors.Wrap(err, "err with Decode access token in AuthVKApp"))
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorInternalServerError"), 301)
		return
	}

	getUserURL := fmt.Sprintf("https://api.vk.com/method/users.get?fields=bdate&access_token=%s&v=5.124", accessToken.AccessToken)

	res, err = http.Get(getUserURL)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with Get getUserURL in AuthVKApp"))
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorInternalServerError"), 301)
		return
	}

	firstName := struct {
		Response []struct {
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
		}
	}{}

	if err = json.NewDecoder(res.Body).Decode(&firstName); err != nil {
		logger.LogError(errors.Wrap(err, "err with Decode first name in AuthVKApp"))
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorInternalServerError"), 301)
		return
	}

	auth := &types.RegisterOrAuthSocial{Email: accessToken.Email}
	if len(firstName.Response) > 0 {
		auth.Firstname = firstName.Response[0].FirstName
		auth.Lastname = firstName.Response[0].LastName
	}

	token, err := h.srv.AuthSocial(auth)
	if err != nil {
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", err), 301)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?jwt=%v", token.Token), 301)
}

func (h *Handlers) RegisterOrAuthStudentOK(w http.ResponseWriter, r *http.Request) {

	//https://connect.ok.ru/oauth/authorize?client_id=512000641091&scope=VALUABLE_ACCESS;GET_EMAIL&response_type=code&redirect_uri=https://lk.bashdiktant.ru/api/ok/callback

	errOK := r.FormValue("error")
	if errOK != "" {
		logger.LogError(fmt.Errorf("err with authOK in RegisterOrAuthStudentOK, err:%s", errOK))
	}

	code := r.FormValue("code")
	if code == "" {
		//apiResponseEncoder(w, infrastruct.ErrorBadRequest)
		http.Redirect(w, r, fmt.Sprintf("%s?err=%v", types.RedirectForLocal, "ErrorBadRequest"), 301)
		return
	}
	accessURL := fmt.Sprintf("https://api.ok.ru/oauth/token.do?code=%s&client_id=%s&client_secret=%s&redirect_uri=%s&grant_type=authorization_code", code, h.soc.OKAppID, h.soc.OKSecretKey, h.soc.OKCallbackURL)
	res, err := http.Post(accessURL, "", nil)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with http.Get accessURL  in RegisterOrAuthStudentOK"))
		//apiErrorEncode(w, infrastruct.ErrorInternalServerError)
		http.Redirect(w, r, fmt.Sprintf("%s?err=%v", types.RedirectForLocal, "ErrorInternalServerError"), 301)
		return
	}

	accessToken := types.AccessOK{}
	if err = json.NewDecoder(res.Body).Decode(&accessToken); err != nil {
		logger.LogError(errors.Wrap(err, "err with Decode access token  in RegisterOrAuthStudentOK"))
		//apiErrorEncode(w, infrastruct.ErrorInternalServerError)
		http.Redirect(w, r, fmt.Sprintf("%s?err=%v", types.RedirectForLocal, "ErrorInternalServerError"), 301)
		return
	}

	if accessToken.Error != "" {
		logger.LogError(fmt.Errorf("err with authOK in RegisterOrAuthStudentOK, err:%s, %s", accessToken.Error, accessToken.ErrorDescription))
	}

	sig := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s%s", accessToken.AccessToken, h.soc.OKSecretKey))))
	data := []byte(fmt.Sprintf("application_key=%sfields=FIRST_NAME,LAST_NAME,EMAILformat=jsonmethod=users.getCurrentUser%s", h.soc.OKPublicKey, sig))
	sig = fmt.Sprintf("%x", md5.Sum(data))

	getUserURL := fmt.Sprintf("https://api.ok.ru/fb.do?application_key=%s&fields=FIRST_NAME,LAST_NAME,EMAIL&format=json&method=users.getCurrentUser&sig=%s&access_token=%s", h.soc.OKPublicKey, sig, accessToken.AccessToken)

	res, err = http.Post(getUserURL, "", nil)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with Get getUserURL  in RegisterOrAuthStudentOK"))
		//apiErrorEncode(w, infrastruct.ErrorInternalServerError)
		http.Redirect(w, r, fmt.Sprintf("%s?err=%v", types.RedirectForLocal, "ErrorInternalServerError"), 301)
		return
	}

	authOK := types.RegOrAuthOK{}
	if err = json.NewDecoder(res.Body).Decode(&authOK); err != nil {
		logger.LogError(errors.Wrap(err, "err with Decode first name  in RegisterOrAuthStudentOK"))
		//apiErrorEncode(w, infrastruct.ErrorInternalServerError)
		http.Redirect(w, r, fmt.Sprintf("%s?err=%v", types.RedirectForLocal, "ErrorInternalServerError"), 301)
		return
	}

	if authOK.Error != "" {
		logger.LogError(fmt.Errorf("err with authOK in RegisterOrAuthStudentOK, err:%s", authOK.Error))
	}

	auth := types.RegisterOrAuthSocial{Email: authOK.Email, Firstname: authOK.Firstname, Lastname: authOK.LastName}
	token, err := h.srv.RegisterOrAuthStudentBySocial(&auth)
	if err != nil {
		//apiErrorEncode(w, err)
		http.Redirect(w, r, fmt.Sprintf("%s?err=%v", types.RedirectForLocal, err), 301)
		return
	}

	//apiResponseEncoder(w, token)
	http.Redirect(w, r, fmt.Sprintf("%s?jwt=%v", types.RedirectForLocal, token.Token), 301)
}

func (h *Handlers) RegisterOrAuthTeacherOK(w http.ResponseWriter, r *http.Request) {

	//https://connect.ok.ru/oauth/authorize?client_id=512000641091&scope=VALUABLE_ACCESS;GET_EMAIL&response_type=code&redirect_uri=https://lk.bashdiktant.ru/api/ok/callback/teacher

	errOK := r.FormValue("error")
	if errOK != "" {
		logger.LogError(fmt.Errorf("err with authOK OK, err:%s", errOK))
	}

	code := r.FormValue("code")
	if code == "" {
		//apiResponseEncoder(w, infrastruct.ErrorBadRequest)
		http.Redirect(w, r, fmt.Sprintf("%s?err=%v", types.RedirectForLocal, "ErrorBadRequest"), 301)
		return
	}
	accessURL := fmt.Sprintf("https://api.ok.ru/oauth/token.do?code=%s&client_id=%s&client_secret=%s&redirect_uri=%s/teacher&grant_type=authorization_code", code, h.soc.OKAppID, h.soc.OKSecretKey, h.soc.OKCallbackURL)
	res, err := http.Post(accessURL, "", nil)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with http.Get accessURL  in RegisterOrAuthTeacherOK"))
		//apiErrorEncode(w, infrastruct.ErrorInternalServerError)
		http.Redirect(w, r, fmt.Sprintf("%s?err=%v", types.RedirectForLocal, "ErrorInternalServerError"), 301)
		return
	}

	accessToken := types.AccessOK{}
	if err = json.NewDecoder(res.Body).Decode(&accessToken); err != nil {
		logger.LogError(errors.Wrap(err, "err with Decode access token  in RegisterOrAuthTeacherOK"))
		//apiErrorEncode(w, infrastruct.ErrorInternalServerError)
		http.Redirect(w, r, fmt.Sprintf("%s?err=%v", types.RedirectForLocal, "ErrorInternalServerError"), 301)
		return
	}

	if accessToken.Error != "" {
		logger.LogError(fmt.Errorf("err with RegisterOrAuthTeacherOK vk, err:%s, %s", accessToken.Error, accessToken.ErrorDescription))
	}

	sig := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s%s", accessToken.AccessToken, h.soc.OKSecretKey))))
	data := []byte(fmt.Sprintf("application_key=%sfields=FIRST_NAME,LAST_NAME,EMAILformat=jsonmethod=users.getCurrentUser%s", h.soc.OKPublicKey, sig))
	sig = fmt.Sprintf("%x", md5.Sum(data))

	getUserURL := fmt.Sprintf("https://api.ok.ru/fb.do?application_key=%s&fields=FIRST_NAME,LAST_NAME,EMAIL&format=json&method=users.getCurrentUser&sig=%s&access_token=%s", h.soc.OKPublicKey, sig, accessToken.AccessToken)

	res, err = http.Post(getUserURL, "", nil)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with http.Post  in RegisterOrAuthTeacherOK"))
		//apiErrorEncode(w, infrastruct.ErrorInternalServerError)
		http.Redirect(w, r, fmt.Sprintf("%s?err=%v", types.RedirectForLocal, "ErrorInternalServerError"), 301)
		return
	}

	authOK := types.RegOrAuthOK{}
	if err = json.NewDecoder(res.Body).Decode(&authOK); err != nil {
		logger.LogError(errors.Wrap(err, "err with Decode first name  in RegisterOrAuthTeacherOK"))
		//apiErrorEncode(w, infrastruct.ErrorInternalServerError)
		http.Redirect(w, r, fmt.Sprintf("%s?err=%v", types.RedirectForLocal, "ErrorInternalServerError"), 301)
		return
	}

	if authOK.Error != "" {
		logger.LogError(fmt.Errorf("err with RegisterOrAuthTeacherOK, err:%s", authOK.Error))
	}

	auth := types.RegisterOrAuthSocial{Email: authOK.Email, Firstname: authOK.Firstname, Lastname: authOK.LastName}
	token, err := h.srv.RegisterOrAuthTeacherBySocial(&auth)
	if err != nil {
		//apiErrorEncode(w, err)
		http.Redirect(w, r, fmt.Sprintf("%s?err=%v", types.RedirectForLocal, err), 301)
		return
	}

	//apiResponseEncoder(w, token)
	http.Redirect(w, r, fmt.Sprintf("%s?jwt=%v", types.RedirectForLocal, token.Token), 301)
}

func (h *Handlers) AuthOK(w http.ResponseWriter, r *http.Request) {

	//https://connect.ok.ru/oauth/authorize?client_id=512000641091&scope=VALUABLE_ACCESS;GET_EMAIL&response_type=code&redirect_uri=https://lk.bashdiktant.ru/api/ok/callback/auth

	errOK := r.FormValue("error")
	if errOK != "" {
		logger.LogError(fmt.Errorf("err with authOK OK, err:%s", errOK))
	}

	code := r.FormValue("code")
	if code == "" {
		//apiResponseEncoder(w, infrastruct.ErrorBadRequest)
		http.Redirect(w, r, fmt.Sprintf("%s?err=%v", types.RedirectForLocal, "ErrorBadRequest"), 301)
		return
	}
	accessURL := fmt.Sprintf("https://api.ok.ru/oauth/token.do?code=%s&client_id=%s&client_secret=%s&redirect_uri=%s/auth&grant_type=authorization_code", code, h.soc.OKAppID, h.soc.OKSecretKey, h.soc.OKCallbackURL)
	res, err := http.Post(accessURL, "", nil)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with http.Get accessURL  in AuthOK"))
		//apiErrorEncode(w, infrastruct.ErrorInternalServerError)
		http.Redirect(w, r, fmt.Sprintf("%s?err=%v", types.RedirectForLocal, "ErrorInternalServerError"), 301)
		return
	}

	accessToken := types.AccessOK{}
	if err = json.NewDecoder(res.Body).Decode(&accessToken); err != nil {
		logger.LogError(errors.Wrap(err, "err with Decode access token  in AuthOK"))
		//apiErrorEncode(w, infrastruct.ErrorInternalServerError)
		http.Redirect(w, r, fmt.Sprintf("%s?err=%v", types.RedirectForLocal, "ErrorInternalServerError"), 301)
		return
	}

	if accessToken.Error != "" {
		logger.LogError(fmt.Errorf("err with AuthOK, err:%s, %s", accessToken.Error, accessToken.ErrorDescription))
	}

	sig := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s%s", accessToken.AccessToken, h.soc.OKSecretKey))))
	data := []byte(fmt.Sprintf("application_key=%sfields=FIRST_NAME,LAST_NAME,EMAILformat=jsonmethod=users.getCurrentUser%s", h.soc.OKPublicKey, sig))
	sig = fmt.Sprintf("%x", md5.Sum(data))

	getUserURL := fmt.Sprintf("https://api.ok.ru/fb.do?application_key=%s&fields=FIRST_NAME,LAST_NAME,EMAIL&format=json&method=users.getCurrentUser&sig=%s&access_token=%s", h.soc.OKPublicKey, sig, accessToken.AccessToken)

	res, err = http.Post(getUserURL, "", nil)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with Get getUserURL  in AuthOK"))
		//apiErrorEncode(w, infrastruct.ErrorInternalServerError)
		http.Redirect(w, r, fmt.Sprintf("%s?err=%v", types.RedirectForLocal, "ErrorInternalServerError"), 301)
		return
	}

	authOK := types.RegOrAuthOK{}
	if err = json.NewDecoder(res.Body).Decode(&authOK); err != nil {
		logger.LogError(errors.Wrap(err, "err with Decode first name  in AuthOK"))
		//apiErrorEncode(w, infrastruct.ErrorInternalServerError)
		http.Redirect(w, r, fmt.Sprintf("%s?err=%v", types.RedirectForLocal, "ErrorInternalServerError"), 301)
		return
	}

	if authOK.Error != "" {
		logger.LogError(fmt.Errorf("err with AuthOK, err:%s", authOK.Error))
	}

	auth := types.RegisterOrAuthSocial{
		Firstname: authOK.Firstname,
		Lastname:  authOK.LastName,
		Email:     authOK.Email,
	}

	token, err := h.srv.AuthSocial(&auth)
	if err != nil {
		//apiErrorEncode(w, err)
		http.Redirect(w, r, fmt.Sprintf("%s?err=%v", types.RedirectForLocal, "ErrorInternalServerError"), 301)
		return
	}

	//apiResponseEncoder(w, token)
	http.Redirect(w, r, fmt.Sprintf("%s?jwt=%v", types.RedirectForLocal, token.Token), 301)
}

func (h *Handlers) RegisterOrAuthStudentOKApp(w http.ResponseWriter, r *http.Request) {

	//https://connect.ok.ru/oauth/authorize?client_id=512000641091&scope=VALUABLE_ACCESS;GET_EMAIL&response_type=code&redirect_uri=https://lk.bashdiktant.ru/api/ok/callback/app

	errOK := r.FormValue("error")
	if errOK != "" {
		logger.LogError(fmt.Errorf("err with authOK OK, err:%s", errOK))
	}

	code := r.FormValue("code")
	if code == "" {
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorBadRequest"), 301)
		return
	}
	accessURL := fmt.Sprintf("https://api.ok.ru/oauth/token.do?code=%s&client_id=%s&client_secret=%s&redirect_uri=%s/app&grant_type=authorization_code", code, h.soc.OKAppID, h.soc.OKSecretKey, h.soc.OKCallbackURL)
	res, err := http.Post(accessURL, "", nil)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with http.Get accessURL  in RegisterOrAuthStudentOKApp"))
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorInternalServerError"), 301)
		return
	}

	accessToken := types.AccessOK{}
	if err = json.NewDecoder(res.Body).Decode(&accessToken); err != nil {
		logger.LogError(errors.Wrap(err, "err with Decode access token  in RegisterOrAuthStudentOKApp"))
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorInternalServerError"), 301)
		return
	}

	if accessToken.Error != "" {
		logger.LogError(fmt.Errorf("err with RegisterOrAuthStudentOKApp, err:%s, %s", accessToken.Error, accessToken.ErrorDescription))
	}

	sig := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s%s", accessToken.AccessToken, h.soc.OKSecretKey))))
	data := []byte(fmt.Sprintf("application_key=%sfields=FIRST_NAME,LAST_NAME,EMAILformat=jsonmethod=users.getCurrentUser%s", h.soc.OKPublicKey, sig))
	sig = fmt.Sprintf("%x", md5.Sum(data))

	getUserURL := fmt.Sprintf("https://api.ok.ru/fb.do?application_key=%s&fields=FIRST_NAME,LAST_NAME,EMAIL&format=json&method=users.getCurrentUser&sig=%s&access_token=%s", h.soc.OKPublicKey, sig, accessToken.AccessToken)

	res, err = http.Post(getUserURL, "", nil)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with Get getUserURL  in RegisterOrAuthStudentOKApp"))
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorInternalServerError"), 301)
		return
	}

	authOK := types.RegOrAuthOK{}
	if err = json.NewDecoder(res.Body).Decode(&authOK); err != nil {
		logger.LogError(errors.Wrap(err, "err with Decode first name  in RegisterOrAuthStudentOKApp"))
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorInternalServerError"), 301)
		return
	}

	if authOK.Error != "" {
		logger.LogError(fmt.Errorf("err with RegisterOrAuthStudentOKApp, err:%s", authOK.Error))
	}

	auth := types.RegisterOrAuthSocial{Email: authOK.Email, Firstname: authOK.Firstname, Lastname: authOK.LastName}
	token, err := h.srv.RegisterOrAuthStudentBySocial(&auth)
	if err != nil {
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", err), 301)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?jwt=%v", token.Token), 301)
}

func (h *Handlers) RegisterOrAuthTeacherOKApp(w http.ResponseWriter, r *http.Request) {

	//https://connect.ok.ru/oauth/authorize?client_id=512000641091&scope=VALUABLE_ACCESS;GET_EMAIL&response_type=code&redirect_uri=https://lk.bashdiktant.ru/api/ok/callback/teacher/app

	errOK := r.FormValue("error")
	if errOK != "" {
		logger.LogError(fmt.Errorf("err with authOK OK, err:%s", errOK))
	}

	code := r.FormValue("code")
	if code == "" {
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorBadRequest"), 301)
		return
	}
	accessURL := fmt.Sprintf("https://api.ok.ru/oauth/token.do?code=%s&client_id=%s&client_secret=%s&redirect_uri=%s/teacher/app&grant_type=authorization_code", code, h.soc.OKAppID, h.soc.OKSecretKey, h.soc.OKCallbackURL)
	res, err := http.Post(accessURL, "", nil)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with http.Get accessURL  in RegisterOrAuthTeacherOKApp"))
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorInternalServerError"), 301)
		return
	}

	accessToken := types.AccessOK{}
	if err = json.NewDecoder(res.Body).Decode(&accessToken); err != nil {
		logger.LogError(errors.Wrap(err, "err with Decode access token  in RegisterOrAuthTeacherOKApp"))
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorInternalServerError"), 301)
		return
	}

	if accessToken.Error != "" {
		logger.LogError(fmt.Errorf("err with RegisterOrAuthTeacherOKApp, err:%s, %s", accessToken.Error, accessToken.ErrorDescription))
	}

	sig := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s%s", accessToken.AccessToken, h.soc.OKSecretKey))))
	data := []byte(fmt.Sprintf("application_key=%sfields=FIRST_NAME,LAST_NAME,EMAILformat=jsonmethod=users.getCurrentUser%s", h.soc.OKPublicKey, sig))
	sig = fmt.Sprintf("%x", md5.Sum(data))

	getUserURL := fmt.Sprintf("https://api.ok.ru/fb.do?application_key=%s&fields=FIRST_NAME,LAST_NAME,EMAIL&format=json&method=users.getCurrentUser&sig=%s&access_token=%s", h.soc.OKPublicKey, sig, accessToken.AccessToken)

	res, err = http.Post(getUserURL, "", nil)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with Get getUserURL  in RegisterOrAuthTeacherOKApp"))
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorInternalServerError"), 301)
		return
	}

	authOK := types.RegOrAuthOK{}
	if err = json.NewDecoder(res.Body).Decode(&authOK); err != nil {
		logger.LogError(errors.Wrap(err, "err with Decode first name  in RegisterOrAuthTeacherOKApp"))
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorInternalServerError"), 301)
		return
	}

	if authOK.Error != "" {
		logger.LogError(fmt.Errorf("err with RegisterOrAuthTeacherOKApp, err:%s", authOK.Error))
	}

	auth := types.RegisterOrAuthSocial{Email: authOK.Email, Firstname: authOK.Firstname, Lastname: authOK.LastName}
	token, err := h.srv.RegisterOrAuthTeacherBySocial(&auth)
	if err != nil {
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", err), 301)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?jwt=%v", token.Token), 301)
}

func (h *Handlers) AuthOKApp(w http.ResponseWriter, r *http.Request) {

	//https://connect.ok.ru/oauth/authorize?client_id=512000641091&scope=VALUABLE_ACCESS;GET_EMAIL&response_type=code&redirect_uri=https://lk.bashdiktant.ru/api/ok/callback/auth/app

	errOK := r.FormValue("error")
	if errOK != "" {
		logger.LogError(fmt.Errorf("err with authOK OK, err:%s", errOK))
	}

	code := r.FormValue("code")
	if code == "" {
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorBadRequest"), 301)
		return
	}
	accessURL := fmt.Sprintf("https://api.ok.ru/oauth/token.do?code=%s&client_id=%s&client_secret=%s&redirect_uri=%s/auth/app&grant_type=authorization_code", code, h.soc.OKAppID, h.soc.OKSecretKey, h.soc.OKCallbackURL)
	res, err := http.Post(accessURL, "", nil)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with http.Get accessURL  in AuthOKApp"))
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorInternalServerError"), 301)
		return
	}

	accessToken := types.AccessOK{}
	if err = json.NewDecoder(res.Body).Decode(&accessToken); err != nil {
		logger.LogError(errors.Wrap(err, "err with Decode access token  in AuthOKApp"))
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorInternalServerError"), 301)
		return
	}

	if accessToken.Error != "" {
		logger.LogError(fmt.Errorf("err with AuthOKApp, err:%s, %s", accessToken.Error, accessToken.ErrorDescription))
	}

	sig := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s%s", accessToken.AccessToken, h.soc.OKSecretKey))))
	data := []byte(fmt.Sprintf("application_key=%sfields=FIRST_NAME,LAST_NAME,EMAILformat=jsonmethod=users.getCurrentUser%s", h.soc.OKPublicKey, sig))
	sig = fmt.Sprintf("%x", md5.Sum(data))

	getUserURL := fmt.Sprintf("https://api.ok.ru/fb.do?application_key=%s&fields=FIRST_NAME,LAST_NAME,EMAIL&format=json&method=users.getCurrentUser&sig=%s&access_token=%s", h.soc.OKPublicKey, sig, accessToken.AccessToken)

	res, err = http.Post(getUserURL, "", nil)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with Get getUserURL  in AuthOKApp"))
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorInternalServerError"), 301)
		return
	}

	authOK := types.RegOrAuthOK{}
	if err = json.NewDecoder(res.Body).Decode(&authOK); err != nil {
		logger.LogError(errors.Wrap(err, "err with Decode first name  in AuthOKApp"))
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorInternalServerError"), 301)
		return
	}

	if authOK.Error != "" {
		logger.LogError(fmt.Errorf("err with AuthOKApp, err:%s", authOK.Error))
	}

	auth := types.RegisterOrAuthSocial{
		Firstname: authOK.Firstname,
		Lastname:  authOK.LastName,
		Email:     authOK.Email,
	}

	token, err := h.srv.AuthSocial(&auth)
	if err != nil {
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", err), 301)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?jwt=%v", token.Token), 301)
}

func (h *Handlers) RegisterOrAuthStudentFB(w http.ResponseWriter, r *http.Request) {

	//https://www.facebook.com/v9.0/dialog/oauth?client_id=677170169585903&state="{st=fsdfdge3,ds=236776345}"&scope=email&redirect_uri=https://lk.bashdiktant.ru/api/fb/callback

	errFB := r.FormValue("error_reason")
	if errFB != "" {
		logger.LogError(fmt.Errorf("err with auth OK, err:%s", errFB))
	}

	code := r.FormValue("code")
	if code == "" {
		//apiResponseEncoder(w, infrastruct.ErrorBadRequest)
		http.Redirect(w, r, fmt.Sprintf("%s?err=%v", types.RedirectForLocal, "ErrorBadRequest"), 301)
		return
	}
	accessURL := fmt.Sprintf("https://graph.facebook.com/v9.0/oauth/access_token?client_id=%s&redirect_uri=%s&client_secret=%s&code=%s", h.soc.FBAppID, h.soc.FBCallbackURL, h.soc.FBSecretKey, code)
	res, err := http.Post(accessURL, "", nil)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with http.Get accessURL  in RegisterOrAuthStudentFB"))
		//apiErrorEncode(w, infrastruct.ErrorInternalServerError)
		http.Redirect(w, r, fmt.Sprintf("%s?err=%v", types.RedirectForLocal, "ErrorInternalServerError"), 301)
		return
	}

	accessToken := types.AccessFB{}
	if err = json.NewDecoder(res.Body).Decode(&accessToken); err != nil {
		logger.LogError(errors.Wrap(err, "err with Decode access token  in RegisterOrAuthStudentFB"))
		//apiErrorEncode(w, infrastruct.ErrorInternalServerError)
		http.Redirect(w, r, fmt.Sprintf("%s?err=%v", types.RedirectForLocal, "ErrorInternalServerError"), 301)
		return
	}

	if accessToken.Error != "" {
		logger.LogError(fmt.Errorf("err with auth vk, err:%s, %s", accessToken.Error, accessToken.ErrorDescription))
	}
	getUserURL := fmt.Sprintf("https://graph.facebook.com/me?access_token=%s&fields=email,first_name,last_name", accessToken.AccessToken)

	res, err = http.Post(getUserURL, "", nil)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with Get getUserURL  in RegisterOrAuthStudentFB"))
		//apiErrorEncode(w, infrastruct.ErrorInternalServerError)
		http.Redirect(w, r, fmt.Sprintf("%s?err=%v", types.RedirectForLocal, "ErrorInternalServerError"), 301)
		return
	}

	auth := types.RegisterOrAuthSocial{}
	if err = json.NewDecoder(res.Body).Decode(&auth); err != nil {
		logger.LogError(errors.Wrap(err, "err with Decode first name  in RegisterOrAuthStudentFB"))
		//apiErrorEncode(w, infrastruct.ErrorInternalServerError)
		http.Redirect(w, r, fmt.Sprintf("%s?err=%v", types.RedirectForLocal, "ErrorInternalServerError"), 301)
		return
	}

	token, err := h.srv.RegisterOrAuthStudentBySocial(&auth)
	if err != nil {
		//apiErrorEncode(w, err)
		http.Redirect(w, r, fmt.Sprintf("%s?err=%v", types.RedirectForLocal, err), 301)
		return
	}

	//apiResponseEncoder(w, token)
	http.Redirect(w, r, fmt.Sprintf("%s?jwt=%v", types.RedirectForLocal, token.Token), 301)
}

func (h *Handlers) RegisterOrAuthTeacherFB(w http.ResponseWriter, r *http.Request) {

	//https://www.facebook.com/v9.0/dialog/oauth?client_id=677170169585903&state="{st=fsdfdge3,ds=236776345}"&scope=email&redirect_uri=https://lk.bashdiktant.ru/api/fb/callback/teacher

	errFB := r.FormValue("error_reason")
	if errFB != "" {
		logger.LogError(fmt.Errorf("err with auth OK, err:%s", errFB))
	}

	code := r.FormValue("code")
	if code == "" {
		//apiResponseEncoder(w, infrastruct.ErrorBadRequest)
		http.Redirect(w, r, fmt.Sprintf("%s?err=%v", types.RedirectForLocal, "ErrorInternalServerError"), 301)
		return
	}
	accessURL := fmt.Sprintf("https://graph.facebook.com/v9.0/oauth/access_token?client_id=%s&redirect_uri=%s/teacher&client_secret=%s&code=%s", h.soc.FBAppID, h.soc.FBCallbackURL, h.soc.FBSecretKey, code)
	res, err := http.Post(accessURL, "", nil)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with http.Get accessURL  in RegisterOrAuthTeacherFB"))
		//apiErrorEncode(w, infrastruct.ErrorInternalServerError)
		http.Redirect(w, r, fmt.Sprintf("%s?err=%v", types.RedirectForLocal, "ErrorInternalServerError"), 301)
		return
	}

	accessToken := types.AccessFB{}
	if err = json.NewDecoder(res.Body).Decode(&accessToken); err != nil {
		logger.LogError(errors.Wrap(err, "err with Decode access token  in RegisterOrAuthTeacherFB"))
		//apiErrorEncode(w, infrastruct.ErrorInternalServerError)
		http.Redirect(w, r, fmt.Sprintf("%s?err=%v", types.RedirectForLocal, "ErrorInternalServerError"), 301)
		return
	}

	if accessToken.Error != "" {
		logger.LogError(fmt.Errorf("err with auth vk, err:%s, %s", accessToken.Error, accessToken.ErrorDescription))
	}
	getUserURL := fmt.Sprintf("https://graph.facebook.com/me?access_token=%s&fields=email,first_name,last_name", accessToken.AccessToken)

	res, err = http.Post(getUserURL, "", nil)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with Get getUserURL  in RegisterOrAuthTeacherFB"))
		//apiErrorEncode(w, infrastruct.ErrorInternalServerError)
		http.Redirect(w, r, fmt.Sprintf("%s?err=%v", types.RedirectForLocal, "ErrorInternalServerError"), 301)
		return
	}

	auth := types.RegisterOrAuthSocial{}
	if err = json.NewDecoder(res.Body).Decode(&auth); err != nil {
		logger.LogError(errors.Wrap(err, "err with Decode first name  in RegisterOrAuthTeacherFB"))
		//apiErrorEncode(w, infrastruct.ErrorInternalServerError)
		http.Redirect(w, r, fmt.Sprintf("%s?err=%v", types.RedirectForLocal, "ErrorInternalServerError"), 301)
		return
	}

	token, err := h.srv.RegisterOrAuthTeacherBySocial(&auth)
	if err != nil {
		//apiErrorEncode(w, err)
		http.Redirect(w, r, fmt.Sprintf("%s?err=%v", types.RedirectForLocal, err), 301)
		return
	}

	//apiResponseEncoder(w, token)
	http.Redirect(w, r, fmt.Sprintf("%s?jwt=%v", types.RedirectForLocal, token.Token), 301)
}

func (h *Handlers) AuthFB(w http.ResponseWriter, r *http.Request) {

	//https://www.facebook.com/v9.0/dialog/oauth?client_id=677170169585903&state="{st=fsdfdge3,ds=236776345}"&scope=email&redirect_uri=https://lk.bashdiktant.ru/api/fb/callback/auth

	errFB := r.FormValue("error_reason")
	if errFB != "" {
		logger.LogError(fmt.Errorf("err with auth OK, err:%s", errFB))
	}

	code := r.FormValue("code")
	if code == "" {
		//apiResponseEncoder(w, infrastruct.ErrorBadRequest)
		http.Redirect(w, r, fmt.Sprintf("%s?err=%v", types.RedirectForLocal, "ErrorBadRequest"), 301)
		return
	}
	accessURL := fmt.Sprintf("https://graph.facebook.com/v9.0/oauth/access_token?client_id=%s&redirect_uri=%s/auth&client_secret=%s&code=%s", h.soc.FBAppID, h.soc.FBCallbackURL, h.soc.FBSecretKey, code)
	res, err := http.Post(accessURL, "", nil)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with http.Get accessURL  in RegisterOrAuthTeacherFB"))
		//apiErrorEncode(w, infrastruct.ErrorInternalServerError)
		http.Redirect(w, r, fmt.Sprintf("%s?err=%v", types.RedirectForLocal, "ErrorInternalServerError"), 301)
		return
	}

	accessToken := types.AccessFB{}
	if err = json.NewDecoder(res.Body).Decode(&accessToken); err != nil {
		logger.LogError(errors.Wrap(err, "err with Decode access token  in RegisterOrAuthTeacherFB"))
		//apiErrorEncode(w, infrastruct.ErrorInternalServerError)
		http.Redirect(w, r, fmt.Sprintf("%s?err=%v", types.RedirectForLocal, "ErrorInternalServerError"), 301)
		return
	}

	if accessToken.Error != "" {
		logger.LogError(fmt.Errorf("err with auth vk, err:%s, %s", accessToken.Error, accessToken.ErrorDescription))
	}
	getUserURL := fmt.Sprintf("https://graph.facebook.com/me?access_token=%s&fields=email,first_name,last_name", accessToken.AccessToken)

	res, err = http.Post(getUserURL, "", nil)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with Get getUserURL  in RegisterOrAuthTeacherFB"))
		//apiErrorEncode(w, infrastruct.ErrorInternalServerError)
		http.Redirect(w, r, fmt.Sprintf("%s?err=%v", types.RedirectForLocal, "ErrorInternalServerError"), 301)
		return
	}

	auth := types.RegisterOrAuthSocial{}
	if err = json.NewDecoder(res.Body).Decode(&auth); err != nil {
		logger.LogError(errors.Wrap(err, "err with Decode first name  in RegisterOrAuthTeacherFB"))
		//apiErrorEncode(w, infrastruct.ErrorInternalServerError)
		http.Redirect(w, r, fmt.Sprintf("%s?err=%v", types.RedirectForLocal, "ErrorInternalServerError"), 301)
		return
	}

	token, err := h.srv.AuthSocial(&auth)
	if err != nil {
		//apiErrorEncode(w, err)
		http.Redirect(w, r, fmt.Sprintf("%s?err=%v", types.RedirectForLocal, err), 301)
		return
	}

	//apiResponseEncoder(w, token)
	http.Redirect(w, r, fmt.Sprintf("%s?jwt=%v", types.RedirectForLocal, token.Token), 301)
}

func (h *Handlers) RegisterOrAuthStudentFBApp(w http.ResponseWriter, r *http.Request) {

	//https://www.facebook.com/v9.0/dialog/oauth?client_id=677170169585903&state="{st=fsdfdge3,ds=236776345}"&scope=email&redirect_uri=https://lk.bashdiktant.ru/api/fb/callback/app

	errFB := r.FormValue("error_reason")
	if errFB != "" {
		logger.LogError(fmt.Errorf("err with auth OK, err:%s", errFB))
	}

	code := r.FormValue("code")
	if code == "" {
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorBadRequest"), 301)
		return
	}
	accessURL := fmt.Sprintf("https://graph.facebook.com/v9.0/oauth/access_token?client_id=%s&redirect_uri=%s/app&client_secret=%s&code=%s", h.soc.FBAppID, h.soc.FBCallbackURL, h.soc.FBSecretKey, code)
	res, err := http.Post(accessURL, "", nil)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with http.Get accessURL  in RegisterOrAuthStudentFB"))
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorInternalServerError"), 301)
		return
	}

	accessToken := types.AccessFB{}
	if err = json.NewDecoder(res.Body).Decode(&accessToken); err != nil {
		logger.LogError(errors.Wrap(err, "err with Decode access token  in RegisterOrAuthStudentFB"))
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorInternalServerError"), 301)
		return
	}

	if accessToken.Error != "" {
		logger.LogError(fmt.Errorf("err with auth vk, err:%s, %s", accessToken.Error, accessToken.ErrorDescription))
	}
	getUserURL := fmt.Sprintf("https://graph.facebook.com/me?access_token=%s&fields=email,first_name,last_name", accessToken.AccessToken)

	res, err = http.Post(getUserURL, "", nil)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with Get getUserURL  in RegisterOrAuthStudentFB"))
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorInternalServerError"), 301)
		return
	}

	auth := types.RegisterOrAuthSocial{}
	if err = json.NewDecoder(res.Body).Decode(&auth); err != nil {
		logger.LogError(errors.Wrap(err, "err with Decode first name  in RegisterOrAuthStudentFB"))
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorInternalServerError"), 301)
		return
	}

	token, err := h.srv.RegisterOrAuthStudentBySocial(&auth)
	if err != nil {
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", err), 301)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?jwt=%s", token.Token), 301)
}

func (h *Handlers) RegisterOrAuthTeacherFBApp(w http.ResponseWriter, r *http.Request) {

	//https://www.facebook.com/v9.0/dialog/oauth?client_id=677170169585903&state="{st=fsdfdge3,ds=236776345}"&scope=email&redirect_uri=https://lk.bashdiktant.ru/api/fb/callback/teacher/app

	errFB := r.FormValue("error_reason")
	if errFB != "" {
		logger.LogError(fmt.Errorf("err with auth OK, err:%s", errFB))
	}

	code := r.FormValue("code")
	if code == "" {
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorBadRequest"), 301)
		return
	}
	accessURL := fmt.Sprintf("https://graph.facebook.com/v9.0/oauth/access_token?client_id=%s&redirect_uri=%s/teacher&client_secret=%s&code=%s", h.soc.FBAppID, h.soc.FBCallbackURL, h.soc.FBSecretKey, code)
	res, err := http.Post(accessURL, "", nil)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with http.Get accessURL  in RegisterOrAuthTeacherFB"))
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorInternalServerError"), 301)
		return
	}

	accessToken := types.AccessFB{}
	if err = json.NewDecoder(res.Body).Decode(&accessToken); err != nil {
		logger.LogError(errors.Wrap(err, "err with Decode access token  in RegisterOrAuthTeacherFB"))
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorInternalServerError"), 301)
		return
	}

	if accessToken.Error != "" {
		logger.LogError(fmt.Errorf("err with auth vk, err:%s, %s", accessToken.Error, accessToken.ErrorDescription))
	}
	getUserURL := fmt.Sprintf("https://graph.facebook.com/me?access_token=%s&fields=email,first_name,last_name", accessToken.AccessToken)

	res, err = http.Post(getUserURL, "", nil)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with Get getUserURL  in RegisterOrAuthTeacherFB"))
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorInternalServerError"), 301)
		return
	}

	auth := types.RegisterOrAuthSocial{}
	if err = json.NewDecoder(res.Body).Decode(&auth); err != nil {
		logger.LogError(errors.Wrap(err, "err with Decode first name  in RegisterOrAuthTeacherFB"))
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorInternalServerError"), 301)
		return
	}

	token, err := h.srv.RegisterOrAuthTeacherBySocial(&auth)
	if err != nil {
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", err), 301)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?jwt=%s", token.Token), 301)
}

func (h *Handlers) AuthFBApp(w http.ResponseWriter, r *http.Request) {

	//https://www.facebook.com/v9.0/dialog/oauth?client_id=677170169585903&state="{st=fsdfdge3,ds=236776345}"&scope=email&redirect_uri=https://lk.bashdiktant.ru/api/fb/callback/auth/app

	errFB := r.FormValue("error_reason")
	if errFB != "" {
		logger.LogError(fmt.Errorf("err with auth OK, err:%s", errFB))
	}

	code := r.FormValue("code")
	if code == "" {
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorBadRequest"), 301)
		return
	}
	accessURL := fmt.Sprintf("https://graph.facebook.com/v9.0/oauth/access_token?client_id=%s&redirect_uri=%s/auth/app&client_secret=%s&code=%s", h.soc.FBAppID, h.soc.FBCallbackURL, h.soc.FBSecretKey, code)
	res, err := http.Post(accessURL, "", nil)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with http.Get accessURL  in RegisterOrAuthTeacherFB"))
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorInternalServerError"), 301)
		return
	}

	accessToken := types.AccessFB{}
	if err = json.NewDecoder(res.Body).Decode(&accessToken); err != nil {
		logger.LogError(errors.Wrap(err, "err with Decode access token  in RegisterOrAuthTeacherFB"))
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorInternalServerError"), 301)
		return
	}

	if accessToken.Error != "" {
		logger.LogError(fmt.Errorf("err with auth vk, err:%s, %s", accessToken.Error, accessToken.ErrorDescription))
	}
	getUserURL := fmt.Sprintf("https://graph.facebook.com/me?access_token=%s&fields=email,first_name,last_name", accessToken.AccessToken)

	res, err = http.Post(getUserURL, "", nil)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with Get getUserURL  in RegisterOrAuthTeacherFB"))
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorInternalServerError"), 301)
		return
	}

	auth := types.RegisterOrAuthSocial{}
	if err = json.NewDecoder(res.Body).Decode(&auth); err != nil {
		logger.LogError(errors.Wrap(err, "err with Decode first name  in RegisterOrAuthTeacherFB"))
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorInternalServerError"), 301)
		return
	}

	token, err := h.srv.AuthSocial(&auth)
	if err != nil {
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", err), 301)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?jwt=%s", token.Token), 301)
}

func (h *Handlers) RegisterOrAuthStudentAppleApp(w http.ResponseWriter, r *http.Request) {

	//https://appleid.apple.com/auth/authorize?client_id=app.bashdictant.com&redirect_uri=https://lk.bashdiktant.ru/api/apple/callback&response_type=code&scope=name+email&state=sxdcfvgbhn&response_mode=form_post

	code := r.FormValue("code")
	if code == "" {
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorBadRequest"), 301)
		return
	}

	claimsForApple := jwt.StandardClaims{
		Issuer:    h.soc.AppleTeamID,
		Audience:  "https://appleid.apple.com",
		Subject:   h.soc.AppleClientID,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
	}
	clientSecret := jwt.NewWithClaims(jwt.SigningMethodES256, claimsForApple)
	clientSecret.Header["kid"] = h.soc.AppleKeyID
	bytes, err := ioutil.ReadFile(h.soc.AppleSecretFile)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with ioutil.ReadFile  in RegisterOrAuthStudentAppleApp"))
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorInternalServerError"), 301)
		return
	}
	block, _ := pem.Decode(bytes)
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with ParsePKCS8PrivateKey  in RegisterOrAuthStudentAppleApp"))
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorInternalServerError"), 301)
		return
	}
	jwtForApple, err := clientSecret.SignedString(key)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with WriteField3  in RegisterOrAuthStudentAppleApp"))
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorInternalServerError"), 301)
		return
	}

	data := url.Values{}
	data.Set("client_id", h.soc.AppleClientID)
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("client_secret", jwtForApple)

	req, err := http.NewRequest("POST", "https://appleid.apple.com/auth/token", strings.NewReader(data.Encode()))
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with NewRequest  in RegisterOrAuthStudentAppleApp"))
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorInternalServerError"), 301)
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with DefaultClient  in RegisterOrAuthStudentAppleApp"))
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorInternalServerError"), 301)
		return
	}

	info := struct {
		Info string `json:"id_token"`
	}{}

	if err := json.NewDecoder(res.Body).Decode(&info); err != nil {
		logger.LogError(errors.Wrap(err, "err with NewDecoder(res2.Body)  in RegisterOrAuthStudentAppleApp"))
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorInternalServerError"), 301)
		return
	}

	tokenString := info.Info
	cl := jwt.MapClaims{}
	var email string
	_, err = jwt.ParseWithClaims(tokenString, cl, func(token *jwt.Token) (interface{}, error) {
		return []byte(""), nil
	})
	for key, val := range cl {
		if key == "email" {
			email = fmt.Sprintf("%v", val)
			break
		}
	}

	auth := types.RegisterOrAuthSocial{
		Email: email,
	}
	token, err := h.srv.RegisterOrAuthStudentBySocial(&auth)
	if err != nil {
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", err), 301)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?jwt=%s", token.Token), 301)
}

func (h *Handlers) AuthAppleApp(w http.ResponseWriter, r *http.Request) {

	//https://appleid.apple.com/auth/authorize?client_id=app.bashdictant.com&redirect_uri=https://lk.bashdiktant.ru/api/apple/callback/auth&response_type=code&scope=name%20email&state=sxdcfvgbhn&response_mode=form_post

	code := r.FormValue("code")
	if code == "" {
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorBadRequest"), 301)
		return
	}

	claimsForApple := jwt.StandardClaims{
		Issuer:    h.soc.AppleTeamID,
		Audience:  "https://appleid.apple.com",
		Subject:   h.soc.AppleClientID,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
	}
	clientSecret := jwt.NewWithClaims(jwt.SigningMethodES256, claimsForApple)
	clientSecret.Header["kid"] = h.soc.AppleKeyID
	bytes, err := ioutil.ReadFile(h.soc.AppleSecretFile)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with ioutil.ReadFile  in AuthAppleApp"))
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorInternalServerError"), 301)
		return
	}
	block, _ := pem.Decode(bytes)
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with ParsePKCS8PrivateKey  in AuthAppleApp"))
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorInternalServerError"), 301)
		return
	}
	jwtForApple, err := clientSecret.SignedString(key)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with WriteField3  in AuthAppleApp"))
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorInternalServerError"), 301)
		return
	}

	data := url.Values{}
	data.Set("client_id", h.soc.AppleClientID)
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("client_secret", jwtForApple)

	req, err := http.NewRequest("POST", "https://appleid.apple.com/auth/token", strings.NewReader(data.Encode()))
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with NewRequest  in AuthAppleApp"))
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorInternalServerError"), 301)
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with DefaultClient  in AuthAppleApp"))
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorInternalServerError"), 301)
		return
	}

	info := struct {
		Info string `json:"id_token"`
	}{}

	if err := json.NewDecoder(res.Body).Decode(&info); err != nil {
		logger.LogError(errors.Wrap(err, "err with NewDecoder(res2.Body)  in AuthAppleApp"))
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorInternalServerError"), 301)
		return
	}

	tokenString := info.Info
	cl := jwt.MapClaims{}
	var email string
	_, err = jwt.ParseWithClaims(tokenString, cl, func(token *jwt.Token) (interface{}, error) {
		return []byte(""), nil
	})
	for key, val := range cl {
		if key == "email" {
			email = fmt.Sprintf("%v", val)
			break
		}
	}

	auth := types.RegisterOrAuthSocial{
		Email: email,
	}
	token, err := h.srv.AuthSocial(&auth)
	if err != nil {
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", err), 301)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?jwt=%s", token.Token), 301)
}

func (h *Handlers) RegisterOrAuthTeacherAppleApp(w http.ResponseWriter, r *http.Request) {

	//https://appleid.apple.com/auth/authorize?client_id=app.bashdictant.com&redirect_uri=https://lk.bashdiktant.ru/api/apple/callback/teacher&response_type=code&scope=name%20email&state=sxdcfvgbhn&response_mode=form_post

	code := r.FormValue("code")
	if code == "" {
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorBadRequest"), 301)
		return
	}

	claimsForApple := jwt.StandardClaims{
		Issuer:    h.soc.AppleTeamID,
		Audience:  "https://appleid.apple.com",
		Subject:   h.soc.AppleClientID,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
	}
	clientSecret := jwt.NewWithClaims(jwt.SigningMethodES256, claimsForApple)
	clientSecret.Header["kid"] = h.soc.AppleKeyID
	bytes, err := ioutil.ReadFile(h.soc.AppleSecretFile)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with ioutil.ReadFile  in RegisterOrAuthStudentAppleApp"))
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorInternalServerError"), 301)
		return
	}
	block, _ := pem.Decode(bytes)
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with ParsePKCS8PrivateKey  in RegisterOrAuthStudentAppleApp"))
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorInternalServerError"), 301)
		return
	}
	jwtForApple, err := clientSecret.SignedString(key)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with WriteField3  in RegisterOrAuthStudentAppleApp"))
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorInternalServerError"), 301)
		return
	}

	data := url.Values{}
	data.Set("client_id", h.soc.AppleClientID)
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("client_secret", jwtForApple)

	req, err := http.NewRequest("POST", "https://appleid.apple.com/auth/token", strings.NewReader(data.Encode()))
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with NewRequest  in RegisterOrAuthStudentAppleApp"))
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorInternalServerError"), 301)
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with DefaultClient  in RegisterOrAuthStudentAppleApp"))
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorInternalServerError"), 301)
		return
	}

	info := struct {
		Info string `json:"id_token"`
	}{}

	if err := json.NewDecoder(res.Body).Decode(&info); err != nil {
		logger.LogError(errors.Wrap(err, "err with NewDecoder(res2.Body)  in RegisterOrAuthStudentAppleApp"))
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", "ErrorInternalServerError"), 301)
		return
	}

	tokenString := info.Info
	cl := jwt.MapClaims{}
	var email string
	_, err = jwt.ParseWithClaims(tokenString, cl, func(token *jwt.Token) (interface{}, error) {
		return []byte(""), nil
	})
	for key, val := range cl {
		if key == "email" {
			email = fmt.Sprintf("%v", val)
			break
		}
	}

	auth := types.RegisterOrAuthSocial{
		Email: email,
	}
	token, err := h.srv.RegisterOrAuthTeacherBySocial(&auth)
	if err != nil {
		http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?error=%s", err), 301)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("bashdikt://redirect?jwt=%s", token.Token), 301)
}
