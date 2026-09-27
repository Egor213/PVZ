package usecase

type ExampleRepo interface {
}

// Реализация usecase
// func (uc ) Do() {}

/*
func (uc CreateOrder) Do(ctx context.Context, userID string, items []domain.Item) (domain.Order, error) {
	order := domain.Order{ID: uc.IDs.New(), UserID: userID, Items: items}

	if err := order.Validate(); err != nil { // правило проверяет domain
		return domain.Order{}, err
	}

	if err := uc.Repo.Save(ctx, order); err != nil { // вызов через интерфейс
		return domain.Order{}, fmt.Errorf("сохранить заказ: %w", err)
	}

	if err := uc.Notifier.OrderCreated(ctx, order); err != nil {
		return order, nil // заказ создан, уведомление не
	}

	return order, nil
}
*/
