package todo

import (
	"errors"
	"log"
)

var (
	ErrInvalidTaskName = errors.New("invalid todo")
	ErrInvalidId       = errors.New("invalid id")
)

type Service interface {
	GetTodos() []Todo
	AddTodo(taskName string) (Todo, error)
	UpdateTaskName(id, taskName string) (Todo, error)
	UpdateCompleted(id string, completed bool) (Todo, error)
	DeleteTodo(id string) (bool, error)
}

type Repository interface {
	GetAllTodos() ([]Todo, error)
	AddTodo(taskName string) (Todo, error)
	UpdateTaskName(id, taskName string) (Todo, error)
	UpdateCompleted(id string, completed bool) (Todo, error)
	DeleteTodo(id string) (bool, error)
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{
		r: r,
	}
}

func (s *service) AddTodo(t string) (Todo, error) {
	if err := ValidateTaskName(t); err != nil {
		return Todo{}, err
	}

	return s.r.AddTodo(t)
}

func (s *service) UpdateTaskName(i, t string) (Todo, error) {
	if err := ValidateId(i); err != nil {
		return Todo{}, err
	}

	if err := ValidateTaskName(t); err != nil {
		return Todo{}, err
	}

	return s.r.UpdateTaskName(i, t)
}

func (s *service) UpdateCompleted(i string, c bool) (Todo, error) {
	if err := ValidateId(i); err != nil {
		return Todo{}, err
	}

	return s.r.UpdateCompleted(i, c)
}

func (s *service) GetTodos() []Todo {
	todos, err := s.r.GetAllTodos()
	if err != nil {
		log.Println(err)

		return []Todo{}
	}

	return todos
}

func (s *service) DeleteTodo(i string) (bool, error) {
	if err := ValidateId(i); err != nil {
		return false, err
	}

	return s.r.DeleteTodo(i)
}

// validations
func ValidateTaskName(t string) error {
	if t == "" {
		return ErrInvalidTaskName
	}

	return nil
}

func ValidateId(i string) error {
	if i == "" {
		return ErrInvalidId
	}

	return nil
}
