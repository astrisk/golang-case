package main

/**
1. errorgroup 不在go的基础类库里，需要用第三饭库方式引入
2. 如果go get golang.org/x/sync/errgroup直接下载不了，可以采用github clone，把代码放在go path对应路径下
3. go get install
*/

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type webHandler struct {
	ctx  context.Context
	name string
}

// resposnse webserver name
func (wh *webHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, wh.name)
}

// new http server
func newHTTPServer(c context.Context, name string, port int) *http.Server {
	mux := http.NewServeMux()
	handler := &webHandler{ctx: c, name: name}
	mux.Handle("/", handler)
	httpServer := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: mux}

	return httpServer
}

func main() {

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(quit)

	group, ctxError := errgroup.WithContext(ctx)

	// start web server
	webServer := newHTTPServer(ctx, "WebServer", 8080)
	group.Go(func() error {
		log.Println("Web Server listening on port 8080")
		if err := webServer.ListenAndServe(); err != http.ErrServerClosed {
			return err
		}
		return nil
	})

	// start monitor server
	monServer := newHTTPServer(ctx, "MonServer", 8081)
	group.Go(func() error {
		log.Println("Monitor Server listening on port 8081")
		if err := monServer.ListenAndServe(); err != http.ErrServerClosed {
			return err
		}
		return nil
	})

	// termination
	select {
	case <-quit:
		break
	case <-ctxError.Done():
		break
	}

	// gracefully shutdown
	cancel()
	timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer timeoutCancel()

	webServer.Shutdown(timeoutCtx)
	monServer.Shutdown(timeoutCtx)

	if err := group.Wait(); err != nil {
		log.Fatal(err)
	}

	log.Println("Server bye")

}
