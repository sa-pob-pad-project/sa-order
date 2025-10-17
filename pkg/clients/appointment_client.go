package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	client_dto "order-service/pkg/clients/dto"
	contextUtils "order-service/pkg/context"

	"github.com/google/uuid"
)

type AppointmentClient struct {
	baseUrl string
	hc      *http.Client
}

func NewAppointmentClient(baseUrl string) *AppointmentClient {
	return &AppointmentClient{
		baseUrl: baseUrl,
		hc: &http.Client{
			Timeout: http.DefaultClient.Timeout,
		},
	}
}

func (c *AppointmentClient) doRequest(ctx context.Context, method, path string, body interface{}, response interface{}) error {
	accessToken := contextUtils.GetAccessToken(ctx)
	if accessToken == "" {
		return fmt.Errorf("access token is empty")
	}

	url := fmt.Sprintf("%s%s", c.baseUrl, path)

	var req *http.Request
	var err error

	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		req, err = http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}
	}

	req.AddCookie(&http.Cookie{
		Name:  "access_token",
		Value: accessToken,
	})

	fmt.Println("Requesting URL:", url)
	fmt.Println("With Method:", method)
	if body != nil {
		fmt.Println("With Body:", body)
	}

	resp, err := c.hc.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if response != nil {
		if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

func (c *AppointmentClient) GetLatestAppointmentByPatientID(ctx context.Context, patientID uuid.UUID) (*client_dto.GetLatestAppointmentResponseDto, error) {

	var appointment client_dto.GetLatestAppointmentResponseDto
	if err := c.doRequest(ctx, http.MethodGet, "/v1/patient/history/latest", nil, &appointment); err != nil {
		return nil, err
	}

	return &appointment, nil
}
