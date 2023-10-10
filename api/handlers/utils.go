package handlers

func jsonResponse(message string) *map[string]string {
	response := map[string]string{"status": message}
	return &response
}
