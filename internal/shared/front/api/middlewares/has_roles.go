package middlewares


func HasRole(roles []string, role string) bool { 
	for _, r := range roles { 
		if r == role { 
			return true 
		} 
	} 
	return false 
}