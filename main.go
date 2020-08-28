package main

import (
	_ "beego/routers"
	"github.com/astaxie/beego"
)

func main() {
	beego.AddFuncMap("ShowPrePage", HandlePrePage)
	beego.AddFuncMap("ShowNextPage", HandleNextPage)
	beego.Run()
}
func HandlePrePage(data int) int {
	return data - 1
}
func HandleNextPage(data int) int {
	return data + 1
}
