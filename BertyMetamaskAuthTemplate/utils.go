package main

import "berty.tech/berty/v2/go/pkg/bertybot"

func Pingpong(ctx bertybot.Context) {
	if ctx.IsReplay || !ctx.IsNew {
		return
	}
	_ = ctx.ReplyString("pong")
}
