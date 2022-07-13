package page

import (
	"github.com/canack/telebots/services/page_reader/page/target"
	"github.com/canack/telebots/services/page_reader/types"
)

func configureReddit(entryUrl string) *target.Reddit {
	return target.NewReddit(
		entryUrl,
		"//div[@data-testid=\"post-container\"]",
		"//h1[@class=\"_eYtD2XCVieq6emjKBH3m\"]",

		"_1YCqQVO-9r-Up6QPB9H6_4 _1YCqQVO-9r-Up6QPB9H6_4",
		"//div[@class=\"_3sf33-9rVAO_v4y0pIW_CH\"]",

		"p",
		"", // not used
		0,  // not used
		64, // image height equal or larger than 64
		10, // limit for images and entries
	)
}

func CrawlReddit(entryUrl string) (*[]types.ImageAndText, error) {
	a := configureReddit(entryUrl)
	content, err := a.Crawl()

	if err != nil {
		return nil, err
	}

	return content, nil

}
