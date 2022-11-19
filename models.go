package main

import "fmt"

// model is the model to use for the completion.
type model string

const (
	// Enumerate supported models

	// GPT-3
	ModelTextDaVinci model = "text-davinci-002"
	ModelTextCurie   model = "text-curie-001"
	ModelTextBabbage model = "text-babbage-001"
	ModelTextAda     model = "text-ada-001"

	// Codex
	ModelCodeDaVinci model = "code-davinci-002"
	ModelCodeCushman model = "code-cushman-001"
)

func (m model) validate() error {
	switch m {
	case ModelTextDaVinci,
		ModelCodeDaVinci,
		ModelCodeCushman,
		ModelTextCurie,
		ModelTextBabbage,
		ModelTextAda:
		return nil
	default:
		return fmt.Errorf("invalid model: %s", m)
	}
}

// CreateCompletionRequest is the request to create a completion.
type CreateCompletionRequest struct {
	Model            model    `json:"model"`
	Prompt           []string `json:"prompt"`
	MaxTokens        int      `json:"max_tokens"`
	Temperature      float32  `json:"temperature"`
	FrequencyPenalty float32  `json:"frequency_penalty"`
	PresencePenalty  float32  `json:"presence_penalty"`
	Stop             []string `json:"stop"`
}

func (r *CreateCompletionRequest) validate() error {
	if err := r.Model.validate(); err != nil {
		return err
	}

	if r.MaxTokens < 1 {
		return fmt.Errorf("max_tokens must be >= 1")
	}

	if r.Temperature < 0 || r.Temperature > 1 {
		return fmt.Errorf("temperature must be between 0 and 1")
	}
	return nil
}

// CreateCompletionResponse is the response from creating a completion.
type CreateCompletionResponse struct {
	ID      string          `json:"id"`
	Object  string          `json:"object"`
	Created uint64          `json:"created"`
	Model   string          `json:"model"`
	Choices []Choice        `json:"choices"`
	Usage   CompletionUsage `json:"usage"`
}

type (
	LogprobResult struct {
		Tokens        []string             `json:"tokens"`
		TokenLogprobs []float32            `json:"token_logprobs"`
		TopLogprobs   []map[string]float32 `json:"top_logprobs"`
		TextOffset    []int                `json:"text_offset"`
	}

	Choice struct {
		Text         string        `json:"text"`
		Index        int           `json:"index"`
		FinishReason string        `json:"finish_reason"`
		LogProbs     LogprobResult `json:"logprobs"`
	}

	CompletionUsage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	}

	ErrorResponse struct {
		Error *struct {
			Code    *int    `json:"code,omitempty"`
			Message string  `json:"message"`
			Param   *string `json:"param,omitempty"`
			Type    string  `json:"type"`
		} `json:"error,omitempty"`
	}
)
