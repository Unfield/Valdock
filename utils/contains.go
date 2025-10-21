package utils

func Contains(slice []int, port int) bool {
	for _, p := range slice {
		if p == port {
			return true
		}
	}
	return false
}
