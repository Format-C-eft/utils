package redis

import (
	"github.com/Format-C-eft/utils/cache/codec"
)

type cacheOptions struct {
	Codec codec.Codec
}

type Options func(option *cacheOptions)

func WithCodec(codec codec.Codec) Options {
	return func(option *cacheOptions) {
		option.Codec = codec
	}
}
