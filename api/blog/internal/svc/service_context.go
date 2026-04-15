// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"time"
	"zero-admin/api/blog/internal/config"
	"zero-admin/api/blog/internal/middleware"
	"zero-admin/api/blog/internal/models"
	"zero-admin/pkg/orm"

	captchaUtil "zero-admin/pkg/captcha"

	"github.com/mojocn/base64Captcha"
	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/syncx"
	"github.com/zeromicro/go-zero/rest"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config        config.Config
	JwtExpireAuth rest.Middleware

	// 缓存
	LocalCache *collection.Cache
	Redis      *redis.Redis
	DB         *gorm.DB

	Barrier syncx.SingleFlight
	Captcha *base64Captcha.Captcha

	ArticleModel         *models.ArticleModel
	ArticleCategoryModel *models.ArticleCategoryModel
	ArticleCommentModel  *models.ArticleCommentModel
	UserModel            *models.UserModel
	NoticeModel          *models.NoticeModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	db := orm.MustNewMysql(&c.Mysql)
	// cache
	redisCli := redis.MustNewRedis(c.Redis)
	localCache, err := collection.NewCache(time.Minute*30, collection.WithName("cache"))
	if err != nil {
		logx.Must(err)
	}

	// captcha
	captcha := captchaUtil.NewCaptchaDriverWithStore(captchaUtil.String, captchaUtil.NewRedisStore(redisCli, captchaUtil.RedisPrefix, captchaUtil.RedisExpire))

	// singleFlight
	barrier := syncx.NewSingleFlight()
	return &ServiceContext{
		Config:        c,
		JwtExpireAuth: middleware.NewJwtExpireAuthMiddleware(redisCli, localCache).Handle,

		LocalCache: localCache,
		Redis:      redisCli,
		DB:         db,
		Barrier:    barrier,
		Captcha:    captcha,

		ArticleModel:         models.NewArticleModel(db),
		ArticleCategoryModel: models.NewArticleCategoryModel(db),
		ArticleCommentModel:  models.NewArticleCommentModel(db),
		UserModel:            models.NewUserModel(db),
		NoticeModel:          models.NewNoticeModel(db),
	}
}
