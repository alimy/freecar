// Code generated by hertz generator.

package trip

import (
	"github.com/alimy/freecar/app/api/biz/router/common"
	"github.com/alimy/freecar/library/cor/consts"
	"github.com/cloudwego/hertz/pkg/app"
)

func rootMw() []app.HandlerFunc {
	return common.CommonMW()
}

func _adminMw() []app.HandlerFunc {
	return []app.HandlerFunc{
		common.PasetoAuth(consts.Admin),
	}
}

func _tripMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _get_lltripsMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _getsometripsMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _createtripMw() []app.HandlerFunc {
	return []app.HandlerFunc{
		common.PasetoAuth(consts.User),
	}
}

func _gettripMw() []app.HandlerFunc {
	return []app.HandlerFunc{
		common.PasetoAuth(consts.User),
	}
}

func _gettripsMw() []app.HandlerFunc {
	return []app.HandlerFunc{
		common.PasetoAuth(consts.User),
	}
}

func _updatetripMw() []app.HandlerFunc {
	return []app.HandlerFunc{
		common.PasetoAuth(consts.User),
	}
}

func _deletetripMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _getalltripsMw() []app.HandlerFunc {
	// your code...
	return nil
}
