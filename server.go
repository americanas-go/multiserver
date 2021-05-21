package multiserver

import (
	"context"
	"sync"

	"github.com/americanas-go/errors"
)

type d struct {
	Ok  bool
	Srv Server
}

var s []d

type Server interface {
	Serve(context.Context)
	Shutdown(context.Context)
}

func Serve(ctx context.Context, srvs ...Server) {

	switch len(srvs) {
	case 0:
		panic("no servers configured")
	case 1:
		srvs[0].Serve(ctx)
	default:
		wg := new(sync.WaitGroup)
		wg.Add(len(srvs))

		for i, srv := range srvs {
			i := i
			srv := srv
			s = append(s, d{
				Ok:  true,
				Srv: srv,
			})
			go func() {
				srv.Serve(ctx)
				s[i].Ok = false
				wg.Done()
			}()
		}

		wg.Wait()
	}

}

func Check(ctx context.Context) error {
	if len(s) == 0 {
		panic("no servers configured")
	}

	for _, a := range s {
		if !a.Ok {
			return errors.ServiceUnavailablef("one of servers is down")
		}
	}

	return nil
}

func Shutdown(ctx context.Context) {
	if len(s) == 0 {
		panic("no servers configured")
	}

	for _, a := range s {
		a.Srv.Shutdown(ctx)
	}
}
