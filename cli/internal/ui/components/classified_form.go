package components

import (
	"fmt"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/wesellis/terminal-news/cli/internal/models"
	"github.com/wesellis/terminal-news/cli/internal/ui/styles"
)

// ClassifiedFormMode represents form state
type ClassifiedFormMode int

const (
	ClassifiedModeEdit ClassifiedFormMode = iota
	ClassifiedModePreview
	ClassifiedModeSubmitting
)

// ClassifiedForm is a form for posting classifieds
type ClassifiedForm struct {
	mode            ClassifiedFormMode
	activeField     int
	width           int
	height          int
	styles          *styles.Styles

	// Input fields
	titleInput      textinput.Model
	descriptionArea textarea.Model
	priceInput      textinput.Model
	locationInput   textinput.Model
	emailInput      textinput.Model
	phoneInput      textinput.Model

	// Selection fields
	category        string
	categories      []string
	contactMethod   string
	contactMethods  []string
	isPremium       bool

	// State
	errorMessage    string
	charCounts      map[string]int
}

// NewClassifiedForm creates a new classified form
func NewClassifiedForm(styles *styles.Styles) *ClassifiedForm {
	// Title input
	titleInput := textinput.New()
	titleInput.Placeholder = "MacBook Pro M3 - Excellent Condition"
	titleInput.CharLimit = 200
	titleInput.Width = 60
	titleInput.Focus()

	// Description textarea
	descriptionArea := textarea.New()
	descriptionArea.Placeholder = "Describe your item in detail..."
	descriptionArea.CharLimit = 1000
	descriptionArea.SetWidth(60)
	descriptionArea.SetHeight(6)

	// Price input
	priceInput := textinput.New()
	priceInput.Placeholder = "2000"
	priceInput.CharLimit = 10
	priceInput.Width = 20

	// Location input
	locationInput := textinput.New()
	locationInput.Placeholder = "San Francisco, CA"
	locationInput.CharLimit = 100
	locationInput.Width = 40

	// Email input
	emailInput := textinput.New()
	emailInput.Placeholder = "your@email.com"
	emailInput.CharLimit = 100
	emailInput.Width = 40

	// Phone input
	phoneInput := textinput.New()
	phoneInput.Placeholder = "(555) 123-4567"
	phoneInput.CharLimit = 20
	phoneInput.Width = 30

	return &ClassifiedForm{
		mode:            ClassifiedModeEdit,
		activeField:     0,
		styles:          styles,
		titleInput:      titleInput,
		descriptionArea: descriptionArea,
		priceInput:      priceInput,
		locationInput:   locationInput,
		emailInput:      emailInput,
		phoneInput:      phoneInput,
		categories:      []string{"For Sale", "Jobs", "Housing", "Services", "Events", "Gigs"},
		category:        "For Sale",
		contactMethods:  []string{"Email", "Phone", "Direct Message"},
		contactMethod:   "Email",
		isPremium:       false,
		charCounts:      make(map[string]int),
	}
}

// SetSize updates dimensions
func (cf *ClassifiedForm) SetSize(width, height int) {
	cf.width = width
	cf.height = height
}

// Update handles input
func (cf *ClassifiedForm) Update(msg tea.Msg) (*ClassifiedForm, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch cf.mode {
		case ClassifiedModeEdit:
			switch msg.String() {
			case "tab":
				cf.nextField()
				return cf, nil

			case "shift+tab":
				cf.prevField()
				return cf, nil

			case "ctrl+p":
				// Preview before posting
				cf.mode = ClassifiedModePreview
				return cf, nil

			case "ctrl+s":
				// Submit
				return cf, cf.submit()

			case "ctrl+c":
				// Toggle category
				cf.nextCategory()
				return cf, nil

			case "ctrl+m":
				// Toggle contact method
				cf.nextContactMethod()
				return cf, nil

			case "ctrl+r":
				// Toggle premium
				cf.isPremium = !cf.isPremium
				return cf, nil

			case "esc":
				// Cancel
				return cf, cf.cancel()
			}

		case ClassifiedModePreview:
			switch msg.String() {
			case "e":
				// Back to edit
				cf.mode = ClassifiedModeEdit
				return cf, nil

			case "enter", "ctrl+s":
				// Confirm and submit
				return cf, cf.submit()

			case "esc":
				// Back to edit
				cf.mode = ClassifiedModeEdit
				return cf, nil
			}
		}
	}

	// Update active field
	var cmd tea.Cmd
	switch cf.activeField {
	case 0:
		cf.titleInput, cmd = cf.titleInput.Update(msg)
		cmds = append(cmds, cmd)
	case 1:
		cf.descriptionArea, cmd = cf.descriptionArea.Update(msg)
		cmds = append(cmds, cmd)
	case 2:
		cf.priceInput, cmd = cf.priceInput.Update(msg)
		cmds = append(cmds, cmd)
	case 3:
		cf.locationInput, cmd = cf.locationInput.Update(msg)
		cmds = append(cmds, cmd)
	case 4:
		if cf.contactMethod == "Email" {
			cf.emailInput, cmd = cf.emailInput.Update(msg)
			cmds = append(cmds, cmd)
		} else if cf.contactMethod == "Phone" {
			cf.phoneInput, cmd = cf.phoneInput.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	// Update character counts
	cf.charCounts["title"] = len(cf.titleInput.Value())
	cf.charCounts["description"] = len(cf.descriptionArea.Value())

	return cf, tea.Batch(cmds...)
}

// nextField moves to next field
func (cf *ClassifiedForm) nextField() {
	cf.blurAll()
	cf.activeField++
	if cf.activeField > 4 {
		cf.activeField = 0
	}
	cf.focusActive()
}

// prevField moves to previous field
func (cf *ClassifiedForm) prevField() {
	cf.blurAll()
	cf.activeField--
	if cf.activeField < 0 {
		cf.activeField = 4
	}
	cf.focusActive()
}

// blurAll unfocuses all fields
func (cf *ClassifiedForm) blurAll() {
	cf.titleInput.Blur()
	cf.descriptionArea.Blur()
	cf.priceInput.Blur()
	cf.locationInput.Blur()
	cf.emailInput.Blur()
	cf.phoneInput.Blur()
}

// focusActive focuses active field
func (cf *ClassifiedForm) focusActive() {
	switch cf.activeField {
	case 0:
		cf.titleInput.Focus()
	case 1:
		cf.descriptionArea.Focus()
	case 2:
		cf.priceInput.Focus()
	case 3:
		cf.locationInput.Focus()
	case 4:
		if cf.contactMethod == "Email" {
			cf.emailInput.Focus()
		} else if cf.contactMethod == "Phone" {
			cf.phoneInput.Focus()
		}
	}
}

// nextCategory cycles to next category
func (cf *ClassifiedForm) nextCategory() {
	for i, cat := range cf.categories {
		if cat == cf.category {
			cf.category = cf.categories[(i+1)%len(cf.categories)]
			return
		}
	}
}

// nextContactMethod cycles to next contact method
func (cf *ClassifiedForm) nextContactMethod() {
	for i, method := range cf.contactMethods {
		if method == cf.contactMethod {
			cf.contactMethod = cf.contactMethods[(i+1)%len(cf.contactMethods)]
			return
		}
	}
}

// validate validates form data
func (cf *ClassifiedForm) validate() error {
	title := strings.TrimSpace(cf.titleInput.Value())
	description := strings.TrimSpace(cf.descriptionArea.Value())
	location := strings.TrimSpace(cf.locationInput.Value())

	if title == "" {
		return fmt.Errorf("title is required")
	}

	if len(title) < 10 {
		return fmt.Errorf("title must be at least 10 characters")
	}

	if description == "" {
		return fmt.Errorf("description is required")
	}

	if len(description) < 20 {
		return fmt.Errorf("description must be at least 20 characters")
	}

	if location == "" {
		return fmt.Errorf("location is required")
	}

	if cf.contactMethod == "Email" {
		email := strings.TrimSpace(cf.emailInput.Value())
		if email == "" || !strings.Contains(email, "@") {
			return fmt.Errorf("valid email is required")
		}
	} else if cf.contactMethod == "Phone" {
		phone := strings.TrimSpace(cf.phoneInput.Value())
		if phone == "" {
			return fmt.Errorf("phone number is required")
		}
	}

	return nil
}

// submit submits the form
func (cf *ClassifiedForm) submit() tea.Cmd {
	if err := cf.validate(); err != nil {
		cf.errorMessage = err.Error()
		return nil
	}

	cf.mode = ClassifiedModeSubmitting

	return func() tea.Msg {
		// Create classified object
		priceVal, _ := strconv.ParseFloat(cf.priceInput.Value(), 64)
		price := &priceVal

		classified := &models.Classified{
			Title:        strings.TrimSpace(cf.titleInput.Value()),
			Description:  strings.TrimSpace(cf.descriptionArea.Value()),
			Price:        price,
			Category:     cf.category,
			City:         cf.locationInput.Value(), // Parse city/state properly
			ContactEmail: cf.emailInput.Value(),
			ContactPhone: cf.phoneInput.Value(),
			IsPremium:    cf.isPremium,
		}

		return ClassifiedSubmitMsg{Classified: classified}
	}
}

// cancel cancels the form
func (cf *ClassifiedForm) cancel() tea.Cmd {
	return func() tea.Msg {
		return ClassifiedCancelMsg{}
	}
}

// View renders the form
func (cf *ClassifiedForm) View() string {
	switch cf.mode {
	case ClassifiedModeEdit:
		return cf.renderEditMode()
	case ClassifiedModePreview:
		return cf.renderPreview()
	case ClassifiedModeSubmitting:
		return cf.renderSubmitting()
	}
	return ""
}

// renderEditMode renders edit form
func (cf *ClassifiedForm) renderEditMode() string {
	var sections []string

	// Header
	header := cf.styles.Header.Render("POST CLASSIFIED")
	sections = append(sections, header, "")

	// Category selector
	categoryLine := fmt.Sprintf("Category: [%s ▾]  (Ctrl+C to change)", cf.category)
	sections = append(sections, cf.renderField("", categoryLine, 0, false))

	// Title field
	titleLabel := fmt.Sprintf("Title: (%d/200 characters)", cf.charCounts["title"])
	sections = append(sections, cf.renderField(titleLabel, cf.titleInput.View(), 0, cf.activeField == 0))

	// Description field
	descLabel := fmt.Sprintf("Description: (%d/1000 characters)", cf.charCounts["description"])
	sections = append(sections, cf.renderField(descLabel, cf.descriptionArea.View(), 1, cf.activeField == 1))

	// Price field
	sections = append(sections, cf.renderField("Price: $", cf.priceInput.View(), 2, cf.activeField == 2))

	// Location field
	sections = append(sections, cf.renderField("Location:", cf.locationInput.View(), 3, cf.activeField == 3))

	// Contact method
	contactLine := fmt.Sprintf("Contact Method: [%s ▾]  (Ctrl+M to change)", cf.contactMethod)
	sections = append(sections, cf.renderField("", contactLine, 0, false))

	// Contact field based on method
	if cf.contactMethod == "Email" {
		sections = append(sections, cf.renderField("Email:", cf.emailInput.View(), 4, cf.activeField == 4))
	} else if cf.contactMethod == "Phone" {
		sections = append(sections, cf.renderField("Phone:", cf.phoneInput.View(), 4, cf.activeField == 4))
	} else {
		sections = append(sections, cf.renderField("Contact:", "Via Direct Message", 4, false))
	}

	// Premium option
	premiumBox := " "
	if cf.isPremium {
		premiumBox = "✓"
	}
	premiumSection := fmt.Sprintf(`
┌─ Premium Listing ($10) ────────────────────────────────────────────┐
│ [%s] Make this a premium listing  (Ctrl+R to toggle)               │
│     • Highlighted with ⭐ in feed                                   │
│     • 3x visibility                                                 │
│     • 60-day duration (vs 30)                                       │
│     • Featured placement                                            │
└─────────────────────────────────────────────────────────────────────┘
	`, premiumBox)
	sections = append(sections, premiumSection)

	// Error message
	if cf.errorMessage != "" {
		sections = append(sections, "", cf.styles.ErrorMessage.Render("✗ "+cf.errorMessage))
	}

	// Help text
	help := cf.styles.HelpText.Render("[Tab] Next  [Shift+Tab] Prev  [Ctrl+P] Preview  [Ctrl+S] Submit  [Esc] Cancel")
	sections = append(sections, "", help)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// renderField renders a form field
func (cf *ClassifiedForm) renderField(label, value string, fieldNum int, isActive bool) string {
	if label != "" {
		label = cf.styles.Label.Render(label)
	}

	var valueStyle lipgloss.Style
	if isActive {
		valueStyle = cf.styles.ActiveInput
	} else {
		valueStyle = cf.styles.Input
	}

	if value != "" {
		value = valueStyle.Render(value)
	}

	if label != "" {
		return lipgloss.JoinVertical(lipgloss.Left, label, value, "")
	}
	return value + "\n"
}

// renderPreview renders preview mode
func (cf *ClassifiedForm) renderPreview() string {
	preview := fmt.Sprintf(`
┌─ PREVIEW ───────────────────────────────────────────────────────────┐
│                                                                      │
│  %s                                                                  │
│                                                                      │
│  Category: %s                                                        │
│  Price: $%s                                                          │
│  Location: %s                                                        │
│                                                                      │
│  %s                                                                  │
│                                                                      │
│  Contact: %s via %s                                                  │
│                                                                      │
└──────────────────────────────────────────────────────────────────────┘

[Enter] Submit  [E] Edit  [Esc] Cancel
	`,
		cf.titleInput.Value(),
		cf.category,
		cf.priceInput.Value(),
		cf.locationInput.Value(),
		cf.descriptionArea.Value(),
		cf.getContactValue(),
		cf.contactMethod,
	)

	return preview
}

// renderSubmitting renders submitting state
func (cf *ClassifiedForm) renderSubmitting() string {
	return cf.styles.InfoMessage.Render("\n  ⠋ Posting classified...\n")
}

// getContactValue returns the contact value based on method
func (cf *ClassifiedForm) getContactValue() string {
	switch cf.contactMethod {
	case "Email":
		return cf.emailInput.Value()
	case "Phone":
		return cf.phoneInput.Value()
	default:
		return "Direct Message"
	}
}

// Message types
type ClassifiedSubmitMsg struct {
	Classified *models.Classified
}

type ClassifiedCancelMsg struct{}
