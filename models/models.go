package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type User struct {
	Id       int
	UserName string     `orm:"unique"`
	Passwd   string     `orm:"size(20)"`
	Articles []*Article `orm:"rel(m2m)"` //设置多对多关系
}

type Article struct {
	Id          int          `orm:"pk;auto"`
	Title       string       `orm:"size(20)"`                         //标题
	Content     string       `orm:"size(500)"`                        //内容
	Img         string       `orm:"size(50);null"`                    //图片路径
	Time        time.Time    `orm:"type(datatime);auto_now_add"`      //发布事件
	Count       int          `orm:"default(0)"`                       //阅读量
	ArticleType *ArticleType `orm:"rel(fk);null;on_delete(set_null)"` //设置一对多关系
	Users       []*User      `orm:"reverse(many)"`                    //设置多对多的反向关系
}

type ArticleType struct {
	Id       int
	TypeName string     `orm:"size(20)"`
	Articles []*Article `orm:"reverse(many)"` //设置一对多的反向关系
}

func init() {
	orm.RegisterDataBase("default", "mysql", "root:root@tcp(127.0.0.1:3306)/itcast?charset=utf8")
	orm.RegisterModel(new(User), new(Article), new(ArticleType))
	//orm.RunSyncdb("default", false, true)
}
