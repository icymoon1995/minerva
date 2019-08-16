package logic

import "github.com/gorilla/sessions"

var CookieSession = sessions.NewCookieStore([]byte("cookie_secret"))
