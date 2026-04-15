// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package home

import (
	"context"
	"strings"
	"zero-admin/api/blog/internal/svc"
	"zero-admin/api/blog/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchSuggestLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchSuggestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchSuggestLogic {
	return &SearchSuggestLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchSuggestLogic) SearchSuggest(req *types.SearchSuggestReq) (resp *types.SearchSuggestResp, err error) {
	if req.Keyword == "" {
		return &types.SearchSuggestResp{
			Keywords: []string{},
		}, nil
	}

	limit := int(req.Limit)
	if limit <= 0 {
		limit = 10
	}

	articles, err := l.svcCtx.ArticleModel.SearchArticlesByKeyword(l.ctx, req.Keyword, limit)
	if err != nil {
		l.Errorf("搜索联想失败: %v", err)
		return nil, err
	}

	keywordsMap := make(map[string]bool)
	var keywords []string

	for _, article := range articles {
		if strings.Contains(strings.ToLower(article.Title), strings.ToLower(req.Keyword)) {
			if !keywordsMap[article.Title] {
				keywordsMap[article.Title] = true
				keywords = append(keywords, article.Title)
			}
		}

		if article.Tags != "" {
			tags := strings.Split(article.Tags, ",")
			for _, tag := range tags {
				tag = strings.TrimSpace(tag)
				if tag != "" && strings.Contains(strings.ToLower(tag), strings.ToLower(req.Keyword)) {
					if !keywordsMap[tag] {
						keywordsMap[tag] = true
						keywords = append(keywords, tag)
					}
				}
			}
		}

		if len(keywords) >= limit {
			break
		}
	}

	if keywords == nil {
		keywords = []string{}
	}

	return &types.SearchSuggestResp{
		Keywords: keywords,
	}, nil
}
