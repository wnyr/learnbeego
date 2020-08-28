package controllers

import (
	"beego/models"
	"bytes"
	"encoding/gob"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/gomodule/redigo/redis"
	"math"
	"path"
	"strconv"
	"time"
)

type ArticleController struct {
	beego.Controller
}

//文章列表页
func (this *ArticleController) ShowArticleList() {
	/*userName := this.GetSession("userName")
	if userName == nil {
		this.Redirect("/login", 302)
		return
	}*/
	o := orm.NewOrm()

	qs := o.QueryTable("Article")
	var articles []models.Article
	//qs.All(&articles)
	//beego.Info(articles[0])

	////分页start
	//起始页
	//pageIndex := 1
	pageIndex, err := this.GetInt("pageIndex")
	if err != nil {
		pageIndex = 1
	}

	//获取select数据
	typeName := this.GetString("select")
	var count int64
	/*count, err := qs.RelatedSel("ArticleType").Count()
	if err != nil {
		beego.Info("Count查询出错")
		return
	}*/
	if typeName == "" {
		count, err = qs.RelatedSel("ArticleType").Count()
		if err != nil {
			beego.Info("Count查询出错")
			return
		}
	} else {
		//存在下来框数据
		count, err = qs.RelatedSel("ArticleType").Filter("ArticleType__TypeName", typeName).Count()
		if err != nil {
			beego.Info("Count查询出错")
			return
		}
	}

	//每页数量
	pageSize := 2
	//获取总页数,向上取整
	pageCount := math.Ceil(float64(count) / float64(pageSize))
	start := pageSize * (pageIndex - 1)
	//qs.Limit(pageSize, start).RelatedSel("ArticleType").All(&articles)

	////增加下来框选择
	//获取select数据
	//typeName := this.GetString("select")

	//处理数据
	if typeName == "" {
		beego.Info("下拉框数据为空")
		qs.Limit(pageSize, start).RelatedSel("ArticleType").All(&articles)
	} else {
		//存在下来框数据
		qs.Limit(pageSize, start).RelatedSel("ArticleType").Filter("ArticleType__TypeName", typeName).All(&articles)
	}
	this.Data["typeName"] = typeName
	////

	FirstPage := false
	LastPage := false
	if pageIndex == 1 {
		FirstPage = true
	}
	if pageIndex >= int(pageCount) {
		LastPage = true
	}

	this.Data["count"] = count
	this.Data["pageCount"] = pageCount
	this.Data["pageIndex"] = pageIndex
	this.Data["FirstPage"] = FirstPage
	this.Data["LastPage"] = LastPage
	////分页end

	////获取类型start
	var types []models.ArticleType
	////redis
	conn, _ := redis.Dial("tcp", ":6379")
	rel, err := redis.Bytes(conn.Do("get", "types"))
	if err != nil {
		beego.Info("获取redis数据失败")
	}
	//反序列化
	dec := gob.NewDecoder(bytes.NewReader(rel))
	dec.Decode(&types)

	if len(types) == 0 {
		_, err = o.QueryTable("ArticleType").All(&types)
		if err != nil {
			beego.Info("列表页获取类型失败:", err)
		}

		//序列化
		var buffer bytes.Buffer
		enc := gob.NewEncoder(&buffer)
		enc.Encode(types)
		if err != nil {
			beego.Info("redis数据库连接错误")
			return
		}
		////添加到redis中
		_, err = conn.Do("set", "types", buffer.Bytes())
		if err != nil {
			beego.Info("types写入redis错误", err)
			return
		}
		beego.Info("从数据库中获取类型")

	}

	this.Data["types"] = types

	////获取类型end
	this.Data["articles"] = articles

	//欢迎xxx
	userName := this.GetSession("userName")
	this.Data["userName"] = userName

	this.Layout = "layout.html"
	this.TplName = "index.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Scripts"] = "indexScript.html"
}

func (this *ArticleController) HandleSelect() {
	typeName := this.GetString("select")
	beego.Info(typeName)
	//处理数据
	if typeName == "" {
		beego.Info("下拉框数据为空")
		return
	}
	//查询数据
	o := orm.NewOrm()
	var articles []models.Article
	o.QueryTable("Article").RelatedSel("ArticleType").Filter("ArticleType__TypeName", typeName).All(&articles)
	beego.Info(articles)
	this.Ctx.WriteString(typeName)
}

func (this *ArticleController) ShowAddArticle() {
	////获取文章类型s
	o := orm.NewOrm()
	var types []models.ArticleType
	_, err := o.QueryTable("ArticleType").All(&types)
	if err != nil {
		beego.Info("添加文章页获取类型错误:", err)
	}
	this.Data["types"] = types
	////获取文章类型end
	this.Layout = "layout.html"
	this.TplName = "add.html"
}
func (this *ArticleController) HandleAddArticle() {
	//1.拿数据
	articleName := this.GetString("articleName")
	typeName := this.GetString("select")
	content := this.GetString("content")

	if typeName == "" {
		beego.Info("下来框数据为空")
		return
	}

	beego.Info(articleName, content)

	//上传图片处理
	//获取图片
	f, h, err := this.GetFile("uploadname")
	defer f.Close()
	if err != nil {
		beego.Info("图片上传失败", err)
		return
	}
	//判断格式
	ext := path.Ext(h.Filename)
	beego.Info(h.Filename, ext)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		beego.Info("上传文件格式不对")
		return
	}
	//文件大小
	if h.Size > 3*1204*1204 {
		beego.Info("上传文件太大")
		return
	}
	//不不重名
	filename := time.Now().Format("2006-01-02-15-04-05")
	this.SaveToFile("uploadname", "./static/img/"+filename+ext)

	// 3.插入数据到数据库
	o := orm.NewOrm()

	article := models.Article{}
	article.Title = articleName
	article.Content = content
	article.Img = "/static/img/" + filename + ext

	//给Article对象赋值
	articleType := models.ArticleType{TypeName: typeName}
	err = o.Read(&articleType, "TypeName")
	if err != nil {
		beego.Info("获取类型错误:", err)
		return
	}
	article.ArticleType = &articleType

	_, err = o.Insert(&article)
	if err != nil {
		beego.Info("article插入数据错误:", err)
	}

	//返回试图
	this.Redirect("/Article/ShowArticle", 302)

}

// 查看详情

func (this *ArticleController) ShowContent() {
	id := this.GetString("id")
	beego.Info(id)
	o := orm.NewOrm()
	id2, _ := strconv.Atoi(id)
	article := models.Article{Id: id2}
	err := o.Read(&article)
	if err != nil {
		beego.Info("ShowContent查询数据为空")
		return
	}
	//更新阅读量
	article.Count += 1
	o.Update(&article) //没有指定更新那一页，会自己查

	//多对多插入读者
	//1.获取orm对象
	//article := models.Article{Id: id2}
	//2.获取多对多操作对象
	m2m := o.QueryM2M(&article, "Users")
	//3.获取插入对象
	name := this.GetSession("userName")
	user := models.User{UserName: name.(string)}
	o.Read(&user, "UserName")
	//4.多对多插入
	m2m.Add(&user)

	//多对多查询
	//查询那些用户浏览了这篇文章
	//方法一
	//o.LoadRelated(&article,"Users)
	//方法二
	var users []models.User
	o.QueryTable("User").Filter("Articles__Article__id", id2).Distinct().All(&users)
	this.Data["users"] = users

	this.Data["article"] = article
	this.Layout = "layout.html"
	this.TplName = "content.html"
}

//删除文章
func (this *ArticleController) ShowDeleteArticle() {
	id, _ := this.GetInt("id")
	o := orm.NewOrm()
	article := models.Article{Id: id}
	_, err := o.Delete(&article)
	if err != nil {
		beego.Info("文章删除失败:", err)
		return
	}
	this.Redirect("/Article/ShowArticle", 302)
}

//更新文章页面显示
func (this *ArticleController) ShowUpdateArticle() {
	//id, _ := this.GetInt("id")
	id := this.GetString("id")
	//判断
	if id == "" {
		beego.Info("连接错误")
		return
	}
	o := orm.NewOrm()
	id2, _ := strconv.Atoi(id)
	article := models.Article{Id: id2}
	err := o.Read(&article)
	if err != nil {
		beego.Info("更新文章查询错误:", err)
		return
	}
	this.Data["article"] = article
	this.Layout = "layout.html"
	this.TplName = "update.html"
}

//更新文章
func (this *ArticleController) UpdateContent() {
	id, _ := this.GetInt("id")
	articleName := this.GetString("articleName")
	content := this.GetString("content")

	//判断
	if articleName == "" || content == "" {
		beego.Info("更新文章标题或内容不能为空")
		return
	}

	//上传图片
	f, h, err := this.GetFile("uploadname")
	var filename string
	if f != nil {
		defer f.Close()
		if err != nil {
			beego.Info("更新文章上传图片错误:", err)
			return
		}
		ext := path.Ext(h.Filename)
		if ext != ".jpg" && ext != ".jpeg" && ext != "png" {
			beego.Info("更新上传图片格式不支持！")
			return
		}

		//文件大小
		if h.Size > 3*1204*1204 {
			beego.Info("上传文件太大")
			return
		}
		//不不重名
		filename = time.Now().Format("2006-01-02-15-04-05") + ext
		this.SaveToFile("uploadname", "./static/img/"+filename)
	}

	o := orm.NewOrm()
	article := models.Article{Id: id}
	//读取数据库
	//判断文章ID是否存在
	err = o.Read(&article)
	if err != nil {
		beego.Info("要更新的文章不存在")
		return
	}

	article.Title = articleName
	article.Content = content
	if f != nil {
		article.Img = "/static/img/" + filename
	}

	_, err = o.Update(&article)
	if err != nil {
		beego.Info("文章更新出错:", err)
		return
	}
	this.Redirect("/Article/ArticleContent/?id="+strconv.Itoa(id), 302)
}

//显示文章类型
func (this *ArticleController) ShowAddType() {
	var articleTypes []models.ArticleType
	o := orm.NewOrm()
	_, err := o.QueryTable("ArticleType").All(&articleTypes)
	if err != nil {
		beego.Info("文章类型查询失败:", err)
	}
	this.Data["types"] = articleTypes
	this.Layout = "layout.html"
	this.TplName = "addType.html"
}

//处理添加类型
func (this *ArticleController) HandleAddType() {
	typeName := this.GetString("typeName")
	if typeName == "" {
		beego.Info("获取数据错误")
		return
	}
	o := orm.NewOrm()
	articleType := models.ArticleType{}
	articleType.TypeName = typeName
	_, err := o.Insert(&articleType)
	if err != nil {
		beego.Info("类型写入错误:", err)
		return
	}
	this.Redirect("/Article/AddArticleType", 302)
}

//退出登陆
func (this *ArticleController) Logout() {
	this.DelSession("userName")
	this.Redirect("/Article/login", 302)
}
