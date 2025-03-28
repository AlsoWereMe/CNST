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

	// "github.com/cloudwego/biz-demo/gomall/app/order/biz/dal/mysql"
	// "github.com/cloudwego/biz-demo/gomall/app/order/biz/model"
	// "github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/cart"
	"github.com/cloudwego/biz-demo/gomall/app/order/biz/dal/mysql"
	"github.com/cloudwego/biz-demo/gomall/app/order/biz/model"
	order "github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/order"
	"github.com/cloudwego/kitex/pkg/kerrors"
)

type ListOrderService struct {
	ctx context.Context
} // NewListOrderService new ListOrderService
func NewListOrderService(ctx context.Context) *ListOrderService {
	return &ListOrderService{ctx: ctx}
}

// Run create note info
func (s *ListOrderService) Run(req *order.ListOrderReq) (resp *order.ListOrderResp, err error) {
	// TODO 请实现ListOrder的业务逻辑，从数据库中的order表和order_item表中查询数据
	// 参数校验
	if req.GetUserId() == 0 {
		return nil, kerrors.NewBizStatusError(400, "invalid user id")
	}

	// 根据uid,从数据库order表中获取订单列表orders
	var orders []model.Order
	orders, err = model.ListOrder(mysql.DB, s.ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}

	// 处理orders格式为协议层后返回resp
	var resOrders []*order.Order
	for _, o := range orders {
		resOrders = append(resOrders, model.ToProtoOrder(o))
	}
	return &order.ListOrderResp{Orders: resOrders}, nil
}
