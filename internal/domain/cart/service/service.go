package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/azka-zaydan/synapsis-test/configs"
	"github.com/azka-zaydan/synapsis-test/infras"
	"github.com/azka-zaydan/synapsis-test/internal/domain/cart/model"
	"github.com/azka-zaydan/synapsis-test/internal/domain/cart/model/dto"
	"github.com/azka-zaydan/synapsis-test/internal/domain/cart/repository"
	orderModel "github.com/azka-zaydan/synapsis-test/internal/domain/order/model"
	orderRepo "github.com/azka-zaydan/synapsis-test/internal/domain/order/repository"
	paymentModel "github.com/azka-zaydan/synapsis-test/internal/domain/payment/model"
	paymentRepo "github.com/azka-zaydan/synapsis-test/internal/domain/payment/repository"
	productRepo "github.com/azka-zaydan/synapsis-test/internal/domain/product/repository"

	"github.com/gofrs/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type CartService interface {
	ListItems(ctx context.Context, userID uuid.UUID) (res dto.ListItemsResponse, err error)
	CreateCart(ctx context.Context, req dto.CreateCartRequest) (res dto.CartResponse, err error)
	AddItems(ctx context.Context, req dto.AddItemsRequest, userID uuid.UUID) (res dto.ListItemsResponse, err error)
	DeleteItems(ctx context.Context, req dto.DeleteItemsRequest, userID uuid.UUID) (res dto.ListItemsResponse, err error)
	Checkout(ctx context.Context, req dto.CheckoutRequest, userID uuid.UUID) (res dto.CheckoutResponse, err error)
}

type CartServiceImpl struct {
	Repo        repository.CartRepository
	Redis       *infras.Redis
	config      *configs.Config
	OrderRepo   orderRepo.OrderRepository
	PaymentRepo paymentRepo.PaymentRepository
	ProductRepo productRepo.ProductRepository
}

func ProvideCartServiceImpl(repo repository.CartRepository, redis *infras.Redis, config *configs.Config, orderRepo orderRepo.OrderRepository, paymentRepo paymentRepo.PaymentRepository, productRepo productRepo.ProductRepository) *CartServiceImpl {
	return &CartServiceImpl{
		Redis:       redis,
		Repo:        repo,
		config:      config,
		OrderRepo:   orderRepo,
		PaymentRepo: paymentRepo,
		ProductRepo: productRepo,
	}
}

func (s *CartServiceImpl) ListItems(ctx context.Context, userID uuid.UUID) (res dto.ListItemsResponse, err error) {
	exist, err := s.isListItemsCacheAvailable(ctx, userID.String())
	if err != nil {
		log.Error().Err(err).Msg("[ListItems] Failed isListItemsCacheAvailable")
		return
	}
	if exist {
		log.Info().Msg("[ListItems] Using From Cache")
		res, err = s.getListItemsCache(ctx, userID.String())
		if err != nil {
			log.Error().Err(err).Msg("[ListItems] Failed getListItemsCache")
			return
		}
		return
	}
	cart, err := s.Repo.GetCartByUserID(ctx, userID.String())
	if err != nil {
		log.Error().Err(err).Msg("[ListItems] Failed GetCartByUserID")
		return
	}
	items, err := s.Repo.GetCartItemsByCartID(ctx, cart.ID.String())
	if err != nil {
		log.Error().Err(err).Msg("[GetCartItemsByCartID] Failed GetCartByUserID")
		return
	}

	res = dto.NewListItemsResponse(cart, items)

	err = s.setListItemsCache(ctx, userID.String(), res)
	if err != nil {
		log.Error().Err(err).Msg("[setListItemsCache] Failed setListItemsCache")
		return
	}

	return dto.NewListItemsResponse(cart, items), nil
}

func (s *CartServiceImpl) isListItemsCacheAvailable(ctx context.Context, userId string) (exist bool, err error) {
	_, err = s.Redis.Client.Get(ctx, fmt.Sprintf("cart:{%s}", userId)).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		log.Error().Err(err).Msg("[isListItemsCacheAvailable] Failed Getting From Redis")
		return false, err
	}
	return true, nil
}

func (s *CartServiceImpl) getListItemsCache(ctx context.Context, userId string) (res dto.ListItemsResponse, err error) {
	data, err := s.Redis.Client.Get(ctx, fmt.Sprintf("cart:{%s}", userId)).Result()
	if err != nil {
		log.Error().Err(err).Msg("[getListItemsCache] Failed Get Cache")
		return
	}
	err = json.Unmarshal([]byte(data), &res)
	if err != nil {
		log.Error().Err(err).Msg("[setListItemsCache] Failed Unmarshal Response")
		return
	}
	return
}

func (s *CartServiceImpl) setListItemsCache(ctx context.Context, userId string, res dto.ListItemsResponse) (err error) {
	marshaled, err := json.Marshal(res)
	if err != nil {
		log.Error().Err(err).Msg("[setListItemsCache] Failed Marshal Response")
		return
	}
	_, err = s.Redis.Client.Set(ctx, fmt.Sprintf("cart:{%s}", userId), marshaled, s.config.Cache.Cart.ExpiresIn).Result()
	if err != nil {
		log.Error().Err(err).Msg("[setListItemsCache] Failed Set Cache")
		return
	}
	return
}

func (s *CartServiceImpl) deleteListItemsCache(ctx context.Context, userId string) (err error) {
	_, err = s.Redis.Client.Del(ctx, fmt.Sprintf("cart:{%s}", userId)).Result()
	if err != nil {
		log.Error().Err(err).Msg("[deleteListItemsCache] Failed Del Cache")
		return
	}
	return
}

func (s *CartServiceImpl) CreateCart(ctx context.Context, req dto.CreateCartRequest) (res dto.CartResponse, err error) {
	cart, err := req.ToModel()
	if err != nil {
		log.Error().Err(err).Msg("[CreateCart] Failed Creating Model")
		return
	}
	err = s.Repo.CreateCart(ctx, &cart)
	if err != nil {
		log.Error().Err(err).Msg("[CreateCart] Failed Creating Cart")
		return
	}
	return dto.NewCartResponse(cart), nil
}

func (s *CartServiceImpl) GetCartByUserID(ctx context.Context, userId uuid.UUID) (res dto.CartResponse, err error) {
	cart, err := s.Repo.GetCartByUserID(ctx, userId.String())
	if err != nil {
		log.Error().Err(err).Msg("[GetCartByUserID] Failed GetCartByUserID")
	}
	return dto.NewCartResponse(cart), nil
}

func (s *CartServiceImpl) AddItems(ctx context.Context, req dto.AddItemsRequest, userID uuid.UUID) (res dto.ListItemsResponse, err error) {
	cart, err := s.Repo.GetCartByUserID(ctx, userID.String())
	if err != nil {
		log.Error().Err(err).Msg("[AddItems] Failed GetCartByUserID")
		return
	}

	existingItems, err := s.Repo.GetCartItemsByCartID(ctx, cart.ID.String())
	if err != nil {
		log.Error().Err(err).Msg("[AddItems] Failed GetCartItemsByCartID")
	}

	newItems, newCart, err := s.addOrUpdateItems(ctx, existingItems, req, cart)
	if err != nil {
		log.Error().Err(err).Msg("[AddItems] Failed addOrUpdateItems")
		return
	}
	err = s.deleteListItemsCache(ctx, userID.String())
	if err != nil {
		log.Error().Err(err).Msg("[AddItems] Failed deleteListItemsCache")
		return
	}
	return dto.NewListItemsResponse(newCart, newItems), nil
}

func (s *CartServiceImpl) addOrUpdateItems(ctx context.Context, existingItems []model.CartItem, req dto.AddItemsRequest, cart model.Cart) (newItems []model.CartItem, updatedCart model.Cart, err error) {
	existingItemsMap := make(map[string]*model.CartItem)
	for i, item := range existingItems {
		existingItemsMap[item.ProductID.String()] = &existingItems[i]
	}
	var mu sync.Mutex
	var wg sync.WaitGroup
	errCh := make(chan error, len(req))
	newItemsCh := make(chan model.CartItem, len(req))
	for _, newItem := range req {
		wg.Add(1)
		go func(newItem dto.ItemRequest) {
			defer wg.Done()

			mu.Lock()
			if existingItem, found := existingItemsMap[newItem.ProductID]; found {
				existingItem.Quantity += newItem.Quantity
				existingItem.TotalPrice += newItem.Price
				mu.Unlock()

				err = s.Repo.UpdateCartItem(ctx, existingItem)
				if err != nil {
					log.Error().Err(err).Msg("[addOrUpdateItems] Failed Update Cart Item")
					errCh <- err
				}

				mu.Lock()
				cart.TotalPrice += existingItem.TotalPrice
				newItemsCh <- *existingItem
				mu.Unlock()
			} else {
				mu.Unlock()
				newItemDto := dto.CartItemCreateRequest{
					CartID:    cart.ID.String(),
					ProductID: newItem.ProductID,
					Quantity:  newItem.Quantity,
					Price:     newItem.Price,
				}
				newItem, err := newItemDto.ToModel()
				if err != nil {
					log.Error().Err(err).Msg("[addOrUpdateItems] Failed Create Model")
					errCh <- err
				}
				err = s.Repo.CreateCartItem(ctx, &newItem)
				if err != nil {
					log.Error().Err(err).Msg("[addOrUpdateItems] Failed Create Cart Item")
					errCh <- err
				}
				mu.Lock()
				cart.TotalPrice += newItem.TotalPrice
				newItemsCh <- newItem
				mu.Unlock()
			}
		}(newItem)

	}
	wg.Wait()
	close(errCh)
	close(newItemsCh)

	for err := range errCh {
		if err != nil {
			log.Error().Err(err).Msg("[addOrUpdateItems] Error Occured")
			return make([]model.CartItem, 0), model.Cart{}, err
		}
	}

	for item := range newItemsCh {
		newItems = append(newItems, item)
	}

	cart.TotalItems = len(newItems)
	if err := s.Repo.UpdateCart(ctx, &cart); err != nil {
		log.Error().Err(err).Msg("[addOrUpdateItems] Failed Update Cart")
		return make([]model.CartItem, 0), model.Cart{}, err
	}

	return newItems, cart, nil
}

func (s *CartServiceImpl) DeleteItems(ctx context.Context, req dto.DeleteItemsRequest, userID uuid.UUID) (res dto.ListItemsResponse, err error) {
	cart, err := s.Repo.GetCartByUserID(ctx, userID.String())
	if err != nil {
		log.Error().Err(err).Msg("[DeleteItems] Failed GetCartByUserID")
		return
	}

	existingItems, err := s.Repo.GetCartItemsByCartID(ctx, cart.ID.String())
	if err != nil {
		log.Error().Err(err).Msg("[DeleteItems] Failed GetCartItemsByCartID")
	}

	if len(existingItems) == 0 {
		return
	}

	newItems, newCart, err := s.removeOrUpdateItems(ctx, existingItems, req, cart)
	if err != nil {
		log.Error().Err(err).Msg("[DeleteItems] Failed removeOrUpdateItems")
		return
	}
	err = s.deleteListItemsCache(ctx, userID.String())
	if err != nil {
		log.Error().Err(err).Msg("[DeleteItems] Failed deleteListItemsCache")
		return
	}
	return dto.NewListItemsResponse(newCart, newItems), nil
}

func (s *CartServiceImpl) removeOrUpdateItems(ctx context.Context, existingItems []model.CartItem, req dto.DeleteItemsRequest, cart model.Cart) (updatedItems []model.CartItem, updatedCart model.Cart, err error) {
	existingItemsMap := make(map[string]*model.CartItem)
	for i, item := range existingItems {
		existingItemsMap[item.ID.String()] = &existingItems[i]
	}
	var mu sync.Mutex
	var wg sync.WaitGroup
	errCh := make(chan error, len(req))
	updatedItemsCh := make(chan model.CartItem, len(req))

	for _, v := range req {
		wg.Add(1)
		go func(v dto.DeleteItemRequest) {
			defer wg.Done()

			mu.Lock()
			existingItem, found := existingItemsMap[v.ItemId]
			mu.Unlock()
			if !found {
				return
			}
			mu.Lock()
			if v.Quantity >= existingItem.Quantity {
				err := s.Repo.DeleteCartItem(ctx, existingItem.ID.String())
				if err != nil {
					errCh <- err
					mu.Unlock()
					return
				}
				cart.TotalPrice -= existingItem.TotalPrice
				cart.TotalItems -= 1
			} else {
				existingItem.Quantity -= v.Quantity
				existingItem.TotalPrice -= v.Price

				err := s.Repo.UpdateCartItem(ctx, existingItem)
				if err != nil {
					errCh <- err
					mu.Unlock()
					return
				}
				cart.TotalPrice -= v.Price
				updatedItemsCh <- *existingItem
			}
			mu.Unlock()
		}(v)
	}

	wg.Wait()
	close(errCh)
	close(updatedItemsCh)

	for err := range errCh {
		if err != nil {
			log.Error().Err(err).Msg("[removeOrUpdateItems] Error Occured")
			return make([]model.CartItem, 0), model.Cart{}, err
		}
	}

	for item := range updatedItemsCh {
		updatedItems = append(updatedItems, item)
	}

	cart.TotalItems = len(updatedItems)
	err = s.Repo.UpdateCart(ctx, &cart)
	if err != nil {
		log.Error().Err(err).Msg("[removeOrUpdateItems] Failed Update Cart")
		return make([]model.CartItem, 0), model.Cart{}, err
	}
	return updatedItems, cart, nil
}

func (s *CartServiceImpl) Checkout(ctx context.Context, req dto.CheckoutRequest, userID uuid.UUID) (res dto.CheckoutResponse, err error) {
	cart, err := s.Repo.GetCartByUserID(ctx, userID.String())
	if err != nil {
		log.Error().Err(err).Msg("[Checkout] Failed GetCartByUserID")
		return
	}
	items, err := s.Repo.GetCartItemsByCartID(ctx, cart.ID.String())
	if err != nil {
		log.Error().Err(err).Msg("[Checkout] Failed GetCartItemsByCartID")
		return
	}
	if len(items) == 0 {
		return
	}

	order, totalItems, err := s.parseCheckoutItems(ctx, items, req, cart)
	if err != nil {
		log.Error().Err(err).Msg("[Checkout] Failed parseCheckoutItems")
		return
	}

	return dto.NewCheckoutResponse(order.ID.String(), order.OrderAt, totalItems, order.TotalPrice), nil
}

func (s *CartServiceImpl) parseCheckoutItems(ctx context.Context, existingItems []model.CartItem, req dto.CheckoutRequest, cart model.Cart) (res orderModel.Order, totalItems int, err error) {
	existingItemsMap := make(map[string]*model.CartItem)
	for i, item := range existingItems {
		existingItemsMap[item.ID.String()] = &existingItems[i]
	}
	var mu sync.Mutex
	var wg sync.WaitGroup
	errCh := make(chan error, len(req))

	var order orderModel.Order
	var payment paymentModel.Payment
	orderId, _ := uuid.NewV4()
	paymentId, _ := uuid.NewV4()
	order.ID = orderId
	order.UserID = cart.UserID
	order.Status = int(orderModel.OrderPlacedStatus)
	order.PaymentID = uuid.NullUUID{UUID: paymentId, Valid: true}
	order.CreatedBy = cart.UserID
	order.UpdatedBy = cart.UserID
	order.OrderAt = time.Now()
	payment.ID = paymentId
	payment.OrderID = orderId
	payment.CreatedBy = cart.UserID
	payment.UpdatedBy = cart.UserID

	for _, v := range req {
		wg.Add(1)
		go func(v dto.CheckoutItem) {
			defer wg.Done()

			mu.Lock()
			existingItem, found := existingItemsMap[v.ItemId]
			mu.Unlock()
			if !found {
				return
			}
			prod, err := s.ProductRepo.GetProductByID(ctx, v.ProductId)
			if err != nil {
				errCh <- err
				return
			}

			mu.Lock()
			if v.Quantity == existingItem.Quantity && prod.Stock >= existingItem.Quantity {
				id, _ := uuid.NewV4()
				detail := orderModel.OrderDetail{
					ID:                   id,
					OrderID:              order.ID,
					ProductID:            prod.ID,
					TotalItems:           v.Quantity,
					SubtotalProductPrice: v.Price,
					CreatedBy:            order.ID,
					UpdatedBy:            order.ID,
				}
				err := s.Repo.DeleteCartItem(ctx, existingItem.ID.String())
				if err != nil {
					errCh <- err
					mu.Unlock()
					return
				}
				err = s.OrderRepo.CreateOrderDetail(ctx, &detail)
				if err != nil {
					errCh <- err
					mu.Unlock()
					return
				}
				prod.Stock -= v.Quantity
				order.TotalPrice += v.Price
				cart.TotalPrice -= existingItem.TotalPrice
				cart.TotalItems -= 1
				totalItems += 1
			}
			mu.Unlock()
		}(v)
	}
	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			log.Error().Err(err).Msg("[parseCheckoutItems] Error Occured")
			return res, 0, err
		}
	}

	payment.TotalPrice = order.TotalPrice
	payment.Status = int(paymentModel.Unpaid)
	payment.PaymentMethod = ""
	payment.UserID = cart.UserID

	err = s.OrderRepo.CreateOrder(ctx, &order)
	if err != nil {
		log.Error().Err(err).Msg("[parseCheckoutItems] Failed Create Order")
		return
	}

	err = s.PaymentRepo.CreatePayment(ctx, &payment)
	if err != nil {
		log.Error().Err(err).Msg("[parseCheckoutItems] Failed Create Payment")
		return
	}
	err = s.Repo.UpdateCart(ctx, &cart)
	if err != nil {
		log.Error().Err(err).Msg("[parseCheckoutItems] Failed Update Cart")
		return
	}
	return order, totalItems, nil
}
