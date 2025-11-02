package configtemplates

import "fmt"

func formatMemory(memory int) string {
	return fmt.Sprintf("%dmb", memory)
}
