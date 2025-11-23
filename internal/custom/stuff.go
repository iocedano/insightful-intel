package custom

import "fmt"

type CustomClient struct {
	Client    *CustomClient
	connected *CustomClient
	Stored    func(any) error
	Retrieved func() (any, error)
	// Deleted func() error
	// Updated func(any) error
	// Created func(any) error
}

type StuffInterface interface {
	Connect(conn *CustomClient)
	Store(data any) error
	Retrieve() (any, error)
	Search(query string) (any, error)
}

func NewCustomClient() *CustomClient {
	return &CustomClient{
		Client: NewCustomClient(),
	}
}

func (s *CustomClient) Connect(conn *CustomClient) {
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

func (s *CustomClient) Search(query string) (any, error) {
	if s.Search == nil {
		return nil, fmt.Errorf("search function is not set")
	}
	return s.Search(query)
}

type CustomPathMap struct {
	BaseURL string
	Paths   map[string]string
}

func (pm *CustomPathMap) GetURLFrom(path string) string {
	if p, ok := pm.Paths[path]; ok {
		return pm.BaseURL + p
	}
	return ""
}
