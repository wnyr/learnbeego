package controllers

import (
	"beego/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type RegController struct {
	beego.Controller
}

func (this *RegController) ShowReg() {
	this.TplName = "register.html"
}

//注册业务
/*
1.接收浏览器传进的数据
2.数据处理
3.插入数据库(数据库表User)
4.返回视图
*/
func (this *RegController) HandleReg() {
	//1.接收浏览器传进的数据
	name := this.GetString("userName")
	passwd := this.GetString("password")
	//2.数据处理
	if name == "" || passwd == "" {
		beego.Info("用户名或者密码不能为空")
		this.TplName = "register.html"
		return
	}
	//3.插入数据库(数据库表User)
	o := orm.NewOrm()
	user := models.User{}
	user.UserName = name
	user.Passwd = passwd
	id, err := o.Insert(&user)
	if err != nil {
		beego.Info("user插入失败")
		return
	}
	beego.Info(name, "插入成功，ID:", id)

	//4.返回登陆
	//this.Ctx.WriteString("注册成功")
	this.Redirect("/login", 302)

}

//登陆
type LoginController struct {
	beego.Controller
}

func (this *LoginController) ShowLogin() {
	name := this.Ctx.GetCookie("userName")
	if name != "" {
		this.Data["name"] = name
		this.Data["check"] = "checked"
	}
	this.TplName = "login.html"
}

//登陆业务
/*
1.拿到浏览器数据
2.数据处理
3.查找数据库
4.返回视图
*/
func (this *LoginController) HandleLogin() {
	//1.拿到浏览器数据
	name := this.GetString("userName")
	passwd := this.GetString("password")
	beego.Info(name, passwd)
	//2.数据处理
	if name == "" || passwd == "" {
		beego.Info("用户名或者密码不能为空")
		this.TplName = "login.html"
		return
	}
	//3.查找数据库
	o := orm.NewOrm()
	user := models.User{}
	user.UserName = name
	//user.Passwd=passwd

	err := o.Read(&user, "Username")
	if err != nil {
		beego.Info("用户名查询失败")
		this.TplName = "login.html"
		return
	}
	//判断密码
	if user.Passwd != passwd {
		beego.Info("密码错误")
		this.TplName = "login.html"
		return
	}

	//记住用户名
	check := this.GetString("remember")
	if check == "on" {
		this.Ctx.SetCookie("userName", name, time.Second*3600)
	} else {
		this.Ctx.SetCookie("userName", "", -1)
	}

	//登陆session
	this.SetSession("userName", name)
	//4.返回视图
	//this.Ctx.WriteString("登陆成功")
	this.Redirect("/Article/ShowArticle", 302)
}
