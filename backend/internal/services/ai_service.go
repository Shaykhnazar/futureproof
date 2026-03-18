package services

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/shaykhnazar/futureproof/internal/models"
	"github.com/shaykhnazar/futureproof/pkg/cache"
)

const (
	analysisCacheTTL = 24 * time.Hour
)

// AIService handles AI-powered career analysis
type AIService struct {
	apiKey string
	model  string
	cache  *cache.Redis
	logger *zap.Logger
}

// NewAIService creates a new AI service
func NewAIService(apiKey, model string, cache *cache.Redis, logger *zap.Logger) *AIService {
	return &AIService{
		apiKey: apiKey,
		model:  model,
		cache:  cache,
		logger: logger,
	}
}

// AnalyzeCareer performs AI-powered career risk analysis
func (s *AIService) AnalyzeCareer(ctx context.Context, req models.AnalysisRequest) (*models.AnalysisResult, error) {
	// Generate request hash for caching
	requestHash := s.generateRequestHash(req)
	cacheKey := fmt.Sprintf("analysis:%s", requestHash)

	// Check cache first
	cached, err := s.cache.Get(ctx, cacheKey)
	if err == nil && cached != "" {
		var result models.AnalysisResult
		if err := json.Unmarshal([]byte(cached), &result); err == nil {
			s.logger.Debug("Analysis retrieved from cache", zap.String("profession", req.ProfessionSlug))
			return &result, nil
		}
	}

	// Call Claude API
	result, err := s.callClaudeAPI(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze career: %w", err)
	}

	// Cache the result
	data, _ := json.Marshal(result)
	_ = s.cache.Set(ctx, cacheKey, data, analysisCacheTTL)

	s.logger.Info("Career analysis completed",
		zap.String("profession", req.ProfessionSlug),
		zap.String("location", req.Location),
	)

	return result, nil
}

// callClaudeAPI makes the actual API call to Claude
func (s *AIService) callClaudeAPI(ctx context.Context, req models.AnalysisRequest) (*models.AnalysisResult, error) {
	// Build the prompt
	prompt := s.buildAnalysisPrompt(req)

	// For now, return a mock response (replace with actual Anthropic SDK call)
	// TODO: Implement actual Anthropic API integration
	result := &models.AnalysisResult{
		ProfessionSlug:  req.ProfessionSlug,
		ProfessionTitle: strings.ReplaceAll(strings.Title(req.ProfessionSlug), "-", " "),
		AIRiskScore:     45,
		RiskLevel:       "Medium",
		Summary:         fmt.Sprintf("Analysis for %s profession shows moderate AI automation risk.", req.ProfessionSlug),
		Threats: []string{
			"Routine tasks may be automated by AI systems",
			"Increased competition from AI-augmented professionals",
			"Some traditional skills becoming less valuable",
		},
		Opportunities: []string{
			"AI tools can enhance productivity and capabilities",
			"New roles emerging in AI implementation and oversight",
			"Opportunity to specialize in AI-resistant aspects of the role",
		},
		RecommendedPivots: []models.PivotSuggestion{
			{
				TargetProfession: "AI Ethics Officer",
				TargetSlug:       "ai-ethics-officer",
				MatchScore:       85,
				Reason:           "Your domain expertise is valuable for AI governance",
				TimeToTransition: "6-9 months with targeted learning",
			},
		},
		Timeline:    "5-10 years for significant automation impact",
		SkillsToLearn: []string{
			"AI/ML fundamentals",
			"Data analysis",
			"Strategic thinking",
			"Ethical AI practices",
		},
		GeneratedAt: time.Now(),
	}

	s.logger.Debug("Claude API called", zap.String("prompt_length", fmt.Sprintf("%d", len(prompt))))
	return result, nil
}

// buildAnalysisPrompt constructs the prompt for Claude
func (s *AIService) buildAnalysisPrompt(req models.AnalysisRequest) string {
	return fmt.Sprintf(`You are a career advisor specializing in AI-era career planning. Analyze the following profession:

Profession: %s
Location: %s
Years of Experience: %d
Current Skills: %s

Provide a comprehensive career risk analysis in JSON format with:
1. AI Risk Score (0-100, where 100 = highest automation risk)
2. Risk Level (Low/Medium/High)
3. Summary (2-3 sentences)
4. Key Threats (list of 3-5 specific threats)
5. Opportunities (list of 3-5 opportunities)
6. Recommended Career Pivots (3-5 suggestions with match scores)
7. Timeline (when significant changes are expected)
8. Skills to Learn (5-7 future-proof skills)

Be realistic, data-driven, and actionable.`,
		req.ProfessionSlug,
		req.Location,
		req.YearsExp,
		strings.Join(req.CurrentSkills, ", "),
	)
}

// generateRequestHash creates a unique hash for caching
func (s *AIService) generateRequestHash(req models.AnalysisRequest) string {
	data := fmt.Sprintf("%s|%s|%d|%s",
		req.ProfessionSlug,
		req.Location,
		req.YearsExp,
		strings.Join(req.CurrentSkills, ","),
	)
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)
}

// StreamChat handles streaming chat with AI career coach
func (s *AIService) StreamChat(ctx context.Context, message string, history []string) (chan string, error) {
	// TODO: Implement streaming chat with Anthropic API
	responseChan := make(chan string)

	go func() {
		defer close(responseChan)
		// Mock streaming response
		response := "I'm here to help you navigate your career in the AI era. How can I assist you today?"
		for _, char := range response {
			select {
			case <-ctx.Done():
				return
			case responseChan <- string(char):
				time.Sleep(10 * time.Millisecond)
			}
		}
	}()

	return responseChan, nil
}
