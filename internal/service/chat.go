package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"ledgerA/internal/dto"
	"ledgerA/internal/model"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// ChatService processes natural language messages.
type ChatService interface {
	Process(ctx context.Context, userID uuid.UUID, message string) (*dto.ChatResponse, error)
}

type chatService struct {
	transactionService TransactionService
	accountService     AccountService
	categoryService    CategoryService
	ollamaURL          string
}

// NewChatService creates a ChatService.
func NewChatService(ts TransactionService, as AccountService, cs CategoryService, ollamaURL string) ChatService {
	return &chatService{
		transactionService: ts,
		accountService:     as,
		categoryService:    cs,
		ollamaURL:          ollamaURL,
	}
}

func (s *chatService) Process(ctx context.Context, userID uuid.UUID, message string) (*dto.ChatResponse, error) {
	msg := strings.TrimSpace(message)
	lower := strings.ToLower(msg)

	// Try Ollama-based parsing first if configured
	if s.ollamaURL != "" {
		parsed, err := s.parseWithOllama(ctx, msg)
		if err == nil && parsed != nil {
			return s.executeAction(ctx, userID, parsed)
		}
		log.Debug().Err(err).Msg("ollama parse failed, falling back to pattern matching")
	}

	// Pattern matching fallback
	if action := s.parseAddTransaction(lower, msg); action != nil {
		return s.executeAction(ctx, userID, action)
	}

	if strings.Contains(lower, "balance") || strings.Contains(lower, "how much") {
		return s.handleBalanceQuery(ctx, userID, lower)
	}

	if strings.Contains(lower, "recent") || strings.Contains(lower, "last") || strings.Contains(lower, "transactions") {
		return s.handleRecentTransactions(ctx, userID)
	}

	if lower == "help" || strings.Contains(lower, "what can you") || strings.Contains(lower, "how to") {
		return &dto.ChatResponse{
			Reply:   "I can help you with:\n• **Add transactions**: \"Spent 500 on groceries\" or \"Got 10000 salary\"\n• **Check balances**: \"What's my balance?\" or \"Balance of Savings\"\n• **Recent transactions**: \"Show recent transactions\"\n\nTry natural phrases like \"Paid 200 for coffee from HDFC\" or \"Received 5000 freelance income\".",
			Success: true,
		}, nil
	}

	return &dto.ChatResponse{
		Reply:   "I didn't understand that. Try something like:\n• \"Spent 500 on groceries\"\n• \"Got 10000 salary\"\n• \"What's my balance?\"\n• \"Show recent transactions\"\n\nType **help** for more examples.",
		Success: true,
	}, nil
}

// Patterns for expense:
// "spent 500 on groceries"
// "paid 200 for coffee"
// "bought lunch for 150"
// "expense 300 rent"
// Patterns for income:
// "got 10000 salary"
// "received 5000 freelance"
// "earned 2000"
// "income 5000 salary"

var (
	spentOnPattern   = regexp.MustCompile(`(?i)(?:spent|paid|expense)\s+(\d+(?:\.\d+)?)\s+(?:on|for)\s+(.+?)(?:\s+from\s+(.+))?$`)
	boughtForPattern = regexp.MustCompile(`(?i)bought\s+(.+?)\s+for\s+(\d+(?:\.\d+)?)(?:\s+from\s+(.+))?$`)
	incomePattern    = regexp.MustCompile(`(?i)(?:got|received|earned|income)\s+(\d+(?:\.\d+)?)\s*(.+?)(?:\s+(?:in|to|into)\s+(.+))?$`)
	simpleExpense    = regexp.MustCompile(`(?i)(?:spent|paid|expense)\s+(\d+(?:\.\d+)?)\s+(.+?)(?:\s+from\s+(.+))?$`)
)

func (s *chatService) parseAddTransaction(lower, original string) *dto.ChatAction {
	// "spent 500 on groceries" / "paid 200 for coffee"
	if m := spentOnPattern.FindStringSubmatch(original); m != nil {
		amount, _ := strconv.ParseFloat(m[1], 64)
		return &dto.ChatAction{
			Type:     "expense",
			Amount:   amount,
			Name:     strings.TrimSpace(m[2]),
			Account:  strings.TrimSpace(m[3]),
			Date:     time.Now().Format("2006-01-02"),
		}
	}

	// "bought lunch for 150"
	if m := boughtForPattern.FindStringSubmatch(original); m != nil {
		amount, _ := strconv.ParseFloat(m[2], 64)
		return &dto.ChatAction{
			Type:     "expense",
			Amount:   amount,
			Name:     strings.TrimSpace(m[1]),
			Account:  strings.TrimSpace(m[3]),
			Date:     time.Now().Format("2006-01-02"),
		}
	}

	// "spent 500 groceries" (without "on/for")
	if m := simpleExpense.FindStringSubmatch(original); m != nil {
		amount, _ := strconv.ParseFloat(m[1], 64)
		return &dto.ChatAction{
			Type:     "expense",
			Amount:   amount,
			Name:     strings.TrimSpace(m[2]),
			Account:  strings.TrimSpace(m[3]),
			Date:     time.Now().Format("2006-01-02"),
		}
	}

	// "got 10000 salary" / "received 5000 freelance"
	if m := incomePattern.FindStringSubmatch(original); m != nil {
		amount, _ := strconv.ParseFloat(m[1], 64)
		return &dto.ChatAction{
			Type:    "income",
			Amount:  amount,
			Name:    strings.TrimSpace(m[2]),
			Account: strings.TrimSpace(m[3]),
			Date:    time.Now().Format("2006-01-02"),
		}
	}

	return nil
}

func (s *chatService) executeAction(ctx context.Context, userID uuid.UUID, action *dto.ChatAction) (*dto.ChatResponse, error) {
	switch action.Type {
	case "expense", "income":
		return s.createTransactionFromAction(ctx, userID, action)
	default:
		return &dto.ChatResponse{
			Reply:   "I understood your request but don't know how to handle that action type.",
			Success: false,
		}, nil
	}
}

func (s *chatService) createTransactionFromAction(ctx context.Context, userID uuid.UUID, action *dto.ChatAction) (*dto.ChatResponse, error) {
	// Find or use first account
	accounts, _, err := s.accountService.List(ctx, userID)
	if err != nil || len(accounts) == 0 {
		return &dto.ChatResponse{
			Reply:   "You need at least one account before adding transactions. Create one in the Accounts page first.",
			Success: false,
		}, nil
	}

	var targetAccount *model.Account
	if action.Account != "" {
		for i := range accounts {
			if strings.EqualFold(accounts[i].Name, strings.TrimSpace(action.Account)) {
				targetAccount = &accounts[i]
				break
			}
		}
	}
	if targetAccount == nil {
		targetAccount = &accounts[0]
	}

	// Find or use first category
	categories, _, err := s.categoryService.List(ctx, userID)
	if err != nil || len(categories) == 0 {
		return &dto.ChatResponse{
			Reply:   "You need at least one category before adding transactions. Create one in the Categories section first.",
			Success: false,
		}, nil
	}

	var targetCategory *model.Category
	if action.Category != "" {
		for i := range categories {
			if strings.EqualFold(categories[i].Name, strings.TrimSpace(action.Category)) {
				targetCategory = &categories[i]
				break
			}
		}
	}
	if targetCategory == nil {
		targetCategory = &categories[0]
	}

	amount := action.Amount
	if action.Type == "expense" {
		amount = -amount
	}

	txReq := dto.CreateTransactionRequest{
		AccountID:       targetAccount.ID,
		CategoryID:      targetCategory.ID,
		Name:            action.Name,
		Amount:          amount,
		TransactionDate: action.Date,
	}

	tx, err := s.transactionService.Create(ctx, userID, txReq)
	if err != nil {
		return &dto.ChatResponse{
			Reply:   "Failed to create the transaction. Please try again.",
			Success: false,
		}, nil
	}

	typeLabel := "Expense"
	if action.Type == "income" {
		typeLabel = "Income"
	}

	reply := fmt.Sprintf("Done! %s of **%.2f** for **%s** recorded to account **%s** (category: %s).",
		typeLabel, action.Amount, tx.Name, targetAccount.Name, targetCategory.Name)

	return &dto.ChatResponse{
		Reply:   reply,
		Action:  action,
		Success: true,
	}, nil
}

func (s *chatService) handleBalanceQuery(ctx context.Context, userID uuid.UUID, lower string) (*dto.ChatResponse, error) {
	accounts, _, err := s.accountService.List(ctx, userID)
	if err != nil || len(accounts) == 0 {
		return &dto.ChatResponse{
			Reply:   "No accounts found. Create one first.",
			Success: true,
		}, nil
	}

	// Check if user asked about specific account
	for _, acc := range accounts {
		if strings.Contains(lower, strings.ToLower(acc.Name)) {
			reply := fmt.Sprintf("**%s** balance: **%.2f**", acc.Name, acc.CurrentBalance)
			return &dto.ChatResponse{Reply: reply, Success: true}, nil
		}
	}

	// Show all balances
	var sb strings.Builder
	sb.WriteString("Account balances:\n")
	var total float64
	for _, acc := range accounts {
		sb.WriteString(fmt.Sprintf("• **%s**: %.2f\n", acc.Name, acc.CurrentBalance))
		total += acc.CurrentBalance
	}
	sb.WriteString(fmt.Sprintf("\n**Total**: %.2f", total))

	return &dto.ChatResponse{Reply: sb.String(), Success: true}, nil
}

func (s *chatService) handleRecentTransactions(ctx context.Context, userID uuid.UUID) (*dto.ChatResponse, error) {
	filters := dto.TransactionFilters{
		SortBy:  "transaction_date",
		SortDir: "desc",
		Page:    1,
		PerPage: 5,
	}
	txns, _, err := s.transactionService.List(ctx, userID, filters)
	if err != nil || len(txns) == 0 {
		return &dto.ChatResponse{
			Reply:   "No recent transactions found.",
			Success: true,
		}, nil
	}

	var sb strings.Builder
	sb.WriteString("Recent transactions:\n")
	for _, tx := range txns {
		typeEmoji := "🔴"
		if tx.Amount > 0 {
			typeEmoji = "🟢"
		}
		sb.WriteString(fmt.Sprintf("• %s **%s** — %.2f (%s)\n", typeEmoji, tx.Name, tx.Amount, tx.TransactionDate.Format("2006-01-02")))
	}

	return &dto.ChatResponse{Reply: sb.String(), Success: true}, nil
}

// Ollama integration

type ollamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type ollamaResponse struct {
	Response string `json:"response"`
}

func (s *chatService) parseWithOllama(ctx context.Context, message string) (*dto.ChatAction, error) {
	prompt := fmt.Sprintf(`You are a financial assistant that parses user messages into structured transaction data.
Given the user message, extract one of these actions:
1. "expense" - user wants to record an expense
2. "income" - user wants to record income
3. "query" - user wants to ask about their finances

Respond ONLY with valid JSON (no markdown, no explanation):
{"type": "expense|income|query", "name": "description", "amount": 123.45, "account": "optional account name", "category": "optional category"}

If you cannot parse it, respond with: {"type": "unknown"}

User message: %s`, message)

	body, _ := json.Marshal(ollamaRequest{
		Model:  "llama3.2",
		Prompt: prompt,
		Stream: false,
	})

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.ollamaURL+"/api/generate", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ollamaResp ollamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return nil, err
	}

	// Parse the JSON from Ollama's response
	responseText := strings.TrimSpace(ollamaResp.Response)
	// Try to extract JSON from the response
	start := strings.Index(responseText, "{")
	end := strings.LastIndex(responseText, "}")
	if start < 0 || end < 0 || end <= start {
		return nil, fmt.Errorf("no JSON found in ollama response")
	}
	jsonStr := responseText[start : end+1]

	var action dto.ChatAction
	if err := json.Unmarshal([]byte(jsonStr), &action); err != nil {
		return nil, err
	}

	if action.Type == "unknown" || action.Type == "query" {
		return nil, fmt.Errorf("ollama returned non-actionable type: %s", action.Type)
	}

	if action.Date == "" {
		action.Date = time.Now().Format("2006-01-02")
	}

	return &action, nil
}
