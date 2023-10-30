package todo

import (
	"log"
	"reflect"
	"testing"
)

type mockRepository struct {
	todos []Todo
}

func (mr *mockRepository) GetAllTodos() ([]Todo, error) {
	return mr.todos, nil
}

func (mr *mockRepository) AddTodo(t string) (Todo, error) {
	return Todo{TaskName: t}, nil
}

func (mr *mockRepository) UpdateTaskName(i, t string) (Todo, error) {
	return Todo{Id: i, TaskName: t}, nil
}

func (mr *mockRepository) UpdateCompleted(i string, c bool) (Todo, error) {
	return Todo{Id: i, Completed: c}, nil
}

func (mr *mockRepository) DeleteTodo(id string) (bool, error) {
	index := -1

	for i, v := range mr.todos {
		if v.Id == id {
			index = i
		}
	}

	if index > -1 {
		mr.todos = append(mr.todos[:index], mr.todos[index+1:]...)

		return true, nil
	}

	return false, nil
}

func TestGetTodos(t *testing.T) {
	expected := []Todo{
		{
			Id:        "1",
			TaskName:  "Todo 1",
			UpdatedAt: "00:00",
			Completed: true,
		},
		{
			Id:        "2",
			TaskName:  "Todo 2",
			UpdatedAt: "00:00",
			Completed: false,
		},
	}

	mr := mockRepository{
		todos: expected,
	}

	s := NewService(&mr)

	if got := s.GetTodos(); !reflect.DeepEqual(expected, got) {
		log.Fatalf("TestGetTodos expected %v, got %v", expected, got)
	}
}

func TestAddTodo(t *testing.T) {
	expected := Todo{
		TaskName: "Todo 1",
	}

	mr := mockRepository{}

	s := NewService(&mr)

	todo, err := s.AddTodo(expected.TaskName)
	if err != nil {
		log.Fatalf("TestAddTodo expected %v, got %v", nil, err)
	}

	if todo != expected {
		log.Fatalf("TestAddTodo expected %v, got %v", expected, todo)
	}

	if _, err := s.AddTodo(""); err != ErrInvalidTaskName {
		log.Fatalf("TestAddTodo expected %v, got %v", ErrInvalidTaskName, err)
	}
}

func TestUpdateTaskName(t *testing.T) {
	expected := Todo{
		Id:       "123",
		TaskName: "Todo 1",
	}

	mr := mockRepository{}

	s := NewService(&mr)

	todo, err := s.UpdateTaskName(expected.Id, expected.TaskName)
	if err != nil {
		log.Fatalf("TestUpdateTodo expected %v, got %v", nil, err)
	}

	if todo != expected {
		log.Fatalf("TestUpdateTodo expected %v, got %v", expected, todo)
	}

	if _, err := s.UpdateTaskName("1234", ""); err != ErrInvalidTaskName {
		log.Fatalf("TestUpdateTaskName expected %v, got %v", ErrInvalidTaskName, err)
	}

	if _, err := s.UpdateTaskName("", "Todo 1"); err != ErrInvalidId {
		log.Fatalf("TestUpdateTaskName expected %v, got %v", ErrInvalidId, err)
	}
}

func TestUpdateCompleted(t *testing.T) {
	expected := Todo{
		Id:        "123",
		Completed: true,
	}

	mr := mockRepository{}

	s := NewService(&mr)

	todo, err := s.UpdateCompleted(expected.Id, expected.Completed)
	if err != nil {
		log.Fatalf("TestUpdateTodo expected %v, got %v", nil, err)
	}

	if todo != expected {
		log.Fatalf("TestUpdateTodo expected %v, got %v", expected, todo)
	}

	if _, err := s.UpdateCompleted("", expected.Completed); err != ErrInvalidId {
		log.Fatalf("TestUpdateCompleted expected %v, got %v", ErrInvalidId, err)
	}
}

func TestDeleteTodo(t *testing.T) {
	todos := []Todo{
		{
			Id:        "1",
			TaskName:  "Todo 1",
			UpdatedAt: "00:00",
			Completed: true,
		},
		{
			Id:        "2",
			TaskName:  "Todo 2",
			UpdatedAt: "00:00",
			Completed: false,
		},
	}

	mr := mockRepository{
		todos: todos,
	}

	s := NewService(&mr)

	if got, _ := s.DeleteTodo("1"); got != true {
		log.Fatalf("TestDeleteTodo expected %v, got %v", true, got)
	}

	expected := todos[1:]
	if !reflect.DeepEqual(expected, mr.todos) {
		log.Fatalf("TestDeleteTodo expected %v, got %v", expected, mr.todos)
	}
}

func TestValidateTaskName(t *testing.T) {
	tests := []struct {
		expected error
		input    string
	}{
		{
			expected: nil,
			input:    "Todo 1",
		},
		{
			expected: ErrInvalidTaskName,
			input:    "",
		},
	}

	for _, test := range tests {
		if got := ValidateTaskName(test.input); got != test.expected {
			log.Fatalf("TestValidateTaskName expected %v, got %v", test.expected, got)
		}
	}
}

func TestValidateId(t *testing.T) {
	tests := []struct {
		expected error
		input    string
	}{
		{
			expected: nil,
			input:    "1234",
		},
		{
			expected: ErrInvalidId,
			input:    "",
		},
	}

	for _, test := range tests {
		if got := ValidateId(test.input); got != test.expected {
			log.Fatalf("TestValidateId expected %v, got %v", test.expected, got)
		}
	}
}
