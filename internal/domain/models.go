package domain

// CustomerType определяет тип клиента
type CustomerType string

const (
	//CustomerTypeNew - клиент, который еще не совершал покупки
	CustomerTypeNew CustomerType = "new"
	//CustomerTypeOld - клиент, который уже совершал покупки
	CustomerTypeOld CustomerType = "old"
)

type (
	Book struct {
		ID    string
		Price float64
	}
	Discount struct {
		CartAmount      float64 // Стоимость корзины
		DiscountPercent float64 // Размер скидки (0,1 или 0,05)
		TotalCost       float64 // Итоговая стоимость корзины
		DiscountAmount  float64 // Итоговая сумма скидки
	}
)
