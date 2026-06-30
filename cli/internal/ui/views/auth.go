package views

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/wesellis/terminal-news/cli/internal/api"
	"github.com/wesellis/terminal-news/cli/internal/cache"
	"github.com/wesellis/terminal-news/cli/internal/models"
	"github.com/wesellis/terminal-news/cli/internal/ui/styles"
)

// AuthMode represents the current auth screen
type AuthMode int

const (
	AuthModeLogin AuthMode = iota
	AuthModeRegister
	AuthModeLoading
	AuthModeSuccess
)

// AuthView handles login and registration
type AuthView struct {
	BaseView
	mode            AuthMode
	activeField     int
	usernameInput   textinput.Model
	emailInput      textinput.Model
	passwordInput   textinput.Model
	password2Input  textinput.Model
	errorMessage    string
	successMessage  string
	onAuthSuccess   func(*models.User, string) // Callback with user and token
}

// NewAuthView creates a new auth view
func NewAuthView(apiClient *api.Client, cache *cache.Cache, styles *styles.Styles) *AuthView {
	// Username input
	usernameInput := textinput.New()
	usernameInput.Placeholder = "username"
	usernameInput.CharLimit = 50
	usernameInput.Width = 40

	// Email input
	emailInput := textinput.New()
	emailInput.Placeholder = "email@example.com"
	emailInput.CharLimit = 100
	emailInput.Width = 40

	// Password input
	passwordInput := textinput.New()
	passwordInput.Placeholder = "password"
	passwordInput.EchoMode = textinput.EchoPassword
	passwordInput.EchoCharacter = '•'
	passwordInput.CharLimit = 100
	passwordInput.Width = 40

	// Confirm password input
	password2Input := textinput.New()
	password2Input.Placeholder = "confirm password"
	password2Input.EchoMode = textinput.EchoPassword
	password2Input.EchoCharacter = '•'
	password2Input.CharLimit = 100
	password2Input.Width = 40

	// Focus first field
	usernameInput.Focus()

	return &AuthView{
		BaseView: BaseView{
			apiClient: apiClient,
			cache:     cache,
			styles:    styles,
		},
		mode:           AuthModeLogin,
		activeField:    0,
		usernameInput:  usernameInput,
		emailInput:     emailInput,
		passwordInput:  passwordInput,
		password2Input: password2Input,
	}
}

// SetOnAuthSuccess sets the callback for successful auth
func (av *AuthView) SetOnAuthSuccess(callback func(*models.User, string)) {
	av.onAuthSuccess = callback
}

// Update handles messages
func (av *AuthView) Update(msg tea.Msg) (*AuthView, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch av.mode {
		case AuthModeLogin, AuthModeRegister:
			switch msg.String() {
			case "tab", "down":
				av.nextField()
				return av, nil

			case "shift+tab", "up":
				av.prevField()
				return av, nil

			case "enter":
				if av.mode == AuthModeLogin {
					return av, av.doLogin()
				} else {
					return av, av.doRegister()
				}

			case "ctrl+t":
				// Toggle between login and register
				av.toggleMode()
				return av, nil

			case "esc":
				// Go back or clear
				if av.errorMessage != "" {
					av.errorMessage = ""
				}
				return av, nil
			}

		case AuthModeSuccess:
			// Any key after success returns to app
			return av, nil
		}

	case LoginSuccessMsg:
		av.mode = AuthModeSuccess
		av.successMessage = fmt.Sprintf("Welcome back, @%s!", msg.User.Username)
		if av.onAuthSuccess != nil {
			av.onAuthSuccess(&msg.User, msg.Token)
		}
		return av, nil

	case RegisterSuccessMsg:
		av.mode = AuthModeSuccess
		av.successMessage = fmt.Sprintf("Account created! Welcome, @%s!", msg.Username)
		// Auto-login after registration
		return av, av.doLogin()

	case ErrorMsg:
		av.mode = AuthModeLogin // Back to form
		av.errorMessage = string(msg)
		av.loading = false
		return av, nil
	}

	// Update active input field
	var cmd tea.Cmd
	switch av.activeField {
	case 0:
		av.usernameInput, cmd = av.usernameInput.Update(msg)
		cmds = append(cmds, cmd)
	case 1:
		if av.mode == AuthModeRegister {
			av.emailInput, cmd = av.emailInput.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			av.passwordInput, cmd = av.passwordInput.Update(msg)
			cmds = append(cmds, cmd)
		}
	case 2:
		if av.mode == AuthModeRegister {
			av.passwordInput, cmd = av.passwordInput.Update(msg)
			cmds = append(cmds, cmd)
		}
	case 3:
		if av.mode == AuthModeRegister {
			av.password2Input, cmd = av.password2Input.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return av, tea.Batch(cmds...)
}

// nextField moves to next input field
func (av *AuthView) nextField() {
	av.blurAll()

	maxField := 1
	if av.mode == AuthModeRegister {
		maxField = 3
	}

	av.activeField++
	if av.activeField > maxField {
		av.activeField = 0
	}

	av.focusActive()
}

// prevField moves to previous input field
func (av *AuthView) prevField() {
	av.blurAll()

	maxField := 1
	if av.mode == AuthModeRegister {
		maxField = 3
	}

	av.activeField--
	if av.activeField < 0 {
		av.activeField = maxField
	}

	av.focusActive()
}

// blurAll unfocuses all inputs
func (av *AuthView) blurAll() {
	av.usernameInput.Blur()
	av.emailInput.Blur()
	av.passwordInput.Blur()
	av.password2Input.Blur()
}

// focusActive focuses the active field
func (av *AuthView) focusActive() {
	switch av.activeField {
	case 0:
		av.usernameInput.Focus()
	case 1:
		if av.mode == AuthModeRegister {
			av.emailInput.Focus()
		} else {
			av.passwordInput.Focus()
		}
	case 2:
		av.passwordInput.Focus()
	case 3:
		av.password2Input.Focus()
	}
}

// toggleMode switches between login and register
func (av *AuthView) toggleMode() {
	if av.mode == AuthModeLogin {
		av.mode = AuthModeRegister
	} else {
		av.mode = AuthModeLogin
	}

	av.activeField = 0
	av.errorMessage = ""
	av.blurAll()
	av.usernameInput.Focus()
}

// doLogin performs login
func (av *AuthView) doLogin() tea.Cmd {
	username := strings.TrimSpace(av.usernameInput.Value())
	password := av.passwordInput.Value()

	if username == "" || password == "" {
		av.errorMessage = "Username and password are required"
		return nil
	}

	av.mode = AuthModeLoading
	av.loading = true

	return func() tea.Msg {
		resp, err := av.apiClient.Login(username, password)
		if err != nil {
			return ErrorMsg(fmt.Sprintf("Login failed: %v", err))
		}

		// Save token to cache
		av.cache.SetSetting("auth_token", resp.Token)
		av.cache.SetSetting("username", resp.User.Username)

		return LoginSuccessMsg{
			User:  resp.User,
			Token: resp.Token,
		}
	}
}

// doRegister performs registration
func (av *AuthView) doRegister() tea.Cmd {
	username := strings.TrimSpace(av.usernameInput.Value())
	email := strings.TrimSpace(av.emailInput.Value())
	password := av.passwordInput.Value()
	password2 := av.password2Input.Value()

	// Validation
	if username == "" || email == "" || password == "" {
		av.errorMessage = "All fields are required"
		return nil
	}

	if len(username) < 3 {
		av.errorMessage = "Username must be at least 3 characters"
		return nil
	}

	if len(password) < 8 {
		av.errorMessage = "Password must be at least 8 characters"
		return nil
	}

	if password != password2 {
		av.errorMessage = "Passwords do not match"
		return nil
	}

	if !strings.Contains(email, "@") {
		av.errorMessage = "Invalid email address"
		return nil
	}

	av.mode = AuthModeLoading
	av.loading = true

	return func() tea.Msg {
		err := av.apiClient.Register(username, email, password)
		if err != nil {
			return ErrorMsg(fmt.Sprintf("Registration failed: %v", err))
		}

		return RegisterSuccessMsg{
			Username: username,
		}
	}
}

// View renders the auth view
func (av *AuthView) View() string {
	switch av.mode {
	case AuthModeLoading:
		return av.renderLoading()
	case AuthModeSuccess:
		return av.renderSuccess()
	case AuthModeLogin:
		return av.renderLogin()
	case AuthModeRegister:
		return av.renderRegister()
	}
	return ""
}

// renderLogin renders login form
func (av *AuthView) renderLogin() string {
	title := av.styles.Header.Render("LOGIN TO TERMINAL NEWS")

	usernameLabel := av.styles.Label.Render("Username:")
	passwordLabel := av.styles.Label.Render("Password:")

	var usernameStyle, passwordStyle lipgloss.Style
	if av.activeField == 0 {
		usernameStyle = av.styles.ActiveInput
	} else {
		usernameStyle = av.styles.Input
	}

	if av.activeField == 1 {
		passwordStyle = av.styles.ActiveInput
	} else {
		passwordStyle = av.styles.Input
	}

	usernameView := usernameStyle.Render(av.usernameInput.View())
	passwordView := passwordStyle.Render(av.passwordInput.View())

	form := lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		usernameLabel,
		usernameView,
		"",
		passwordLabel,
		passwordView,
		"",
	)

	// Error message
	errorMsg := ""
	if av.errorMessage != "" {
		errorMsg = av.styles.ErrorMessage.Render("✗ " + av.errorMessage)
	}

	// Help text
	help := av.styles.HelpText.Render("[Enter] Login  [Tab] Next Field  [Ctrl+T] Register  [Esc] Cancel")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		form,
		errorMsg,
		"",
		help,
	)

	return av.styles.Border.Render(content)
}

// renderRegister renders registration form
func (av *AuthView) renderRegister() string {
	title := av.styles.Header.Render("CREATE ACCOUNT")

	usernameLabel := av.styles.Label.Render("Username:")
	emailLabel := av.styles.Label.Render("Email:")
	passwordLabel := av.styles.Label.Render("Password:")
	password2Label := av.styles.Label.Render("Confirm Password:")

	// Apply active styling
	inputs := []textinput.Model{av.usernameInput, av.emailInput, av.passwordInput, av.password2Input}
	views := make([]string, 4)
	labels := []string{usernameLabel, emailLabel, passwordLabel, password2Label}

	for i, input := range inputs {
		var style lipgloss.Style
		if av.activeField == i {
			style = av.styles.ActiveInput
		} else {
			style = av.styles.Input
		}

		views[i] = lipgloss.JoinVertical(
			lipgloss.Left,
			labels[i],
			style.Render(input.View()),
			"",
		)
	}

	form := lipgloss.JoinVertical(lipgloss.Left, views...)

	// Error message
	errorMsg := ""
	if av.errorMessage != "" {
		errorMsg = av.styles.ErrorMessage.Render("✗ " + av.errorMessage)
	}

	// Help text
	help := av.styles.HelpText.Render("[Enter] Register  [Tab] Next Field  [Ctrl+T] Login  [Esc] Cancel")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		form,
		errorMsg,
		"",
		help,
	)

	return av.styles.Border.Render(content)
}

// renderLoading renders loading state
func (av *AuthView) renderLoading() string {
	loading := `
	⠋ Processing...
	`
	return av.styles.InfoMessage.Render(loading)
}

// renderSuccess renders success state
func (av *AuthView) renderSuccess() string {
	success := fmt.Sprintf(`
╔════════════════════════════════════════════════════════════╗
║                                                            ║
║                    ✓ Success!                             ║
║                                                            ║
║              %s              ║
║                                                            ║
║              Press any key to continue...                  ║
║                                                            ║
╚════════════════════════════════════════════════════════════╝
	`, av.successMessage)

	return av.styles.SuccessMessage.Render(success)
}

// Message types
type LoginSuccessMsg struct {
	User  models.User
	Token string
}

type RegisterSuccessMsg struct {
	Username string
}
