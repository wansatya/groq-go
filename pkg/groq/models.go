package groq

// ModelsList contains the available models, for all supported models visit https://console.groq.com/docs/models.
var ModelsList = []string{
	"llama3-8b-8192",
	"llama-3.2-1b-preview",
	"llama3-groq-8b-8192-tool-use-preview",
	"llama-3.2-11b-text-preview",
	"llava-v1.5-7b-4096-preview",
	"gemma2-9b-it",
	"llama3-groq-70b-8192-tool-use-preview",
	"llama-3.2-90b-text-preview",
	"gemma-7b-it",
	"llama-3.2-3b-preview",
	"llama-guard-3-8b",
	"llama-3.1-8b-instant",
	"llama3-70b-8192",
	"mixtral-8x7b-32768",
	"llama-3.1-70b-versatile"
}

// IsValidModel checks if the given model name is valid
func IsValidModel(model string) bool {
	for _, m := range ModelsList {
		if m == model {
			return true
		}
	}
	return false
}
