package bytedance

import (
	"context"
	"fmt"
	"net/http"
)

// Template 模版
type Template struct {
	TemplateID  int        `json:"template_id"`
	UserVersion string     `json:"user_version"`
	UserDesc    string     `json:"user_desc"`
	CreateTime  *Timestamp `json:"create_time"`
}

// Templates 获取模版列表返回值
type Templates struct {
	TemplateList []*Template `json:"template_list"`
}

// GetTemplates 获取第三方应用的所有模版
func (s *ThirdPartyService) GetTemplates(ctx context.Context, componentAppID, componentAccessToken string) (
	*Templates, *http.Response, error) {
	u := fmt.Sprintf(
		"v1/tp/template/get_tpl_list?component_appid=%v&component_access_token=%v",
		componentAppID,
		componentAccessToken,
	)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}
	templates := new(Templates)
	resp, err := s.client.Do(ctx, req, templates)
	if err != nil {
		return nil, resp, err
	}
	return templates, resp, nil
}

// Draft 草稿
type Draft struct {
	DraftID     int        `json:"draft_id"`
	UserVersion string     `json:"user_version"`
	UserDesc    string     `json:"user_desc"`
	CreateTime  *Timestamp `json:"create_time"`
}

// Drafts 获取草稿列表返回值
type Drafts struct {
	DraftList []*Draft `json:"draft_list"`
}

// GetDrafts 获取第三方应用的草稿
func (s *ThirdPartyService) GetDrafts(ctx context.Context, componentAppID, componentAccessToken string) (
	*Drafts, *http.Response, error) {
	u := fmt.Sprintf(
		"v1/tp/template/get_draft_list?component_appid=%v&component_access_token=%v",
		componentAppID,
		componentAccessToken,
	)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}
	drafts := new(Drafts)
	resp, err := s.client.Do(ctx, req, drafts)
	if err != nil {
		return nil, resp, err
	}
	return drafts, resp, nil
}

// AddTemplateRequest 添加模版请求
type AddTemplateRequest struct {
	DraftID int `json:"draft_id"`
}

// AddTemplate adds a template from draft.
// 将临时草稿设置为持久的代码模板。每个第三方应用的模板上限为200个。
func (s *ThirdPartyService) AddTemplate(ctx context.Context, componentAppID, componentAccessToken string,
	body *AddTemplateRequest) (*http.Response, error) {
	u := fmt.Sprintf("v1/tp/template/add_tpl?component_appid=%v&component_access_token=%v",
		componentAccessToken, componentAppID)

	req, err := s.client.NewRequest(http.MethodPost, u, body)
	if err != nil {
		return nil, err
	}
	return s.client.Do(ctx, req, nil)
}

// DeleteTemplateRequest 删除模版请求
type DeleteTemplateRequest struct {
	TemplateID int `json:"template_id"`
}

// DeleteTemplate deletes a template from draft.
// 将临时草稿设置为持久的代码模板。每个第三方应用的模板上限为200个。
func (s *ThirdPartyService) DeleteTemplate(ctx context.Context, componentAppID, componentAccessToken string,
	body *DeleteTemplateRequest) (*http.Response, error) {
	u := fmt.Sprintf("v1/tp/template/del_tpl?component_appid=%v&component_access_token=%v",
		componentAccessToken, componentAppID)

	req, err := s.client.NewRequest(http.MethodPost, u, body)
	if err != nil {
		return nil, err
	}
	return s.client.Do(ctx, req, nil)
}
