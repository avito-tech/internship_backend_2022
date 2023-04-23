package service

import (
	"account-management-service/internal/entity"
	"account-management-service/internal/mocks/repomocks"
	"account-management-service/internal/mocks/webapimocks"
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestOperationService_OperationHistory(t *testing.T) {
	type args struct {
		ctx   context.Context
		input OperationHistoryInput
	}

	type MockBehavior func(o *repomocks.MockOperation, p *repomocks.MockProduct, g *webapimocks.MockGDrive, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		want         []OperationHistoryOutput
		wantErr      bool
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				input: OperationHistoryInput{
					AccountId: 1,
				},
			},
			mockBehavior: func(o *repomocks.MockOperation, p *repomocks.MockProduct, g *webapimocks.MockGDrive, args args) {
				o.EXPECT().OperationsPagination(args.ctx, args.input.AccountId, args.input.SortType, args.input.Offset, args.input.Limit).
					Return([]entity.Operation{
						{
							Id:            1,
							AccountId:     1,
							Amount:        100,
							OperationType: entity.OperationTypeDeposit,
							CreatedAt:     time.UnixMilli(123456),
							ProductId:     nil,
							OrderId:       nil,
							Description:   "",
						},
					}, []string{
						"some product name",
					}, nil)
			},
			want: []OperationHistoryOutput{
				{
					Amount:      100,
					Operation:   "deposit",
					Time:        time.UnixMilli(123456),
					Product:     "some product name",
					Order:       nil,
					Description: "",
				},
			},
			wantErr: false,
		},
		{
			name: "operations pagination error",
			args: args{
				ctx: context.Background(),
				input: OperationHistoryInput{
					AccountId: 1,
				},
			},
			mockBehavior: func(o *repomocks.MockOperation, p *repomocks.MockProduct, g *webapimocks.MockGDrive, args args) {
				o.EXPECT().OperationsPagination(args.ctx, args.input.AccountId, args.input.SortType, args.input.Offset, args.input.Limit).
					Return(nil, nil, errors.New("some error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// init deps
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// init mocks
			operationRepo := repomocks.NewMockOperation(ctrl)
			productRepo := repomocks.NewMockProduct(ctrl)
			gDrive := webapimocks.NewMockGDrive(ctrl)
			tc.mockBehavior(operationRepo, productRepo, gDrive, tc.args)

			// init service
			s := NewOperationService(operationRepo, productRepo, gDrive)

			// run test
			got, err := s.OperationHistory(tc.args.ctx, tc.args.input)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestOperationService_MakeReportFile(t *testing.T) {
	type args struct {
		ctx   context.Context
		month int
		year  int
	}

	type MockBehavior func(o *repomocks.MockOperation, p *repomocks.MockProduct, g *webapimocks.MockGDrive, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		want         []byte
		wantErr      bool
	}{
		{
			name: "OK",
			args: args{
				ctx:   context.Background(),
				month: 1,
				year:  2021,
			},
			mockBehavior: func(o *repomocks.MockOperation, p *repomocks.MockProduct, g *webapimocks.MockGDrive, args args) {
				o.EXPECT().GetAllRevenueOperationsGroupedByProduct(args.ctx, args.month, args.year).
					Return([]string{
						"some product name",
					}, []int{
						100,
					}, nil)
			},
			want:    []byte("some product name,100\n"),
			wantErr: false,
		},
		{
			name: "get all revenue operations grouped by product error",
			args: args{
				ctx:   context.Background(),
				month: 1,
				year:  2021,
			},
			mockBehavior: func(o *repomocks.MockOperation, p *repomocks.MockProduct, g *webapimocks.MockGDrive, args args) {
				o.EXPECT().GetAllRevenueOperationsGroupedByProduct(args.ctx, args.month, args.year).
					Return(nil, nil, errors.New("some error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// init deps
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// init mocks
			operationRepo := repomocks.NewMockOperation(ctrl)
			productRepo := repomocks.NewMockProduct(ctrl)
			gDrive := webapimocks.NewMockGDrive(ctrl)
			tc.mockBehavior(operationRepo, productRepo, gDrive, tc.args)

			// init service
			s := NewOperationService(operationRepo, productRepo, gDrive)

			// run test
			got, err := s.MakeReportFile(tc.args.ctx, tc.args.month, tc.args.year)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestOperationService_MakeReportLink(t *testing.T) {
	type args struct {
		ctx   context.Context
		month int
		year  int
	}

	type MockBehavior func(o *repomocks.MockOperation, p *repomocks.MockProduct, g *webapimocks.MockGDrive, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		want         string
		wantErr      bool
	}{
		{
			name: "OK",
			args: args{
				ctx:   context.Background(),
				month: 1,
				year:  2021,
			},
			mockBehavior: func(o *repomocks.MockOperation, p *repomocks.MockProduct, g *webapimocks.MockGDrive, args args) {
				g.EXPECT().IsAvailable().Return(true)

				o.EXPECT().GetAllRevenueOperationsGroupedByProduct(args.ctx, args.month, args.year).
					Return([]string{
						"some product name",
					}, []int{
						100,
					}, nil)

				g.EXPECT().UploadCSVFile(args.ctx, "report_1_2021.csv", []byte("some product name,100\n")).
					Return("https://example.com", nil)
			},
			want:    "https://example.com",
			wantErr: false,
		},
		{
			name: "gdrive is not available",
			args: args{
				ctx:   context.Background(),
				month: 1,
				year:  2021,
			},
			mockBehavior: func(o *repomocks.MockOperation, p *repomocks.MockProduct, g *webapimocks.MockGDrive, args args) {
				g.EXPECT().IsAvailable().Return(false)
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "get all revenue operations grouped by product error",
			args: args{
				ctx:   context.Background(),
				month: 1,
				year:  2021,
			},
			mockBehavior: func(o *repomocks.MockOperation, p *repomocks.MockProduct, g *webapimocks.MockGDrive, args args) {
				g.EXPECT().IsAvailable().Return(true)

				o.EXPECT().GetAllRevenueOperationsGroupedByProduct(args.ctx, args.month, args.year).
					Return(nil, nil, errors.New("some error"))
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "upload csv file error",
			args: args{
				ctx:   context.Background(),
				month: 1,
				year:  2021,
			},
			mockBehavior: func(o *repomocks.MockOperation, p *repomocks.MockProduct, g *webapimocks.MockGDrive, args args) {
				g.EXPECT().IsAvailable().Return(true)

				o.EXPECT().GetAllRevenueOperationsGroupedByProduct(args.ctx, args.month, args.year).
					Return([]string{
						"some product name",
					}, []int{
						100,
					}, nil)

				g.EXPECT().UploadCSVFile(args.ctx, "report_1_2021.csv", []byte("some product name,100\n")).
					Return("", errors.New("some error"))
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// init deps
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// init mocks
			operationRepo := repomocks.NewMockOperation(ctrl)
			productRepo := repomocks.NewMockProduct(ctrl)
			gDrive := webapimocks.NewMockGDrive(ctrl)
			tc.mockBehavior(operationRepo, productRepo, gDrive, tc.args)

			// init service
			s := NewOperationService(operationRepo, productRepo, gDrive)

			// run test
			got, err := s.MakeReportLink(tc.args.ctx, tc.args.month, tc.args.year)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}
