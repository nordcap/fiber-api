package routes

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nordcap/fiber-api/database"
	"github.com/nordcap/fiber-api/models"
)

// используем структуру как сериализатор
type Order struct {
	ID        uint      `json:"id"`
	User      User      `json: "user"`
	Product   Product   `json: "product"`
	CreatedAt time.Time `json: "order_date"`
}

// почему возникает ошибка - нельзя передавать models.Product?
// разве это не одно и тоже что Product
func CreateResponseOrder(order models.Order, user User, product Product) Order {
	return Order{ID: order.ID, User: user, Product: product, CreatedAt: order.CreatedAt}
}

func CreateOrder(c *fiber.Ctx) error {
	var order models.Order

	if err := c.BodyParser(&order); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	//проверка что user существует
	var user models.User

	if err := findUser(int(order.UserRefer), &user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	//проверка что product существует
	var product models.Product

	if err := findProduct(int(order.ProductRefer), &product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	database.Database.Db.Create(&order)

	responseUser := CreateResponseUser(user)
	responseProduct := CreateResponseProduct(product)
	responseOrder := CreateResponseOrder(order, responseUser, responseProduct)

	return c.Status(fiber.StatusOK).JSON(responseOrder)
}

func GetOrders(c *fiber.Ctx) error {
	orders := []models.Order{}
	database.Database.Db.Find(&orders)

	responseOrders := []Order{}
	for _, order := range orders {
		var user models.User
		var product models.Product

		database.Database.Db.Find(&user, "id = ?", order.UserRefer)
		database.Database.Db.Find(&product, "id = ?", order.ProductRefer)
		responseOrder := CreateResponseOrder(order, CreateResponseUser(user), CreateResponseProduct(product))
		responseOrders = append(responseOrders, responseOrder)
	}

	return c.Status(fiber.StatusOK).JSON(responseOrders)

}

func findOrder(id int, order *models.Order) error {
	database.Database.Db.Find(order, "id = ?", id)
	if order.ID == 0 {
		return errors.New("заказ не существует")
	}
	return nil

}

func GetOrder(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var order models.Order

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("Ожидается параметр :id integer")
	}

	if err := findOrder(id, &order); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	var user models.User
	var product models.Product

	database.Database.Db.First(&user, order.UserRefer)
	database.Database.Db.First(&product, order.ProductRefer)
	responseOrder := CreateResponseOrder(order, CreateResponseUser(user), CreateResponseProduct(product))
	return c.Status(fiber.StatusOK).JSON(responseOrder)
}
