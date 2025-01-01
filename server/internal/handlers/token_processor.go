package handlers

import (
	"encoding/json"
	"relayapi/server/internal/models"
)

// TokenProcessor 处理令牌的扩展信息
type TokenProcessor struct{}

// ExtInfoData 扩展信息的数据结构
type ExtInfoData struct {
	RepM string `json:"rep_m,omitempty"`
}

// ProcessRequestBody 处理请求体，根据扩展信息进行修改
func (p *TokenProcessor) ProcessRequestBody(token *models.Token, requestBody []byte) ([]byte, error) {
	if token.ExtInfo == "" {
		return requestBody, nil
	}

	// 解析扩展信息
	var extInfo ExtInfoData
	if err := json.Unmarshal([]byte(token.ExtInfo), &extInfo); err != nil {
		return nil, err
	}

	// 如果没有替换模型信息，直接返回原始请求体
	if extInfo.RepM == "" {
		return requestBody, nil
	}

	// 尝试解析请求体为 JSON
	var requestData map[string]interface{}
	if err := json.Unmarshal(requestBody, &requestData); err != nil {
		return requestBody, nil
	}

	// 检查并替换模型信息
	requestData["model"] = extInfo.RepM
	// 重新编码为 JSON
	modifiedBody, err := json.Marshal(requestData)
	if err != nil {
		return nil, err
	}
	return modifiedBody, nil

}
