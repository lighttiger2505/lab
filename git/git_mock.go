package git

type MockClient struct {
}

func NewMockClient() *MockClient {
	return &MockClient{}
}

func (g *MockClient) CurrentBranch() (string, error) {
	return "currentBranch", nil
}
