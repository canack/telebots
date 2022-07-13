package target

import (
	"bytes"
	"fmt"
	"github.com/canack/telebots/services/page_reader/types"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/devices"
	"image"
)

type Reddit struct {
	types.Webpage
	EntryUrl          string
	MainEntryPath     string
	MainEntryTextPath string

	SubEntryDivPath  string
	SubEntryPath     string
	SubEntryTextPath string

	ScreenShotPath  string
	ScreenShotSizeX int
	ScreenShotSizeY int
	EntryLimit      int
}

func NewReddit(entryUrl string, mainEntryPath string, mainEntryTextPath string,
	subEntryDivPath string, subEntryPath string, subEntryTextPath string,
	screenShotPath string, screenShotSizeX int, screenShotSizeY int, entryLimit int) *Reddit {
	return &Reddit{
		EntryUrl:          entryUrl,
		MainEntryPath:     mainEntryPath,
		MainEntryTextPath: mainEntryTextPath,
		SubEntryDivPath:   subEntryDivPath,
		SubEntryPath:      subEntryPath,
		SubEntryTextPath:  subEntryTextPath,
		ScreenShotPath:    screenShotPath,
		ScreenShotSizeX:   screenShotSizeX,
		ScreenShotSizeY:   screenShotSizeY,
		EntryLimit:        entryLimit,
	}
}

func (r *Reddit) Crawl() (*[]types.ImageAndText, error) {

	var result []types.ImageAndText

	// To load entry page
	page := rod.New().MustConnect().MustPage(r.EntryUrl)
	page.MustEmulate(devices.LaptopWithTouch)
	page.MustReload().MustWaitLoad()
	page.MustReload().MustWaitLoad()

	mainResult, err := r.mainEntry(page)
	if err != nil {
		return &[]types.ImageAndText{}, err
	}

	result = append(result, *mainResult)

	subResult, err := r.subEntry(page)
	if err != nil {
		return &[]types.ImageAndText{}, err
	}
	result = append(result, *subResult...)

	return &result, nil
}

func (r *Reddit) mainEntry(page *rod.Page) (*types.ImageAndText, error) {
	//get main entry
	var result types.ImageAndText

	m := page.MustElementX(r.MainEntryPath)
	t := page.MustElementX(r.MainEntryTextPath)

	imgBytes := m.MustScreenshot()
	bytesReader := bytes.NewReader(imgBytes)
	img, _, _ := image.Decode(bytesReader)
	result.Image = img
	result.Text = t.MustText()

	//

	return &result, nil
}

func (r *Reddit) subEntry(page *rod.Page) (*[]types.ImageAndText, error) {
	postMain := page.MustSearch(r.SubEntryDivPath)
	post := postMain.MustElementsX(r.SubEntryPath)

	var result []types.ImageAndText
	e := 0

	for _, p := range post {
		if e == r.EntryLimit {
			break
		}

		imgBytes := p.MustScreenshot()
		bytesReader := bytes.NewReader(imgBytes)
		img, _, _ := image.Decode(bytesReader)
		if img.Bounds().Dy() >= r.ScreenShotSizeY {
			//i to str
			//iStr := strconv.Itoa(e)
			p.MustScreenshot() // TODO: edit this line
			e++
			comment := p.MustElement(r.SubEntryTextPath).MustText()
			fmt.Println(comment)

			result = append(result, types.ImageAndText{Image: img, Text: comment})

		}
	}

	return &result, nil
}
