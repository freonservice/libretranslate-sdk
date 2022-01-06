package libretranslate

import (
	"context"
	"testing"
)

func TestClient_GetLanguages(t *testing.T) {
	client := NewLibreTranslate("https://libretranslate.com")
	data, err := client.GetLanguages(context.Background())
	if err != nil {
		t.Error(err)
	}
	if len(data) == 0 {
		t.Error("empty languages count")
	}
}

func TestClient_GetFrontendSetting(t *testing.T) {
	client := NewLibreTranslate("https://libretranslate.com")
	data, err := client.GetFrontendSetting(context.Background())
	if err != nil {
		t.Error(err)
	}
	if data.CharLimit == 0 {
		t.Error("char limit can't be zero")
	}
}
