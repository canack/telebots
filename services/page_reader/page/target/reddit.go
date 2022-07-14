package target

import (
	"bytes"
	"context"
	"github.com/canack/telebots/services/page_reader/types"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/devices"
	"image"
	"time"
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

func (r *Reddit) CrawlWithTimeout() (*[]types.ImageAndText, error) {
	// Eğer 90 saniye içerisinde bu işlemi yapamıyorsa problem var demektir.
	ctx, cancel := context.WithTimeout(context.TODO(), 90*time.Second)
	defer cancel()

	dataChan := make(chan *[]types.ImageAndText)
	errChan := make(chan error)

	go r.Crawl(ctx, dataChan, errChan)

	select {
	case <-ctx.Done():
		return &[]types.ImageAndText{}, ctx.Err()
	case err := <-errChan:
		if err != nil {
			return &[]types.ImageAndText{}, err
		}
		return <-dataChan, nil
	}
}

func (r *Reddit) Crawl(ctx context.Context, dataChan chan *[]types.ImageAndText, errChan chan error) {

	var result []types.ImageAndText

	p := rod.New().Context(ctx).MustConnect().MustPage()
	defer p.Browser().Close()

	rod.Try(func() {
		p.MustNavigate(r.EntryUrl)
		p.MustEmulate(devices.LaptopWithTouch)
		p.MustReload().MustWaitLoad()
		p.MustReload().MustWaitLoad()

		mainResult, err := r.mainEntry(p)
		if err != nil {
			errChan <- err
			dataChan <- &result
		}

		result = append(result, *mainResult)

		subResult, err := r.subEntry(p)
		if err != nil {
			errChan <- err
			dataChan <- &result
		}
		result = append(result, *subResult...)

		errChan <- nil
		dataChan <- &result
	})

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
			p.MustScreenshot()
			e++
			comment := p.MustElement(r.SubEntryTextPath).MustText()
			result = append(result, types.ImageAndText{Image: img, Text: comment})
		}
	}

	return &result, nil
}
