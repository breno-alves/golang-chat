package service

func CheckPasswordHash(hash, password string) bool {
	return hash == password
}
