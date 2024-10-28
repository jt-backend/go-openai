package openai

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrTooManyEmptyStreamMessages = errors.New("stream has sent too many empty messages")
)

type CompletionStream struct {
	*streamReader[CompletionResponse]
}

// CreateCompletionStream — API call to create a completion w/ streaming
// support. It sets whether to stream back partial progress. If set, tokens will be
// sent as data-only server-sent events as they become available, with the
// stream terminated by a data: [DONE] message.
func (c *Client) CreateCompletionStream(
	ctx context.Context,
	request CompletionRequest,
) (stream *CompletionStream, err error) {
	urlSuffix := "/completions"
	if !checkEndpointSupportsModel(urlSuffix, request.Model) {
		err = ErrCompletionUnsupportedModel
		return
	}

	if !checkPromptType(request.Prompt) {
		err = ErrCompletionRequestPromptTypeNotSupported
		return
	}

	request.Stream = true
	req, err := c.newRequest(
		ctx,
		http.MethodPost,
		c.fullURL(urlSuffix, withModel(request.Model)),
		withBody(request),
	)
	if err != nil {
		return nil, err
	}

	resp, err := sendRequestStream[CompletionResponse](c, req)
	if err != nil {
		return
	}
	stream = &CompletionStream{
		streamReader: resp,
	}
	return
}

type RunStream struct {
	*streamReader[Run]
}

// CreateRunStream creates a new run and returns a streaming response.
func (c *Client) CreateRunStream(
	ctx context.Context,
	threadID string,
	request RunRequest,
) (stream *RunStream, err error) {
	urlSuffix := fmt.Sprintf("/threads/%s/runs", threadID)

	// 在這裡設置 Stream 標誌，以確保請求為流式請求
	request.Stream = true

	req, err := c.newRequest(
		ctx,
		http.MethodPost,
		c.fullURL(urlSuffix),
		withBody(request),
		withBetaAssistantVersion(c.config.AssistantVersion),
	)
	if err != nil {
		return nil, err
	}

	// 使用流式請求的方式來發送，並獲取流式響應
	resp, err := sendRequestStream[Run](c, req)
	if err != nil {
		return
	}

	stream = &RunStream{
		streamReader: resp,
	}
	return
}
