package stuff

import "fmt"

type Stuff struct {
	Client    *CustomClient
	connected *Stuff
	Stored    func(any) error
	Retrieved func() (any, error)
	// Deleted func() error
	// Updated func(any) error
	// Created func(any) error
}

type StuffInterface interface {
	Connect(conn *Stuff)
	Store(data any) error
	Retrieve() (any, error)
	Search(query string) (any, error)
}

func NewStuff() *Stuff {
	return &Stuff{
		Client: NewCustomClient(),
	}
}

func (s *Stuff) Connect(conn *Stuff) {
	s.connected = conn
}

// func (s *Stuff) Store(data any) error {
// 	if s.Stored == nil {
// 		return fmt.Errorf("stored function is not set")
// 	}
// 	return s.Stored(data)
// }

// func (s *Stuff) Retrieve() (any, error) {
// 	if s.Retrieved == nil {
// 		return nil, fmt.Errorf("retrieved function is not set")
// 	}
// 	return s.Retrieved()
// }

func (s *Stuff) Search(query string) (any, error) {
	if s.Search == nil {
		return nil, fmt.Errorf("search function is not set")
	}
	return s.Search(query)
}

type PathMap struct {
	BaseURL string
	Paths   map[string]string
}

func (pm PathMap) GetURLFrom(path string) string {
	if p, ok := pm.Paths[path]; ok {
		return pm.BaseURL + p
	}
	return ""
}
