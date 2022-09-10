package controlepeso

type Repository interface {
	Save(entry Entry) (*Entry, error)
}

type Service struct {
	Repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		Repository: repository,
	}
}

func (service *Service) CreateEntry(entry Entry) (*Entry, error) {
	return service.Repository.Save(entry)
}
