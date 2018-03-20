package test

type DiffCalculator struct {
	client RPCClient
}

type Post struct {
	Message    string
	ExternalID string
}

type GetDiffRequest struct {
	SourceID string
	Posts    []Post
}

type GetDiffResponse struct {
}

func (c *DiffCalculator) GetDiff(req GetDiffRequest) (resp GetDiffResponse, err error) {
	return resp, c.client.do("diff-calculator", "GetDiff", req, &resp)
}
