package menu

import (
	"github.com/xinliangnote/go-gin-api/skitii/pkg/core"
	"github.com/xinliangnote/go-gin-api/skitii/repository/mysql"
	"github.com/xinliangnote/go-gin-api/skitii/repository/mysql/menu"
)

type UpdateMenuData struct {
	Name string // 菜单名称
	Link string // 链接地址
	Icon string // 图标
}

func (s *service) Modify(ctx core.Context, id int32, menuData *UpdateMenuData) (err error) {
	data := map[string]interface{}{
		"name":         menuData.Name,
		"link":         menuData.Link,
		"icon":         menuData.Icon,
		"updated_user": ctx.SessionUserInfo().UserName,
	}

	qb := menu.NewQueryBuilder()
	qb.WhereId(mysql.EqualPredicate, id)
	err = qb.Updates(s.db.GetDbW().WithContext(ctx.RequestContext()), data)
	if err != nil {
		return err
	}

	return
}
