package bytedance

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// UploadPackageRequest 上传代码请求.
type UploadPackageRequest struct {
	TemplateID  int    `json:"template_id"`
	UserDesc    string `json:"user_desc"`
	UserVersion string `json:"user_version"`
	ExtJSON     string `json:"ext_json"`
}

// UploadPackage 提交代码
// 为授权小程序提交代码（提交成功后，授权小程序具有测试版本）.
func (s *MicroAppService) UploadPackage(ctx context.Context, componentAppID, componentAccessToken string,
	body *UploadPackageRequest) (*http.Response, error) {
	u := fmt.Sprintf("v1/microapp/package/upload?component_app_id=%v&component_access_token=%v",
		componentAccessToken, componentAppID)

	req, err := s.client.NewRequest(http.MethodPost, u, body)
	if err != nil {
		return nil, err
	}
	return s.client.Do(ctx, req, nil)
}

// PackageAuditHosts 获取可选审核宿主端列表返回.
type PackageAuditHosts struct {
	HostNames         []string `json:"hostNames"`
	ReleasedHostNames []string `json:"releasedHostNames"`
}

// GetPackageAuditHosts 获取可选审核宿主端列表
// 获取可以提审的端，作为参数传入提审代码v2接口中
func (s *MicroAppService) GetPackageAuditHosts(ctx context.Context, componentAppID, componentAccessToken string) (
	*PackageAuditHosts, *http.Response, error) {
	u := fmt.Sprintf("v1/microapp/package/audit_hosts?component_app_id=%v&component_access_token=%v",
		componentAppID, componentAccessToken)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	auditHosts := new(PackageAuditHosts)
	resp, err := s.client.Do(ctx, req, auditHosts)
	if err != nil {
		return nil, resp, err
	}
	return auditHosts, resp, nil
}

// CommitAuditPackageRequest 提审代码请求
type CommitAuditPackageRequest struct {
	HostNames []string `json:"hostNames"`
}

// CommitAuditPackage 提审代码 v2 支持传入宿主端参数
// 为授权小程序提审代码（审核成功后，授权小程序具有审核版本）
func (s *MicroAppService) CommitAuditPackage(ctx context.Context, componentAppID, componentAccessToken string,
	body *UploadPackageRequest) (*http.Response, error) {
	u := fmt.Sprintf("v2/microapp/package/audit?component_app_id=%v&component_access_token=%v",
		componentAccessToken, componentAppID)

	req, err := s.client.NewRequest(http.MethodPost, u, body)
	if err != nil {
		return nil, err
	}
	return s.client.Do(ctx, req, nil)
}

// ReleasePackage 发布代码
// 为授权小程序发布代码（发布成功后，授权小程序具有线上版本）
func (s *MicroAppService) ReleasePackage(ctx context.Context, componentAppID, componentAccessToken string) (
	*http.Response, error) {
	u := fmt.Sprintf("v2/microapp/package/release?component_app_id=%v&component_access_token=%v",
		componentAccessToken, componentAppID)

	req, err := s.client.NewRequest(http.MethodPost, u, nil)
	if err != nil {
		return nil, err
	}
	return s.client.Do(ctx, req, nil)
}

// RollbackPackage 回退代码版本
// 为授权小程序回退代码版本，此操作可能需要等待一会（如果可以回退，执行成功后，授权小程序将回退至上一个线上版本）
func (s *MicroAppService) RollbackPackage(ctx context.Context, componentAppID, componentAccessToken string) (
	*http.Response, error) {
	u := fmt.Sprintf("v2/microapp/package/rollback?component_app_id=%v&component_access_token=%v",
		componentAccessToken, componentAppID)

	req, err := s.client.NewRequest(http.MethodPost, u, nil)
	if err != nil {
		return nil, err
	}
	return s.client.Do(ctx, req, nil)
}

// Rollback info
type Rollback struct {
	CanRollback bool   `json:"can_rollback"`
	LastVersion string `json:"last_version"`
}

// VersionCommon 通用返回值
type VersionCommon struct {
	Categories      []string  `json:"categories"`
	Ctime           time.Time `json:"ctime"`
	DeveloperAvatar string    `json:"developer_avatar"`
	DeveloperID     string    `json:"developer_id"`
	DeveloperName   string    `json:"developer_name"`
	Summary         string    `json:"summary"`
	Version         string    `json:"version"`
}

// Audit 审核版本返回值
type Audit struct {
	VersionCommon

	ApprovedApps     []int           `json:"approvedApps"`
	AttachInfo       json.RawMessage `json:"attachInfo"`
	HasPublish       int             `json:"has_publish"`
	IsIllegalVersion bool            `json:"is_illegal_version"`
	Reason           string          `json:"reason"`
	ReasonDetail     json.RawMessage `json:"reason_detail"`
	Status           int             `json:"status"`
}

// Current 线上版本返回值
type Current struct {
	VersionCommon

	ApprovedApps    []int           `json:"approvedApps"`
	AttachInfo      json.RawMessage `json:"attachInfo"`
	HasDown         int             `json:"has_down"`
	NotApprovedApps []string        `json:"notApprovedApps"`
	Reason          string          `json:"reason"`
	ReasonDetail    json.RawMessage `json:"reason_detail"`
	Rollback        *Rollback       `json:"rollback"`
	LastVersion     string          `json:"last_version"`
	UID             string          `json:"uid,omitempty"`
}

// Latest 测试版本返回值
type Latest struct {
	VersionCommon

	HasAudit   int    `json:"has_audit"`
	ScreenShot string `json:"screen_shot"`
}

// PackageVersions 获取小程序版本列表信息返回
type PackageVersions struct {
	Audit   *Audit   `json:"audit"`
	Current *Current `json:"current"`
	Latest  *Latest  `json:"latest"`
}

// GetPackageVersions 获取小程序版本列表信息
// 返回的结果列表中一共会展示三种状态的小程序代码版本信息，包括测试版本、审核版本、线上版本。
func (s *MicroAppService) GetPackageVersions(ctx context.Context, componentAppID, componentAccessToken string) (
	*PackageVersions, *http.Response, error) {
	u := fmt.Sprintf("v1/microapp/package/versions?component_app_id=%v&component_access_token=%v",
		componentAppID, componentAccessToken)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	versions := new(PackageVersions)
	resp, err := s.client.Do(ctx, req, versions)
	if err != nil {
		return nil, resp, err
	}
	return versions, resp, nil
}
