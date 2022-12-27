package captcha

import (
	"context"
	_ "embed"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/aidenesco/anticaptcha"
	"golang.org/x/sync/errgroup"
)

var (
	//go:embed captcha.html
	captchaHtml []byte
)

func New(ctx context.Context, anticaptchaKey string) (string, error) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		return "", err
	}

	g, ctx := errgroup.WithContext(ctx)

	serverCtx, serverCancel := context.WithCancel(ctx)

	g.Go(func() error {
		return serve(serverCtx, l)
	})

	var ret string

	g.Go(func() error {
		defer serverCancel()
		var err error
		ret, err = execute(ctx, anticaptchaKey)
		return err
	})

	return ret, g.Wait()
}

func serve(ctx context.Context, l net.Listener) error {
	g, ctx := errgroup.WithContext(ctx)
	srv := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.Write(captchaHtml)
		}),
	}

	g.Go(func() error {
		err := srv.Serve(l)
		if err == http.ErrServerClosed {
			return nil
		}
		return err
	})

	g.Go(func() error {
		// Handle shutdown
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return srv.Shutdown(ctx)
	})

	return g.Wait()
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
