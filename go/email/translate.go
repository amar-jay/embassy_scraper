package email

import (
	"context"
	"fmt"
	"log"

	translate "cloud.google.com/go/translate"
	"golang.org/x/text/language"
)

type Text struct {
	text string
	lang string
}

func (t Text) String() string {
	return t.text
}

type TranslationClient struct {
	ctx context.Context
	*translate.Client
}

func NewTranslationClient() *TranslationClient {
	ctx := context.Background()

	client, err := translate.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return &TranslationClient{ctx, client}
}

// from is by default english
func (t TranslationClient) TranslateTo(text, to string, from ...string) (string, error) {
	// TODO: translate text
	model := "nmt"
	lang, err := language.Parse(to)
	if err != nil {
		return "", fmt.Errorf("language.Parse: %w", err)
	}
	resp, err := t.Translate(t.ctx, []string{text}, lang, &translate.Options{
		Model: model, // Either "nmt" or "base".
	})

	if err != nil {
		return "", fmt.Errorf("Translate: %w", err)
	}

	if len(resp) == 0 {
		return "", nil
	}
	return resp[0].Text, nil
}
