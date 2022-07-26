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
	dc.DrawRoundedRectangle(110, 5, messageboxWidth, 175, 32)
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

	//font, err := truetype.Parse(goregular.TTF)
	//if err != nil {
	//	log.Fatal(err)
	//}
	sizeUsername := s.determineUsernameFontSize(quoted.Sender.FirstName, smallBubble)

	face := truetype.NewFace(font, &truetype.Options{Size: sizeUsername})
	dc.SetFontFace(face)
	dc.SetRGB(0, 0, 125)
	dc.DrawStringWrapped(quoted.Sender.FirstName, 330, 26, 0.5, 0.5, 400, 1, gg.AlignLeft)
	//dc.DrawString(quoted.Sender.FirstName, 130, 42)

	s.log.Debug("drawing message")
	text := quoted.Text
	//if len(text) > 80 {
	//	text = text[:80] + "..."
	//}

	size := s.determineMessageFontSize(text)

	face = truetype.NewFace(font, &truetype.Options{Size: size})
	dc.SetFontFace(face)
	dc.SetRGB(0, 0, 0)
	dc.DrawStringWrapped(text, 330, 110, 0.5, 0.5, 400, 1, gg.AlignLeft)

	dc.SavePNG("sticker.png")

	/*

		background := image.NewRGBA(image.Rect(0, 0, 320, 240))
		blue := color.RGBA{25, 25, 50, 255}
		draw.Draw(background, background.Bounds(), &image.Uniform{C: blue}, image.Point{}, draw.Src)
		draw.Draw(background, background.Bounds(), resizedUserpic, image.Point{}, draw.Src)

		addLabel(background, 50, 50, message.Sender.FirstName)

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
	*/

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

	if wordBatches < 1 {
		baseSize += 10
	}

	coefficient := 1.0 - (float64(wordBatches) * 0.1)
	s.log.Debug(coefficient)

	return baseSize * math.Max(0.1, coefficient)
}

/*
func addLabel(img *image.RGBA, x, y int, label string) {
	var (
		dpi      = flag.Float64("dpi", 72, "screen resolution in Dots Per Inch")
		fontfile = flag.String("fontfile", "../../testdata/luxisr.ttf", "filename of the ttf font")
		//hinting  = flag.String("hinting", "none", "none | full")
		size     = flag.Float64("size", 12, "font size in points")
		//spacing  = flag.Float64("spacing", 1.5, "line spacing (e.g. 2 means double spaced)")
		//wonb     = flag.Bool("whiteonblack", false, "white text on a black background")
	)

	// Read the font data.
	fontBytes, err := ioutil.ReadFile(*fontfile)
	if err != nil {
		log.Println(err)
		return
	}
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
		return
	}

	// Initialize the context.
	fg, bg := image.Black, image.White
	rgba := image.NewRGBA(image.Rect(0, 0, 640, 480))
	draw.Draw(rgba, rgba.Bounds(), bg, image.Point{}, draw.Src)

	c := freetype.NewContext()
	c.SetDPI(*dpi)
	c.SetFont(f)
	c.SetFontSize(*size)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)


	c.SetDst(img)
	sizeD := 12.0 // font size in pixels
	pt := freetype.Pt(x, y+int(c.PointToFixed(sizeD)>>6))

	if _, err := c.DrawString(label, pt); err != nil {
		// handle error
	}
}
*/
