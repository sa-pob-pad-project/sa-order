package service

import (
	"context"
	"errors"
	"order-service/pkg/apperr"
	"order-service/pkg/clients"
	contextUtils "order-service/pkg/context"
	"order-service/pkg/dto"
	"order-service/pkg/repository"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeliveryService struct {
	db                     *gorm.DB
	deliveryRepository     *repository.DeliveryRepository
	deliveryInfoRepository *repository.DeliveryInformationRepository
	orderRepository        *repository.OrderRepository
	userClient             *clients.UserClient
}

func NewDeliveryService(
	db *gorm.DB,
	deliveryRepo *repository.DeliveryRepository,
	deliveryInfoRepo *repository.DeliveryInformationRepository,
	orderRepo *repository.OrderRepository,
	userClient *clients.UserClient,
) *DeliveryService {
	return &DeliveryService{
		db:                     db,
		deliveryRepository:     deliveryRepo,
		deliveryInfoRepository: deliveryInfoRepo,
		orderRepository:        orderRepo,
		userClient:             userClient,
	}
}

// CreateDeliveryInfo creates a new delivery information record
func (s *DeliveryService) CreateDeliveryInfo(ctx context.Context, req dto.CreateDeliveryInfoRequestDto) (*dto.CreateDeliveryInfoResponseDto, error) {
	userID := contextUtils.GetUserId(ctx)
	deliveryInfo, err := dto.ToDeliveryInformation(userID, req)
	if err != nil {
		return nil, apperr.New(apperr.CodeBadRequest, "Invalid user ID format", err)
	}

	err = s.deliveryInfoRepository.Create(ctx, deliveryInfo)
	if err != nil {
		return nil, apperr.New(apperr.CodeInternal, "Failed to create delivery information", err)
	}

	return &dto.CreateDeliveryInfoResponseDto{
		DeliveryInfo: dto.ToDeliveryInfoDto(deliveryInfo),
	}, nil
}

// GetDeliveryInfoByID retrieves a delivery information record by ID
func (s *DeliveryService) GetDeliveryInfoByID(ctx context.Context, id string) (*dto.GetDeliveryInfoResponseDto, error) {
	deliveryInfoID, err := uuid.Parse(id)
	if err != nil {
		return nil, apperr.New(apperr.CodeBadRequest, "Invalid delivery information ID format", err)
	}

	deliveryInfo, err := s.deliveryInfoRepository.FindByID(ctx, deliveryInfoID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.New(apperr.CodeNotFound, "Delivery information not found", err)
		}
		return nil, apperr.New(apperr.CodeInternal, "Failed to retrieve delivery information", err)
	}

	return &dto.GetDeliveryInfoResponseDto{
		DeliveryInfo: dto.ToDeliveryInfoDto(deliveryInfo),
	}, nil
}

// GetDeliveryInfosByUserID retrieves all delivery information records for a user
func (s *DeliveryService) GetDeliveryInfosByUserID(ctx context.Context, userID string) (*dto.GetAllDeliveryInfosResponseDto, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, apperr.New(apperr.CodeBadRequest, "Invalid user ID format", err)
	}

	deliveryInfos, err := s.deliveryInfoRepository.FindByUserID(ctx, userUUID)
	if err != nil {
		return nil, apperr.New(apperr.CodeInternal, "Failed to retrieve delivery information", err)
	}

	return &dto.GetAllDeliveryInfosResponseDto{
		DeliveryInfos: dto.ToDeliveryInfoDtoList(deliveryInfos),
	}, nil
}

// GetAllDeliveryInfos retrieves all delivery information records
func (s *DeliveryService) GetAllDeliveryInfos(ctx context.Context) (*dto.GetAllDeliveryInfosResponseDto, error) {
	deliveryInfos, err := s.deliveryInfoRepository.FindAll(ctx)
	if err != nil {
		return nil, apperr.New(apperr.CodeInternal, "Failed to retrieve delivery information", err)
	}

	return &dto.GetAllDeliveryInfosResponseDto{
		DeliveryInfos: dto.ToDeliveryInfoDtoList(deliveryInfos),
	}, nil
}

// UpdateDeliveryInfo updates an existing delivery information record
func (s *DeliveryService) UpdateDeliveryInfo(ctx context.Context, req dto.UpdateDeliveryInfoRequestDto) (*dto.UpdateDeliveryInfoResponseDto, error) {
	deliveryInfoID, err := uuid.Parse(req.ID)
	if err != nil {
		return nil, apperr.New(apperr.CodeBadRequest, "Invalid delivery information ID format", err)
	}

	// Check if delivery info exists
	existingInfo, err := s.deliveryInfoRepository.FindByID(ctx, deliveryInfoID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.New(apperr.CodeNotFound, "Delivery information not found", err)
		}
		return nil, apperr.New(apperr.CodeInternal, "Failed to retrieve delivery information", err)
	}

	// Update fields
	existingInfo.Address = req.Address
	existingInfo.PhoneNumber = req.PhoneNumber
	existingInfo.DeliveryMethod = req.DeliveryMethod
	existingInfo.Version = existingInfo.Version + 1

	// Save updates
	err = s.deliveryInfoRepository.Update(ctx, existingInfo)
	if err != nil {
		return nil, apperr.New(apperr.CodeInternal, "Failed to update delivery information", err)
	}

	return &dto.UpdateDeliveryInfoResponseDto{
		DeliveryInfo: dto.ToDeliveryInfoDto(existingInfo),
	}, nil
}

// DeleteDeliveryInfo deletes a delivery information record
func (s *DeliveryService) DeleteDeliveryInfo(ctx context.Context, id string) (*dto.DeleteDeliveryInfoResponseDto, error) {
	deliveryInfoID, err := uuid.Parse(id)
	if err != nil {
		return nil, apperr.New(apperr.CodeBadRequest, "Invalid delivery information ID format", err)
	}

	// Check if delivery info exists
	_, err = s.deliveryInfoRepository.FindByID(ctx, deliveryInfoID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.New(apperr.CodeNotFound, "Delivery information not found", err)
		}
		return nil, apperr.New(apperr.CodeInternal, "Failed to retrieve delivery information", err)
	}

	// Delete the record
	err = s.deliveryInfoRepository.Delete(ctx, deliveryInfoID)
	if err != nil {
		return nil, apperr.New(apperr.CodeInternal, "Failed to delete delivery information", err)
	}

	return &dto.DeleteDeliveryInfoResponseDto{
		ID:        id,
		DeletedAt: time.Now().Format(time.RFC3339),
	}, nil
}
