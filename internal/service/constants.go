package service

// Discount percentages
const (
	// discountPercentNewCustomer2to5 - скидка 20% для нового клиента при покупке 2-5 книг
	discountPercentNewCustomer2to5 = 0.20

	// discountPercentOldCustomer2to5 - скидка 10% для старого клиента при покупке 2-5 книг
	discountPercentOldCustomer2to5 = 0.10

	// discountPercentOldCustomer6to10 - скидка 5% для старого клиента при покупке 6-10 книг
	discountPercentOldCustomer6to10 = 0.05

	// discountPercentOldCustomerOver10 - скидка 2% для старого клиента при покупке более 10 книг
	discountPercentOldCustomerOver10 = 0.02

	// discountPercentFriday - скидка 5% в пятницу для всех клиентов
	discountPercentFriday = 0.05

	// discountPercentBulkBook - скидка 40% на каждую вторую книгу при покупке более 5 экземпляров одной книги
	discountPercentBulkBook = 0.40
)

// Book quantity thresholds
const (
	// minBooksForNewCustomerDiscount - минимальное количество книг для скидки новому клиенту
	minBooksForNewCustomerDiscount = 2

	// maxBooksForNewCustomerDiscount - максимальное количество книг для скидки новому клиенту
	maxBooksForNewCustomerDiscount = 5

	// minBooksForOldCustomerSmallDiscount - минимальное количество книг для малой скидки старому клиенту
	minBooksForOldCustomerSmallDiscount = 2

	// maxBooksForOldCustomerSmallDiscount - максимальное количество книг для малой скидки старому клиенту
	maxBooksForOldCustomerSmallDiscount = 5

	// minBooksForOldCustomerMediumDiscount - минимальное количество книг для средней скидки старому клиенту
	minBooksForOldCustomerMediumDiscount = 6

	// maxBooksForOldCustomerMediumDiscount - максимальное количество книг для средней скидки старому клиенту
	maxBooksForOldCustomerMediumDiscount = 10

	// minBooksForOldCustomerLargeDiscount - минимальное количество книг для большой скидки старому клиенту
	minBooksForOldCustomerLargeDiscount = 11

	// minBooksForBulkDiscount - минимальное количество экземпляров одной книги для применения оптовой скидки
	minBooksForBulkDiscount = 6
)
