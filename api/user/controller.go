package user

import (
	"net/http"
	userBusiness "roby-backend-golang/business/user"
	"roby-backend-golang/utils"

	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	service userBusiness.Service
}

func NewController(service userBusiness.Service) *Controller {
	return &Controller{
		service: service,
	}
}

// ShowAccount godoc
// @Summary      Show an account
// @Description  get string by ID
// @Tags         account
// @Accept       json
// @Produce      json
// @Param	     AuthLogin body user.AuthLogin true "Login"
// @Success      200 {object} user.AuthLogin
// @Failure      400 {object} user.AuthLogin
// @Router 		/user/login [post]
func (Controller *Controller) Login(c *fiber.Ctx) error {
	var auth userBusiness.AuthLogin
	if err := c.BodyParser(&auth); err != nil {
		return err
	}
	res, err := Controller.service.Login(auth)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    res.Token.AccessToken,
		MaxAge:   res.Token.AccessTokenExpired,
		HTTPOnly: true,
		Secure:   true,
	})
	return c.JSON(fiber.Map{
		"code":    200,
		"message": "success login",
		"result":  res,
	})
}

func (Controller *Controller) Register(c *fiber.Ctx) error {
	var data userBusiness.Register
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var err error
	data.File, err = c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"code":    400,
			"message": "please upload foto profile",
		})
	}

	err = Controller.service.RegisterUser(data)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"code":    400,
			"message": err.Error(),
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"code":    200,
		"message": "success register",
	})
}

func (Controller *Controller) Logout(c *fiber.Ctx) error {
	c.ClearCookie()
	return c.Status(200).JSON(fiber.Map{
		"code":    200,
		"message": "success logout",
	})
}

func (Controller *Controller) GetMe(c *fiber.Ctx) error {
	id := c.Locals("id").(string)
	res, err := Controller.service.GetMe(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"code":    http.StatusBadRequest,
			"message": "invalid token",
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"code":    200,
		"message": "success get data",
		"result":  res,
	})
}

func (Controller *Controller) GetRandomUser(c *fiber.Ctx) error {
	id := c.Locals("id").(string)
	res, err := Controller.service.GetRandomUser(id)
	if err != nil {
		return c.Status(utils.GetStatusCode(err)).JSON(err)
	}
	return c.Status(200).JSON(fiber.Map{
		"code":    200,
		"message": "success get data",
		"result":  res,
	})
}

func (Controller *Controller) SwipeUser(c *fiber.Ctx) error {
	id := c.Locals("id").(string)
	var input userBusiness.SwipeUser
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"code":    400,
			"message": err.Error(),
		})
	}
	err := Controller.service.SwipeUser(id, input)
	if err != nil {
		return c.Status(utils.GetStatusCode(err)).JSON(err)
	}
	return c.Status(200).JSON(fiber.Map{
		"code":    200,
		"message": "success swipe",
	})
}

func (Controller *Controller) GetListPackage(c *fiber.Ctx) error {
	res, err := Controller.service.GetListPackage()
	if err != nil {
		return c.Status(utils.GetStatusCode(err)).JSON(err)
	}
	return c.Status(200).JSON(fiber.Map{
		"code":    200,
		"message": "success get data",
		"result":  res,
	})
}

func (Controller *Controller) PurchasePackage(c *fiber.Ctx) error {
	id := c.Locals("id").(string)
	var input userBusiness.Purchase
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"code":    400,
			"message": err.Error(),
		})
	}
	err := Controller.service.PurchasePackage(id, input.ID)
	if err != nil {
		return c.Status(utils.GetStatusCode(err)).JSON(err)
	}
	return c.Status(200).JSON(fiber.Map{
		"code":    200,
		"message": "success purchase",
	})
}

func (Controller *Controller) GetPackageByID(c *fiber.Ctx) error {
	id := c.Query("id")
	res, err := Controller.service.GetPackageByID(id)
	if err != nil {
		return c.Status(utils.GetStatusCode(err)).JSON(err)
	}
	return c.Status(200).JSON(fiber.Map{
		"code":    200,
		"message": "success get data",
		"result":  res,
	})
}
