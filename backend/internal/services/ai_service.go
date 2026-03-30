package services

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"go.uber.org/zap"

	"github.com/shaykhnazar/futureproof/internal/models"
	"github.com/shaykhnazar/futureproof/pkg/cache"
)

const analysisCacheTTL = 24 * time.Hour

// ChatMessage represents a single turn in a conversation
type ChatMessage struct {
	Role    string // "user" or "assistant"
	Content string
}

type AIService struct {
	client *anthropic.Client
	model  string
	cache  *cache.Redis
	logger *zap.Logger
}

func NewAIService(apiKey, model string, cache *cache.Redis, logger *zap.Logger) *AIService {
	c := anthropic.NewClient(option.WithAPIKey(apiKey))
	return &AIService{client: &c, model: model, cache: cache, logger: logger}
}

// AnalyzeCareer performs AI-powered career risk analysis with 24h caching
func (s *AIService) AnalyzeCareer(ctx context.Context, req models.AnalysisRequest) (*models.AnalysisResult, error) {
	cacheKey := fmt.Sprintf("analysis:%s", s.generateRequestHash(req))

	if cached, err := s.cache.Get(ctx, cacheKey); err == nil && cached != "" {
		var result models.AnalysisResult
		if err := json.Unmarshal([]byte(cached), &result); err == nil {
			s.logger.Debug("Analysis cache hit", zap.String("profession", req.ProfessionSlug))
			return &result, nil
		}
	}

	result, err := s.callClaude(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("career analysis failed: %w", err)
	}

	data, _ := json.Marshal(result)
	_ = s.cache.Set(ctx, cacheKey, data, analysisCacheTTL)
	s.logger.Info("Career analysis completed",
		zap.String("profession", req.ProfessionSlug),
		zap.String("location", req.Location),
	)
	return result, nil
}

func (s *AIService) callClaude(ctx context.Context, req models.AnalysisRequest) (*models.AnalysisResult, error) {
	msg, err := s.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.Model(s.model),
		MaxTokens: 2048,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(s.buildAnalysisPrompt(req))),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("Anthropic API error: %w", err)
	}
	if len(msg.Content) == 0 {
		return nil, fmt.Errorf("empty response from Claude")
	}

	raw := msg.Content[0].Text
	raw = extractJSON(raw)

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &data); err != nil {
		return nil, fmt.Errorf("failed to parse Claude response: %w\nraw: %s", err, raw)
	}

	title := strings.Title(strings.ReplaceAll(req.ProfessionSlug, "-", " "))
	result := &models.AnalysisResult{
		ProfessionSlug:  req.ProfessionSlug,
		ProfessionTitle: title,
		AIRiskScore:     int(toFloat64(data["ai_risk_score"])),
		RiskLevel:       toString(data["risk_level"]),
		Summary:         toString(data["summary"]),
		Timeline:        toString(data["timeline"]),
		Threats:         toStringSlice(data["threats"]),
		Opportunities:   toStringSlice(data["opportunities"]),
		SkillsToLearn:   toStringSlice(data["skills_to_learn"]),
		GeneratedAt:     time.Now(),
	}

	if pivots, ok := data["recommended_pivots"].([]interface{}); ok {
		for _, p := range pivots {
			if pm, ok := p.(map[string]interface{}); ok {
				name := toString(pm["target_profession"])
				result.RecommendedPivots = append(result.RecommendedPivots, models.PivotSuggestion{
					TargetProfession: name,
					TargetSlug:       toSlug(name),
					MatchScore:       int(toFloat64(pm["match_score"])),
					Reason:           toString(pm["reason"]),
					TimeToTransition: toString(pm["time_to_transition"]),
				})
			}
		}
	}

	return result, nil
}

func (s *AIService) buildAnalysisPrompt(req models.AnalysisRequest) string {
	skills := "not specified"
	if len(req.CurrentSkills) > 0 {
		skills = strings.Join(req.CurrentSkills, ", ")
	}
	return fmt.Sprintf(`You are a world-class career strategist specializing in AI-era disruption (2025+).

Analyze this professional profile and return ONLY valid JSON (no markdown, no explanation outside the JSON):

Profession: %s
Location: %s
Years of Experience: %d
Current Skills: %s

Return exactly this JSON structure:
{
  "ai_risk_score": <integer 0-100>,
  "risk_level": "<Low|Medium|High>",
  "summary": "<2-3 sentence honest assessment>",
  "threats": ["<threat 1>", "<threat 2>", "<threat 3>"],
  "opportunities": ["<opportunity 1>", "<opportunity 2>", "<opportunity 3>"],
  "recommended_pivots": [
    {"target_profession": "<job title>", "match_score": <0-100>, "reason": "<why>", "time_to_transition": "<e.g. 3-6 months>"}
  ],
  "skills_to_learn": ["<skill 1>", "<skill 2>", "<skill 3>", "<skill 4>", "<skill 5>"],
  "timeline": "<when significant disruption is expected>"
}`,
		req.ProfessionSlug, req.Location, req.YearsExp, skills)
}

// StreamChat streams a career coaching conversation
func (s *AIService) StreamChat(ctx context.Context, message string, history []ChatMessage) (chan string, error) {
	ch := make(chan string, 64)

	messages := []anthropic.MessageParam{}
	for _, h := range history {
		if h.Content == "" {
			continue
		}
		if h.Role == "assistant" {
			messages = append(messages, anthropic.NewAssistantMessage(anthropic.NewTextBlock(h.Content)))
		} else {
			messages = append(messages, anthropic.NewUserMessage(anthropic.NewTextBlock(h.Content)))
		}
	}
	messages = append(messages, anthropic.NewUserMessage(anthropic.NewTextBlock(message)))

	go func() {
		defer close(ch)
		stream := s.client.Messages.NewStreaming(ctx, anthropic.MessageNewParams{
			Model:     anthropic.Model(s.model),
			MaxTokens: 1024,
			System: []anthropic.TextBlockParam{
				{Text: `You are FutureProof AI — a career intelligence advisor for the post-AI era.
Help people navigate AI-driven career disruption with honesty, empathy, and actionable advice.
Be direct, specific, and realistic. Structure advice around: immediate actions, 6-month plan, 2-year vision.`},
			},
			Messages: messages,
		})

		for stream.Next() {
			event := stream.Current()
			switch e := event.AsAny().(type) {
			case anthropic.ContentBlockDeltaEvent:
				if delta, ok := e.Delta.AsAny().(anthropic.TextDelta); ok {
					select {
					case <-ctx.Done():
						return
					case ch <- delta.Text:
					}
				}
			}
		}
	}()

	return ch, nil
}

func (s *AIService) generateRequestHash(req models.AnalysisRequest) string {
	data := fmt.Sprintf("%s|%s|%d|%s", req.ProfessionSlug, req.Location, req.YearsExp, strings.Join(req.CurrentSkills, ","))
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)
}

func extractJSON(raw string) string {
	if idx := strings.Index(raw, "```json"); idx != -1 {
		raw = raw[idx+7:]
		if end := strings.Index(raw, "```"); end != -1 {
			raw = raw[:end]
		}
	} else if idx := strings.Index(raw, "{"); idx != -1 {
		if last := strings.LastIndex(raw, "}"); last != -1 {
			raw = raw[idx : last+1]
		}
	}
	return strings.TrimSpace(raw)
}

func toFloat64(v interface{}) float64 {
	switch n := v.(type) {
	case float64:
		return n
	case int:
		return float64(n)
	}
	return 0
}

func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	return fmt.Sprint(v)
}

func toStringSlice(v interface{}) []string {
	arr, ok := v.([]interface{})
	if !ok {
		return nil
	}
	out := make([]string, 0, len(arr))
	for _, item := range arr {
		out = append(out, fmt.Sprint(item))
	}
	return out
}

func toSlug(s string) string {
	return strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(s, " ", "-"), "/", "-"))
}
