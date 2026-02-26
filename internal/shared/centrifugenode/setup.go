package centrifugenode

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	pkgjwt "starter-boilerplate/pkg/jwt"

	gocentrifuge "github.com/centrifugal/centrifuge"
)

type Init struct{}

// Setup wires JWT authentication, channel authorization, recovery, and WebSocket transport
// onto the centrifuge Node. Returns early if node is nil (standalone mode).
func Setup(node *gocentrifuge.Node, mux *http.ServeMux, jwtManager *pkgjwt.Manager) Init {
	if node == nil {
		return Init{}
	}

	node.OnConnecting(func(ctx context.Context, e gocentrifuge.ConnectEvent) (gocentrifuge.ConnectReply, error) {
		if e.Token == "" {
			return gocentrifuge.ConnectReply{}, errors.New("empty token")
		}

		claims, err := jwtManager.ValidateAccessToken(e.Token)
		if err != nil {
			return gocentrifuge.ConnectReply{}, err
		}

		return gocentrifuge.ConnectReply{
			Credentials: &gocentrifuge.Credentials{
				UserID: claims.UserID,
			},
		}, nil
	})

	node.OnConnect(func(client *gocentrifuge.Client) {
		client.OnSubscribe(func(e gocentrifuge.SubscribeEvent, cb gocentrifuge.SubscribeCallback) {
			if err := authorizeChannel(client.UserID(), e.Channel); err != nil {
				cb(gocentrifuge.SubscribeReply{}, err)
				return
			}
			cb(gocentrifuge.SubscribeReply{
				Options: gocentrifuge.SubscribeOptions{
					EnableRecovery: true,
					RecoveryMode:   gocentrifuge.RecoveryModeCache,
				},
			}, nil)
		})
	})

	wsHandler := gocentrifuge.NewWebsocketHandler(node, gocentrifuge.WebsocketConfig{
		CheckOrigin: func(r *http.Request) bool { return true },
	})
	mux.Handle("GET /connection/websocket", wsHandler)

	slog.Info("centrifuge websocket handler mounted", slog.String("path", "/connection/websocket"))
	return Init{}
}

// authorizeChannel checks that the user is allowed to subscribe to the given channel.
// For "personal:<userID>" channels, only the owner may subscribe.
func authorizeChannel(userID, channel string) error {
	ns, id, ok := strings.Cut(channel, ":")
	if !ok {
		return errors.New("invalid channel format")
	}

	switch ns {
	case "personal":
		if id != userID {
			return errors.New("permission denied")
		}
		return nil
	default:
		return errors.New("unknown channel namespace: " + ns)
	}
}
