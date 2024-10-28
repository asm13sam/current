package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// withDeleted http://127.0.0.1:8080/product_get_all?all

// access constant
const (
	LOGIN  = 0
	LOGOUT = 1 << iota
	USER_READ
	USER_CREATE
	USER_UPDATE
	USER_DELETE
	CONTRAGENT_READ
	CONTRAGENT_CREATE
	CONTRAGENT_UPDATE
	CONTRAGENT_DELETE
	DOC_READ
	DOC_OWNREAD
	DOC_CREATE
	DOC_UPDATE
	DOC_OWNUPDATE
	DOC_DELETE
	CATALOG_READ
	CATALOG_CREATE
	CATALOG_UPDATE
	CATALOG_DELETE
	OWNER_READ
	OWNER_CREATE
	OWNER_UPDATE
	OWNER_DELETE
	WS_CONNECT
	ADMIN
)

var clients = make(map[*websocket.Conn]WsClient)
var broadcast = make(chan Message)

type WsClient struct {
	UserId   int    `json:"user_id"`
	Username string `json:"username"`
}

type Message struct {
	UserId   int    `json:"user_id"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

var wsUpgrader = websocket.Upgrader{
	// CheckOrigin: func(r *http.Request) bool {
	//  return true
	// },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func UpgradeWS(r Req) {
	// Upgrade the HTTP connection to a WebSocket connection.
	// var err error
	// var wsConn *websocket.Conn
	wsConn, err := wsUpgrader.Upgrade(r.W, r.R, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer wsConn.Close()
	userid, username, err := CurrentUser(r.R)
	if err != nil {
		fmt.Println(err)
		return
	}
	clients[wsConn] = WsClient{userid, username}
	// clients[wsConn] = r.UserId
	fmt.Println("Connected", userid)
	fmt.Println("Clients", clients)
	// Read messages from the client.

	for {
		var msg Message
		err := wsConn.ReadJSON(&msg)
		fmt.Println("'userid'", userid)
		if err != nil {
			fmt.Println("OnReadMessage>>", clients[wsConn], err)
			msg = Message{clients[wsConn].UserId, clients[wsConn].Username, "виходить з чату"}
			delete(clients, wsConn)
			broadcast <- msg
			return
		}
		fmt.Println(msg.Message)
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		fmt.Printf("%s >> %s\n", msg.Username, msg.Message)
		fmt.Println("OnWriteMessage>>", clients)
		for wsConn := range clients {
			if msg.UserId != clients[wsConn].UserId {
				err := wsConn.WriteJSON(msg)
				if err != nil {
					fmt.Println("OnWriteMessage>>", clients[wsConn], err)
					wsConn.Close()
					msg = Message{clients[wsConn].UserId, clients[wsConn].Username, "виходить з чату"}
					delete(clients, wsConn)
					broadcast <- msg
				}
			}
		}
	}
}

type Req struct {
	W           http.ResponseWriter
	R           *http.Request
	UserId      int
	IntParam    int
	StrParam    string
	Int2Param   int
	Str2Param   string
	WithDeleted bool
	DeletedOnly bool
}

type Result struct {
	Value interface{} `json:"value"`
	Error string      `json:"error"`
}

func (r Req) Respond(payload interface{}, err error) {
	r.W.Header().Set("Content-Type", "application/json")

	res := Result{payload, ""}
	code := http.StatusOK

	if err != nil {
		res = Result{nil, err.Error()}
		code = http.StatusInternalServerError
	}

	response, err := json.Marshal(res)

	if err != nil {
		response, _ = json.Marshal(Result{nil, "Uncorrect payload to JSON marshalling"})
		code = http.StatusInternalServerError
	}
	r.W.WriteHeader(code)
	r.W.Write(response)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func copyFile(src, dest string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	sourceinfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	err = os.Chmod(dest, sourceinfo.Mode())
	if err != nil {
		return err
	}

	return nil
}

func copyDir(src string, dest string) error {
	sourceinfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	err = os.MkdirAll(dest, sourceinfo.Mode())
	if err != nil {
		return err
	}

	directory, err := os.Open(src)
	if err != nil {
		return err
	}

	objects, err := directory.Readdir(-1)
	if err != nil {
		return err
	}

	for _, obj := range objects {
		srcFile := filepath.Join(src, obj.Name())
		destFile := filepath.Join(dest, obj.Name())
		if obj.IsDir() {
			// create sub-directories - recursively
			err = copyDir(srcFile, destFile)
			if err != nil {
				return err
			}
		} else {
			// perform copy
			err = copyFile(srcFile, destFile)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// func copyFiles(src, dst string) error {
//  bytesRead, err := ioutil.ReadFile(src)

//  if err != nil {
//      return err
//  }

//  err = ioutil.WriteFile(dst, bytesRead, 0777)

//  if err != nil {
//      return err
//  }
//  return nil
// }

func WrapAuth(handler func(req Req), access uint64) func(w http.ResponseWriter, r *http.Request) {

	h := func(w http.ResponseWriter, r *http.Request) {
		req := Req{R: r, W: w}

		if (!IsLoggedIn(r) || !HasBaseAccess(access, r)) && access != LOGIN {
			log.Print("try to do some without auth", access)
			req.Respond(nil, errors.New("access denied"))
			return
		}

		q := r.URL.Query()
		all := q.Get("all")
		req.WithDeleted = (all == "all")
		req.DeletedOnly = (all == "deleted")
		vars := mux.Vars(r)
		var err error
		num, ok := vars["id"]
		if ok {
			req.IntParam, err = strconv.Atoi(num)
			if err != nil {
				req.Respond(nil, errors.New("invalid integer parameter"))
				return
			}
		} else {
			req.IntParam = 0
		}

		fs, ok := vars["fs"]
		if ok {
			req.StrParam = fs

		} else {
			req.StrParam = ""
		}

		num2, ok := vars["id2"]
		if ok {
			req.Int2Param, err = strconv.Atoi(num2)
			if err != nil {
				req.Respond(nil, errors.New("invalid integer parameter"))
				return
			}
		} else {
			req.Int2Param = 0
		}

		fs2, ok := vars["fs2"]
		if ok {
			req.Str2Param = fs2

		} else {
			req.Str2Param = ""
		}

		handler(req)
	}
	return h
}

// // Send a message back to the client.
// func sendWsMessage(text string) {
//  err := wsConn.WriteMessage(websocket.TextMessage, []byte(text))
//  if err != nil {
//      fmt.Println(err)
//  }
// }

// User access functions for session
func UserPassword(login string) (int, string, error) {
	var pass string
	var id int
	var is_active bool
	var name string
	err := db.QueryRow("SELECT id, password, is_active, name FROM user WHERE login=?", login).Scan(&id, &pass, &is_active, &name)
	if err != nil {
		return id, pass, err
	}
	if !is_active {
		return id, "", nil
	}
	msg := Message{id, name, fmt.Sprintf("%s приєднався до чату", name)}
	broadcast <- msg

	return id, pass, nil
}

func UserBaseAccess(id int) uint64 {
	var a uint64
	err := db.QueryRow("SELECT base_access FROM user WHERE id=?", id).Scan(&a)
	if err != nil {
		a = 0
	}
	return a
}

func UserAddAccess(id int) uint64 {
	var a uint64
	err := db.QueryRow("SELECT add_access FROM user WHERE id=?", id).Scan(&a)
	if err != nil {
		a = 0
	}
	return a
}

// Users handlers for login and logout

func Login(r Req) {
	//fmt.Println("get request")
	sess, err := Store.Get(r.R, "sess")
	//fmt.Println("get sess cookie")
	if err != nil {
		log.Println("Error identifying session", err)
		r.Respond(nil, err)
		return
	}
	//fmt.Println("session identificated")
	usr, err := DecodeUser(r)
	if err != nil {
		log.Println("Invalid request payload", err)
		r.Respond(nil, err)
		return
	}
	//fmt.Println("user decoded")
	login := usr.Login
	pass := usr.Password
	//fmt.Println("login pass decomposition")
	bdid, bdpass, err := UserPassword(login)
	//fmt.Println("userPassword works")
	if err != nil {
		r.Respond(nil, err)
	}

	fmt.Printf("login: %s pass %s bdpass %s\n", login, pass, bdpass)
	if bdpass == pass && bdpass != "" {
		sess.Values["loggedin"] = "true"
		sess.Values["username"] = login
		sess.Values["userid"] = bdid
		sess.Save(r.R, r.W)
		log.Print("user ", login, " is authenticated")
		//fmt.Fprintf(w, "%s", login)

		//jstr, err := json.Marshal(map[string]string{"user": login})
		//if err != nil {
		//  fmt.Println(err)
		//}
		//fmt.Println(jstr)
		r.Respond(map[string]string{"user": login}, nil)
		return
	}
	log.Print("Invalid user " + login)
	r.Respond(nil, errors.New("невірний логін та/або пароль - спробуйте ще раз"))
}

func Logout(r Req) {
	sess, err := Store.Get(r.R, "sess")
	if err == nil {
		if sess.Values["loggedin"] != "false" {
			sess.Values["loggedin"] = "false"
			sess.Save(r.R, r.W)
			log.Print("User ", sess.Values["username"], " is logout.")
		}
		fmt.Fprintf(r.W, "200")
	} else {
		fmt.Fprintf(r.W, "500")
	}
}

func UploadFile(req Req) {
	o, err := OrderingGet(req.IntParam, nil)
	if err != nil {
		req.Respond(nil, err)
		return
	}

	c, err := ContragentGet(o.ContragentId, nil)
	if err != nil {
		req.Respond(nil, err)
		return
	}

	numberDir := fmt.Sprintf("%d", o.Id)

	req.R.ParseMultipartForm(10 << 30)
	file, handler, err := req.R.FormFile("botFile")
	if err != nil {
		req.Respond(nil, err)
		return
	}
	defer file.Close()
	ext := "upload-*" + filepath.Ext(handler.Filename)
	//ext := fmt.Sprintf("upload-*%s", filepath.Ext(handler.Filename))
	newPath := filepath.Join(Cfg.MaketsPath, c.DirName, numberDir, "pix")
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	// fmt.Printf("File Size: %+v\n", handler.Size)
	// fmt.Printf("MIME Header: %+v\n", handler.Header)

	tempFile, err := os.CreateTemp(newPath, ext)
	if err != nil {
		req.Respond(nil, err)
	}
	defer tempFile.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		req.Respond(nil, err)
	}
	tempFile.Write(fileBytes)
	req.Respond(map[string]string{"message": "Successfully Uploaded File\n"}, nil)
}

func CopyBase(r Req) {
	t := time.Now()
	createdAt := t.Format("2006-01-02T15-04-05")
	base_name := r.StrParam + "_" + createdAt
	newPath := filepath.Join(Cfg.BckpPath, base_name)
	err := copyFile(Cfg.DBFile, newPath)
	if err != nil {
		r.Respond(nil, err)
		return
	}
	r.Respond(map[string]string{"base_name": base_name}, nil)
}

func GetBackupBases(r Req) {
	directory, err := os.Open(Cfg.BckpPath)
	if err != nil {
		r.Respond(nil, err)
		return
	}

	objects, err := directory.ReadDir(-1)
	if err != nil {
		r.Respond(nil, err)
		return
	}
	res := []string{}
	for _, obj := range objects {
		res = append(res, obj.Name())
	}
	r.Respond(map[string][]string{"base_names": res}, nil)
}

func RestoreBaseFromBackup(r Req) {
	err := DbClose()
	if err != nil {
		r.Respond(nil, err)
		return
	}
	base_name := r.StrParam
	oldPath := filepath.Join(Cfg.BckpPath, base_name)
	err = copyFile(oldPath, Cfg.DBFile)
	if err != nil {
		r.Respond(nil, err)
		return
	}
	err = DBconnect(Cfg.DBFile)
	if err != nil {
		r.Respond(nil, err)
		return
	}
	r.Respond(map[string]string{"message": "Successfully Restore From Backup\n"}, nil)
}

func DeleteBackupBase(r Req) {
	base_name := r.StrParam
	path := filepath.Join(Cfg.BckpPath, base_name)
	err := os.Remove(path)
	if err != nil {
		r.Respond(nil, err)
		return
	}
	r.Respond(map[string]string{"message": "Successfully Remove Backup File\n"}, nil)
}

func CopyProject(r Req) {
	p, err := ProjectGet(r.IntParam, nil)
	if err != nil {
		r.Respond(nil, err)
		return
	}
	c, err := ContragentGet(p.ContragentId, nil)
	if err != nil {
		r.Respond(nil, err)
		return
	}
	var oldPath, newPath string
	if p.ProjectGroupId == 4 || p.ProjectGroupId == 5 {
		oldPath = filepath.Join(Cfg.OldMaketsPath, c.DirName, p.TypeDir, p.NumberDir)
		newPath = filepath.Join(Cfg.NewMaketsPath, c.DirName, p.TypeDir, p.NumberDir)
	} else {
		newPath = filepath.Join(Cfg.OldMaketsPath, c.DirName, p.TypeDir, p.NumberDir)
		oldPath = filepath.Join(Cfg.NewMaketsPath, c.DirName, p.TypeDir, p.NumberDir)
	}
	err = copyDir(oldPath, newPath)
	if err != nil {
		r.Respond(nil, err)
		return
	}
	err = os.RemoveAll(oldPath)
	if err != nil {
		r.Respond(nil, err)
		return
	}
	res := struct {
		new_path string
		old_path string
	}{
		new_path: newPath,
		old_path: oldPath,
	}
	r.Respond(res, nil)
}

func CreateProjectDirs(r Req) {

	o, err := OrderingGet(r.IntParam, nil)
	if err != nil {
		r.Respond(nil, err)
		return
	}

	c, err := ContragentGet(o.ContragentId, nil)
	if err != nil {
		r.Respond(nil, err)
		return
	}
	numberDir := fmt.Sprintf("%d", o.Id)
	newPath := filepath.Join(Cfg.MaketsPath, c.DirName, numberDir)
	err = os.MkdirAll(newPath, 0777)
	if err != nil {
		r.Respond(nil, err)
		return
	}
	for _, d := range Cfg.MaketDirs {
		p := filepath.Join(newPath, d)
		err := os.MkdirAll(p, 0777)
		if err != nil {
			r.Respond(nil, err)
			return
		}
	}
	maketFile := filepath.Join(Cfg.MaketsPath, "maket.cdr")
	newMaketFile := filepath.Join(newPath, fmt.Sprintf("maket_%d.cdr", o.Id))

	if fileExists(newMaketFile) {
		r.Respond(nil, errors.New("file exist"))
		return
	}
	err = copyFile(maketFile, newMaketFile)
	if err != nil {
		r.Respond(nil, err)
		return
	}

	res := struct {
		path string
		file string
	}{
		path: newPath,
		file: newMaketFile,
	}

	r.Respond(res, nil)
}

func CreateMatherialToWhsInToNumber(m *MatherialToWhsIn, tx *sql.Tx) error {
	var sql_reg string
	var wmc_number WmcNumber

	whs_in, err := WhsInGet(m.WhsInId, tx)
	if err != nil {
		return err
	}
	sql_reg = `SELECT * FROM wmc_number WHERE whs_id = ? AND matherial_id = ? AND color_id = ?;`
	row := tx.QueryRow(sql_reg, whs_in.WhsId, m.MatherialId, m.ColorId)

	err = row.Scan(
		&wmc_number.Id,
		&wmc_number.WhsId,
		&wmc_number.MatherialId,
		&wmc_number.ColorId,
		&wmc_number.Total,
		&wmc_number.IsActive,
	)
	if err != nil {
		wmc_number.Id = 0
		wmc_number.WhsId = whs_in.WhsId
		wmc_number.MatherialId = m.MatherialId
		wmc_number.ColorId = m.ColorId
		wmc_number.Total = m.Number
		wmc_number.IsActive = true
		wmc_number, err = WmcNumberCreate(wmc_number, tx)
		if err != nil {
			return err
		}
	} else {
		wmc_number.Total += m.Number
		wmc_number, err = WmcNumberUpdate(wmc_number, tx)
		if err != nil {
			return err
		}
	}
	return nil

}

func DeleteMatherialToWhsInToNumber(m *MatherialToWhsIn, tx *sql.Tx) error {
	var sql_reg string
	var wmc_number WmcNumber

	whs_in, err := WhsInGet(m.WhsInId, tx)
	if err != nil {
		return err
	}
	sql_reg = `SELECT * FROM wmc_number WHERE whs_id = ? AND matherial_id = ? AND color_id = ?;`
	row := tx.QueryRow(sql_reg, whs_in.WhsId, m.MatherialId, m.ColorId)

	err = row.Scan(
		&wmc_number.Id,
		&wmc_number.WhsId,
		&wmc_number.MatherialId,
		&wmc_number.ColorId,
		&wmc_number.Total,
		&wmc_number.IsActive,
	)
	if err != nil {
		return err
	} else {
		wmc_number.Total -= m.Number
		wmc_number, err = WmcNumberUpdate(wmc_number, tx)
		if err != nil {
			return err
		}
	}
	return nil

}

func UpdateMatherialToWhsInToNumber(m *MatherialToWhsIn, old_number float64, tx *sql.Tx) error {
	var sql_reg string
	var wmc_number WmcNumber

	whs_in, err := WhsInGet(m.WhsInId, tx)
	if err != nil {
		return err
	}
	sql_reg = `SELECT * FROM wmc_number WHERE whs_id = ? AND matherial_id = ? AND color_id = ?;`
	row := tx.QueryRow(sql_reg, whs_in.WhsId, m.MatherialId, m.ColorId)

	err = row.Scan(
		&wmc_number.Id,
		&wmc_number.WhsId,
		&wmc_number.MatherialId,
		&wmc_number.ColorId,
		&wmc_number.Total,
		&wmc_number.IsActive,
	)
	if err != nil {
		return err
	} else {
		wmc_number.Total = wmc_number.Total - old_number + m.Number
		wmc_number, err = WmcNumberUpdate(wmc_number, tx)
		if err != nil {
			return err
		}
	}
	return nil

}

func CreateMatherialToWhsOutToNumber(m *MatherialToWhsOut, tx *sql.Tx) error {
	var sql_reg string
	var wmc_number WmcNumber

	whs_out, err := WhsOutGet(m.WhsOutId, tx)
	if err != nil {
		return err
	}
	sql_reg = `SELECT * FROM wmc_number WHERE whs_id = ? AND matherial_id = ? AND color_id = ?;`
	row := tx.QueryRow(sql_reg, whs_out.WhsId, m.MatherialId, m.ColorId)

	err = row.Scan(
		&wmc_number.Id,
		&wmc_number.WhsId,
		&wmc_number.MatherialId,
		&wmc_number.ColorId,
		&wmc_number.Total,
		&wmc_number.IsActive,
	)
	if err != nil {
		wmc_number.Id = 0
		wmc_number.WhsId = whs_out.WhsId
		wmc_number.MatherialId = m.MatherialId
		wmc_number.ColorId = m.ColorId
		wmc_number.Total = -m.Number
		wmc_number.IsActive = true
		wmc_number, err = WmcNumberCreate(wmc_number, tx)
		if err != nil {
			return err
		}
	} else {
		wmc_number.Total -= m.Number
		wmc_number, err = WmcNumberUpdate(wmc_number, tx)
		if err != nil {
			return err
		}
	}
	return nil

}

func DeleteMatherialToWhsOutToNumber(m *MatherialToWhsOut, tx *sql.Tx) error {
	var sql_reg string
	var wmc_number WmcNumber

	whs_out, err := WhsOutGet(m.WhsOutId, tx)
	if err != nil {
		return err
	}
	sql_reg = `SELECT * FROM wmc_number WHERE whs_id = ? AND matherial_id = ? AND color_id = ?;`
	row := tx.QueryRow(sql_reg, whs_out.WhsId, m.MatherialId, m.ColorId)

	err = row.Scan(
		&wmc_number.Id,
		&wmc_number.WhsId,
		&wmc_number.MatherialId,
		&wmc_number.ColorId,
		&wmc_number.Total,
		&wmc_number.IsActive,
	)
	if err != nil {
		return err
	} else {
		wmc_number.Total += m.Number
		wmc_number, err = WmcNumberUpdate(wmc_number, tx)
		if err != nil {
			return err
		}
	}
	return nil

}

func UpdateMatherialToWhsOutToNumber(m *MatherialToWhsOut, old_number float64, tx *sql.Tx) error {
	var sql_reg string
	var wmc_number WmcNumber

	whs_out, err := WhsOutGet(m.WhsOutId, tx)
	if err != nil {
		return err
	}
	sql_reg = `SELECT * FROM wmc_number WHERE whs_id = ? AND matherial_id = ? AND color_id = ?;`
	row := tx.QueryRow(sql_reg, whs_out.WhsId, m.MatherialId, m.ColorId)

	err = row.Scan(
		&wmc_number.Id,
		&wmc_number.WhsId,
		&wmc_number.MatherialId,
		&wmc_number.ColorId,
		&wmc_number.Total,
		&wmc_number.IsActive,
	)
	if err != nil {
		return err
	} else {
		wmc_number.Total = wmc_number.Total + old_number - m.Number
		wmc_number, err = WmcNumberUpdate(wmc_number, tx)
		if err != nil {
			return err
		}
	}
	return nil

}

func UpdateCost(ordering *Ordering) {
	ordering.Cost = ordering.Price * (1 + ordering.Persent/100)
}
