package chatDB

import (
	"crypto/sha256"
	"database/sql"
	"log"
	"reflect"

	_ "github.com/go-sql-driver/mysql"
)

//Account DB上ではアカウントのnameとハッシュ後のパスワードを持つ
type Account struct {
	name     string
	password []byte
}

//Chat チャットの構造体
type Chat struct {
	message string //メッセージ内容
	name    string //送り主
	date    string //年(西暦):月:日:時:分:秒 年以外のすべてを二桁(1は01など)とする
}

const (
	salt     string = "salt"
	accessDB string = "root:password@tcp(localhost:3306)/gosql"
)

//SaveChat chat内容をDBに保存
func SaveChat(chat Chat) {

	//使用するDBに接続
	db, err := sql.Open("mysql", accessDB)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	ins, err := db.Prepare("INSERT INTO Chat(message,name,date) VALUES(?,?,?)")
	if err != nil {
		log.Println(err)
	}
	defer ins.Close()

	ins.Exec(chat.message, chat.name, chat.date)
}

//AddAccount 入力したアカウントを登録
func AddAccount(nameIn string, passIn string) {

	//入力パスワードのハッシュ化
	hash := sha256.Sum256([]byte(salt + passIn))

	//使用するDBに接続
	db, err := sql.Open("mysql", accessDB)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	ins, err := db.Prepare("INSERT INTO Account(name,pass) VALUES(?,?)")
	if err != nil {
		log.Println(err)
	}
	defer ins.Close()

	ins.Exec(nameIn, byteToStr(hash))

}

//CollationAccount は入力されたnameとpassが登録済みか
func CollationAccount(nameIn string, passIn string) bool {

	//使用するDBに接続
	db, err := sql.Open("mysql", accessDB)

	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	//アカウントのパスワードを取得
	sel, err := db.Prepare("SELECT pass FROM Account WHERE name=?")
	if err != nil {
		log.Println(err)
	}
	defer sel.Close()

	var pass string

	err = sel.QueryRow(nameIn).Scan(&pass)
	if err != nil {
		panic(err)
	}

	//入力パスワードのハッシュ化
	hash := sha256.Sum256([]byte(salt + passIn))

	//入力パスワードと比較
	if reflect.DeepEqual(byteToStr(hash), pass) {
		log.Println("Login Success")
		return true
	} else {
		log.Println("Login Faild")
		return false
	}
}

func byteToStr(hash [32]byte) string {

	var hexStr string

	for _, hex := range hash {
		hexStr = hexStr + string(hex)
	}
	return hexStr
}
