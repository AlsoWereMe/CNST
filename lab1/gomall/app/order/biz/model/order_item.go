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

	"gorm.io/gorm"
)

type OrderItem struct {
	Base
	ProductId    uint32
	OrderIdRefer string `gorm:"size:256;index"`
	Quantity     int32
	Cost         float32
}

func (oi OrderItem) TableName() string {
	return "order_item"
}

func CreateOrderItem(db *gorm.DB, ctx context.Context, item OrderItem) error {
	return db.WithContext(ctx).
		Model(&OrderItem{}).
		Create(&item).
		Error
}
