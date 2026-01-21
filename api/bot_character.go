package api

import (
	"encoding/json"
	"sync"
	"time"
)

// Character card API skeleton for future SealDice integration
// These APIs follow the protocol defined in docs/sealchat-protocol.md

// CharacterPendingRequest stores a pending character API request
type CharacterPendingRequest struct {
	Echo      string
	API       string
	Data      any
	CreatedAt time.Time
	Response  chan *CharacterAPIResponse
}

// CharacterAPIResponse represents a response from SealDice
type CharacterAPIResponse struct {
	OK    bool            `json:"ok"`
	Data  json.RawMessage `json:"data,omitempty"`
	Error string          `json:"error,omitempty"`
}

// characterPendingRequests stores pending character API requests by echo ID
var characterPendingRequests = &sync.Map{}

// characterRequestTimeout is the timeout for character API requests
const characterRequestTimeout = 5 * time.Second

// apiCharacterGet handles character.get requests
// This is a SealChat → SealDice API that retrieves character card data
func apiCharacterGet(ctx *ChatContext, msg []byte) {
	data := struct {
		Echo string `json:"echo"`
		Data struct {
			GroupID string `json:"group_id"` // channel_id
			UserID  string `json:"user_id"`
		} `json:"data"`
	}{}
	if err := json.Unmarshal(msg, &data); err != nil {
		sendCharacterError(ctx, data.Echo, "请求解析失败")
		return
	}

	// Find connected BOT for this channel
	botConn := findBotConnectionForChannel(ctx, data.Data.GroupID)
	if botConn == nil {
		sendCharacterError(ctx, data.Echo, "未找到可用的 BOT 连接")
		return
	}

	// Forward request to BOT and wait for response
	resp := forwardCharacterRequest(botConn, "character.get", data.Echo, data.Data)
	if resp == nil {
		sendCharacterError(ctx, data.Echo, "请求超时")
		return
	}

	sendCharacterResponse(ctx, data.Echo, resp)
}

// apiCharacterSet handles character.set requests
// This is a SealChat → SealDice API that writes character card data
func apiCharacterSet(ctx *ChatContext, msg []byte) {
	data := struct {
		Echo string `json:"echo"`
		Data struct {
			GroupID string                 `json:"group_id"` // channel_id
			UserID  string                 `json:"user_id"`
			Name    string                 `json:"name"`
			Attrs   map[string]interface{} `json:"attrs"`
		} `json:"data"`
	}{}
	if err := json.Unmarshal(msg, &data); err != nil {
		sendCharacterError(ctx, data.Echo, "请求解析失败")
		return
	}

	botConn := findBotConnectionForChannel(ctx, data.Data.GroupID)
	if botConn == nil {
		sendCharacterError(ctx, data.Echo, "未找到可用的 BOT 连接")
		return
	}

	resp := forwardCharacterRequest(botConn, "character.set", data.Echo, data.Data)
	if resp == nil {
		sendCharacterError(ctx, data.Echo, "请求超时")
		return
	}

	sendCharacterResponse(ctx, data.Echo, resp)
}

// apiCharacterList handles character.list requests
// This is a SealChat → SealDice API that lists user's character cards
func apiCharacterList(ctx *ChatContext, msg []byte) {
	data := struct {
		Echo string `json:"echo"`
		Data struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}{}
	if err := json.Unmarshal(msg, &data); err != nil {
		sendCharacterError(ctx, data.Echo, "请求解析失败")
		return
	}

	// For character.list, we need to find any connected BOT
	botConn := findAnyBotConnection(ctx)
	if botConn == nil {
		sendCharacterError(ctx, data.Echo, "未找到可用的 BOT 连接")
		return
	}

	resp := forwardCharacterRequest(botConn, "character.list", data.Echo, data.Data)
	if resp == nil {
		sendCharacterError(ctx, data.Echo, "请求超时")
		return
	}

	sendCharacterResponse(ctx, data.Echo, resp)
}

// findBotConnectionForChannel finds a BOT WebSocket connection for a specific channel
// TODO: Implement proper BOT routing based on channel membership
func findBotConnectionForChannel(ctx *ChatContext, channelID string) *WsSyncConn {
	// Skeleton: iterate through userId2ConnInfo to find BOT connections
	// that are associated with the given channel
	_ = channelID
	_ = ctx

	// For now, return nil - to be implemented when BOT connection tracking is added
	return nil
}

// findAnyBotConnection finds any available BOT WebSocket connection
func findAnyBotConnection(ctx *ChatContext) *WsSyncConn {
	_ = ctx
	// Skeleton: iterate through userId2ConnInfo to find any BOT connection
	return nil
}

// forwardCharacterRequest forwards a character API request to a BOT
func forwardCharacterRequest(botConn *WsSyncConn, api, echo string, data any) *CharacterAPIResponse {
	if botConn == nil {
		return nil
	}

	// Create pending request with response channel
	respChan := make(chan *CharacterAPIResponse, 1)
	pending := &CharacterPendingRequest{
		Echo:      echo,
		API:       api,
		Data:      data,
		CreatedAt: time.Now(),
		Response:  respChan,
	}
	characterPendingRequests.Store(echo, pending)
	defer characterPendingRequests.Delete(echo)

	// Send request to BOT
	req := map[string]any{
		"api":  api,
		"echo": echo,
		"data": data,
	}
	if err := botConn.WriteJSON(req); err != nil {
		return nil
	}

	// Wait for response with timeout
	select {
	case resp := <-respChan:
		return resp
	case <-time.After(characterRequestTimeout):
		return nil
	}
}

// HandleCharacterResponse processes a character API response from BOT
// This should be called when receiving a response with empty "api" field
func HandleCharacterResponse(echo string, data json.RawMessage) bool {
	pending, ok := characterPendingRequests.Load(echo)
	if !ok {
		return false
	}

	req := pending.(*CharacterPendingRequest)

	var resp CharacterAPIResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		resp = CharacterAPIResponse{OK: false, Error: "响应解析失败"}
	}

	select {
	case req.Response <- &resp:
	default:
	}

	return true
}

func sendCharacterError(ctx *ChatContext, echo, errMsg string) {
	resp := map[string]any{
		"api":  "",
		"echo": echo,
		"data": map[string]any{
			"ok":    false,
			"error": errMsg,
		},
	}
	_ = ctx.Conn.WriteJSON(resp)
}

func sendCharacterResponse(ctx *ChatContext, echo string, resp *CharacterAPIResponse) {
	result := map[string]any{
		"api":  "",
		"echo": echo,
		"data": map[string]any{
			"ok": resp.OK,
		},
	}
	if resp.Error != "" {
		result["data"].(map[string]any)["error"] = resp.Error
	}
	if resp.Data != nil {
		result["data"].(map[string]any)["data"] = resp.Data
	}
	_ = ctx.Conn.WriteJSON(result)
}
