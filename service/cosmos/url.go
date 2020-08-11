package cosmos

import (
	"fmt"
	"regexp"

	"github.com/sorcererxw/jikeview-bot/util/log"
)

type Url struct {
	ID string
}

func ParseUrl(url string) *Url {
	re := regexp.MustCompile("^(https?://)?(www\\.)?xiaoyuzhoufm\\.com/episodes?/([0-9a-z]+)/?(\\?.*)?$")
	if re.MatchString(url) {
		return &Url{re.FindStringSubmatch(url)[3]}
	}
	log.Print("not cosmos")
	return nil
}

func (p Episode) generateUrl() string {
	return fmt.Sprintf("https://www.xiaoyuzhoufm.com/episode/%s", p.Eid)
}
