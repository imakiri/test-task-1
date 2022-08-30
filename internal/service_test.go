package internal

import (
	"context"
	"fmt"
	proto "github.com/imakiri/test-task-1/internal/proto/v1"
	"testing"
)

func TestGet(t *testing.T) {
	var service, _ = NewService()

	var resp, err = service.Get(context.Background(), &proto.Request{Inn: 8622000931})
	if err != nil {
		t.Error(err)
	}
	if resp.GetINN() != 8622000931 {
		t.Error("1")
	}

	fmt.Println(resp.GetName())
	fmt.Println(resp.GetINN())
	fmt.Println(resp.GetKPP())
	fmt.Println(resp.GetDirector())
	fmt.Println(err)

	resp, err = service.Get(context.Background(), &proto.Request{Inn: 862200093})
	if err != nil {
		t.Error(err)
	}
	if resp.GetINN() != 8622000931 {
		t.Error("2")
	}

	resp, err = service.Get(context.Background(), &proto.Request{Inn: 86220009})
	if err == nil {
		t.Error("4")
	}

	resp, err = service.Get(context.Background(), &proto.Request{Inn: 86220009311})
	if err == nil {
		t.Error("5")
	}
}
