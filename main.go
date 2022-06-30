package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	"github.com/Tunakank/go-todo/my"
	"github.com/gorilla/sessions"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// DB関連
var dbDriver = "sqlite3"
var dbName = "data.sqlite3"

// Cookie関連
var sesName = "tunakank-go-todo-session"
var cs = sessions.NewCookieStore([]byte("kokoniha-sukina-moziwo-ireru"))

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		index(w, req)
	})
	http.HandleFunc("/home", func(w http.ResponseWriter, req *http.Request) {
		home(w, req)
	})
	http.HandleFunc("/detail", func(w http.ResponseWriter, req *http.Request) {
		detail(w, req)
	})
	http.HandleFunc("/list", func(w http.ResponseWriter, req *http.Request) {
		list(w, req)
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, req *http.Request) {
		login(w, req)
	})
	http.HandleFunc("/logout", func(w http.ResponseWriter, req *http.Request) {
		logout(w, req)
	})
	http.HandleFunc("/new-account", func(w http.ResponseWriter, req *http.Request) {
		newAccount(w, req)
	})
	http.HandleFunc("/edit", func(w http.ResponseWriter, req *http.Request) {
		editMemo(w, req)
	})
	http.HandleFunc("/delete", func(w http.ResponseWriter, req *http.Request) {
		deleteMemo(w, req)
	})

	http.ListenAndServe("", nil)
}

// ログインされているかのチェック
func checkLogin(w http.ResponseWriter, req *http.Request) *my.User {
	ses, _ := cs.Get(req, sesName)
	if ses.Values["login"] == nil || !ses.Values["login"].(bool) {
		// 未ログインなので`/login`へリダイレクト
		http.Redirect(w, req, "/login", 302)
	}

	ac := ""
	if ses.Values["account"] != nil {
		ac = ses.Values["account"].(string)
	}

	var user my.User
	db, _ := gorm.Open(dbDriver, dbName)
	defer db.Close()

	db.Where("account = ?", ac).First(&user)

	return &user
}

func notemp() *template.Template {
	tmp, _ := template.New("index").Parse("NO PAGE.")
	return tmp
}

func page(fname string) *template.Template {
	tmps, _ := template.ParseFiles("templates/"+fname+".html", "templates/head.html", "templates/foot.html")
	return tmps
}

func index(w http.ResponseWriter, req *http.Request) {
	user := checkLogin(w, req)

	db, _ := gorm.Open(dbDriver, dbName)
	defer db.Close()

	var pl []my.Memo
	db.Where("user_id = ?", user.Model.ID).Order("created_at desc").Limit(10).Find(&pl)
	var gl []my.List
	db.Where("user_id = ?", user.Model.ID).Order("created_at desc").Limit(10).Find(&gl)

	res := struct {
		Title   string
		Message string
		Name    string
		Account string
		Plist   []my.Memo
		Glist   []my.List
	}{
		Title:   "Top Page",
		Message: user.Name + "さん、ようこそ。",
		Name:    user.Name,
		Account: user.Account,
		Plist:   pl,
		Glist:   gl,
	}

	er := page("index").Execute(w, res)
	if er != nil {
		log.Fatal(er)
	}
}

func detail(w http.ResponseWriter, req *http.Request) {
	user := checkLogin(w, req)

	mid := req.FormValue("mid")
	db, _ := gorm.Open(dbDriver, dbName)
	defer db.Close()

	if req.Method == "POST" {
		msg := req.PostFormValue("message")
		mId, _ := strconv.Atoi(mid)
		cmt := my.Comment{
			UserId:  int(user.Model.ID),
			MemoId:  mId,
			Message: msg,
		}
		db.Create(&cmt)
	}

	var pst my.Memo
	var cmts []my.CommentJoin

	db.Where("id = ?", mid).First(&pst)
	db.Table("comments").Select("comments.*, users.id, users.name").Joins("join users on users.id =comments.user_id").Where("comments.memo_id = ?", mid).Order("created_at desc").Find(&cmts)

	res := struct {
		Title   string
		Message string
		Name    string
		Account string
		Memo    my.Memo
		Clist   []my.CommentJoin
	}{
		Title:   "動画",
		Message: "",
		Name:    user.Name,
		Account: user.Account,
		Memo:    pst,
		Clist:   cmts,
	}
	er := page("post").Execute(w, res)
	if er != nil {
		log.Fatal(er)
	}
}

func home(w http.ResponseWriter, req *http.Request) {
	user := checkLogin(w, req)

	db, _ := gorm.Open(dbDriver, dbName)
	defer db.Close()

	if req.Method == "POST" {
		switch req.PostFormValue("form") {
		case "memo":
			ad := req.PostFormValue("address")
			ti := req.PostFormValue("title")
			ad = strings.TrimSpace(ad)
			if strings.HasPrefix(ad, "https://youtu.be/") {
				ad = strings.TrimPrefix(ad, "https://youtu.be/")
			}

			pt := my.Memo{
				UserId:  int(user.Model.ID),
				Address: ad,
				Title:   ti,
				Message: req.PostFormValue("message"),
			}
			db.Create(&pt)
		case "list":
			gp := my.List{
				UserId:  int(user.Model.ID),
				Name:    req.PostFormValue("name"),
				Message: req.PostFormValue("message"),
			}
			db.Create(&gp)
		}
	}

	var pts []my.Memo
	var gps []my.List

	db.Where("user_id=?", user.ID).Order("created_at desc").Limit(10).Find(&pts)
	db.Where("user_id=?", user.ID).Order("created_at desc").Limit(10).Find(&gps)

	itm := struct {
		Title   string
		Message string
		Name    string
		Account string
		Plist   []my.Memo
		Glist   []my.List
	}{
		Title:   "マイページ",
		Message: "メールアドレス=\"" + user.Account + "\"<br>ニックネーム=\"" + user.Name + "\"",
		Name:    user.Name,
		Account: user.Account,
		Plist:   pts,
		Glist:   gps,
	}
	er := page("home").Execute(w, itm)
	if er != nil {
		log.Fatal(er)
	}
}

func list(w http.ResponseWriter, req *http.Request) {
	user := checkLogin(w, req)

	gid := req.FormValue("gid")
	db, _ := gorm.Open(dbDriver, dbName)
	defer db.Close()

	if req.Method == "POST" {
		ti := req.PostFormValue("title")
		ad := req.PostFormValue("address")
		ad = strings.TrimSpace(ad)
		if strings.HasPrefix(ad, "https://youtu.be/") {
			ad = strings.TrimPrefix(ad, "https://youtu.be/")
		}
		lId, _ := strconv.Atoi(gid)
		pt := my.Memo{
			UserId:  int(user.Model.ID),
			Address: ad,
			Title:   ti,
			Message: req.PostFormValue("message"),
			ListId:  lId,
		}
		db.Create(&pt)
	}

	var lrp my.List
	var pts []my.Memo

	db.Where("id=?", gid).First(&lrp)
	db.Order("created_at desc").Model(&lrp).Related(&pts)

	itm := struct {
		Title   string
		Message string
		Name    string
		Account string
		List    my.List
		Plist   []my.Memo
	}{
		Title:   "リスト",
		Message: "",
		Name:    user.Name,
		Account: user.Account,
		List:    lrp,
		Plist:   pts,
	}
	er := page("list").Execute(w, itm)
	if er != nil {
		log.Fatal(er)
	}
}

// ログイン処理
func login(w http.ResponseWriter, req *http.Request) {
	// 戻り値
	res := struct {
		Title   string
		Message string
		Account string
	}{
		Title:   "Login",
		Message: "ログインしてください",
		Account: "",
	}

	if req.Method == "GET" {
		er := page("login").Execute(w, res)
		if er != nil {
			log.Fatal(er)
		}
		return
	}

	if req.Method == "POST" {
		db, _ := gorm.Open(dbDriver, dbName)
		defer db.Close()

		usr := req.PostFormValue("account")
		pass := req.PostFormValue("pass")
		res.Account = usr

		var re int
		var user my.User

		db.Where("account = ? and password = ?", usr, pass).Find(&user).Count(&re)

		if re <= 0 {
			// ログイン失敗
			res.Message = "Wrong account or password."
			page("login").Execute(w, res)
			return
		}

		// ログイン成功
		ses, _ := cs.Get(req, sesName)
		ses.Values["login"] = true
		ses.Values["account"] = usr
		ses.Values["name"] = user.Name
		ses.Save(req, w)
		http.Redirect(w, req, "/", 302)
	}

	er := page("login").Execute(w, res)
	if er != nil {
		log.Fatal(er)
	}
}

// セッション削除
func resetSession(w http.ResponseWriter, req *http.Request) {
	ses, _ := cs.Get(req, sesName)
	ses.Values["login"] = nil
	ses.Values["account"] = nil
	ses.Save(req, w)
}

// ログアウト
func logout(w http.ResponseWriter, req *http.Request) {
	resetSession(w, req)
	http.Redirect(w, req, "/login", 302)
}

// アカウント作成
func newAccount(w http.ResponseWriter, req *http.Request) {
	resetSession(w, req)

	// 戻り値
	res := struct {
		Title   string
		Message string
	}{
		Title:   "New Account",
		Message: "新規アカウント作成",
	}

	if req.Method == "GET" {
		er := page("newAccount").Execute(w, res)
		if er != nil {
			log.Fatal(er)
		}
		return
	}

	if req.Method == "POST" {
		db, _ := gorm.Open(dbDriver, dbName)
		defer db.Close()

		pUsr := req.PostFormValue("account")
		pNickname := req.PostFormValue("nickname")
		pPass := req.PostFormValue("pass")

		usr := my.User{
			Account:  pUsr,
			Name:     pNickname,
			Password: pPass,
			Message:  pNickname + "'s account.",
		}
		db.Create(&usr)

		ses, _ := cs.Get(req, sesName)
		ses.Values["login"] = true
		ses.Values["account"] = pUsr
		ses.Values["name"] = pNickname
		ses.Save(req, w)
		http.Redirect(w, req, "/", 302)
	}
}

func editMemo(w http.ResponseWriter, req *http.Request) {
	user := checkLogin(w, req)

	db, _ := gorm.Open(dbDriver, dbName)
	defer db.Close()

	mid := req.FormValue("mid")

	var mm my.Memo
	db.Where("id = ?", mid).Find(&mm)

	if req.Method == "POST" {
		mm.Title = req.PostFormValue("title")
		mm.Message = req.PostFormValue("message")
		db.Save(&mm)
		http.Redirect(w, req, "/detail?mid="+mid, 302)
	}

	itm := struct {
		Title   string
		Message string
		Name    string
		Account string
		Memo    my.Memo
	}{
		Title:   "編集",
		Message: "",
		Name:    user.Name,
		Account: user.Account,
		Memo:    mm,
	}
	er := page("edit").Execute(w, itm)
	if er != nil {
		log.Fatal(er)
	}
}

func deleteMemo(w http.ResponseWriter, req *http.Request) {
	checkLogin(w, req)

	db, _ := gorm.Open(dbDriver, dbName)
	defer db.Close()

	mid := req.FormValue("mid")

	db.Delete(&my.Memo{}, mid)
	http.Redirect(w, req, "/", 302)
}
