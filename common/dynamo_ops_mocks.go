package common

import (
	"context"
	"errors"
	"time"
)

type DynamoOpsMock struct {
	LastContentId            int32
	FCSIdForContentIdResults *[]string
	FCSIdForContentError     error

	IdMappingIndexQueried      string
	IdMappingKeyFieldQueried   string
	IdMappingSearchTermQueried interface{}
	IdMappingResult            IdMappingRecord
	IdMappingError             error
}

func (ops *DynamoOpsMock) QueryFCSIdForContentId(ctx context.Context, contentId int32) (*[]string, error) {
	ops.LastContentId = contentId
	if ops.FCSIdForContentError != nil {
		return nil, ops.FCSIdForContentError
	} else {
		return ops.FCSIdForContentIdResults, nil
	}
}

func (ops *DynamoOpsMock) QueryEncodingsForFCSId(ctx context.Context, fcsid string) ([]*Encoding, error) {
	return nil, errors.New("not implemented")
}

func (ops *DynamoOpsMock) QueryEncodingsForContentId(ctx context.Context, contentid int32, maybeSince *time.Time) ([]*Encoding, error) {
	return nil, errors.New("not implemented")
}

func (ops *DynamoOpsMock) QueryIdMappings(ctx context.Context, indexName string, keyFieldName string, searchTerm interface{}) (*IdMappingRecord, error) {
	ops.IdMappingIndexQueried = indexName
	ops.IdMappingKeyFieldQueried = keyFieldName
	ops.IdMappingSearchTermQueried = searchTerm
	if ops.IdMappingError != nil {
		return nil, ops.IdMappingError
	} else {
		coped := ops.IdMappingResult
		return &coped, nil
	}
}
