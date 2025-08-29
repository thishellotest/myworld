package adobesign

import "context"

type LibraryDocumentsService service

func (s *LibraryDocumentsService) LibraryDocuments(ctx context.Context) (interface{}, error) {

	req, err := s.client.NewRequest("GET", "libraryDocuments", nil)
	if err != nil {
		return nil, err
	}

	var response interface{}
	if _, err := s.client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	return response, nil
}
