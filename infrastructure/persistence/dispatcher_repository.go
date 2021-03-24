package persistence

import (
	"bytes"
	"config-manager/domain"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type MockDoType func(req *http.Request) (*http.Response, error)

type ClientMock struct {
	MockDo MockDoType
}

func (m *ClientMock) Do(req *http.Request) (*http.Response, error) {
	return m.MockDo(req)
}

func SetupMockDispatcherClient(expectedResponse string) *ClientMock {
	r := ioutil.NopCloser(bytes.NewReader([]byte(expectedResponse)))

	client := &ClientMock{
		MockDo: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 207,
				Body:       r,
			}, nil
		},
	}

	return client
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type DispatcherClient struct {
	DispatcherHost string
	DispatcherPSK  string
	Client         HTTPClient
}

func (r *DispatcherClient) Dispatch(
	ctx context.Context,
	inputs []domain.DispatcherInput,
) ([]domain.DispatcherResponse, error) {
	fmt.Println("Sending request to playbook dispatcher")

	reqBody, err := json.Marshal(inputs)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", r.DispatcherHost+"/internal/dispatch", bytes.NewBuffer(reqBody))
	req.Header.Set("Authorization", fmt.Sprintf("PSK %s", r.DispatcherPSK))
	req.Header.Set("Content-Type", "application/json")

	res, err := r.Client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	var dRes []domain.DispatcherResponse
	err = json.NewDecoder(res.Body).Decode(&dRes)
	return dRes, err
}
