package clients

import (
	client_dto "order-service/pkg/clients/dto"
	contextUtils "order-service/pkg/context"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type UserClient struct {
	baseUrl string
	hc      *http.Client
}

func NewUserClient(baseUrl string) *UserClient {
	return &UserClient{
		baseUrl: baseUrl,
		hc: &http.Client{
			Timeout: http.DefaultClient.Timeout,
		},
	}
}

func (c *UserClient) doRequest(ctx context.Context, method, path string, body interface{}, response interface{}) error {
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

func (c *UserClient) GetDoctorByIds(ctx context.Context, doctorID []string) (*[]client_dto.GetDoctorProfileResponseDto, error) {
	reqBody := client_dto.GetDoctorsByIDsRequestDto{
		DoctorIDs: doctorID,
	}

	var doctorProfiles []client_dto.GetDoctorProfileResponseDto
	if err := c.doRequest(ctx, http.MethodPost, "/v1/doctors", reqBody, &doctorProfiles); err != nil {
		return nil, err
	}

	return &doctorProfiles, nil
}

func (c *UserClient) GetDoctorById(ctx context.Context, doctorID string) (*client_dto.GetDoctorProfileResponseDto, error) {
	var doctorProfile client_dto.GetDoctorProfileResponseDto
	if err := c.doRequest(ctx, http.MethodPost, "/v1/doctors", client_dto.GetDoctorsByIDsRequestDto{DoctorIDs: []string{doctorID}}, &doctorProfile); err != nil {
		return nil, err
	}

	return &doctorProfile, nil
}

func (c *UserClient) GetPatientByIds(ctx context.Context, patientIDs []string) (*[]client_dto.GetPatientProfileResponseDto, error) {
	reqBody := client_dto.GetPatientsByIDsRequestDto{
		PatientIDs: patientIDs,
	}

	var patientProfiles []client_dto.GetPatientProfileResponseDto
	if err := c.doRequest(ctx, http.MethodPost, "/v1/patients", reqBody, &patientProfiles); err != nil {
		return nil, err
	}

	return &patientProfiles, nil
}
