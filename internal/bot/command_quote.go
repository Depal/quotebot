package bot

import (
	"github.com/Depal/quotebot/internal/bot/static"
	"github.com/fogleman/gg"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/nfnt/resize"
	"github.com/tucnak/telebot"
	"image"
	_ "image/jpeg"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strings"
)

const MaxTextLengthSymbols = 400
const MaxWordBatches = 7

func (s *Service) handleQuote(message *telebot.Message) {
	s.announceCommand(static.CommandQuote, message)

	payload := message.Payload
	if payload != "" {
		_, err := s.bot.Reply(message, "Сейчас работает только /q без аргументов (@Depal, сделай)")
		if err != nil {
			s.log.Fatal(err)
			return
		}
		return
	}

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
	var userpic image.Image

	userpics, err := s.bot.ProfilePhotosOf(quoted.Sender)
	if err != nil || len(userpics) < 1 {
		userpic, err = s.loadEmptyUserpic()
	} else {
		userpic, err = s.getUserpic(userpics[0].FileID)
	}
	if err != nil {
		s.log.Error(err)
		return
	}

	resizedUserpic := resize.Resize(96, 96, userpic, resize.Lanczos3)

	width := 575.0
	height := 175.0

	dc := gg.NewContext(int(width), int(height))
	dc.SetRGBA(1, 1, 1, 0)
	dc.Clear()

	s.log.Debug("drawing avatar")
	dc.DrawCircle(53, 53, 48)
	dc.Clip()
	dc.DrawImage(resizedUserpic, 5, 5)
	dc.ResetClip()

	s.log.Debug("drawing messagebox")
	var messageboxWidth float64
	var smallBubble bool
	if len([]rune(quoted.Text)) < 5 {
		s.log.Debug("small text")
		messageboxWidth = 255
		smallBubble = true
	} else {
		s.log.Debug("big text")
		messageboxWidth = 455
	}

	dc.SetRGB(20, 20, 20)
	dc.DrawRoundedRectangle(110, 5, messageboxWidth, 165, 32)
	dc.Fill()

	s.log.Debug("drawing username")

	fontBytes, err := ioutil.ReadFile("fonts/lucidagrande.ttf")
	if err != nil {
		s.log.Error(err)
		return
	}

	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		s.log.Error(err)
		return
	}

	sizeUsername := s.determineUsernameFontSize(quoted.Sender.FirstName, smallBubble)

	face := truetype.NewFace(font, &truetype.Options{Size: sizeUsername})
	dc.SetFontFace(face)
	dc.SetRGB(0, 0, 125)
	dc.DrawStringWrapped(quoted.Sender.FirstName, 330, 26, 0.5, 0.5, 400, 1, gg.AlignLeft)
	//dc.DrawString(quoted.Sender.FirstName, 130, 42)

	s.log.Debug("drawing message")
	text := quoted.Text
	if len([]rune(text)) > MaxTextLengthSymbols {
		text = text[:MaxTextLengthSymbols] + "..."
	}

	size := s.determineMessageFontSize(text)

	face = truetype.NewFace(font, &truetype.Options{Size: size})
	dc.SetFontFace(face)
	dc.SetRGB(0, 0, 0)
	dc.DrawStringWrapped(text, 330, 105, 0.5, 0.5, 400, 1, gg.AlignLeft)

	dc.SavePNG("sticker.png")

	file := &telebot.Sticker{File: telebot.FromDisk("sticker.png")}

	s.log.Debug("replying")

	_, err = s.bot.Reply(message, file)
	if err != nil {
		s.log.Error(err)
	}

	s.finishCommand(static.CommandQuote)
}

func (s *Service) loadEmptyUserpic() (userpic image.Image, err error) {
	file, err := os.Open("images/default_avatar.png")
	if err != nil {
		s.log.Error(err)
		return userpic, err
	}

	userpic, _, err = image.Decode(file)
	if err != nil {
		s.log.Error(err)
		return userpic, err
	}

	return userpic, nil
}

func (s *Service) getUserpic(fileID string) (userpic image.Image, err error) {
	userpicUrl, err := s.bot.FileURLByID(fileID)
	if err != nil {
		s.log.Error(err)
		return userpic, err
	}

	res, err := http.Get(userpicUrl)
	if err != nil {
		s.log.Error(err)
		return userpic, err
	}
	defer res.Body.Close()

	userpic, _, err = image.Decode(res.Body)
	if err != nil {
		s.log.Error(err)
		return userpic, err
	}

	return userpic, err
}

func (s *Service) determineUsernameFontSize(username string, isSmallBubble bool) (fontSize float64) {
	if isSmallBubble {
		if len([]rune(username)) > 11 {
			return 12
		} else {
			return 36
		}
	} else {
		return 36
	}
}

func (s *Service) determineMessageFontSize(text string) (fontSize float64) {
	words := len(strings.Split(text, " "))

	baseSize := float64(48)

	wordBatches := words / 3
	s.log.Debug(wordBatches)
	if wordBatches > MaxWordBatches {
		wordBatches = MaxWordBatches
	}

	if wordBatches < 1 {
		baseSize += 10
	}

	coefficient := 1.0 - (float64(wordBatches) * 0.09)
	s.log.Debug(coefficient)

	return baseSize * math.Max(0.1, coefficient)
}
