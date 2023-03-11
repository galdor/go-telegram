package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

type Response struct {
	Ok          bool               `json:"ok"`
	Description string             `json:"description"`
	Result      json.RawMessage    `json:"result"`
	ErrorCode   int                `json:"error_code"`
	Parameters  ResponseParameters `json:"parameters"`
}

type ResponseParameters struct {
	MigrateToChatId int `json:"migrate_to_chat_id"`
	RetryAfter      int `json:"retry_after"` // seconds
}

type MethodError struct {
	Description string
	ErrorCode   int
	Parameters  ResponseParameters
}

func (err MethodError) Error() string {
	var buf bytes.Buffer

	buf.WriteString("method error")

	if err.ErrorCode != 0 {
		buf.WriteByte(' ')
		buf.WriteString(strconv.Itoa(err.ErrorCode))
	}

	buf.WriteString(": ")
	buf.WriteString(err.Description)

	return buf.String()
}

func DecodeResponse(data []byte, result interface{}) error {
	var response Response
	if err := json.Unmarshal(data, &response); err != nil {
		return fmt.Errorf("cannot decode response: %w", err)
	}

	if !response.Ok {
		err := MethodError{
			Description: response.Description,
			ErrorCode:   response.ErrorCode,
			Parameters:  response.Parameters,
		}

		return &err
	}

	if err := json.Unmarshal(response.Result, result); err != nil {
		return fmt.Errorf("cannot decode result: %w", err)
	}

	return nil
}

type User struct {
	Id                      int64  `json:"id"`
	IsBot                   bool   `json:"is_bot"`
	FirstName               string `json:"first_name"`
	LastName                string `json:"last_name,omitempty"`
	Username                string `json:"username,omitempty"`
	LanguageCode            string `json:"language_code"`
	IsPremium               bool   `json:"is_premium,omitempty"`
	AddedToAttachmentMenu   bool   `json:"added_to_attachment_menu,omitempty"`
	CanJoinGroups           bool   `json:"can_join_groups,omitempty"`
	CanReadAllGroupMessages bool   `json:"can_read_all_group_messages,omitempty"`
	SupportsInlineQueries   bool   `json:"supports_inline_queries,omitempty"`
}
