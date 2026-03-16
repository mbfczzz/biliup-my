package service

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"biliup/internal/config"
	"biliup/internal/model"

	qrcode "github.com/skip2/go-qrcode"
)

// UploadService Bilibili 上传服务
type UploadService struct {
	cfg        *config.Config
	httpClient *http.Client
}

const (
	biliAPIBase   = "https://api.bilibili.com"
	biliPassport  = "https://passport.bilibili.com"
	biliMember    = "https://member.bilibili.com"
	biliPreupload = "https://member.bilibili.com/preupload"
)

func NewUploadService(cfg *config.Config) *UploadService {
	return &UploadService{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// ============ QR 码登录 ============

// GetQRCode 获取B站登录二维码
func (s *UploadService) GetQRCode() (*model.QRCodeResponse, error) {
	req, err := http.NewRequest("GET", biliPassport+"/x/passport-login/web/qrcode/generate", nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求二维码失败: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result struct {
		Code int `json:"code"`
		Data struct {
			URL       string `json:"url"`
			QRCodeKey string `json:"qrcode_key"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if result.Code != 0 {
		return nil, fmt.Errorf("获取二维码失败, code: %d", result.Code)
	}

	// 生成二维码图片 (Base64 PNG)
	png, err := qrcode.Encode(result.Data.URL, qrcode.Medium, 256)
	if err != nil {
		return nil, fmt.Errorf("生成二维码图片失败: %w", err)
	}

	imageBase64 := base64.StdEncoding.EncodeToString(png)

	return &model.QRCodeResponse{
		URL:       result.Data.URL,
		QRCodeKey: result.Data.QRCodeKey,
		Image:     fmt.Sprintf("data:image/png;base64,%s", imageBase64),
	}, nil
}

// LoginByQRCode 通过二维码登录
func (s *UploadService) LoginByQRCode(qrcodeKey string) (*model.LoginResponse, error) {
	params := url.Values{}
	params.Set("qrcode_key", qrcodeKey)

	reqURL := fmt.Sprintf("%s/x/passport-login/web/qrcode/poll?%s", biliPassport, params.Encode())
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("轮询登录状态失败: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result struct {
		Code int `json:"code"`
		Data struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
			URL     string `json:"url"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	loginResp := &model.LoginResponse{
		Code:    result.Data.Code,
		Message: result.Data.Message,
	}

	// code=0 表示登录成功
	if result.Data.Code == 0 {
		var cookies []string
		for _, cookie := range resp.Cookies() {
			cookies = append(cookies, cookie.Name+"="+cookie.Value)
		}

		// 从 URL 参数也提取 cookies
		if result.Data.URL != "" {
			parsedURL, _ := url.Parse(result.Data.URL)
			if parsedURL != nil {
				for _, kv := range strings.Split(parsedURL.RawQuery, "&") {
					parts := strings.SplitN(kv, "=", 2)
					if len(parts) == 2 {
						name := parts[0]
						if name == "SESSDATA" || name == "bili_jct" || name == "DedeUserID" || name == "DedeUserID__ckMd5" {
							cookies = append(cookies, kv)
						}
					}
				}
			}
		}

		loginResp.Cookie = strings.Join(cookies, "; ")
	}

	return loginResp, nil
}

// ============ 用户信息 ============

// GetMyInfo 获取B站用户信息
func (s *UploadService) GetMyInfo(cookie string) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", biliAPIBase+"/x/space/myinfo", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Cookie", cookie)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)
	return result, nil
}

// ============ 视频上传 ============

// UploadVideo 上传视频到B站
func (s *UploadService) UploadVideo(filePath string, template *model.UploadTemplate, cookie string) error {
	log.Printf("开始上传视频: %s", filePath)

	// Step 1: 预上传
	upInfo, err := s.preUpload(filePath, cookie)
	if err != nil {
		return fmt.Errorf("预上传失败: %w", err)
	}

	// Step 2: 分块上传
	if err := s.uploadChunks(filePath, upInfo, cookie); err != nil {
		return fmt.Errorf("分块上传失败: %w", err)
	}

	// Step 3: 提交投稿
	if err := s.submitVideo(upInfo, template, cookie); err != nil {
		return fmt.Errorf("投稿提交失败: %w", err)
	}

	log.Printf("视频上传成功: %s", filePath)
	return nil
}

type uploadInfo struct {
	UploadURL string
	Filename  string
	BizID     int64
	ChunkSize int64
	Auth      string
	UPOS_URI  string
	Endpoint  string
}

func (s *UploadService) preUpload(filePath string, cookie string) (*uploadInfo, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Set("name", filepath.Base(filePath))
	params.Set("size", strconv.FormatInt(fileInfo.Size(), 10))
	params.Set("r", "upos")
	params.Set("profile", "ugcfx/bup")
	params.Set("ssl", "0")
	params.Set("version", "2.14.0")
	params.Set("build", "2140000")
	params.Set("upcdn", "bda2")
	params.Set("probe_version", "20221109")

	reqURL := fmt.Sprintf("%s?%s", biliPreupload, params.Encode())
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Cookie", cookie)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result struct {
		OK        int      `json:"OK"`
		Auth      string   `json:"auth"`
		BizID     int64    `json:"biz_id"`
		ChunkSize int64    `json:"chunk_size"`
		Endpoint  string   `json:"endpoint"`
		Endpoints []string `json:"endpoints"`
		UposURI   string   `json:"upos_uri"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if result.OK != 1 {
		return nil, fmt.Errorf("预上传失败: %s", string(body))
	}

	endpoint := result.Endpoint
	if len(result.Endpoints) > 0 {
		endpoint = result.Endpoints[0]
	}

	uposPath := strings.TrimPrefix(result.UposURI, "upos://")

	return &uploadInfo{
		UploadURL: fmt.Sprintf("https://%s/%s", endpoint, uposPath),
		Filename:  filepath.Base(filePath),
		BizID:     result.BizID,
		ChunkSize: result.ChunkSize,
		Auth:      result.Auth,
		UPOS_URI:  result.UposURI,
		Endpoint:  endpoint,
	}, nil
}

func (s *UploadService) uploadChunks(filePath string, info *uploadInfo, cookie string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()
	chunkSize := info.ChunkSize
	if chunkSize == 0 {
		chunkSize = 4 * 1024 * 1024
	}
	totalChunks := (fileSize + chunkSize - 1) / chunkSize

	// 初始化上传
	initURL := fmt.Sprintf("%s?uploads&output=json", info.UploadURL)
	req, _ := http.NewRequest("POST", initURL, nil)
	req.Header.Set("X-Upos-Auth", info.Auth)
	req.Header.Set("Cookie", cookie)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("初始化上传失败: %w", err)
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	var initResult struct {
		UploadID string `json:"upload_id"`
	}
	json.Unmarshal(body, &initResult)

	log.Printf("上传初始化完成, uploadID: %s, 总分块: %d", initResult.UploadID, totalChunks)

	// 逐块上传
	buf := make([]byte, chunkSize)
	for i := int64(0); i < totalChunks; i++ {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return fmt.Errorf("读取文件失败: %w", err)
		}

		start := i * chunkSize
		end := start + int64(n)

		chunkURL := fmt.Sprintf("%s?partNumber=%d&uploadId=%s&chunk=%d&chunks=%d&size=%d&start=%d&end=%d&total=%d",
			info.UploadURL, i+1, initResult.UploadID, i, totalChunks, int64(n), start, end, fileSize)

		chunkReq, _ := http.NewRequest("PUT", chunkURL, bytes.NewReader(buf[:n]))
		chunkReq.Header.Set("X-Upos-Auth", info.Auth)
		chunkReq.Header.Set("Content-Type", "application/octet-stream")
		chunkReq.Header.Set("Content-Length", strconv.Itoa(n))
		chunkReq.Header.Set("Cookie", cookie)

		chunkResp, err := s.httpClient.Do(chunkReq)
		if err != nil {
			return fmt.Errorf("上传分块 %d 失败: %w", i, err)
		}
		chunkResp.Body.Close()

		log.Printf("上传进度: %d/%d (%.1f%%)", i+1, totalChunks, float64(i+1)/float64(totalChunks)*100)
	}

	// 完成上传
	finishURL := fmt.Sprintf("%s?output=json&name=%s&profile=ugcfx%%2Fbup&uploadId=%s&biz_id=%d",
		info.UploadURL, info.Filename, initResult.UploadID, info.BizID)

	var parts []map[string]interface{}
	for i := int64(0); i < totalChunks; i++ {
		parts = append(parts, map[string]interface{}{
			"partNumber": i + 1,
			"eTag":       "etag",
		})
	}
	partsJSON, _ := json.Marshal(map[string]interface{}{"parts": parts})

	finishReq, _ := http.NewRequest("POST", finishURL, bytes.NewReader(partsJSON))
	finishReq.Header.Set("X-Upos-Auth", info.Auth)
	finishReq.Header.Set("Content-Type", "application/json")
	finishReq.Header.Set("Cookie", cookie)

	finishResp, err := s.httpClient.Do(finishReq)
	if err != nil {
		return fmt.Errorf("完成上传失败: %w", err)
	}
	finishResp.Body.Close()

	log.Printf("文件上传完成: %s", info.Filename)
	return nil
}

func (s *UploadService) submitVideo(info *uploadInfo, template *model.UploadTemplate, cookie string) error {
	csrf := extractCookieValue(cookie, "bili_jct")

	// 解析 upos_uri 中的文件名（去除路径和扩展名）
	uposBase := filepath.Base(info.UPOS_URI)
	uposName := strings.TrimSuffix(uposBase, filepath.Ext(uposBase))

	submitData := map[string]interface{}{
		"copyright":      template.Copyright,
		"source":         template.CopyrightSource,
		"tid":            template.TID,
		"cover":          template.CoverPath,
		"title":          template.Title,
		"desc_format_id": 0,
		"desc":           template.Description,
		"dynamic":        template.Dynamic,
		"tag":            template.Tags,
		"videos": []map[string]interface{}{
			{
				"filename": uposName,
				"title":    "P1",
				"desc":     "",
			},
		},
		"no_reprint": template.NoReprint,
		"csrf":       csrf,
	}

	if template.Dolby == 1 {
		submitData["dolby"] = 1
	}
	if template.Dtime > 0 {
		submitData["dtime"] = template.Dtime
	}

	jsonData, _ := json.Marshal(submitData)

	reqURL := fmt.Sprintf("%s/x/vu/web/add/v3?csrf=%s", biliMember, csrf)
	req, err := http.NewRequest("POST", reqURL, bytes.NewReader(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Cookie", cookie)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			Aid  int64  `json:"aid"`
			Bvid string `json:"bvid"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("解析投稿响应失败: %w", err)
	}

	if result.Code != 0 {
		return fmt.Errorf("投稿失败: [%d] %s", result.Code, result.Message)
	}

	log.Printf("投稿成功! aid=%d, bvid=%s", result.Data.Aid, result.Data.Bvid)
	return nil
}

// ============ 分区信息 ============

// GetArchivePre 获取投稿预处理信息
func (s *UploadService) GetArchivePre(cookie string) (map[string]interface{}, error) {
	req, _ := http.NewRequest("GET", biliMember+"/x/vupre/web/archive/pre", nil)
	req.Header.Set("Cookie", cookie)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)
	return result, nil
}

// ============ 辅助函数 ============

func extractCookieValue(cookie, name string) string {
	for _, part := range strings.Split(cookie, ";") {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, name+"=") {
			return strings.TrimPrefix(part, name+"=")
		}
	}
	return ""
}

func fileMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// UploadCover 上传封面图片
func UploadCover(coverPath string, cookie string) (string, error) {
	file, err := os.Open(coverPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", filepath.Base(coverPath))
	io.Copy(part, file)
	writer.Close()

	req, _ := http.NewRequest("POST", biliMember+"/x/vu/web/cover/up", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Cookie", cookie)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	var result struct {
		Data struct {
			URL string `json:"url"`
		} `json:"data"`
	}
	json.Unmarshal(respBody, &result)
	return result.Data.URL, nil
}
