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

	// Build context for LLM: accounts & categories
	accounts, _, _ := s.accountService.List(ctx, userID)
	categories, _, _ := s.categoryService.List(ctx, userID)

	// Try Ollama-based parsing first if configured
	if s.ollamaURL != "" {
		parsed, err := s.parseWithOllama(ctx, msg, accounts, categories)
		if err == nil && parsed != nil {
			return s.executeAction(ctx, userID, parsed, accounts, categories)
		}
		log.Debug().Err(err).Msg("ollama parse failed, falling back to pattern matching")
	}

	// Pattern matching fallback
	if action := s.parseAddTransaction(lower, msg); action != nil {
		return s.executeAction(ctx, userID, action, accounts, categories)
	}

	if strings.Contains(lower, "transfer") {
		return s.handleTransferQuery(ctx, userID, lower, msg, accounts, categories)
	}

	if strings.Contains(lower, "balance") || strings.Contains(lower, "how much") {
		return s.handleBalanceQuery(ctx, userID, lower, accounts)
	}

	if strings.Contains(lower, "recent") || strings.Contains(lower, "last") || strings.Contains(lower, "transactions") {
		return s.handleRecentTransactions(ctx, userID)
	}

	if lower == "help" || strings.Contains(lower, "what can you") || strings.Contains(lower, "how to") {
		return &dto.ChatResponse{
			Reply:   "I can help you with:\n• **Add expenses**: \"Spent 500 on groceries from HDFC\"\n• **Add income**: \"Received 10000 salary in Savings\"\n• **Transfer money**: \"Transfer 5000 from HDFC to SBI\"\n• **Check balances**: \"What's my balance?\" or \"Balance of Savings\"\n• **Recent transactions**: \"Show recent transactions\"\n• **Delete transactions**: \"Delete the last transaction\"\n\nTry natural phrases — I'll figure out the rest!",
			Success: true,
		}, nil
	}

	return &dto.ChatResponse{
		Reply:   "I didn't understand that. Try something like:\n• \"Spent 500 on groceries\"\n• \"Got 10000 salary\"\n• \"Transfer 2000 from HDFC to SBI\"\n• \"What's my balance?\"\n\nType **help** for more examples.",
		Success: true,
	}, nil
}

// ─── Pattern matching ───

var (
	spentOnPattern   = regexp.MustCompile(`(?i)(?:spent|paid|expense)\s+(\d+(?:\.\d+)?)\s+(?:on|for)\s+(.+?)(?:\s+from\s+(.+))?$`)
	boughtForPattern = regexp.MustCompile(`(?i)bought\s+(.+?)\s+for\s+(\d+(?:\.\d+)?)(?:\s+from\s+(.+))?$`)
	incomePattern    = regexp.MustCompile(`(?i)(?:got|received|earned|income)\s+(\d+(?:\.\d+)?)\s*(.+?)(?:\s+(?:in|to|into)\s+(.+))?$`)
	simpleExpense    = regexp.MustCompile(`(?i)(?:spent|paid|expense)\s+(\d+(?:\.\d+)?)\s+(.+?)(?:\s+from\s+(.+))?$`)
	transferPattern  = regexp.MustCompile(`(?i)transfer\s+(\d+(?:\.\d+)?)\s+from\s+(.+?)\s+to\s+(.+?)$`)
)

func (s *chatService) parseAddTransaction(lower, original string) *dto.ChatAction {
	if m := spentOnPattern.FindStringSubmatch(original); m != nil {
		amount, _ := strconv.ParseFloat(m[1], 64)
		return &dto.ChatAction{Type: "expense", Amount: amount, Name: strings.TrimSpace(m[2]), Account: strings.TrimSpace(m[3]), Date: time.Now().Format("2006-01-02")}
	}
	if m := boughtForPattern.FindStringSubmatch(original); m != nil {
		amount, _ := strconv.ParseFloat(m[2], 64)
		return &dto.ChatAction{Type: "expense", Amount: amount, Name: strings.TrimSpace(m[1]), Account: strings.TrimSpace(m[3]), Date: time.Now().Format("2006-01-02")}
	}
	if m := simpleExpense.FindStringSubmatch(original); m != nil {
		amount, _ := strconv.ParseFloat(m[1], 64)
		return &dto.ChatAction{Type: "expense", Amount: amount, Name: strings.TrimSpace(m[2]), Account: strings.TrimSpace(m[3]), Date: time.Now().Format("2006-01-02")}
	}
	if m := incomePattern.FindStringSubmatch(original); m != nil {
		amount, _ := strconv.ParseFloat(m[1], 64)
		return &dto.ChatAction{Type: "income", Amount: amount, Name: strings.TrimSpace(m[2]), Account: strings.TrimSpace(m[3]), Date: time.Now().Format("2006-01-02")}
	}
	return nil
}

// ─── Action execution ───

func (s *chatService) executeAction(ctx context.Context, userID uuid.UUID, action *dto.ChatAction, accounts []model.Account, categories []model.Category) (*dto.ChatResponse, error) {
	switch action.Type {
	case "expense", "income":
		return s.createTransactionFromAction(ctx, userID, action, accounts, categories)
	case "transfer":
		return s.executeTransfer(ctx, userID, action, accounts, categories)
	case "delete":
		return s.executeDelete(ctx, userID, action)
	case "balance":
		return s.handleBalanceQuery(ctx, userID, strings.ToLower(action.Account), accounts)
	case "recent":
		return s.handleRecentTransactions(ctx, userID)
	default:
		return &dto.ChatResponse{Reply: "I understood your request but don't know how to handle that action type.", Success: false}, nil
	}
}

func (s *chatService) createTransactionFromAction(ctx context.Context, userID uuid.UUID, action *dto.ChatAction, accounts []model.Account, categories []model.Category) (*dto.ChatResponse, error) {
	if len(accounts) == 0 {
		return &dto.ChatResponse{Reply: "You need at least one account before adding transactions.", Success: false}, nil
	}
	if len(categories) == 0 {
		return &dto.ChatResponse{Reply: "You need at least one category before adding transactions.", Success: false}, nil
	}

	targetAccount := findAccountByName(accounts, action.Account)
	if targetAccount == nil {
		targetAccount = &accounts[0]
	}

	targetCategory := findCategoryByName(categories, action.Category)
	if targetCategory == nil {
		targetCategory = &categories[0]
	}

	amount := action.Amount
	if action.Type == "expense" {
		amount = -amount
	}

	date := action.Date
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	txReq := dto.CreateTransactionRequest{
		AccountID:       targetAccount.ID,
		CategoryID:      targetCategory.ID,
		Name:            action.Name,
		Amount:          amount,
		TransactionDate: date,
	}

	tx, err := s.transactionService.Create(ctx, userID, txReq)
	if err != nil {
		return &dto.ChatResponse{Reply: "Failed to create the transaction. Please try again.", Success: false}, nil
	}

	typeLabel := "Expense"
	if action.Type == "income" {
		typeLabel = "Income"
	}

	reply := fmt.Sprintf("✅ %s of **₹%.2f** for **%s** recorded to account **%s** (category: %s).",
		typeLabel, action.Amount, tx.Name, targetAccount.Name, targetCategory.Name)

	return &dto.ChatResponse{Reply: reply, Action: action, Success: true}, nil
}

func (s *chatService) executeTransfer(ctx context.Context, userID uuid.UUID, action *dto.ChatAction, accounts []model.Account, categories []model.Category) (*dto.ChatResponse, error) {
	if len(accounts) < 2 {
		return &dto.ChatResponse{Reply: "You need at least two accounts for a transfer.", Success: false}, nil
	}
	if len(categories) == 0 {
		return &dto.ChatResponse{Reply: "You need at least one category for a transfer.", Success: false}, nil
	}

	fromAccount := findAccountByName(accounts, action.Account)
	toAccount := findAccountByName(accounts, action.ToAccount)

	if fromAccount == nil || toAccount == nil || fromAccount.ID == toAccount.ID {
		return &dto.ChatResponse{Reply: "I couldn't identify two different accounts for the transfer. Please specify the account names.", Success: false}, nil
	}

	targetCategory := findCategoryByName(categories, action.Category)
	if targetCategory == nil {
		targetCategory = &categories[0]
	}

	date := action.Date
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	name := action.Name
	if name == "" {
		name = fmt.Sprintf("Transfer: %s → %s", fromAccount.Name, toAccount.Name)
	}

	transferReq := dto.TransferRequest{
		FromAccountID:   fromAccount.ID,
		ToAccountID:     toAccount.ID,
		CategoryID:      targetCategory.ID,
		Amount:          action.Amount,
		TransactionDate: date,
		Name:            name,
	}

	_, _, err := s.transactionService.Transfer(ctx, userID, transferReq)
	if err != nil {
		return &dto.ChatResponse{Reply: "Failed to process the transfer. Please try again.", Success: false}, nil
	}

	reply := fmt.Sprintf("✅ Transferred **₹%.2f** from **%s** to **%s**.", action.Amount, fromAccount.Name, toAccount.Name)
	return &dto.ChatResponse{Reply: reply, Action: action, Success: true}, nil
}

func (s *chatService) executeDelete(ctx context.Context, userID uuid.UUID, action *dto.ChatAction) (*dto.ChatResponse, error) {
	if action.TransactionID == "" {
		// Delete the most recent transaction
		filters := dto.TransactionFilters{SortBy: "created_at", SortDir: "desc", Page: 1, PerPage: 1}
		txns, _, err := s.transactionService.List(ctx, userID, filters)
		if err != nil || len(txns) == 0 {
			return &dto.ChatResponse{Reply: "No transactions found to delete.", Success: false}, nil
		}
		action.TransactionID = txns[0].ID.String()
	}

	id, err := uuid.Parse(action.TransactionID)
	if err != nil {
		return &dto.ChatResponse{Reply: "Invalid transaction ID.", Success: false}, nil
	}

	if err := s.transactionService.Delete(ctx, userID, id); err != nil {
		return &dto.ChatResponse{Reply: "Failed to delete the transaction.", Success: false}, nil
	}

	return &dto.ChatResponse{Reply: "✅ Transaction deleted successfully.", Action: action, Success: true}, nil
}

func (s *chatService) handleTransferQuery(ctx context.Context, userID uuid.UUID, lower, original string, accounts []model.Account, categories []model.Category) (*dto.ChatResponse, error) {
	if m := transferPattern.FindStringSubmatch(original); m != nil {
		amount, _ := strconv.ParseFloat(m[1], 64)
		action := &dto.ChatAction{
			Type:      "transfer",
			Amount:    amount,
			Account:   strings.TrimSpace(m[2]),
			ToAccount: strings.TrimSpace(m[3]),
			Date:      time.Now().Format("2006-01-02"),
		}
		return s.executeTransfer(ctx, userID, action, accounts, categories)
	}
	return &dto.ChatResponse{Reply: "To transfer, say something like: \"Transfer 5000 from HDFC to SBI\"", Success: true}, nil
}

func (s *chatService) handleBalanceQuery(ctx context.Context, userID uuid.UUID, lower string, accounts []model.Account) (*dto.ChatResponse, error) {
	if len(accounts) == 0 {
		return &dto.ChatResponse{Reply: "No accounts found. Create one first.", Success: true}, nil
	}

	for _, acc := range accounts {
		if strings.Contains(lower, strings.ToLower(acc.Name)) {
			reply := fmt.Sprintf("**%s** balance: **₹%.2f**", acc.Name, acc.CurrentBalance)
			return &dto.ChatResponse{Reply: reply, Success: true}, nil
		}
	}

	var sb strings.Builder
	sb.WriteString("💰 Account balances:\n")
	var total float64
	for _, acc := range accounts {
		sb.WriteString(fmt.Sprintf("• **%s**: ₹%.2f\n", acc.Name, acc.CurrentBalance))
		total += acc.CurrentBalance
	}
	sb.WriteString(fmt.Sprintf("\n**Total**: ₹%.2f", total))

	return &dto.ChatResponse{Reply: sb.String(), Success: true}, nil
}

func (s *chatService) handleRecentTransactions(ctx context.Context, userID uuid.UUID) (*dto.ChatResponse, error) {
	filters := dto.TransactionFilters{SortBy: "transaction_date", SortDir: "desc", Page: 1, PerPage: 5}
	txns, _, err := s.transactionService.List(ctx, userID, filters)
	if err != nil || len(txns) == 0 {
		return &dto.ChatResponse{Reply: "No recent transactions found.", Success: true}, nil
	}

	var sb strings.Builder
	sb.WriteString("📋 Recent transactions:\n")
	for _, tx := range txns {
		emoji := "🔴"
		if tx.Amount > 0 {
			emoji = "🟢"
		}
		sb.WriteString(fmt.Sprintf("• %s **%s** — ₹%.2f (%s)\n", emoji, tx.Name, tx.Amount, tx.TransactionDate.Format("2006-01-02")))
	}

	return &dto.ChatResponse{Reply: sb.String(), Success: true}, nil
}

// ─── Ollama Integration ───

type ollamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type ollamaResponse struct {
	Response string `json:"response"`
}

func (s *chatService) parseWithOllama(ctx context.Context, message string, accounts []model.Account, categories []model.Category) (*dto.ChatAction, error) {
	// Build account/category context for the LLM
	var accountNames []string
	for _, a := range accounts {
		accountNames = append(accountNames, fmt.Sprintf("%s (balance: %.2f)", a.Name, a.CurrentBalance))
	}
	var categoryNames []string
	for _, c := range categories {
		categoryNames = append(categoryNames, c.Name)
	}

	prompt := fmt.Sprintf(`You are a financial assistant for a personal expense tracker app. Parse the user's message into a structured JSON action.

The user has these accounts: %s
The user has these categories: %s
Today's date is: %s

Supported action types:
- "expense": Record an expense (negative transaction). Fields: type, name, amount (positive number), account (optional), category (optional), date (YYYY-MM-DD, default today)
- "income": Record income (positive transaction). Fields: type, name, amount (positive number), account (optional), category (optional), date (YYYY-MM-DD, default today)  
- "transfer": Transfer money between accounts. Fields: type, amount (positive), account (source), to_account (destination), name (optional), date
- "delete": Delete a transaction. Fields: type, transaction_id (if known, otherwise empty to delete the most recent)
- "balance": Check account balance. Fields: type, account (optional, empty = all accounts)
- "recent": Show recent transactions. Fields: type

Match account and category names closely to the user's existing ones (case-insensitive).

Respond ONLY with valid JSON, no markdown, no explanation:
{"type": "expense|income|transfer|delete|balance|recent", "name": "...", "amount": 0, "account": "...", "to_account": "...", "category": "...", "date": "YYYY-MM-DD", "transaction_id": "..."}

If you cannot parse the message, respond: {"type": "unknown"}

User message: %s`, strings.Join(accountNames, ", "), strings.Join(categoryNames, ", "), time.Now().Format("2006-01-02"), message)

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

	responseText := strings.TrimSpace(ollamaResp.Response)
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

	if action.Type == "unknown" {
		return nil, fmt.Errorf("ollama returned unknown type")
	}

	if action.Date == "" {
		action.Date = time.Now().Format("2006-01-02")
	}

	return &action, nil
}

// ─── Helpers ───

func findAccountByName(accounts []model.Account, name string) *model.Account {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil
	}
	for i := range accounts {
		if strings.EqualFold(accounts[i].Name, name) {
			return &accounts[i]
		}
	}
	// Partial match fallback
	lower := strings.ToLower(name)
	for i := range accounts {
		if strings.Contains(strings.ToLower(accounts[i].Name), lower) {
			return &accounts[i]
		}
	}
	return nil
}

func findCategoryByName(categories []model.Category, name string) *model.Category {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil
	}
	for i := range categories {
		if strings.EqualFold(categories[i].Name, name) {
			return &categories[i]
		}
	}
	lower := strings.ToLower(name)
	for i := range categories {
		if strings.Contains(strings.ToLower(categories[i].Name), lower) {
			return &categories[i]
		}
	}
	return nil
}
