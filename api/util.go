package api

import "github.com/rs/xid"

func generateID() string {
	return xid.New().String()
}
