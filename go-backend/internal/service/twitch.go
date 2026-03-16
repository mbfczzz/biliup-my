package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"biliup/internal/config"
	"biliup/internal/model"
)

// TwitchService Twitch 直播检测和流获取
type TwitchService struct {
	cfg        *config.Config
	httpClient *http.Client
	clientID   string
}

const (
	twitchGQLEndpoint = "https://gql.twitch.tv/gql"
	twitchClientID    = "kimne78kx3ncx6brgo4mv6wki5h1ko"
)

// GraphQL 请求/响应结构
type gqlRequest struct {
	OperationName string      `json:"operationName"`
	Variables     interface{} `json:"variables"`
	Extensions    interface{} `json:"extensions"`
}

type streamMetadataVars struct {
	ChannelLogin string `json:"channelLogin"`
}

type accessTokenVars struct {
	IsLive     bool   `json:"isLive"`
	Login      string `json:"login"`
	IsVod      bool   `json:"isVod"`
	VodID      string `json:"vodID"`
	PlayerType string `json:"playerType"`
	Platform   string `json:"platform"`
}

type accessTokenExtensions struct {
	PersistedQuery struct {
		Version    int    `json:"version"`
		Sha256Hash string `json:"sha256Hash"`
	} `json:"persistedQuery"`
}

func NewTwitchService(cfg *config.Config) *TwitchService {
	return &TwitchService{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
		clientID: twitchClientID,
	}
}

// CheckStream 检测频道是否在线
func (s *TwitchService) CheckStream(channelName string) (*model.TwitchStreamInfo, error) {
	// 使用 GraphQL 查询频道状态
	reqBody := []gqlRequest{
		{
			OperationName: "StreamMetadata",
			Variables: streamMetadataVars{
				ChannelLogin: channelName,
			},
			Extensions: accessTokenExtensions{
				PersistedQuery: struct {
					Version    int    `json:"version"`
					Sha256Hash string `json:"sha256Hash"`
				}{
					Version:    1,
					Sha256Hash: "252a46e3f5b1ddc431b396e688331d8d020daec27079893ac7d4e6db759a7402",
				},
			},
		},
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	req, err := http.NewRequest("POST", twitchGQLEndpoint, bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Client-ID", s.clientID)
	req.Header.Set("Content-Type", "application/json")

	// 如果配置了 Twitch auth token
	if token := s.cfg.Get("twitch_cookie"); token != "" {
		req.Header.Set("Authorization", "OAuth "+token)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求 Twitch API 失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var results []map[string]interface{}
	if err := json.Unmarshal(body, &results); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	info := &model.TwitchStreamInfo{
		UserLogin: channelName,
		IsLive:    false,
	}

	if len(results) == 0 {
		return info, nil
	}

	dataMap, ok := results[0]["data"].(map[string]interface{})
	if !ok {
		return info, nil
	}

	user, ok := dataMap["user"].(map[string]interface{})
	if !ok {
		return info, nil
	}

	if displayName, ok := user["displayName"].(string); ok {
		info.UserName = displayName
	}

	stream, ok := user["stream"].(map[string]interface{})
	if !ok || stream == nil {
		return info, nil
	}

	info.IsLive = true
	if title, ok := stream["title"].(string); ok {
		info.Title = title
	}
	if vc, ok := stream["viewersCount"].(float64); ok {
		info.ViewerCount = int(vc)
	}
	if game, ok := stream["game"].(map[string]interface{}); ok {
		if name, ok := game["name"].(string); ok {
			info.GameName = name
		}
	}

	// 获取预览图
	if previewURL, ok := stream["previewImageURL"].(string); ok {
		info.ThumbnailURL = previewURL
	}

	return info, nil
}

// GetStreamURL 获取直播流 URL
func (s *TwitchService) GetStreamURL(channelName string) (string, error) {
	// 获取 playback access token
	reqBody := []gqlRequest{
		{
			OperationName: "PlaybackAccessToken",
			Variables: accessTokenVars{
				IsLive:     true,
				Login:      channelName,
				IsVod:      false,
				VodID:      "",
				PlayerType: "site",
				Platform:   "web",
			},
			Extensions: accessTokenExtensions{
				PersistedQuery: struct {
					Version    int    `json:"version"`
					Sha256Hash string `json:"sha256Hash"`
				}{
					Version:    1,
					Sha256Hash: "ed230aa1e33e07eebb8928504583da78a5173989fadfb1ac94be06a04f3cdbe9",
				},
			},
		},
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %w", err)
	}

	req, err := http.NewRequest("POST", twitchGQLEndpoint, bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Client-ID", s.clientID)
	req.Header.Set("Content-Type", "application/json")
	if token := s.cfg.Get("twitch_cookie"); token != "" {
		req.Header.Set("Authorization", "OAuth "+token)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	var results []map[string]interface{}
	if err := json.Unmarshal(body, &results); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	if len(results) == 0 {
		return "", fmt.Errorf("空响应")
	}

	dataMap, ok := results[0]["data"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("响应格式错误")
	}

	tokenData, ok := dataMap["streamPlaybackAccessToken"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("无法获取 access token")
	}

	token, _ := tokenData["value"].(string)
	sig, _ := tokenData["signature"].(string)

	// 构建 HLS URL
	streamURL := fmt.Sprintf(
		"https://usher.ttvnw.net/api/channel/hls/%s.m3u8?sig=%s&token=%s&allow_source=true&allow_audio_only=true&fast_bread=true&p=%d",
		channelName, sig, token, time.Now().UnixMilli()%999999+1,
	)

	return streamURL, nil
}

// BatchCheckStreams 批量检测频道状态
func (s *TwitchService) BatchCheckStreams(channels []string) (map[string]*model.TwitchStreamInfo, error) {
	result := make(map[string]*model.TwitchStreamInfo)

	// Twitch GraphQL 支持批量查询，每次最多30个
	batchSize := 30
	for i := 0; i < len(channels); i += batchSize {
		end := i + batchSize
		if end > len(channels) {
			end = len(channels)
		}
		batch := channels[i:end]

		var reqBody []gqlRequest
		for _, ch := range batch {
			reqBody = append(reqBody, gqlRequest{
				OperationName: "StreamMetadata",
				Variables: streamMetadataVars{
					ChannelLogin: ch,
				},
				Extensions: accessTokenExtensions{
					PersistedQuery: struct {
						Version    int    `json:"version"`
						Sha256Hash string `json:"sha256Hash"`
					}{
						Version:    1,
						Sha256Hash: "252a46e3f5b1ddc431b396e688331d8d020daec27079893ac7d4e6db759a7402",
					},
				},
			})
		}

		data, err := json.Marshal(reqBody)
		if err != nil {
			continue
		}

		req, err := http.NewRequest("POST", twitchGQLEndpoint, bytes.NewReader(data))
		if err != nil {
			continue
		}

		req.Header.Set("Client-ID", s.clientID)
		req.Header.Set("Content-Type", "application/json")
		if token := s.cfg.Get("twitch_cookie"); token != "" {
			req.Header.Set("Authorization", "OAuth "+token)
		}

		resp, err := s.httpClient.Do(req)
		if err != nil {
			continue
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		var results []map[string]interface{}
		if err := json.Unmarshal(body, &results); err != nil {
			continue
		}

		for j, res := range results {
			if j >= len(batch) {
				break
			}
			ch := batch[j]
			info := &model.TwitchStreamInfo{
				UserLogin: ch,
				IsLive:    false,
			}

			if dataMap, ok := res["data"].(map[string]interface{}); ok {
				if user, ok := dataMap["user"].(map[string]interface{}); ok {
					if displayName, ok := user["displayName"].(string); ok {
						info.UserName = displayName
					}
					if stream, ok := user["stream"].(map[string]interface{}); ok && stream != nil {
						info.IsLive = true
						if title, ok := stream["title"].(string); ok {
							info.Title = title
						}
						if vc, ok := stream["viewersCount"].(float64); ok {
							info.ViewerCount = int(vc)
						}
						if game, ok := stream["game"].(map[string]interface{}); ok {
							if name, ok := game["name"].(string); ok {
								info.GameName = name
							}
						}
					}
				}
			}
			result[ch] = info
		}
	}

	return result, nil
}

// ExtractChannelName 从 URL 提取频道名
func ExtractChannelName(url string) string {
	url = strings.TrimRight(url, "/")
	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return url
}

// ParseTwitchURL 验证并解析 Twitch URL
func ParseTwitchURL(url string) (string, error) {
	url = strings.TrimSpace(url)
	if !strings.Contains(url, "twitch.tv/") {
		return "", fmt.Errorf("不是有效的 Twitch URL: %s", url)
	}
	channel := ExtractChannelName(url)
	if channel == "" {
		return "", fmt.Errorf("无法提取频道名: %s", url)
	}
	log.Printf("解析 Twitch URL: %s -> 频道: %s", url, channel)
	return channel, nil
}
