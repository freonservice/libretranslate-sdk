package libretranslate

type FrontendSettingLanguage struct {
	Source Language `json:"source"`
	Target Language `json:"target"`
}

type FrontendSetting struct {
	Keys                 bool                    `json:"apiKeys"`
	KeyRequired          bool                    `json:"keyRequired"`
	Suggestions          bool                    `json:"suggestions"`
	CharLimit            int                     `json:"charLimit"`
	FrontendTimeout      int                     `json:"frontendTimeout"`
	Language             FrontendSettingLanguage `json:"language"`
	SupportedFilesFormat []string                `json:"supportedFilesFormat"`
}

type Language struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type TranslateRequest struct {
	Q      string `json:"q"`
	Source string `json:"source"`
	Target string `json:"target"`
	Format string `json:"format"`
	Key    string `json:"api_key"`
}

type TranslatedText struct {
	Text string `json:"translatedText"`
}
