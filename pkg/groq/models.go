package groq

// ModelsList contains the available models
var ModelsList = []string{
	"mixtral-8x7b-32768",
	"llama2-70b-4096",
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
