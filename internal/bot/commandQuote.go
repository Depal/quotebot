package bot

import (
	"github.com/Depal/quotebot/internal/bot/static"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/nfnt/resize"
	"github.com/tucnak/telebot"
	"golang.org/x/image/font/gofont/goregular"
	"image"
	_ "image/jpeg"
	"log"
	"net/http"
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

	res, err := http.Get(userpicUrl)
	if err != nil {
		s.log.Fatal(err)
		return
	}
	defer res.Body.Close()

	s.log.Debug(res.Body)

	userpic, _, err := image.Decode(res.Body)
	if err != nil {
		s.log.Fatal(err)
		return
	}

	resizedUserpic := resize.Resize(96, 96, userpic, resize.Lanczos3)

	width := 500.0
	height := 200.0

	dc := gg.NewContext(int(width), int(height))
	dc.SetRGBA(1, 1, 1, 0)
	dc.Clear()
	dc.SetRGBA(0, 0, 0, 0)

	s.log.Debug("drawing avatar")
	dc.DrawCircle(53, 53, 48)
	dc.Clip()
	dc.DrawImage(resizedUserpic, 5, 5)
	dc.ResetClip()

	s.log.Debug("drawing messagebox")
	dc.SetRGB(50, 50, 75)
	dc.DrawRoundedRectangle(110, 5, 380, 190, 32)
	dc.Fill()

	s.log.Debug("drawing username")
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal(err)
	}
	face := truetype.NewFace(font, &truetype.Options{Size: 24})
	dc.SetFontFace(face)
	dc.SetRGB(0, 0, 0)
	dc.DrawStringAnchored(quoted.Sender.FirstName, 170, 30, 0.5, 0.5)

	s.log.Debug("drawing message")
	face = truetype.NewFace(font, &truetype.Options{Size: 20})
	dc.SetFontFace(face)
	dc.SetRGB(0, 0, 0)
	dc.DrawStringWrapped(quoted.Text, 300, 80, 0.5, 0.5, 300, 1, gg.AlignLeft)

	dc.SavePNG("test.png")

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

	file := &telebot.Photo{File: telebot.FromDisk("test.png")}

	s.log.Debug("replying")

	_, err = s.bot.Send(message.Sender, file)
	if err != nil {
		s.log.Error(err)
	}

	s.finishCommand(static.CommandQuote)
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
