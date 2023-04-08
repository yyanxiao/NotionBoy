package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"

	"notionboy/api/pb"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/logger"
	"notionboy/internal/server/handler"
	"notionboy/internal/server/middleware"
	"notionboy/webui"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
)

func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	})
	return h2c.NewHandler(handler, &http2.Server{})
}

func grpcSvc() *grpc.Server {
	svc := grpc.NewServer(
		// midddlewares
		grpc.StreamInterceptor(grpc_auth.StreamServerInterceptor(middleware.NewAuthFunc())),
		grpc.UnaryInterceptor(grpc_auth.UnaryServerInterceptor(middleware.NewAuthFunc())),
	)
	pb.RegisterServiceServer(svc, handler.NewServer())
	reflection.Register(svc)
	return svc
}

func registerHttpHandlers(ctx context.Context, mux *http.ServeMux) {
	initNotion(mux)
	initWx(mux)
	mux.HandleFunc("/files/tg/", proxyTelegramFile)
	mux.HandleFunc("/files/ipfs/", proxyIpfs)
	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		renderHtml(w, "ok", http.StatusOK)
	})
	// proxy for openai chat API
	mux.HandleFunc("/v1/chat/completions", corsMiddleware(completions))
	mux.HandleFunc("/v1/models", corsMiddleware(proxyOpenAI))

	// wechat pay callback
	mux.HandleFunc("/v1/wechatpay/callback", wechatPayCallback)

	webui.RegisterHandlers(mux)
}

func Serve() {
	cfg := config.GetConfig().Service
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := http.NewServeMux()
	// register http handlers
	registerHttpHandlers(ctx, mux)

	gwmux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(CustomHeaderMatcher),
		runtime.WithMetadata(MetadataInjector),
	)
	// register grpc-gateway
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := pb.RegisterServiceHandlerFromEndpoint(ctx, gwmux, addr, opts); err != nil {
		logger.SugaredLogger.Fatalw("failed to register gateway", "error", err)
	}
	mux.Handle("/", gwmux)

	// wxgzh need use / as path, so we need to handle it manually
	_ = gwmux.HandlePath("POST", "/", wxProcessMsg)
	_ = gwmux.HandlePath("GET", "/", wxProcessMsg)

	grpcServer := grpcSvc()
	srv := &http.Server{
		Addr:    addr,
		Handler: grpcHandlerFunc(grpcServer, mux),
	}
	conn, err := net.Listen("tcp", addr)
	if err != nil {
		logger.SugaredLogger.Fatalw("failed to listen", "error", err)
	}
	logger.SugaredLogger.Infow("Server started", "addr", addr)
	err = srv.Serve(conn)
	if err != nil {
		logger.SugaredLogger.Fatalw("failed to serve", "error", err)
	}
}

var headersMap = map[string]struct{}{
	config.AUTH_HEADER_X_API_KEY: {},
	config.AUTH_HEADER_COOKIE:    {},
}

// CustomHeaderMatcher is a custom header matcher for gRPC-Gateway.
// https://grpc-ecosystem.github.io/grpc-gateway/docs/mapping/customizing_your_gateway/#mapping-from-http-request-headers-to-grpc-client-metadata
func CustomHeaderMatcher(key string) (string, bool) {
	key = strings.ToLower(key)
	if _, ok := headersMap[key]; ok {
		return key, true
	}
	return runtime.DefaultHeaderMatcher(key)
}

func MetadataInjector(ctx context.Context, req *http.Request) metadata.MD {
	md := make(map[string]string)
	// if method, ok := runtime.RPCMethod(ctx); ok {
	// 	md["method"] = method
	// }
	if pattern, ok := runtime.HTTPPathPattern(ctx); ok {
		md[config.AUTH_HEADER_PATH] = pattern
	}
	return metadata.New(md)
}
