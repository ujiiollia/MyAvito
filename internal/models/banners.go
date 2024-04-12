package models

import "context"

type BannerService interface {
	Ping(context.Context) error
}

type BannerRepository interface {
	Ping(context.Context) error
}
