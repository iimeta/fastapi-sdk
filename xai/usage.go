package xai

import "github.com/iimeta/fastapi-sdk/v2/model"

// 1 tick = 1e-10 USD；系统额度 $1 = 1,000,000 tokens
// SpendTokens = ceil(cost_in_usd_ticks / 10,000)
const ticksPerQuotaToken int64 = 10000

func SpendTokensFromCostTicks(ticks int64) int {
	if ticks <= 0 {
		return 0
	}
	return int((ticks + ticksPerQuotaToken - 1) / ticksPerQuotaToken)
}

// 把 xAI 官方 usage 映射到统一 Usage，对齐计费读取字段：
//
//	输入/输出: PromptTokens / CompletionTokens（及 Input/OutputTokens）
//	文本: Prompt/Input TextTokens、Completion/Output TextTokens
//	图像: Prompt/Input ImageTokens、Completion/Output ImageTokens
//	推理: Completion/Output ReasoningTokens
//	缓存读取: PromptTokensDetails.CachedTokens、InputTokensDetails.CachedTokens、CacheReadInputTokens
//	缓存写入: xAI REST 不返回，保持 0（与 OpenAI 一致）
func ConvertUsage(official *model.XAIUsage) *model.Usage {
	if official == nil {
		return nil
	}

	prompt := firstNonZero(official.PromptTokens, official.InputTokens)
	completion := firstNonZero(official.CompletionTokens, official.OutputTokens)

	var pCached, pText, pImage, pAudio int
	if official.PromptTokensDetails != nil {
		pCached = official.PromptTokensDetails.CachedTokens
		pText = official.PromptTokensDetails.TextTokens
		pImage = official.PromptTokensDetails.ImageTokens
		pAudio = official.PromptTokensDetails.AudioTokens
	}
	var iCached, iText, iImage, iAudio int
	if official.InputTokensDetails != nil {
		iCached = official.InputTokensDetails.CachedTokens
		iText = official.InputTokensDetails.TextTokens
		iImage = official.InputTokensDetails.ImageTokens
		iAudio = official.InputTokensDetails.AudioTokens
	}
	var cText, cImage, cReasoning, cAudio, cAccepted, cRejected int
	if official.CompletionTokensDetails != nil {
		cText = official.CompletionTokensDetails.TextTokens
		cImage = official.CompletionTokensDetails.ImageTokens
		cReasoning = official.CompletionTokensDetails.ReasoningTokens
		cAudio = official.CompletionTokensDetails.AudioTokens
		cAccepted = official.CompletionTokensDetails.AcceptedPredictionTokens
		cRejected = official.CompletionTokensDetails.RejectedPredictionTokens
	}
	var oText, oImage, oReasoning int
	if official.OutputTokensDetails != nil {
		oText = official.OutputTokensDetails.TextTokens
		oImage = official.OutputTokensDetails.ImageTokens
		oReasoning = official.OutputTokensDetails.ReasoningTokens
	}

	cached := firstNonZero(pCached, iCached, official.CachedPromptTextTokens)
	imageIn := firstNonZero(pImage, iImage, official.PromptImageTokens)
	textIn := firstNonZero(pText, iText, official.PromptTextTokens)
	if textIn == 0 {
		textIn = prompt - imageIn
		if textIn < 0 {
			textIn = prompt
		}
	}
	audioIn := firstNonZero(pAudio, iAudio)

	reasoning := firstNonZero(cReasoning, oReasoning, official.ReasoningTokens)
	imageOut := firstNonZero(cImage, oImage)
	textOut := firstNonZero(cText, oText)
	if textOut == 0 {
		textOut = completion
	}

	usage := &model.Usage{
		PromptTokens:         prompt,
		CompletionTokens:     completion,
		InputTokens:          firstNonZero(official.InputTokens, prompt),
		OutputTokens:         firstNonZero(official.OutputTokens, completion),
		TotalTokens:          official.TotalTokens,
		CostInUsdTicks:       official.CostInUsdTicks,
		CacheReadInputTokens: cached,
		PromptTokensDetails: model.PromptTokensDetails{
			AudioTokens:  audioIn,
			CachedTokens: cached,
			TextTokens:   textIn,
			ImageTokens:  imageIn,
		},
		InputTokensDetails: model.InputTokensDetails{
			TextTokens:   textIn,
			ImageTokens:  imageIn,
			CachedTokens: cached,
		},
		CompletionTokensDetails: model.CompletionTokensDetails{
			AudioTokens:              cAudio,
			ReasoningTokens:          reasoning,
			TextTokens:               textOut,
			ImageTokens:              imageOut,
			AcceptedPredictionTokens: cAccepted,
			RejectedPredictionTokens: cRejected,
		},
		OutputTokensDetails: model.OutputTokensDetails{
			TextTokens:      textOut,
			ReasoningTokens: reasoning,
			ImageTokens:     imageOut,
		},
	}
	if usage.TotalTokens == 0 && (prompt > 0 || completion > 0 || reasoning > 0) {
		usage.TotalTokens = prompt + completion + reasoning
	}
	if usage.CostInUsdTicks == 0 && official.CostInNanoUsd != nil && *official.CostInNanoUsd > 0 {
		usage.CostInUsdTicks = *official.CostInNanoUsd * 10
	}
	if official.NumSourcesUsed > 0 {
		usage.SearchTokens = official.NumSourcesUsed
	}
	return usage
}

func firstNonZero(values ...int) int {
	for _, v := range values {
		if v != 0 {
			return v
		}
	}
	return 0
}
