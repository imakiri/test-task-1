package internal

import (
	"context"
	"errors"
	"fmt"
	"github.com/gocolly/colly"
	proto "github.com/imakiri/test-task-1/internal/proto/v1"
	"net/http"
	"strconv"
)

type Service struct {
	proto.UnimplementedServiceServer
	client    *http.Client
	collector *colly.Collector
}

func a(n uint64, a [9]uint8) (r uint32) {
	for i := range a {
		r += uint32(uint64(a[i]) * (n % 10))
		n /= 10
	}
	return
}

func (Service) checkINN(inn uint64) (uint64, bool) {
	if 100000000 > inn || inn >= 10000000000 {
		return 0, false
	}

	var innValue uint64
	var innCheck uint8
	var checkArray = [9]uint8{8, 6, 4, 9, 5, 3, 10, 4, 2}

	if inn < 1000000000 {
		innValue = inn
		innCheck = uint8((a(innValue, checkArray) % 11) % 10)
		return innValue*10 + uint64(innCheck), true
	} else {
		innValue = inn / 10
		innCheck = uint8(inn % 10)
		if innCheck != uint8((a(innValue, checkArray)%11)%10) {
			return 0, false
		} else {
			return inn, true
		}
	}
}

func (s *Service) Get(c context.Context, request *proto.Request) (*proto.Response, error) {
	var inn, ok = s.checkINN(request.GetInn())
	if !ok {
		return nil, errors.New("invalid inn")
	}

	var collector = s.collector.Clone()
	var response = new(proto.Response)

	collector.OnHTML("div[id=main]", func(element *colly.HTMLElement) {
		response.Name = element.ChildText("h1[itemprop=name]")
		response.Director = element.ChildText(".company-info__text > a > span")

		var inn = element.ChildText("span[id=clip_inn]")
		var kpp = element.ChildText("span[id=clip_kpp]")

		var err error
		response.INN, err = strconv.ParseUint(inn, 10, 64)
		if err != nil {
			response.INN = 0
		}

		response.KPP, err = strconv.ParseUint(kpp, 10, 64)
		if err != nil {
			response.KPP = 0
		}
	})

	var err = collector.Visit(fmt.Sprintf("https://www.rusprofile.ru/search?query=%d&type=ul&search_inactive=2", inn))
	if err != nil {
		return nil, err
	}

	return response, nil
}

func NewService() (*Service, error) {
	var service = new(Service)
	service.client = new(http.Client)
	service.collector = colly.NewCollector(colly.AllowURLRevisit())

	return service, nil
}
