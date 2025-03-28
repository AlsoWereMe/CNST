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

package model

import (
	"context"
	"time"

	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/cart"
	order "github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/order"
	"gorm.io/gorm"
)

type Consignee struct {
	Email string

	StreetAddress string
	City          string
	State         string
	Country       string
	ZipCode       int32
}

type Order struct {
	Base
	OrderId      string `gorm:"uniqueIndex;size:256"`
	UserId       uint32
	UserCurrency string
	Consignee    Consignee   `gorm:"embedded"`
	OrderItems   []OrderItem `gorm:"foreignKey:OrderIdRefer;references:OrderId"`
}

func (o Order) TableName() string {
	return "order"
}

func CreateOrder(db *gorm.DB, ctx context.Context, order *Order) error {
	return db.WithContext(ctx).Model(&Order{}).Create(order).Error
}

func ListOrder(db *gorm.DB, ctx context.Context, userId uint32) (orders []Order, err error) {
	err = db.Debug().WithContext(ctx).Model(&Order{}).Where(&Order{UserId: userId}).Preload("OrderItems").Find(&orders).Error
	return
}

// 将数据层模型转为协议层模型
func ToProtoOrder(m Order) *order.Order {
	protoOrder := &order.Order{
		OrderId:      m.OrderId,
		UserId:       m.UserId,
		UserCurrency: m.UserCurrency,
		Email:        m.Consignee.Email,
		CreatedAt:    convertTime(m.Base.CreatedAt),
		Address:      convertAddress(m.Consignee),
		OrderItems:   convertOrderItems(m.OrderItems),
	}
	return protoOrder
}

// 地址转换
func convertAddress(c Consignee) *order.Address {
	return &order.Address{
		StreetAddress: c.StreetAddress,
		City:          c.City,
		State:         c.State,
		Country:       c.Country,
		ZipCode:       c.ZipCode,
	}
}

// 时间转换
func convertTime(t time.Time) int32 {
	if t.IsZero() {
		return 0
	}
	return int32(t.Unix())
}

// OrderItem转换
func convertOrderItems(items []OrderItem) []*order.OrderItem {
	protoItems := make([]*order.OrderItem, 0, len(items))
	for _, item := range items {
		cartItem := cart.CartItem{
			ProductId: item.ProductId,
			Quantity:  item.Quantity,
		}
		protoItems = append(protoItems, &order.OrderItem{
			Item: &cartItem,
			Cost: item.Cost,
		})
	}
	return protoItems
}
