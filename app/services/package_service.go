package services

import (
	"github.com/bearts/go-fiber/app/dbFunctions"
	"github.com/bearts/go-fiber/app/models"
	"github.com/bearts/go-fiber/app/structs"
	"github.com/bearts/go-fiber/app/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func UserCreatePackage(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	_id := user.Claims.(jwt.MapClaims)["id"].(string)
	uid, _ := primitive.ObjectIDFromHex(_id)
	claims := user.Claims.(jwt.MapClaims)
	PhoneAvailable := claims["PhoneAvailable"].(bool)
	if !PhoneAvailable {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Phone number is required",
		})
	}
	var body structs.UserCreatePackage
	if err := c.BodyParser(&body); err != nil {
		return err
	}
	if err := utils.Validate.Struct(body); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	delivery, err := dbFunctions.GetPlaceByCode(body.DeliveryLocation)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	pickup, err := dbFunctions.GetPlaceByCode(body.Package.Location)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	price, err := dbFunctions.GetPriceFromToById(pickup.Id, delivery.Id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	var packageInstance models.Package
	packageInstance.User = uid
	packageInstance.Package.NameOfApp = body.NameOfApp
	packageInstance.Package.Location = pickup.Id
	packageInstance.Package.TrackingId = body.Package.TrackingId
	packageInstance.DeliveryLocation = body.DeliveryLocation
	if body.Package.OTP != nil {
		packageInstance.Package.Otp = body.Package.OTP
	}
	if body.Package.Eta != nil {
		packageInstance.Package.Eta = body.Package.Eta
	}
	packageInstance.Package.Status = body.Package.Status
	packageInstance.Price = price
	packageInstance.Status = "pending"

	_pacakge, err := dbFunctions.CreatePackage(packageInstance)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"package": _pacakge,
	})
}

func UserGetAllPackage(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	_id := user.Claims.(jwt.MapClaims)["id"].(string)
	uid, _ := primitive.ObjectIDFromHex(_id)
	status := c.Query("status")
	packages, err := dbFunctions.GetAllPackagesOfUser(uid, status)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success":  true,
		"packages": packages,
	})
}

func UserGetPackageById(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	_id := user.Claims.(jwt.MapClaims)["id"].(string)
	uid, _ := primitive.ObjectIDFromHex(_id)
	_pid := c.Params("id")
	pid, _ := primitive.ObjectIDFromHex(_pid)

	packageInstance, err := dbFunctions.GetPackageById(pid)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	if packageInstance.User != uid {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "You are not the owner of this package",
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"package": packageInstance,
	})
}

func UserUpdatePackageStatus(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	_id := user.Claims.(jwt.MapClaims)["id"].(string)
	uid, _ := primitive.ObjectIDFromHex(_id)
	_pid := c.Params("id")
	pid, _ := primitive.ObjectIDFromHex(_pid)
	var body structs.UserUpdatePackageStatus
	if err := c.BodyParser(&body); err != nil {
		return err
	}
	if err := utils.Validate.Struct(body); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	packageInstance, err := dbFunctions.UpdatePackageDeliveryStatus(pid, body.Status, &uid)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"package": packageInstance,
	})
}

// Runner

func RunnerGetAllUnAssignedPackage(c *fiber.Ctx) error {
	packages, err := dbFunctions.GetAllUnAssignedPackages()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success":  true,
		"packages": packages,
	})
}

func RunnerGetAllPreviousPackage(c *fiber.Ctx) error {
	status := c.Query("status")
	runner := c.Locals("runner").(*jwt.Token)
	_id := runner.Claims.(jwt.MapClaims)["id"].(string)
	rid, _ := primitive.ObjectIDFromHex(_id)
	if status == "" {
		packages, err := dbFunctions.GetAllPackageByRunner(rid)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		}
		return c.JSON(fiber.Map{
			"success":  true,
			"packages": packages,
		})
	} else {
		packages, err := dbFunctions.GetAllPackageByRunnerByStatus(rid, status)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		}
		return c.JSON(fiber.Map{
			"success":  true,
			"packages": packages,
		})
	}
}

func RunnerGetPackageById(c *fiber.Ctx) error {
	runner := c.Locals("runner").(*jwt.Token)
	_id := runner.Claims.(jwt.MapClaims)["id"].(string)
	rid, _ := primitive.ObjectIDFromHex(_id)
	_pid := c.Params("id")
	pid, _ := primitive.ObjectIDFromHex(_pid)
	packageInstance, err := dbFunctions.GetPackageById(pid)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	if packageInstance.Runner != &rid {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "You are not assigned to this package",
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"package": packageInstance,
	})
}

func RunnerAssignPackage(c *fiber.Ctx) error {
	runner := c.Locals("runner").(*jwt.Token)
	_id := runner.Claims.(jwt.MapClaims)["id"].(string)
	rid, _ := primitive.ObjectIDFromHex(_id)
	_pid := c.Params("id")
	pid, _ := primitive.ObjectIDFromHex(_pid)
	_package, err := dbFunctions.GetPackageById(pid)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	if _package.Status != "pending" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Package is already assigned",
		})
	}
	packageInstance, err := dbFunctions.AssignPackage(pid, rid)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"package": packageInstance,
	})
}

func RunnerUpdatePackage(c *fiber.Ctx) error {
	runner := c.Locals("runner").(*jwt.Token)
	_id := runner.Claims.(jwt.MapClaims)["id"].(string)
	rid, _ := primitive.ObjectIDFromHex(_id)
	_pid := c.Params("id")
	pid, _ := primitive.ObjectIDFromHex(_pid)
	var body structs.RunnerUpdatePackage
	if err := c.BodyParser(&body); err != nil {
		return err
	}
	if err := utils.Validate.Struct(body); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	_package, err := dbFunctions.GetPackageById(pid)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	if _package.Runner != &rid {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "You are not assigned to this package",
		})
	}
	packageInstance, err := dbFunctions.UpdatePackageStatus(pid, body.Status)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"package": packageInstance,
	})
}

func RunnerDeliverPackage(c *fiber.Ctx) error {
	var body structs.RunnerDeliverPackage
	if err := c.BodyParser(&body); err != nil {
		return err
	}
	runner := c.Locals("runner").(*jwt.Token)
	_id := runner.Claims.(jwt.MapClaims)["id"].(string)
	rid, _ := primitive.ObjectIDFromHex(_id)
	_pid := c.Params("id")
	pid, _ := primitive.ObjectIDFromHex(_pid)
	packageInstance, err := dbFunctions.GetPackageById(pid)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	if packageInstance.Runner != &rid {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "You are not assigned to this package",
		})
	}
	if packageInstance.RunnerOtp != body.Otp {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "OTP is incorrect",
		})
	}
	packageUpdate, err := dbFunctions.UpdatePackageStatus(pid, "delivered")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"package": packageUpdate,
	})
}
