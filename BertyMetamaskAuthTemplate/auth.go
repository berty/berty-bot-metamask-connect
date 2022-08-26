package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"berty.tech/berty/v2/go/pkg/bertybot"

	"github.com/Doozers/ETH-Signature/ethsign"
	"gorm.io/gorm"
)

func LinkMetamask(db *gorm.DB) func(ctx bertybot.Context) {
	return func(ctx bertybot.Context) {
		if ctx.IsReplay || !ctx.IsNew {
			return
		}

		nonce := rand.Uint32()

		_ = ctx.ReplyString("Please link your metamask account following this link:" + opts.authAddr + "sign?nonce=" + strconv.Itoa(int(nonce)))
		fmt.Println("Please link your metamask account following this link:" + opts.authAddr + "sign?nonce=" + strconv.Itoa(int(nonce)))
		
		// add to db user
		var user User

		res := db.First(&user, "berty_conv_pub_key = ?", ctx.ConversationPK)
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			user.BertyConvPubKey = ctx.ConversationPK
			user.Nonce = nonce
			db.Create(&user)
		} else {
			db.Model(&user).Update("nonce", nonce)
		}
	}
}

func retrieveInfos(code string) (string, string, error) {
	decode, err := base64.StdEncoding.DecodeString(code)
	if err != nil {
		fmt.Println("decode: ", err)
		return "", "", err
	}

	infos := strings.Split(string(decode), " ")

	if len(infos) != 2 {
		return "", "", errors.New("bad signature")
	}

	return infos[0], strings.Replace(infos[1], " ", "", -1), nil
}

func VerifyMetamask(db *gorm.DB) func(ctx bertybot.Context) {
	return func(ctx bertybot.Context) {
		if ctx.IsReplay || !ctx.IsNew {
			return
		}

		code, err := func(ctx bertybot.Context) (string, error) {
			if len(ctx.CommandArgs) == 2 {
				return ctx.CommandArgs[1], nil
			}
			return "", errors.New("bad arguments")
		}(ctx)
		if err != nil {
			_ = ctx.ReplyString(err.Error())
			return
		}

		// decode base64 encoded code to retrieve pubkey and Nonce
		pubkey, sig, err := retrieveInfos(code)
		if err != nil {
			_ = ctx.ReplyString("retrieve infos: " + err.Error())
			return
		}

		var user User
		db.Where(User{BertyConvPubKey: ctx.ConversationPK}).First(&user)

		nonce := user.Nonce
		nonceStr := strconv.Itoa(int(nonce))

		fmt.Println("nonce: |" + nonceStr + "|")
		fmt.Println("sig: |" + sig + "|")
		fmt.Println("pubkey: |" + pubkey + "|")

		verify, err := ethsign.Verify(nonceStr, sig, pubkey, ethsign.Metamask)
		fmt.Println(verify)
		if err != nil {
			_ = ctx.ReplyString(err.Error())
			return
		}

		if !verify {
			_ = ctx.ReplyString("bad signature")
			return
		}

		db.Where(User{
			BertyConvPubKey: ctx.ConversationPK,
			Nonce:           nonce,
		}).Updates(map[string]interface{}{
			"metamask_address": pubkey,
			"verified":         1,
			"nonce":            0,
		})

		_ = ctx.ReplyString("metamask linked")
	}
}
