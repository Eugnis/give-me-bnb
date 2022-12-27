package captcha

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/aidenesco/anticaptcha"
	"golang.org/x/sync/errgroup"
)

func New(ctx context.Context, anticaptchaKey string) (string, error) {
	g, ctx := errgroup.WithContext(ctx)

	var ret string

	g.Go(func() error {
		var err error
		ret, err = execute(ctx, anticaptchaKey)
		return err
	})

	return ret, g.Wait()
}

func execute(ctx context.Context, anticaptchaKey string) (string, error) {
	client := anticaptcha.NewClient(anticaptchaKey)
	balance, _ := client.GetBalance(context.Background())

	fmt.Println("Anticaptcha balance left:", balance)
	res, err := client.HCaptchaProxyless(ctx, "https://testnet.bnbchain.org/faucet-smart", "d9a9ee67-74da-4601-9f31-efe6a297a5cc")

	if err != nil {
		return "", err
	}
	return res.GRecaptchaResponse, nil
}
