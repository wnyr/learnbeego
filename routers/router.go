package routers

import (
	"beego/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
	beego.InsertFilter("/Article/*", beego.BeforeRouter, FilterFunc)
	beego.Router("/", &controllers.MainController{})
	beego.Router("/register", &controllers.RegController{}, "get:ShowReg;post:HandleReg")
	beego.Router("/login", &controllers.LoginController{}, "get:ShowLogin;post:HandleLogin")
	beego.Router("/Article/ShowArticle", &controllers.ArticleController{}, "get:ShowArticleList;post:HandleSelect")
	beego.Router("/Article/AddArticle", &controllers.ArticleController{}, "get:ShowAddArticle;post:HandleAddArticle")
	beego.Router("/Article/ArticleContent", &controllers.ArticleController{}, "get:ShowContent")
	beego.Router("/Article/DeleteContent", &controllers.ArticleController{}, "get:ShowDeleteArticle")
	beego.Router("/Article/UpdateArticle", &controllers.ArticleController{}, "get:ShowUpdateArticle;post:UpdateContent")
	//添加类型
	beego.Router("/Article/AddArticleType", &controllers.ArticleController{}, "get:ShowAddType;post:HandleAddType")
	//退出登陆
	beego.Router("/Article/Logout", &controllers.ArticleController{}, "get:Logout")

}

var FilterFunc = func(ctx *context.Context) {
	userName := ctx.Input.Session("userName")
	if userName == nil {
		ctx.Redirect(302, "/login")
	}
}
