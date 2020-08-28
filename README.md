### validator-zh

------

Chinese display of [go-playground](https://github.com/go-playground/validator), support mobile phone and idcard.

### how to use

------

```go
package main

import (
	"github.com/gin-gonic/gin"
	zh "github.com/luvinci/validate-zh"
)

type User struct {
	Username string `json:"username" validate:"required,min=5" label:"用户名"`
	Password string `json:"password" validate:"required,gte=6" label:"密码"`
	// Nickname可以为空；不为空时则不能与Username相同
	Nickname string `json:"nickname" validate:"nefield=Username" label:"昵称"`
	Phone    string `json:"phone" validate:"mobile" label:"手机号"`
	IdCard   string `json:"idcard" validate:"required,idcard" label:"身份证号码"`
}

func main() {
	r := gin.Default()
	r.POST("/", func(c *gin.Context) {
		var user User
		_ = c.ShouldBind(&user)
		errMsgs := zh.Validate(user)
		c.JSON(200, errMsgs)
	})
	r.Run(":8080")
}

/*请求参数：
{
	"username": "pd",
	"password": "123",
	"nickname": "pd",
	"phone": "111111",
	"idcard": "22222"
}
*/
/*输出结果：
[
    "用户名长度必须至少为5个字符",
    "密码长度必须至少为6个字符",
    "昵称不能等于Username",
    "手机号格式错误",
    "身份证号码格式错误"
]
*/
```

