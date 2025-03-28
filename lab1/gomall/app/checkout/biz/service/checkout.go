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
	"strconv"

	// "errors"
	// "fmt"
	// "strconv"

	// "github.com/cloudwego/biz-demo/gomall/app/checkout/infra/rpc"
	"github.com/cloudwego/biz-demo/gomall/app/checkout/infra/mq"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"

	cart "github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/cart"
	cartservice "github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/cart/cartservice"
	checkout "github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/checkout"
	email "github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/email"
	order "github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/order"
	orderservice "github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/order/orderservice"
	payment "github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/payment"
	paymentservice "github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/payment/paymentservice"
	product "github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/product"
	productcatalogservice "github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/product/productcatalogservice"
)

type CheckoutService struct {
	ctx context.Context
} // NewCheckoutService new CheckoutService
func NewCheckoutService(ctx context.Context) *CheckoutService {
	return &CheckoutService{ctx: ctx}
}

/*
	Run

1. get cart
2. calculate cart
3. create order
4. empty cart
5. pay
6. change order result
7. finish
*/
func (s *CheckoutService) Run(req *checkout.CheckoutReq) (resp *checkout.CheckoutResp, err error) {
	// 1. get cart (使用RPC调用Cart服务以获得购物车信息)
	userId := req.GetUserId()
	if userId == 0 {
		return nil, kerrors.NewBizStatusError(400, "invalid user id")
	}
	cartClient, _ := cartservice.NewClient(
		"cart",
		client.WithHostPorts("127.0.0.1:8883"),
	)
	cartResp, err := cartClient.GetCart(s.ctx, &cart.GetCartReq{UserId: userId})
	if err != nil {
		return nil, err
	}
	userCart := cartResp.GetCart()

	// 2. calc cart（根据第1步的购物车信息，计算总价和订单项信息）
	cartItems := userCart.GetItems()
	if cartItems == nil {
		return nil, kerrors.NewBizStatusError(400, "invalid cart items")
	}
	productClient, _ := productcatalogservice.NewClient(
		"product",
		client.WithHostPorts("127.0.0.1:8881"),
	)

	// 计算购物车中某商品总价cost与全部总价amount，而后构建订单项
	var orderItems []*order.OrderItem
	var amount float32
	for _, item := range cartItems {
		productResp, err := productClient.GetProduct(s.ctx, &product.GetProductReq{Id: item.GetProductId()})
		if err != nil {
			return nil, err
		}
		price := productResp.GetProduct().GetPrice()
		cost := price * float32(item.GetQuantity())
		amount += cost
		orderItems = append(orderItems, &order.OrderItem{
			Item: item,
			Cost: cost,
		})
	}

	// 3. create order（根据第1步和第2步的信息，创建order.PlaceOrderReq，并使用RPC调用Order服务创建订单）
	orderClient, _ := orderservice.NewClient(
		"order",
		client.WithHostPorts("127.0.0.1:8885"),
	)

	// 将数据层的string型ZipCode转为协议层的int32型
	zipCode, err := strconv.ParseInt(req.GetAddress().GetZipCode(), 10, 32)
	if err != nil {
		return nil, err
	}

	// 将checkout协议层的Address转为order协议层的Address
	address := &order.Address{
		StreetAddress: req.GetAddress().GetStreetAddress(),
		City:          req.GetAddress().GetCity(),
		State:         req.GetAddress().GetState(),
		Country:       req.GetAddress().GetCountry(),
		ZipCode:       int32(zipCode),
	}

	// 创建订单
	orderResp, err := orderClient.PlaceOrder(s.ctx, &order.PlaceOrderReq{
		UserId:       userId,
		UserCurrency: "CNY",
		Address:      address,
		Email:        req.GetEmail(),
		OrderItems:   orderItems,
	})
	if err != nil {
		return nil, err
	}

	// 4. empty cart（使用RPC调用Cart服务清空购物车）
	cartClient.EmptyCart(s.ctx, &cart.EmptyCartReq{UserId: userId})

	// 5. pay（使用RPC调用Payment服务进行支付）
	paymentClient, _ := paymentservice.NewClient(
		"payment",
		client.WithHostPorts("127.0.0.1:8886"),
	)

	chargeResp, err := paymentClient.Charge(s.ctx, &payment.ChargeReq{
		Amount:     amount,
		CreditCard: req.GetCreditCard(),
		OrderId:    orderResp.GetOrder().GetOrderId(),
		UserId:     userId,
	})
	if err != nil {
		return nil, err
	}

	// 6. send email（使用MQ发送邮件通知）
	data, _ := proto.Marshal(&email.EmailReq{
		From:        "22302010051@m.fudan.edu.cn",
		To:          req.Email,
		ContentType: "text/plain",
		Subject:     "You just created an order in CloudWeGo shop",
		Content:     "You just created an order in CloudWeGo shop",
	})
	msg := &nats.Msg{Subject: "email", Data: data}
	err = mq.Nc.PublishMsg(msg)
	if err != nil {
		klog.Errorf("Failed to send message: %v", err)
	}

	// 7. finish（返回订单ID和支付结果）
	return &checkout.CheckoutResp{
		OrderId:       orderResp.GetOrder().GetOrderId(),
		TransactionId: chargeResp.GetTransactionId(),
	}, nil
}
