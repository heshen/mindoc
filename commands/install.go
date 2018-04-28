package commands

import (
	"fmt"
	"os"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/lifei6671/mindoc/conf"
	"github.com/lifei6671/mindoc/models"
	"github.com/lifei6671/mindoc/utils"
)

//系统安装.
func Install() {

	fmt.Println("Initializing...")

	err := orm.RunSyncdb("default", false, true)
	if err == nil {
		initialization()
	} else {
		panic(err.Error())
		os.Exit(1)
	}
	fmt.Println("Install Successfully!")
	os.Exit(0)

}

func Version() {
	if len(os.Args) >= 2 && os.Args[1] == "version" {
		fmt.Println(conf.VERSION)
		os.Exit(0)
	}
}

//修改用户密码
func ModifyPassword(account,password string) {
	if len(os.Args) < 2 {
		fmt.Println("Parameter error.")
		os.Exit(1)
	}

	if account == "" {
		fmt.Println("Account cannot be empty.")
		os.Exit(1)
	}
	if password == "" {
		fmt.Println("Password cannot be empty.")
		os.Exit(1)
	}
	member,err := models.NewMember().FindByAccount(account)

	if err != nil {
		fmt.Println("Failed to change password:",err)
		os.Exit(1)
	}
	pwd,err := utils.PasswordHash(password)

	if err != nil {
		fmt.Println("Failed to change password:",err)
		os.Exit(1)
	}
	member.Password = pwd

	err = member.Update("password")
	if err != nil {
		fmt.Println("Failed to change password:",err)
		os.Exit(1)
	}
	fmt.Println("Successfully modified.")
	os.Exit(0)

}
//初始化数据
func initialization() {

	err := models.NewOption().Init()

	if err != nil {
		panic(err.Error())
		os.Exit(1)
	}

	member, err := models.NewMember().FindByFieldFirst("account", "admin")
	if err == orm.ErrNoRows {

		member.Account = "admin"
		member.Avatar = "/static/images/headimgurl.jpg"
		member.Password = "123456"
		member.AuthMethod = "local"
		member.Role = 0
		member.Email = "admin@iminho.me"

		if err := member.Add(); err != nil {
			panic("Member.Add => " + err.Error())
			os.Exit(0)
		}

		book := models.NewBook()

		book.MemberId = member.MemberId
		book.BookName = "MinDoc演示项目"
		book.Status = 0
		book.Description = "这是一个MinDoc演示项目，该项目是由系统初始化时自动创建。"
		book.CommentCount = 0
		book.PrivatelyOwned = 0
		book.CommentStatus = "closed"
		book.Identify = "mindoc"
		book.DocCount = 0
		book.CommentCount = 0
		book.Version = time.Now().Unix()
		book.Cover = conf.GetDefaultCover()
		book.Editor = "markdown"
		book.Theme = "default"

		if err := book.Insert(); err != nil {
			panic("Book.Insert => " + err.Error())
			os.Exit(0)
		}
	}

	group := models.NewMemberGroup()
	group.GroupId = 1
	group.GroupName = "管理员组"
	group.GroupNumber = 1
	group.CreateTime = time.Now()
	group.CreateAt = 1
	group.IsEnableDelete = false


}
