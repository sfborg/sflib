package sfgaio

type sfgaio struct {
	repo sfga.GitHubRepo
}

func New(gf sfga.GitHubRepo) sfga.SFGA {
	return &sfgaio{}
}

func (s *sfgaio) FetchSchema() ([]byte, error) {
	return nil, nil
}
