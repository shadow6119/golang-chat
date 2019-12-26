package main

import (
	// "./chatDB"

	"fmt"
	"net/http"
	// "log"
	// "os"
	// "io/ioutil"
)

func main() {

	// var message []string

	//チャット受信時の処理
	http.HandleFunc("/sendChat", receiveChat)
	//チャットの更新
	http.HandleFunc("/updateChat", updateChat)

	http.ListenAndServe(":8080", nil)
}

//handler関数を設定
func receiveChat(w http.ResponseWriter, r *http.Request) {

	//ファイル書き込み URL.RawQueryはパラメータ
	fmt.Fprintln(w, r.URL.RawQuery)

}

func updateChat(w http.ResponseWriter, r *http.Request) {

	//ファイル書き込み URL.RawQueryはパラメータ
	fmt.Fprintln(w, r.URL.RawQuery)

}
