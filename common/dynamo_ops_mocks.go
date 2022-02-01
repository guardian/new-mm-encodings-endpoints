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

	EncodingsForFCSIdResults []*Encoding
	EncodingsForFCSIdError   error
	FCSIdQueried             string

	EncodingsForContentIdResults []*Encoding
	EncodingsForContentIdError   error
	ContentIdQueried             int32
	ContentIdSince               *time.Time
	IdMappingIndexQueried        string
	IdMappingKeyFieldQueried     string
	IdMappingSearchTermQueried   interface{}
	IdMappingResult              IdMappingRecord
	IdMappingError               error
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
	ops.FCSIdQueried = fcsid
	if ops.EncodingsForFCSIdError != nil {
		return nil, ops.EncodingsForFCSIdError
	} else {
		return ops.EncodingsForFCSIdResults, nil
	}
}

func (ops *DynamoOpsMock) QueryEncodingsForContentId(ctx context.Context, contentid int32, maybeSince *time.Time) ([]*Encoding, error) {
	ops.ContentIdQueried = contentid
	copiedTime := *maybeSince
	ops.ContentIdSince = &copiedTime
	if ops.EncodingsForContentIdError != nil {
		return nil, ops.EncodingsForContentIdError
	} else {
		return ops.EncodingsForContentIdResults, nil
	}
}

func (ops *DynamoOpsMock) QueryIdMappings(ctx context.Context, indexName string, keyFieldName string, searchTerm interface{}) (*IdMappingRecord, error) {
	ops.IdMappingIndexQueried = indexName
	ops.IdMappingKeyFieldQueried = keyFieldName
	ops.IdMappingSearchTermQueried = searchTerm
	if ops.IdMappingError != nil {
		return nil, ops.IdMappingError
	} else {
		empty := IdMappingRecord{}
		if ops.IdMappingResult == empty {
			return nil, nil
		}
		coped := ops.IdMappingResult
		return &coped, nil
	}
}

func (ops *DynamoOpsMock) GetAllMimeEquivalents(ctx context.Context) ([]*MimeEquivalent, error) {
	return nil, errors.New("not implemented in DynamoOpsMock")
}
