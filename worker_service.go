package adp

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type ADPWorkerService struct {
	Authenticator ADPAuthenticationSystem
}

type ADPWorkerEmployeeResponse struct {
	Workers []*Worker `json:"workers"`
}

func NewADPWorkerService(authenticator ADPAuthenticationSystem) ADPWorkerService {
	return ADPWorkerService{
		Authenticator: authenticator,
	}
}

func (a *ADPWorkerService) GetWorker(aoid string) (*ADPWorkerEmployeeResponse, error) {
	url := fmt.Sprintf("hr/v2/workers/%s", aoid)
	body, err := a.MakeRequest(http.MethodGet, url)
	if err != nil {
		return nil, err
	}
	response := ADPWorkerEmployeeResponse{}
	if err := json.NewDecoder(body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (a *ADPWorkerService) ListWorkersAsync(reciever chan<- *Worker, resultLimit int) {

	errChan := make(chan error)
	closeChan := make(chan struct{})

	defer close(reciever)
	defer close(errChan)
	defer close(closeChan)

	go func() {
		top := resultLimit
		skip := 0
		for {
			response, err := a.makeWorkerAPIRequest(top, skip)
			if err != nil {
				errChan <- err
			}

			if len(response.Workers) == 0 {
				closeChan <- struct{}{}
			}

			for _, worker := range response.Workers {
				reciever <- worker
			}
			skip += resultLimit
		}
	}()

	for {
		select {
		case err := <-errChan:
			log.Fatal(err)
		case <-closeChan:
			return
		}
	}
}

func (a *ADPWorkerService) makeWorkerAPIRequest(top, skip int) (*ADPWorkerEmployeeResponse, error) {
	response := &ADPWorkerEmployeeResponse{}
	body, err := a.MakeRequest(http.MethodGet, fmt.Sprintf("hr/v2/workers?$skip=%d&top=%d", skip, top))
	if err != nil {
		return nil, err
	}
	if err := json.NewDecoder(body).Decode(response); err != nil && err != io.EOF {
		return nil, err
	}
	return response, nil
}

func (a *ADPWorkerService) ListWorkers() (*ADPWorkerEmployeeResponse, error) {
	top := 100
	skip := 0
	workers := ADPWorkerEmployeeResponse{
		Workers: []*Worker{},
	}

	for {
		response, err := a.makeWorkerAPIRequest(top, skip)
		if err != nil {
			return nil, err
		}

		if len(response.Workers) == 0 {
			break
		}

		workers.Workers = append(workers.Workers, response.Workers...)
		skip += 100
	}

	return &workers, nil
}

func (a *ADPWorkerService) MakeRequest(method string, path string) (io.ReadCloser, error) {
	baseUrl := "https://api.adp.com"
	request, err := http.NewRequest(method, fmt.Sprintf("%s/%s", baseUrl, path), nil)
	if err != nil {
		return nil, err
	}
	a.Authenticator.SetRequestAuthorizationHeader(request)
	if method == http.MethodGet {
		request.Header.Set("Content-Type", "application/json")
	}

	client, err := a.Authenticator.NewHttpClient()
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	if !isValidResponseStatusCode(resp) {
		return nil, fmt.Errorf("invalid request. status code: %d", resp.StatusCode)
	}

	return resp.Body, nil
}
