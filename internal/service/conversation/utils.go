package conversation

import (
	"sync"

	"notionboy/internal/pkg/logger"

	"github.com/pkoukk/tiktoken-go"
	"github.com/sashabaranov/go-openai"
)

var (
	tk   *tiktoken.Tiktoken
	once sync.Once
)

func getModelMaxTokens(model string) int {
	switch model {
	case openai.GPT3Dot5Turbo, openai.GPT3Dot5Turbo0301:
		return 4096
	case openai.GPT4, openai.GPT40314:
		return 8092
	case openai.GPT432K, openai.GPT432K0314:
		return 32768
	}
	return 2000
}

func getDefaultModelMaxReturnToken(model string) int {
	return getModelMaxTokens(model) / 2
}

func calculateMaxReturnTokens(userPromptTokens, userRequiredMaxTokens int, model string) int {
	remainTokens := getModelMaxTokens(model) - userPromptTokens - 100
	if remainTokens > userRequiredMaxTokens {
		return userRequiredMaxTokens
	}
	return remainTokens
}

func getTiktoken() *tiktoken.Tiktoken {
	if tk == nil {
		var err error
		once.Do(func() {
			tk, err = tiktoken.GetEncoding("cl100k_base")
			if err != nil {
				logger.SugaredLogger.Errorw("tiktoken.EncodingForModel", "error", err)
			}
		})
	}
	return tk
}

// calculateTotalTokensRequired This function calculates the total number of tokens required based on the model and usage.
func calculateTotalTokensRequired(promptTokens, completionTokens int, model string) int {
	// Initializing totalTokens variable to zero.
	totalTokens := 0
	// The token pricing details can be found at https://openai.com/pricing#faq-token
	// A switch case is used to determine the model and its corresponding token price.
	switch model {
	case openai.GPT3Dot5Turbo, openai.GPT3Dot5Turbo0301:
		// For GPT-3.5, promptTokens and completionTokens are charged at the same rate.
		totalTokens = promptTokens + completionTokens
	case openai.GPT4, openai.GPT40314:
		// For GPT-4 8K, the prompt tokens are charged at 15x the rate of the completion tokens.
		totalTokens = promptTokens*15 + completionTokens*30
	case openai.GPT432K, openai.GPT432K0314:
		// For GPT-4 32K, the prompt tokens are charged at 30x the rate of the completion tokens.
		totalTokens = promptTokens*30 + completionTokens*60
	default:
		// If an unknown model is provided, we assume promptTokens and completionTokens are charged at the same rate.
		totalTokens = promptTokens + completionTokens
	}
	// logger.SugaredLogger.Debugw("calculateTotalTokensRequired", "totalTokens", totalTokens, "model", model, "promptTokens", promptTokens, "completionTokens", completionTokens)
	return totalTokens
}

// calculateTotalTokens calculates the total number of tokens required based on the model and usage.
// if model is unknown, we assume promptTokens and completionTokens are charged at the same rate.
func calculateTotalTokens(historyMessages []openai.ChatCompletionMessage, model, instruction, prompt, response string) int {
	// logger.SugaredLogger.Debugw("calculateTotalTokens", "model", model, "instruction", instruction, "prompt", prompt, "response", response, "historyMessages", historyMessages)
	// get the tiktoken
	tk := getTiktoken()

	// calculateTokens calculates the number of tokens required for a given message
	calculateTokens := func(msg string) int {
		return len(tk.Encode(msg, nil, nil))
	}

	// calculate the number of tokens required for the instruction
	instructionTokens := calculateTokens(instruction)

	// calculate the number of tokens required for each message in the history
	historyTokens := 0
	for _, m := range historyMessages {
		historyTokens += calculateTokens(m.Content)
	}

	// calculate the number of tokens required for the prompt and response
	promptTokens := calculateTokens(prompt)
	responseTokens := calculateTokens(response)

	// return the total number of tokens required
	return calculateTotalTokensRequired(instructionTokens+historyTokens+promptTokens, responseTokens, model)
}

// calculateTotalTokensForMessages calculates the total number of tokens required for a given set of messages.
// It takes in a slice of messages, a model, an instruction, a prompt, and a response.
// It returns the total number of tokens required.
func calculateTotalTokensForMessages(messages []*Message, model, instruction, prompt, response string) int {
	// Convert messages to chat messages
	chatMessages := make([]openai.ChatCompletionMessage, 0)
	if len(messages) > 0 {
		for _, m := range messages {
			chatMessages = append(chatMessages, m.toChatMessage()...)
		}
	}

	// Calculate the total number of tokens required
	return calculateTotalTokens(chatMessages, model, instruction, prompt, response)
}
