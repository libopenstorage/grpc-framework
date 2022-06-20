// Example code
package main

import (
	"context"
	"crypto/x509"
	"flag"
	"fmt"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/libopenstorage/grpc-framework/test/app/api"
)

const (
	Bytes = uint64(1)
	KB    = Bytes * uint64(1024)
	MB    = KB * uint64(1024)
	GB    = MB * uint64(1024)
)

var (
	useTls  = flag.Bool("usetls", false, "Connect to server using TLS. Loads CA from the system")
	token   = flag.String("token", "", "Authorization token if any")
	address = flag.String("address", "127.0.0.1:9009", "Address to server as <address>:<port>")
)

type OpenStorageSdkToken struct{}

func (t OpenStorageSdkToken) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": "bearer " + *token,
	}, nil
}

func (t OpenStorageSdkToken) RequireTransportSecurity() bool {
	return *useTls
}

func main() {
	flag.Parse()

	// There are two ways to setup a token:
	//   - One is to setup a client interceptor which adds the token
	//     to every call automatically using grpc.WithPerRPCCredentials().
	//   - Second way is just to add it to the context directly as follows:
	//   import "google.golang.org/grpc/metadata"
	//   md := metadata.New(map[string]string{
	//		"authorization": "bearer" + token,
	//	 })
	//   ctx := metadata.NewOutgoingContext(context.Background(), md)
	//
	//   We will be using the more complicated first model to show how it can be done
	//
	//   To accomplish this, we first need to create an object that satisfies the
	//   interface needed by grpc.WithPerRPCCredentials(..)
	contextToken := OpenStorageSdkToken{}

	dialOptions := []grpc.DialOption{grpc.WithInsecure()}
	if *useTls {
		// Setup a connection
		capool, err := x509.SystemCertPool()
		if err != nil {
			fmt.Printf("Failed to load system certs: %v\n", err)
			os.Exit(1)
		}
		dialOptions = []grpc.DialOption{grpc.WithTransportCredentials(
			credentials.NewClientTLSFromCert(capool, ""),
		)}
	}

	if len(*token) != 0 {
		// Add token interceptor
		dialOptions = append(dialOptions, grpc.WithPerRPCCredentials(contextToken))
	}

	conn, err := grpc.Dial(*address, dialOptions...)
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

	greeter := api.NewHelloGreeterClient(conn)
	resp, err := greeter.SayHello(context.Background(), &api.HelloGreeterSayHelloRequest{
		Name: "theo",
	})
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
	fmt.Printf("message received: %s\n", resp.GetMessage())
}
