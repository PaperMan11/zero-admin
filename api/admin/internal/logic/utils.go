package logic

import (
	"context"
	"zero-admin/pkg/convert"
)

func GetOperateID(ctx context.Context) int64 {
	return convert.ToInt64(ctx.Value("uid"))
}
