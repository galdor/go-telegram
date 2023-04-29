package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

type Integer int64

func (i Integer) MarshalJSON() ([]byte, error) {
	return json.Marshal(strconv.FormatInt(int64(i), 10))
}

func (i *Integer) UnmarshalJSON(data []byte) error {
	var value json.Number

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	i64, err := value.Int64()
	if err != nil {
		return err
	}

	*i = Integer(i64)
	return nil
}

type Response struct {
	Ok          bool               `json:"ok"`
	Description string             `json:"description"`
	Result      json.RawMessage    `json:"result"`
	ErrorCode   int                `json:"error_code"`
	Parameters  ResponseParameters `json:"parameters"`
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
