package models

import "strconv"

type PathGlobalModelsJSON struct {
	RequestsInSeconds string `json:"requestInSecond"`
	StartTokens       string `json:"startTokens"`
}

type PatchGlobalModelsResult struct {
	ReqInSecond int
	StartTokens int
}

func (pr *PatchGlobalModelsResult) Validate() bool {
	return pr.ReqInSecond > 0
}

func (p *PathGlobalModelsJSON) ToIntegerStruct() (*PatchGlobalModelsResult, error) {
	result := &PatchGlobalModelsResult{}

	var err error
	result.ReqInSecond, err = strconv.Atoi(p.RequestsInSeconds)
	if err != nil {
		return nil, err
	}
	result.StartTokens, err = strconv.Atoi(p.StartTokens)
	if err != nil {
		return nil, err
	}

	return result, nil
}
