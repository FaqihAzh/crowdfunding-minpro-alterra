package payment

import (
	"crowdfunding-minpro-alterra/modules/user"
	"os"
	"strconv"

	midtrans "github.com/veritrans/go-midtrans"
)

type service struct {
}

type Service interface {
	GetPaymentURL(donation Donation, user user.User) (string, error)
}

func NewService() *service {
	return &service{}
}

func (s *service) GetPaymentURL(donation Donation, user user.User) (string, error) {
	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	clientKey := os.Getenv("MIDTRANS_CLIENT_KEY")

	midclient := midtrans.NewClient()
  midclient.ServerKey = serverKey
  midclient.ClientKey = clientKey
  midclient.APIEnvType = midtrans.Sandbox

  snapGateway := midtrans.SnapGateway{
    Client: midclient,
  }

	snapReq := &midtrans.SnapReq{
		CustomerDetail: &midtrans.CustDetail{
			Email: user.Email,
			FName: user.Name,
		},
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(donation.ID),
			GrossAmt: int64(donation.Amount),
		},
	}

	snapTokenResp, err := snapGateway.GetToken(snapReq)

	if err != nil {
		return "", err
	}

	return snapTokenResp.RedirectURL, nil
}