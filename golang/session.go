package main

import (
	"net/http"

	"github.com/gorilla/sessions"
)

// Store the cookie store which is going to store session data in the cookie
var Store = sessions.NewCookieStore([]byte("dfrtyu7654erdsw3213456789iuytre4"))

func init() {
	Store.Options = &sessions.Options{
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
}

// IsLoggedIn will check if the user has an active session and return True
func IsLoggedIn(r *http.Request) bool {
	session, err := Store.Get(r, "sess")
	if err != nil || session.Values["loggedin"] == "false" {
		return false
	}
	return true
}

func CurrentUser(r *http.Request) (int, string, error) {
	sess, err := Store.Get(r, "sess")
	if err != nil {
		return 0, "", err
	}
	uid := sess.Values["userid"].(int)
	login := sess.Values["username"].(string)
	return uid, login, nil
}

func HasBaseAccess(accessMask uint64, r *http.Request) bool {
	sess, err := Store.Get(r, "sess")
	if err != nil {
		return false
	}
	uid, ok := sess.Values["userid"].(int)
	if !ok {
		return false
	}
	userAccess := UserBaseAccess(uid)
	if userAccess&accessMask != 0 {
		return true
	}
	if accessMask == DOC_READ {
		return HasBaseAccess(DOC_OWNREAD, r)
	}
	if accessMask == DOC_UPDATE {
		return HasBaseAccess(DOC_OWNUPDATE, r)
	}
	// if accessMask == WS_CONNECT {
	// 	return true
	// }
	return false
}

func HasAddAccess(accessMask uint64, r *http.Request) bool {
	sess, err := Store.Get(r, "sess")
	if err != nil {
		return false
	}
	uid := sess.Values["userid"].(int)
	userAccess := UserAddAccess(uid)
	return userAccess&accessMask != 0
}
