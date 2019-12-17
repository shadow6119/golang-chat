package main

import (
	"crypto/sha256"
	"database/sql"
	"log"
	"reflect"

	_ "github.com/go-sql-driver/mysql"
)

type Account struct {
	name     string
	password []byte
}

const (
	salt     string = "salt"
	accessDB string = "root:password@tcp(localhost:3306)/gosql"
)

func main() {

	// AddAccount("shadow", "shadow")
	AccountCollation("shadow", "shadow")
}

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

//AccountCollation は入力されたnameとpassが登録済みか
func AccountCollation(nameIn string, passIn string) bool {

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
