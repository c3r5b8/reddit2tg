package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/c3r5b8/r2g/sqlite"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
)

var client *tg.Client

func runBot(ctx context.Context) error {
	client = tg.New(os.Getenv("TELEGRAM_BOT_TOKEN"))
	channelId, err := strconv.Atoi(os.Getenv("CHANNEL_ID"))
	if err != nil {
		fmt.Println(err)
	}
	router := tgb.NewRouter().
		Message(func(ctx context.Context, msg *tgb.MessageUpdate) error {
			fmt.Println(msg.Message.From.ID)
			if msg.Message.Chat.ID == 5170084396 || msg.Message.Chat.ID == 848140054 {
				NextUrl = SubredditUrl
				scroblePage(NextUrl, int(msg.Message.Chat.ID))
				return nil
			}

			fmt.Println(msg.Message.Chat.ID)
			fmt.Println(msg.Message.From.FirstName)
			return nil

		}, tgb.Command("first")).
		Message(func(ctx context.Context, msg *tgb.MessageUpdate) error {
			if msg.Message.Chat.ID == 5170084396 || msg.Message.Chat.ID == 848140054 {
				scroblePage(NextUrl, int(msg.Message.Chat.ID))
				return nil
			}
			fmt.Println(msg.Message.Chat.ID)
			fmt.Println(msg.Message.From.FirstName)
			return nil
		}, tgb.Command("next")).
		// Message(func(ctx context.Context, msg *tgb.MessageUpdate) error {
		// 	fmt.Println(msg.Message.Chat.ID)
		// 	return nil
		// }).
		CallbackQuery(func(ctx context.Context, cbq *tgb.CallbackQueryUpdate) error {
			data := cbq.CallbackQuery.Data
			datas := strings.Split(data, "-")
			chatId, _ := strconv.Atoi(datas[2])
			messegeId, _ := strconv.Atoi(datas[3])
			switch datas[0] {
			case "post":
				_, err := client.CopyMessage(tg.ChatID(channelId), tg.ChatID(chatId), messegeId-1).Do(context.Background())
				if err != nil {
					fmt.Println(err)
				}
				queries.WritePost(context.Background(), sqlite.WritePostParams{ID: datas[1], Shown: true})
				client.DeleteMessage(tg.ChatID(chatId), messegeId-1).DoVoid(context.Background())
				client.DeleteMessage(tg.ChatID(chatId), messegeId).DoVoid(context.Background())
			case "skip":
				client.DeleteMessage(tg.ChatID(chatId), messegeId-1).DoVoid(context.Background())
				client.DeleteMessage(tg.ChatID(chatId), messegeId).DoVoid(context.Background())
			case "del":
				queries.WritePost(context.Background(), sqlite.WritePostParams{ID: datas[1], Shown: true})
				client.DeleteMessage(tg.ChatID(chatId), messegeId-1).DoVoid(context.Background())
				client.DeleteMessage(tg.ChatID(chatId), messegeId).DoVoid(context.Background())
			}
			return nil
		})
	return tgb.NewPoller(
		router,
		client,
	).Run(ctx)
}

func sendPhoto(url string, id string, chatId int) {
	msg, err := client.SendPhoto(
		tg.Chat{ID: tg.ChatID(chatId)},
		tg.NewFileArgURL(url)).
		Do(context.Background())
	if err != nil {
		fmt.Println(err)
		fmt.Println(url)
		return
	}

	sendControls(id, int(msg.Chat.ID), msg.ID+1)
}

func sendVideo(url string, id string, chatId int) {
	msg, err := client.SendVideo(
		tg.Chat{ID: tg.ChatID(chatId)},
		tg.NewFileArgURL(url)).
		Do(context.Background())
	if err != nil {
		fmt.Println(err)
		fmt.Println(url)
		return
	}

	sendControls(id, int(msg.Chat.ID), msg.ID+1)
}

func sendPhotoGallery(fileUrls []string, id string, chatId int) {
	for i := 0; i < len(fileUrls); i++ {
		msg, err := client.SendPhoto(
			tg.Chat{ID: tg.ChatID(chatId)},
			tg.NewFileArgURL(fileUrls[i])).
			Do(context.Background())
		if err != nil {
			fmt.Println(err)
			fmt.Println(fileUrls[i])
			return
		}
		sendControls(id, int(msg.Chat.ID), msg.ID+1)
	}

}

func sendControls(mediaId string, chatId int, messegeId int) {
	client.SendMessage(
		tg.Chat{ID: tg.ChatID(chatId)}, "control for previus port").
		ReplyMarkup(inlineKeyboard(mediaId, chatId, messegeId)).
		Do(context.Background())
}

func inlineKeyboard(mediaId string, chatId int, messegeId int) tg.InlineKeyboardMarkup {
	layout := tg.NewButtonLayout[tg.InlineKeyboardButton](3).Row(
		tg.NewInlineKeyboardButtonCallback("post", "post-"+mediaId+"-"+strconv.Itoa(chatId)+"-"+strconv.Itoa(messegeId)),
		tg.NewInlineKeyboardButtonCallback("skip", "skip-"+mediaId+"-"+strconv.Itoa(chatId)+"-"+strconv.Itoa(messegeId)),
		tg.NewInlineKeyboardButtonCallback("del", "del-"+mediaId+"-"+strconv.Itoa(chatId)+"-"+strconv.Itoa(messegeId)),
	)

	return tg.NewInlineKeyboardMarkup(layout.Keyboard()...)

}
