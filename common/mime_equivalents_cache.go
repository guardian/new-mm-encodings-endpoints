package common

import "context"

type MimeEquivalentsCache interface {
	/*
		EquivalentsFor will return a list of _all_ known equivalents to the given MIME type (including itself).
	*/
	EquivalentsFor(input string) []string
}

type MimeEquivalentsCacheImpl struct {
	loadedData map[string]string
}

func NewMimeEquivalentsCache(ctx context.Context, ops DynamoDbOps) (MimeEquivalentsCache, error) {
	equivs, err := ops.GetAllMimeEquivalents(ctx)
	if err != nil {
		return nil, err
	}

	cache := &MimeEquivalentsCacheImpl{
		loadedData: make(map[string]string, 2*len(equivs)),
	}

	for _, equiv := range equivs {
		cache.loadedData[equiv.RealName] = equiv.MimeEquivalent
		cache.loadedData[equiv.MimeEquivalent] = equiv.RealName
	}
	return cache, nil
}

func (cache *MimeEquivalentsCacheImpl) EquivalentsFor(input string) []string {
	result := []string{input}

	if lookedUpEquivalent, haveEquiv := cache.loadedData[input]; haveEquiv {
		result = append(result, lookedUpEquivalent)
		return result
	} else {
		return result
	}
}

type MimeEquivalentsCacheMock struct {
}

func (cache *MimeEquivalentsCacheMock) EquivalentsFor(input string) []string {
	result := []string{input}
	return result
}
