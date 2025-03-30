// Copyright 2024 CloudWeGo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package service

import (
	"context"
	"errors"

	// "fmt"

	"github.com/cloudwego/biz-demo/gomall/app/order/biz/dal/mysql"
	"github.com/cloudwego/biz-demo/gomall/app/order/biz/model"
	order "github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/order"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/google/uuid"
	// "gorm.io/gorm"
)

type PlaceOrderService struct {
	ctx context.Context
} // NewPlaceOrderService new PlaceOrderService
func NewPlaceOrderService(ctx context.Context) *PlaceOrderService {
	return &PlaceOrderService{ctx: ctx}
}

// Run create note info
func (s *PlaceOrderService) Run(req *order.PlaceOrderReq) (resp *order.PlaceOrderResp, err error) {
	// 业务逻辑：插入数据到数据库中的order表和order_item表，生成一个随机的uuid作为订单号

	// 参数校验
	if req.GetUserId() == 0 {
		return nil, kerrors.NewBizStatusError(400, "invalid user id")
	}
	if len(req.GetOrderItems()) == 0 {
		return nil, kerrors.NewBizStatusError(400, "empty order items")
	}

	// 生成订单ID
	orderId, err := uuid.NewRandom()
	if err != nil {
		return nil, kerrors.NewBizStatusError(500, "generate order id failed")
	}

	// 创建订单项数组，作为主订单的订单项属性
	var orderItems []model.OrderItem
	for _, pbItem := range req.GetOrderItems() {
		// GetOrderItems里的item是protobuf自动生成的结构
		// 用于gRPC，不是OrderItem的数据库表结构，要转化为数据库模型
		if pbItem == nil || pbItem.Item == nil {
			return nil, errors.New("invalid order item")
		}

		// pbItem包含指向cart中物品的指针Item，通过该指针可以读到物品的Id和数量
		orderItems = append(orderItems, model.OrderItem{
			OrderIdRefer: orderId.String(),
			ProductId:    pbItem.Item.GetProductId(),
			Quantity:     pbItem.Item.GetQuantity(),
			Cost:         pbItem.GetCost(),
		})
	}

	// 创建主订单
	err = model.CreateOrder(mysql.DB, s.ctx, &model.Order{
		OrderId:      orderId.String(),
		UserId:       req.GetUserId(),
		UserCurrency: req.GetUserCurrency(),
		Consignee: model.Consignee{
			Email:         req.GetEmail(),
			StreetAddress: req.GetAddress().GetStreetAddress(),
			City:          req.GetAddress().GetCity(),
			State:         req.GetAddress().GetState(),
			Country:       req.GetAddress().GetCountry(),
			ZipCode:       req.GetAddress().GetZipCode(),
		},
		OrderItems: orderItems,
	})
	if err != nil {
		return nil, err
	}

	// // 在主订单后，将订单项写入数据库
	// for _, item := range orderItems {
	// 	err = model.CreateOrderItem(mysql.DB, s.ctx, item)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }
	return &order.PlaceOrderResp{Order: &order.OrderResult{OrderId: orderId.String()}}, nil
}
