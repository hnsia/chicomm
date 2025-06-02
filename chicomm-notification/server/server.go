package server

import (
	"context"
	"sync"
	"time"

	"github.com/hnsia/chicomm/chicomm-grpc/pb"
	"golang.org/x/sync/semaphore"

	gomail "gopkg.in/mail.v2"
)

type Server struct {
	client pb.ChicommClient
}

func NewServer(client pb.ChicommClient) *Server {
	return &Server{client: client}
}

func (s *Server) Run(ctx context.Context) {
	// process notification event every 30 seconds
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		// process notification event

		select {
		case <-ticker.C:
			s.processNotificationEvents(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (s *Server) processNotificationEvents(ctx context.Context) error {
	res, err := s.client.ListNotificationEvents(ctx, &pb.ListNotificationEventsReq{})
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	sem := semaphore.NewWeighted(10)
	for _, ev := range res.Events {
		wg.Add(1)
		if err := sem.Acquire(ctx, 1); err != nil {
			return err
		}

		go func(ev *pb.NotificationEvent) {
			defer sem.Release(1)
			defer wg.Done()
			// send email notification
			// update the notification event/state accordingly
		}(ev)
	}

	go func() {
		wg.Wait()
	}()

	return nil
}

func (s *Server) sendNotification(ctx context.Context, ev *pb.NotificationEvent) error {
	m := gomail.NewMessage()
	return nil
}
