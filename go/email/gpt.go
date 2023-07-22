package email

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	store "scraper/store"
	"strings"

	openai "github.com/sashabaranov/go-openai"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type EmailService struct {
	openai    *openai.Client
	sendgrid  *sendgrid.Client
	translate *TranslationClient
}

type Email struct {
	EmailService
	Ctx      context.Context
	Language string
	Prompt   string // prompt to generate email by gpt-3
	Subject  string // subject of email
	From     string // sender of email
	ToEmail  string
	Message  string
}

// generate email using gpt-3.5-turbo
func (e *Email) generate() (string, error) {

	ctx := context.Background()
	res, err := e.openai.CreateCompletion(ctx, openai.CompletionRequest{
		Model:       openai.GPT432K0613,
		MaxTokens:   200,
		Temperature: 0.7,
		Prompt:      e.Prompt,
	})

	if err != nil {
		return "", err
	}

	e.Message = res.Choices[0].Text
	return res.Choices[0].Text, nil
}

// send email using smtp
func (e *Email) Send() error {
	from := mail.NewEmail("Example User", "test@example.com")
	to := mail.NewEmail("Example User", "test@example.com")
	htmlContent := convertToHTML(e.Message)
	// convert to html

	message := mail.NewSingleEmail(from, e.Subject, to, e.Message, htmlContent)
	res, err := e.sendgrid.Send(message)
	if err != nil {
		log.Println(err)
		return err
	}
	if res.StatusCode != 202 {
		log.Println(res.StatusCode)
		return fmt.Errorf("status code: %d", res.StatusCode)
	}
	return nil
}

func (e *EmailService) NewEmail(st store.AccountInfo, lang string) (Email, error) {
	// verify email exists
	verifiedEmail := verifyEmail(st.Email)
	from := verifyEmail([]string{"abdelmanan.abdelrahman03@gmail.com"})
	phone := "+905074405861"
	addr := "19 Mayis, Istanbul, Turkey"
	verifiedPhone := verifyPhone(st.Phone)
	sb := strings.Builder{}
	// use gpt-3 to generate email
	sb.WriteString(
		fmt.Sprintf("Given a [%s] of a country, ", strings.Join(store.GetKeysWithPrefix(st, "Country"), ", ")),
	)
	sb.WriteString(
		"write an email requesting a coin from that country as part of a personal project/hobby of collecting coins from different countries. The email should be polite, respectful and friendly, and include some details about the sender’s deep respect and interest to [CountryName] and its culture. The email should also offer to send a coin from [SenderResidingCountry] in return, if the recipient wishes. since [CountryName] is the only country the sender does not have a coin from. The email should speak passionately about [CountryName] by [SenderName]. The email is being sent to an ambassador, it should be addressed as “Your Excellency”.\n",
	)
	sb.WriteString(fmt.Sprintf(`Details about the sender:
	- SenderName: Abdelmanan Abdelrahman
	- SenderEmail: %s
	- SenderPhone: %s
	- SenderAddress: %s
	- SenderBio: 
		I am a student at Marmara University, studying Electrical and Electronics Engineering. 
		I am Ghanaian by origin, but I have lived in Turkey for the past 2 years.
		I am interested in learning about the history, language, and culture of different nations through their currencies.
	- CountryName: %s
	- CountryAmbassador: %s
	- CountryEmail: %s
	- CountryPhone: %s
	- CountryAddress: %s
	`, from, phone, addr, st.Name, st.Ambassador, verifiedEmail, verifiedPhone, st.Address),
	)
	return Email{
		EmailService: *e,
		From:         from,
		Subject:      "Request for a coin from " + st.Name,
		ToEmail:      verifiedEmail,
		Language:     lang,
		Prompt:       sb.String(),
		Message:      "",
	}, nil
}
func NewEmailService(openai *openai.Client, sendgrid *sendgrid.Client, translate *TranslationClient) *EmailService {
	return &EmailService{openai, sendgrid, translate}
}

// convert plain text to html
func convertToHTML(text string) string {
	// TODO: convert text to html
	return text
}

// return the first verified email among list of email else empty string
func verifyEmail(emails []string) string {
	if len(emails) == 0 {
		return ""
	}
	// regex to match email
	regex := regexp.MustCompile(`[a-zA-Z.]+@[a-zA-Z.]{3,}`)
	for n, email := range emails {
		if !regex.MatchString(email) {
			// remove emaol from emails
			emails = append(emails[:n], emails[n+1:]...)
		}
	}
	//TODO: check if email is valid
	// Step 1: Validate the email using NeverBounce
	for n, email := range emails {
		neverBounceURL := fmt.Sprintf("https://api.neverbounce.com/v4/single/check?key=%s&email=%s", os.Getenv("NEVERBOUNC_API_KEY"), url.QueryEscape(email))
		response, err := http.Get(neverBounceURL)
		if err != nil {
			log.Fatal("Failed to connect to NeverBounce:", err)
		}
		defer response.Body.Close()

		var nbResponse struct {
			Result    string `json:"result"`
			Email     string `json:"email"`
			Suggested string `json:"suggested"`
		}

		if err := json.NewDecoder(response.Body).Decode(&nbResponse); err != nil {
			log.Fatal("Failed to decode NeverBounce response:", err)
		}

		if nbResponse.Result == "valid" || nbResponse.Result == "catchall" {
			continue
		}
		emails = append(emails[:n], emails[n+1:]...)
	}
	if len(emails) > 0 {
		return emails[0]
	}
	return ""
}

func verifyPhone(phones []string) string {
	regex := regexp.MustCompile(`[+][0-9]{11,}`)
	for n, phone := range phones {
		if !regex.MatchString(phone) {
			phones = append(phones[:n], phones[n+1:]...)
		}
	}

	if len(phones) > 0 {
		return phones[0]
	}
	return ""
}
func main() {
	openai := openai.NewClient("your token")
	sendgrid := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	client := NewTranslationClient()
	defer client.Close()

	e := NewEmailService(openai, sendgrid, client)
	e.NewEmail(store.AccountInfo{
		Name:       "Ghana",
		Email:      []string{"me@me.me"},
		Phone:      []string{"+233123456789"},
		Ambassador: "Kofi Nsiah-Poku",
		Address:    "Maltepe, Istanbul, Turkey",
	}, "en")

}
