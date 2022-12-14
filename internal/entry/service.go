package entry

type Repository interface {
	Save(entry Entry) (*Entry, error)
	ListAll(start, count int) ([]Entry, error)
	FindById(id int) (*Entry, error)
	DeleteById(id int) error
	Update(entry Entry) (*Entry, error)
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

func (service *Service) DeleteEntry(id int) error {
	_, err := service.GetEntry(id)
	if err != nil {
		return err
	}
	return service.Repository.DeleteById(id)
}

func (service *Service) UpdateEntry(id int, entry Entry) (*Entry, error) {
	currentEntry, err := service.GetEntry(id)
	if err != nil {
		return nil, err
	}

	if currentEntry == nil || currentEntry.ID == 0 {
		return nil, nil
	}

	entry.ID = id
	return service.Repository.Update(entry)
}
