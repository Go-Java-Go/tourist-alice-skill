package service

import (
	"tourist-alice-skill/internal/repository"
)

func NewUserService(r repository.UserRepository) *UserService {
	return &UserService{r}
}

func NewChatStateService(r repository.ChatStateRepository) *ChatStateService {
	return &ChatStateService{r}
}

type UserService struct {
	repository.UserRepository
}

type ChatStateService struct {
	repository.ChatStateRepository
}
