package routes

import (
	"altt/internal/entities"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gofiber/fiber/v2"
)

// getNativeBalance gets the native balance of an address. returns 404 if the chain is not supported.
func (s *Server) getNativeBalance(ctx *fiber.Ctx) error {
	chain, err := entities.ChainFromString(ctx.Params("chain"))
	if err != nil {
		return ctx.Status(http.StatusNotFound).SendString(err.Error())
	}
	address := common.HexToAddress(ctx.Params("address"))
	if !checkAddressValid(address) {
		return ctx.Status(http.StatusBadRequest).SendString("invalid address")
	}
	balance, err := s.serviceBalancer.GetNativeBalance(ctx.UserContext(), chain, address)
	if err != nil {
		return err
	}
	return ctx.JSON(balance)
}

// getKnownTokenBalance gets the balance of a known token. returns 404 if the token is not known.
func (s *Server) getKnownTokenBalance(ctx *fiber.Ctx) error {
	chain, err := entities.ChainFromString(ctx.Params("chain"))
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return err
	}
	token, err := entities.TokenFromString(ctx.Params("token"))
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return err
	}
	if _, err = entities.GetTokenAddress(chain, token); err != nil {
		ctx.Status(http.StatusNotFound)
		return err
	}
	address := common.HexToAddress(ctx.Params("address"))
	if !checkAddressValid(address) {
		ctx.Status(http.StatusBadRequest)
		return err
	}
	balance, err := s.serviceBalancer.GetKnownTokenBalance(ctx.UserContext(), token, chain, address)
	if err != nil {
		return err
	}
	return ctx.JSON(balance)
}

func checkAddressValid(address common.Address) bool {
	return address.String() != "0x0000000000000000000000000000000000000000"
}
