package bytedance

import (
	"context"
	"fmt"
	"net/http"
)

// MicroAppService handles communication with the micro app related
// method of the ByteDance API.
type MicroAppService service

// NewNameAuditInfo 名字审核信息
type NewNameAuditInfo struct {
	NewName           string `json:"new_name"`
	RemainingTimes    int    `json:"remaining_times"`
	NewNameAuditState int    `json:"new_name_audit_state"`
	Reason            string `json:"reason"`
	Advice            string `json:"advice"`
}

// NewIntroAuditInfo 介绍审核信息
type NewIntroAuditInfo struct {
	NewIntro           string `json:"new_intro"`
	RemainingTimes     int    `json:"remaining_times"`
	NewIntroAuditState int    `json:"new_intro_audit_state"`
	Reason             string `json:"reason"`
	Advice             string `json:"advice"`
}

// NewIconAuditInfo 图标审核信息
type NewIconAuditInfo struct {
	NewIcon           string `json:"new_icon"`
	RemainingTimes    int    `json:"remaining_times"`
	NewIconAuditState int    `json:"new_icon_audit_state"`
	Reason            string `json:"reason"`
	Advice            string `json:"advice"`
}

// AppCategoriesAuditInfo 类目审核信息
type AppCategoriesAuditInfo struct {
	AppCategory           string `json:"app_category"`
	AppCategoryName       string `json:"app_category_name"`
	AppCategoryAuditState int    `json:"app_category_audit_state"`
	Reason                string `json:"reason"`
}

// SubjectAuditInfo 主体信息
type SubjectAuditInfo struct {
	SubjectNumber     string `json:"subject_number"`
	SubjectName       string `json:"subject_name"`
	SubjectType       int    `json:"subject_type"`
	SubjectAuditState int    `json:"subject_audit_state"`
	Reason            string `json:"reason"`
}

// AppInfo 应用信息
type AppInfo struct {
	AppID                  string                  `json:"app_id"`
	AppType                int                     `json:"app_type"`
	AppState               int                     `json:"app_state"`
	AppName                string                  `json:"app_name"`
	NewNameAuditInfo       *NewNameAuditInfo       `json:"new_name_audit_info"`
	AppIntro               string                  `json:"app_intro"`
	NewIntroAuditInfo      *NewIntroAuditInfo      `json:"new_intro_audit_info"`
	AppIcon                string                  `json:"app_icon"`
	NewIConAuditInfo       *NewIconAuditInfo       `json:"new_i_con_audit_info"`
	AppCategoriesAuditInfo *AppCategoriesAuditInfo `json:"app_categories_audit_info"`
	SubjectAuditInfo       *SubjectAuditInfo       `json:"subject_audit_info"`
}

// AppInfoResponse 应用信息返回
type AppInfoResponse struct {
	Data *AppInfo `json:"data"`
}

// GetAppInfo 获取应用信息
func (s *MicroAppService) GetAppInfo(ctx context.Context, componentAppID, authorizerAccessToken string) (
	*AppInfoResponse, *http.Response, error) {
	u := fmt.Sprintf("v1/microapp/app/info?component_appid=%v&authorizer_access_token=%v",
		componentAppID, authorizerAccessToken)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	appInfo := new(AppInfoResponse)
	resp, err := s.client.Do(ctx, req, appInfo)
	if err != nil {
		return nil, resp, err
	}
	return appInfo, resp, nil
}

// DownloadQrcodeRequest 获取二维码请求
type DownloadQrcodeRequest struct {
	Version string `json:"version"` // 入参：current 或 audit 或 latest current 线上版 audit 审核版 latest 测试版
	Path    string `json:"path,omitempty"`
}

// DownloadQrcode 获取二维码
func (s *MicroAppService) DownloadQrcode(ctx context.Context, componentAppID, authorizerAccessToken string,
	body *DownloadQrcodeRequest) (*http.Response, error) {
	u := fmt.Sprintf("v1/microapp/app/qrcode?component_appid=%v&authorizer_access_token=%v",
		authorizerAccessToken, componentAppID)

	req, err := s.client.NewRequest(http.MethodGet, u, body)
	if err != nil {
		return nil, err
	}
	return s.client.Do(ctx, req, nil)
}

// CheckAppName 小程序名称检测
func (s *MicroAppService) CheckAppName(ctx context.Context, componentAppID, authorizerAccessToken, appName string) (*http.Response, error) {
	u := fmt.Sprintf("v1/microapp/app/check_app_name?component_appid=%v&authorizer_access_token=%v&app_name=%v",
		authorizerAccessToken, componentAppID, appName)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	return s.client.Do(ctx, req, nil)
}

// ModifyAppNameRequest 修改应用名称
type ModifyAppNameRequest struct {
	NewName          string `json:"new_name"`
	MaterialFilePath string `json:"material_file_path,omitempty"`
}

// ModifyAppName 修改小程序名称
func (s *MicroAppService) ModifyAppName(ctx context.Context, componentAppID, authorizerAccessToken string,
	body *ModifyAppNameRequest) (*http.Response, error) {
	u := fmt.Sprintf("v1/microapp/app/modify_app_name?component_appid=%v&authorizer_access_token=%v",
		authorizerAccessToken, componentAppID)

	req, err := s.client.NewRequest(http.MethodGet, u, body)
	if err != nil {
		return nil, err
	}
	return s.client.Do(ctx, req, nil)
}

// ModifyAppIntroRequest 修改应用简介
type ModifyAppIntroRequest struct {
	NewIntro string `json:"new_intro"`
}

// ModifyIntro 修改小程序简介
func (s *MicroAppService) ModifyIntro(ctx context.Context, componentAppID, authorizerAccessToken string,
	body *ModifyAppIntroRequest) (*http.Response, error) {
	u := fmt.Sprintf("v1/microapp/app/modify_app_intro?component_appid=%v&authorizer_access_token=%v",
		authorizerAccessToken, componentAppID)

	req, err := s.client.NewRequest(http.MethodGet, u, body)
	if err != nil {
		return nil, err
	}
	return s.client.Do(ctx, req, nil)
}

// ModifyAppIconRequest 修改应用简介
type ModifyAppIconRequest struct {
	NewIconPath string `json:"new_icon_path"`
}

// ModifyAppIcon 修改小程序图标
func (s *MicroAppService) ModifyAppIcon(ctx context.Context, componentAppID, authorizerAccessToken string,
	body *ModifyAppIconRequest) (*http.Response, error) {
	u := fmt.Sprintf("v1/microapp/app/modify_app_intro?component_appid=%v&authorizer_access_token=%v",
		authorizerAccessToken, componentAppID)

	req, err := s.client.NewRequest(http.MethodGet, u, body)
	if err != nil {
		return nil, err
	}
	return s.client.Do(ctx, req, nil)
}

// ModifyServerDomainRequest 修改服务域名
type ModifyServerDomainRequest struct {
	// 默认将第三方平台配置好的 request、socket、upload、download 域名全部添加到授权小程序
	Action   string   `json:"action"` // add 添加，delete 删除，set 覆盖，get 获取
	Request  []string `json:"request"`
	Socket   []string `json:"socket"`
	Upload   []string `json:"upload"`
	Download []string `json:"download"`
}

// ServerDomain 服务域名
type ServerDomain struct {
	Request  []string `json:"request"`
	Socket   []string `json:"socket"`
	Upload   []string `json:"upload"`
	Download []string `json:"download"`
}

// ModifyServerDomain 修改服务域名
func (s *MicroAppService) ModifyServerDomain(ctx context.Context, componentAppID, authorizerAccessToken string,
	body *ModifyServerDomainRequest) (*http.Response, error) {
	u := fmt.Sprintf("v1/microapp/app/modify_server_domain?component_appid=%v&authorizer_access_token=%v",
		authorizerAccessToken, componentAppID)

	req, err := s.client.NewRequest(http.MethodGet, u, body)
	if err != nil {
		return nil, err
	}
	return s.client.Do(ctx, req, nil)
}

// ModifyWebviewDomainRequest 修改webview域名
type ModifyWebviewDomainRequest struct {
	// 默认将第三方平台配置好的 webview 域名全部添加到授权小程序
	Action  string   `json:"action"` // add 添加，delete 删除，set 覆盖，get 获取
	Webview []string `json:"webview"`
}

// ModifyWebviewDomain 修改webview域名
func (s *MicroAppService) ModifyWebviewDomain(ctx context.Context, componentAppID, authorizerAccessToken string,
	body *ModifyWebviewDomainRequest) (*http.Response, error) {
	u := fmt.Sprintf("v1/microapp/app/modify_webview_domain?component_appid=%v&authorizer_access_token=%v",
		authorizerAccessToken, componentAppID)

	req, err := s.client.NewRequest(http.MethodGet, u, body)
	if err != nil {
		return nil, err
	}
	return s.client.Do(ctx, req, nil)
}

// Session 返回
type Session struct {
	SessionKey      string `json:"session_key"`
	OpenID          string `json:"openid"`
	AnonymousOpenID string `json:"anonymous_openid"`
}

// Code2Session code2session
func (s *MicroAppService) Code2Session(ctx context.Context, componentAppID, authorizerAccessToken,
	code, anonymousCode string) (*Session, *http.Response, error) {
	u := fmt.Sprintf("v1/microapp/code2session?component_appid=%v&authorizer_access_token=%v&code=%v&anonymous_code=%v",
		componentAppID, authorizerAccessToken, code, anonymousCode)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	session := new(Session)
	resp, err := s.client.Do(ctx, req, code)
	if err != nil {
		return nil, resp, err
	}
	return session, resp, nil
}
