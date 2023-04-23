package service

import (
	"account-management-service/internal/repo"
	"account-management-service/internal/webapi"
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"strconv"
)

type OperationService struct {
	operationRepo repo.Operation
	productRepo   repo.Product
	gDrive        webapi.GDrive
}

func NewOperationService(operationRepo repo.Operation, productRepo repo.Product, gDrive webapi.GDrive) *OperationService {
	return &OperationService{
		operationRepo: operationRepo,
		productRepo:   productRepo,
		gDrive:        gDrive,
	}
}

func (s *OperationService) OperationHistory(ctx context.Context, input OperationHistoryInput) ([]OperationHistoryOutput, error) {
	operations, productNames, err := s.operationRepo.OperationsPagination(ctx, input.AccountId, input.SortType, input.Offset, input.Limit)
	if err != nil {
		return nil, err
	}

	output := make([]OperationHistoryOutput, 0, len(operations))
	for i, operation := range operations {

		output = append(output, OperationHistoryOutput{
			Amount:      operation.Amount,
			Operation:   operation.OperationType,
			Time:        operation.CreatedAt,
			Product:     productNames[i],
			Order:       operation.OrderId,
			Description: operation.Description,
		})
	}
	return output, nil
}

func (s *OperationService) MakeReportLink(ctx context.Context, month, year int) (string, error) {
	if !s.gDrive.IsAvailable() {
		return "", errors.New("google drive is not available")
	}

	file, err := s.MakeReportFile(ctx, month, year)
	if err != nil {
		return "", err
	}

	url, err := s.gDrive.UploadCSVFile(ctx, fmt.Sprintf("report_%d_%d.csv", month, year), file)
	if err != nil {
		return "", errors.New("failed to upload csv file")
	}

	return url, nil
}

func (s *OperationService) MakeReportFile(ctx context.Context, month, year int) ([]byte, error) {
	products, amounts, err := s.operationRepo.GetAllRevenueOperationsGroupedByProduct(ctx, month, year)
	if err != nil {
		return nil, errors.New("failed to get revenue operations")
	}

	b := bytes.Buffer{}
	w := csv.NewWriter(&b)

	for i := range products {
		err := w.Write([]string{products[i], strconv.Itoa(amounts[i])})
		if err != nil {
			return nil, errors.New("failed to write csv")
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return nil, errors.New("failed to write csv")
	}

	return b.Bytes(), nil
}
