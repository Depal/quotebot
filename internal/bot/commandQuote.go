package bot

import (
	"github.com/Depal/quotebot/internal/bot/static"
	"github.com/tucnak/telebot"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"net/http"
	"os"
)

func (s *Service) handleQuote(message *telebot.Message) {
	s.announceCommand(static.CommandQuote, message)

	/*

		response := fmt.Sprintf("Will create a sticker of message \"%s\" from %s", quoted.Text, quoted.Sender.FirstName)

		_, err := s.bot.Reply(message, response)
		if err != nil {
			s.log.Error(err)
			return
		}
	*/

	quoted := message.ReplyTo
	if quoted == nil {
		_, err := s.bot.Reply(message, "Please quote a message first")
		if err != nil {
			s.log.Error(err)
			return
		}
		return
	}

	s.log.Debug("getting photos")
	userpics, err := s.bot.ProfilePhotosOf(quoted.Sender)
	if err != nil {
		s.log.Error(err)
	}

	s.log.Debug(len(userpics))
	s.log.Debug(userpics[0].FileID)
	//userpics[0].Send(s.bot, message.Sender, nil)

	userpicUrl, err := s.bot.FileURLByID(userpics[0].FileID)
	if err != nil {
		s.log.Fatal(err)
		return
	}

	s.log.Debug(userpicUrl)

	res, err := http.Get("http://i.imgur.com/m1UIjW1.jpg")
	if err != nil {
		s.log.Fatal(err)
		return
	}

	s.log.Debug(res.Body)

	/*
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			s.log.Fatal(err)
			return
		}
		//fmt.Println(data)

		defer func() {
			err = res.Body.Close()
			if err != nil {
				s.log.Fatal(err)
				return
			}
		}()

		/*
		img, _, err := image.Decode(bytes.NewReader(data))
		if err != nil || img == nil {
			s.log.Fatal(err)
			return
		}
	*/

	background := image.NewRGBA(image.Rect(0, 0, 640, 480))
	blue := color.RGBA{255, 122, 255, 255}
	draw.Draw(background, background.Bounds(), &image.Uniform{blue}, image.Point{}, draw.Src)

	f, err := os.Create("test.png")
	if err != nil {
		s.log.Fatal(err)
		return
	}
	defer f.Close()

	err = png.Encode(f, background)
	if err != nil {
		s.log.Fatal(err)
		return
	}

	file := &telebot.Photo{File: telebot.FromDisk("test.png")}

	s.log.Debug("replying")

	_, err = s.bot.Send(message.Sender, file)
	if err != nil {
		s.log.Error(err)
	}

	s.finishCommand(static.CommandQuote)
}
