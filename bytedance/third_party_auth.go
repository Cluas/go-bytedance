package bytedance

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// ComponentAccessToken is response of API v1/auth/tp/token.
type ComponentAccessToken struct {
	ComponentAccessToken string        `json:"component_access_token"`
	ExpiresIn            time.Duration `json:"expires_in"`
}

// GetComponentAccessToken gets a component_access_token.
// 获取第三方平台 component_access_token
// 每个令牌有效期是 2 小时
func (s *ThirdPartyService) GetComponentAccessToken(ctx context.Context, componentAppID, componentAppSecret,
	componentTicket string) (*ComponentAccessToken, *http.Response, error) {
	u := fmt.Sprintf("v1/auth/tp/token?component_app_id=%v&component_appsecret=%v&component_tiket=%v",
		componentAppID, componentAppSecret, componentTicket)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	v := new(ComponentAccessToken)
	resp, err := s.client.Do(ctx, req, v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// CreatePreAuthCodeRequest 创建预授权码请求body
type CreatePreAuthCodeRequest struct {
	ShareRatio  int `json:"share_ratio"`
	ShareAmount int `json:"share_amount"`
}

// PreAuthCode 创建预授权码返回值
type PreAuthCode struct {
	PreAuthCode string    `json:"pre_auth_code"`
	ExpiresIn   time.Time `json:"expires_in"`
}

// CreatePreAuthCode creates a pre_auth_code for auth.
// 用于获取预授权码，预授权码用于小程序授权时的第三方平台方安全验证。
// 每个预授权码有效期为 10 分钟。
func (s *ThirdPartyService) CreatePreAuthCode(ctx context.Context, componentAccessToken,
	componentAppID string, body *CreatePreAuthCodeRequest) (*PreAuthCode, *http.Response, error) {
	u := fmt.Sprintf("v2/auth/pre_auth_code?component_access_token=%v&component_appid=%v",
		componentAccessToken, componentAppID)

	req, err := s.client.NewRequest(http.MethodPost, u, body)
	if err != nil {
		return nil, nil, err
	}

	preAuthCode := new(PreAuthCode)
	resp, err := s.client.Do(ctx, req, preAuthCode)
	if err != nil {
		return nil, resp, err
	}
	return preAuthCode, resp, nil

}

// AuthorizePermission 授权小程序在授权跳转页勾选的权限
type AuthorizePermission struct {
	ID          int    `json:"id"`
	Category    string `json:"category"`
	Description string `json:"description"`
}

// OAuthToken 获取接口调用凭据
type OAuthToken struct {
	AuthorizeAccessToken  string                 `json:"authorize_access_token"`
	AuthorizeRefreshToken string                 `json:"authorize_refresh_token"`
	ExpiresIn             time.Time              `json:"expires_in"`
	AuthorizerAppID       string                 `json:"authorizer_app_id"`
	AuthorizePermission   []*AuthorizePermission `json:"authorize_permission"`
}

// GetOAuthToken get authorizer_access_token.
// 使用授权码换取小程序的接口调用凭据
// authorizer_access_token 有效期 2 小时
// authorizer_refresh_token 有效期 1 个月，且只可使用一次，使用后失效
func (s *ThirdPartyService) GetOAuthToken(ctx context.Context, componentAppID, componentAccessToken,
	authorizationCode, grantType string) (*OAuthToken, *http.Response, error) {
	u := fmt.Sprintf(
		"v1/oauth/token?component_appid=%v&component_access_token=%v&authorzation_code=%v&grant_type=%v",
		componentAppID,
		componentAccessToken,
		authorizationCode,
		grantType,
	)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}
	oAuthToken := new(OAuthToken)
	resp, err := s.client.Do(ctx, req, oAuthToken)
	if err != nil {
		return nil, resp, err
	}
	return oAuthToken, resp, nil

}

// RefreshOAuthTokenResponse 刷新授权小程序的接口调用凭据返回值
type RefreshOAuthTokenResponse struct {
	AuthorizeAccessToken  string `json:"authorize_access_token"`
	AuthorizeRefreshToken string `json:"authorize_refresh_token"`
	ExpiresIn             string `json:"expires_in"`
}

// RefreshOAuthToken refresh authorizer_access_token
// 刷新授权小程序的接口调用凭据
func (s *ThirdPartyService) RefreshOAuthToken(ctx context.Context, componentAppID, componentAccessToken,
	authorizerRefreshToken, grantType string) (*RefreshOAuthTokenResponse, *http.Response, error) {
	u := fmt.Sprintf(
		"v1/oauth/token?component_appid=%v&component_access_token=%v&authorzation_refresh_token=%v&grant_type=%v",
		componentAppID,
		componentAccessToken,
		authorizerRefreshToken,
		grantType,
	)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}
	oAuthToken := new(RefreshOAuthTokenResponse)
	resp, err := s.client.Do(ctx, req, oAuthToken)
	if err != nil {
		return nil, resp, err
	}
	return oAuthToken, resp, nil

}

// RetrieveAuthorizationCodeResponse 找回授权码
type RetrieveAuthorizationCodeResponse struct {
	AuthorizationCode string    `json:"authorization_code"`
	ExpiresIn         time.Time `json:"expires_in"`
}

// RetrieveAuthorizationCode retrieve authorizer_access_token.
// 找回授权码 补偿机制
func (s *ThirdPartyService) RetrieveAuthorizationCode(ctx context.Context, componentAppID, componentAccessToken,
	authorizationAppID string) (*RetrieveAuthorizationCodeResponse, *http.Response, error) {
	u := fmt.Sprintf(
		"v1/oauth/retrieve?component_appid=%v&component_access_token=%v&authorzation_app_id=%v",
		componentAppID,
		componentAccessToken,
		authorizationAppID,
	)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}
	authorizationCode := new(RetrieveAuthorizationCodeResponse)
	resp, err := s.client.Do(ctx, req, authorizationCode)
	if err != nil {
		return nil, resp, err
	}
	return authorizationCode, resp, nil
}
