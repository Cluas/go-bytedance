package bytedance

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// File wrapped file content
type File struct {
	Name    string
	Content io.Reader
}

// Read proxy Content Read
func (f File) Read(p []byte) (n int, err error) {
	return f.Content.Read(p)
}

// DownloadWebViewFile 下载域名校验文件
func (s *ThirdPartyService) DownloadWebViewFile(ctx context.Context, componentAppID, componentAccessToken string,
	body *DeleteTemplateRequest) (*http.Response, error) {
	u := fmt.Sprintf("v1/tp/download/webview_file?component_appid=%v&component_access_token=%v",
		componentAccessToken, componentAppID)

	req, err := s.client.NewRequest(http.MethodGet, u, body)
	if err != nil {
		return nil, err
	}
	return s.client.Do(ctx, req, nil)
}

// UploadPicMaterialRequest 上传图片请求
type UploadPicMaterialRequest struct {
	MaterialType int   `json:"material_type"`
	MaterialFile *File `json:"material_file"`
}

// Params implements FormRender interface.
func (u UploadPicMaterialRequest) Params() map[string]string {
	return map[string]string{
		"material_type": strconv.Itoa(u.MaterialType),
	}
}

// MultipartParams implements FormRender interface.
func (u UploadPicMaterialRequest) MultipartParams() map[string]io.Reader {
	return map[string]io.Reader{
		"material_file": u.MaterialFile,
	}
}

// UploadPicMaterial 上传图片材料
// 使用修改名称、图标、服务类目等涉及材料证明的接口前，都需要先使用这个图片上传接口拿到返回的图片地址
// 目前只支持bmp、jpeg、jpg、png格式。
func (s *ThirdPartyService) UploadPicMaterial(ctx context.Context, componentAppID, componentAccessToken string,
	body *UploadPicMaterialRequest) (string, *http.Response, error) {
	u := fmt.Sprintf("v1/tp/upload_pic_material?component_appid=%v&component_access_token=%v",
		componentAccessToken, componentAppID)
	var address string
	req, err := s.client.NewRequest(http.MethodPost, u, body)
	if err != nil {
		return "", nil, err
	}
	resp, err := s.client.Do(ctx, req, &address)
	if err != nil {
		return "", resp, err
	}
	return address, resp, nil
}
