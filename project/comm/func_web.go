package comm

import (
	"crypto/md5"
	"fmt"
	"log"
	"lottery/conf"
	"lottery/models"
	"net"
	"net/http"
	"net/url"
	"strconv"
)

func ClientIp(request *http.Request) string {
	host, _, _ := net.SplitHostPort(request.RemoteAddr)
	return host
}

func Redirect(writer http.ResponseWriter, url string) {
	writer.Header().Add("location", url)
	writer.WriteHeader(http.StatusFound)
}

func GetLoginUser(r *http.Request) *models.ObjLoginuser {
	c, err := r.Cookie("lottery_loginUser")
	if err != nil {
		return nil
	}
	params, err := url.ParseQuery(c.Value)
	if err != nil {
		return nil
	}
	uid, err := strconv.Atoi(params.Get("uid"))
	if err != nil || uid < 1 {
		return nil
	}
	now, err := strconv.Atoi(params.Get("now"))
	if err != nil || NowUnix()-now > 86400*30 {
		return nil
	}
	loginUser := &models.ObjLoginuser{
		Uid:      uid,
		Username: params.Get("username"),
		Now:      now,
		Ip:       ClientIp(r),
		Sign:     params.Get("sign"),
	}
	sign := createLoginUserSign(loginUser)
	if sign != loginUser.Sign {
		log.Println("func_web GetLoginUser createLoginUserSign not signed", sign, loginUser.Sign)
		return nil
	}

	return loginUser
}

func SetLoginUser(w http.ResponseWriter, loginUser *models.ObjLoginuser) {
	if loginUser == nil || loginUser.Uid < 1 {
		c := &http.Cookie{
			Name:   "lottery_loginUser",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		}
		http.SetCookie(w, c)
		return
	}
	if loginUser.Sign == "" {
		loginUser.Sign = createLoginUserSign(loginUser)
	}
	params := url.Values{}
	params.Add("uid", strconv.Itoa(loginUser.Uid))
	params.Add("username", loginUser.Username)
	params.Add("now", strconv.Itoa(loginUser.Now))
	params.Add("ip", loginUser.Ip)
	params.Add("sign", loginUser.Sign)
	c := &http.Cookie{
		Name:  "lottery_loginUser",
		Value: params.Encode(),
		Path:  "/",
	}
	http.SetCookie(w, c)
}

func createLoginUserSign(loginUser *models.ObjLoginuser) string {
	str := fmt.Sprintf("uid=%d&username=%s&secret=%s&now=%d", loginUser.Uid, loginUser.Username, conf.CookieSecret, loginUser.Now)
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}
