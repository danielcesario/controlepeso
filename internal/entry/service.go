package entry

type Repository interface {
	Save(entry Entry) (*Entry, error)
	ListAll(start, count int) ([]Entry, error)
	FindById(id int) (*Entry, error)
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

func (service *Service) ListEntries(start, count int) ([]Entry, error) {
	return service.Repository.ListAll(start, count)
}

func (service *Service) GetEntry(id int) (*Entry, error) {
	return service.Repository.FindById(id)
}
