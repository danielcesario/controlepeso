package controlepeso

import "fmt"

type Repository interface {
	Save(entry Entry) error
}

type Service struct {
	Repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		Repository: repository,
	}
}

func (service *Service) CreateEntry(entry Entry) error {
	fmt.Println("CreateEntry")
	return service.Repository.Save(entry)
}
