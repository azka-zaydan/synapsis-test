package service

import (
	"context"
	"sync"

	"github.com/azka-zaydan/synapsis-test/internal/domain/cart/model"
	"github.com/azka-zaydan/synapsis-test/internal/domain/cart/model/dto"
	"github.com/azka-zaydan/synapsis-test/internal/domain/cart/repository"
	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"
)

type CartService interface {
	ListItems(ctx context.Context, userID uuid.UUID) (res dto.ListItemsResponse, err error)
	CreateCart(ctx context.Context, req dto.CreateCartRequest) (res dto.CartResponse, err error)
	AddItems(ctx context.Context, req dto.AddItemsRequest, userID uuid.UUID) (res dto.ListItemsResponse, err error)
	DeleteItems(ctx context.Context, req dto.DeleteItemsRequest, userID uuid.UUID) (res dto.ListItemsResponse, err error)
}

type CartServiceImpl struct {
	Repo repository.CartRepository
}

func ProvideCartServiceImpl(repo repository.CartRepository) *CartServiceImpl {
	return &CartServiceImpl{
		Repo: repo,
	}
}

func (s *CartServiceImpl) ListItems(ctx context.Context, userID uuid.UUID) (res dto.ListItemsResponse, err error) {
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

	return dto.NewListItemsResponse(cart, items), nil
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

	newItems, newCart, err := s.removeOrUpdateItems(ctx, existingItems, req, cart)
	if err != nil {
		log.Error().Err(err).Msg("[DeleteItems] Failed removeOrUpdateItems")
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
