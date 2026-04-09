package services

import (
	"user-service/clients"
	"user-service/models"
	"user-service/repositories"
)

type UserService struct {
	repo     *repositories.UserRepository
	producer *clients.KafkaProducer
}

func NewUserService(repo *repositories.UserRepository, producer *clients.KafkaProducer) *UserService {
	return &UserService{repo: repo, producer: producer}
}

func (s *UserService) CreateUser(input models.CreateUserInput) (*models.User, error) {
	user := &models.User{
		Name:  input.Name,
		Email: input.Email,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	_ = s.producer.PublishUserCreated(models.UserCreatedEvent{
		UserID: user.ID,
		Email:  user.Email,
		Name:   user.Name,
	})

	return user, nil
}

func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	return s.repo.FindByID(id)
}
