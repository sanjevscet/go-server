package utils

import (
	"log"
	"strings"

	"github.com/google/uuid"
)

func GetUUID() string {
	log.Println("get-uuid called")
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}
