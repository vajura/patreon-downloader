package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type PatreonData struct {
	Data     []PatreonDataData     `json:"data"`
	Included []PatreonDataIncluded `json:"included"`
}

type PatreonDataIncluded struct {
	Attributes struct {
		DownloadURL string `json:"download_url"`
		FileName    string `json:"file_name"`
		ImageUrls   struct {
			Default   string `json:"default"`
			Original  string `json:"original"`
			Thumbnail string `json:"thumbnail"`
		} `json:"image_urls"`
		Metadata struct {
			Dimensions struct {
				H int `json:"h"`
				W int `json:"w"`
			} `json:"dimensions"`
		} `json:"metadata"`
	} `json:"attributes,omitempty"`
	ID   string `json:"id"`
	Type string `json:"type"`
}

type PatreonDataData struct {
	Attributes    Attributes    `json:"attributes"`
	ID            string        `json:"id"`
	Relationships Relationships `json:"relationships"`
	Type          string        `json:"type"`
}
type Image struct {
	Height   int    `json:"height"`
	LargeURL string `json:"large_url"`
	ThumbURL string `json:"thumb_url"`
	URL      string `json:"url"`
	Width    int    `json:"width"`
}
type PostFile struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type PostMetadata struct {
	ImageOrder []string `json:"image_order"`
}
type Attributes struct {
	ChangeVisibilityAt       interface{}  `json:"change_visibility_at"`
	CommentCount             int          `json:"comment_count"`
	Content                  string       `json:"content"`
	CurrentUserCanComment    bool         `json:"current_user_can_comment"`
	CurrentUserCanDelete     bool         `json:"current_user_can_delete"`
	CurrentUserCanView       bool         `json:"current_user_can_view"`
	CurrentUserHasLiked      bool         `json:"current_user_has_liked"`
	Embed                    interface{}  `json:"embed"`
	HasTiViolation           bool         `json:"has_ti_violation"`
	Image                    Image        `json:"image"`
	IsPaid                   bool         `json:"is_paid"`
	LikeCount                int          `json:"like_count"`
	MetaImageURL             string       `json:"meta_image_url"`
	MinCentsPledgedToView    interface{}  `json:"min_cents_pledged_to_view"`
	PatreonURL               string       `json:"patreon_url"`
	PledgeURL                string       `json:"pledge_url"`
	PostFile                 PostFile     `json:"post_file"`
	PostMetadata             PostMetadata `json:"post_metadata"`
	PostType                 string       `json:"post_type"`
	PublishedAt              time.Time    `json:"published_at"`
	TeaserText               interface{}  `json:"teaser_text"`
	Title                    string       `json:"title"`
	UpgradeURL               string       `json:"upgrade_url"`
	URL                      string       `json:"url"`
	WasPostedByCampaignOwner bool         `json:"was_posted_by_campaign_owner"`
}
type Data struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}
type AccessRules struct {
	Data []Data `json:"data"`
}
type Attachments struct {
	Data []interface{} `json:"data"`
}
type Audio struct {
	Data interface{} `json:"data"`
}
type Links struct {
	Related string `json:"related"`
}
type Campaign struct {
	Data  Data  `json:"data"`
	Links Links `json:"links"`
}
type Images struct {
	Data []Data `json:"data"`
}
type Media struct {
	Data []Data `json:"data"`
}
type Poll struct {
	Data interface{} `json:"data"`
}
type TiChecks struct {
	Data []interface{} `json:"data"`
}
type User struct {
	Data  Data  `json:"data"`
	Links Links `json:"links"`
}
type UserDefinedTags struct {
	Data []interface{} `json:"data"`
}
type Relationships struct {
	AccessRules     AccessRules     `json:"access_rules"`
	Attachments     Attachments     `json:"attachments"`
	Audio           Audio           `json:"audio"`
	Campaign        Campaign        `json:"campaign"`
	Images          Images          `json:"images"`
	Media           Media           `json:"media"`
	Poll            Poll            `json:"poll"`
	TiChecks        TiChecks        `json:"ti_checks"`
	User            User            `json:"user"`
	UserDefinedTags UserDefinedTags `json:"user_defined_tags"`
}

func GetPatreonPosts(month string, config Config) (PatreonData, error) {
	data := PatreonData{}

	urlStr := fmt.Sprintf("https://www.patreon.com/api/posts?include=campaign%%2Caccess_rules%%2Cattachments%%2Caudio%%2Cimages%%2Cmedia%%2Cpoll.choices%%2Cpoll.current_user_responses.user%%2Cpoll.current_user_responses.choice%%2Cpoll.current_user_responses.poll%%2Cuser%%2Cuser_defined_tags%%2Cti_checks&fields[campaign]=currency%%2Cshow_audio_post_download_links%%2Cavatar_photo_url%%2Cearnings_visibility%%2Cis_nsfw%%2Cis_monthly%%2Cname%%2Curl&fields[post]=change_visibility_at%%2Ccomment_count%%2Ccontent%%2Ccurrent_user_can_comment%%2Ccurrent_user_can_delete%%2Ccurrent_user_can_view%%2Ccurrent_user_has_liked%%2Cembed%%2Cimage%%2Cis_paid%%2Clike_count%%2Cmeta_image_url%%2Cmin_cents_pledged_to_view%%2Cpost_file%%2Cpost_metadata%%2Cpublished_at%%2Cpatreon_url%%2Cpost_type%%2Cpledge_url%%2Cthumbnail_url%%2Cteaser_text%%2Ctitle%%2Cupgrade_url%%2Curl%%2Cwas_posted_by_campaign_owner%%2Chas_ti_violation&fields[post_tag]=tag_type%%2Cvalue&fields[user]=image_url%%2Cfull_name%%2Curl&fields[access_rule]=access_rule_type%%2Camount_cents&fields[media]=id%%2Cimage_urls%%2Cdownload_url%%2Cmetadata%%2Cfile_name&filter[campaign_id]=%s&filter[contains_exclusive_posts]=true&filter[is_draft]=false&filter[month]=%s&sort=-published_at&json-api-version=1.0", config.CampaignId, month)
	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return data, err
	}

	req.Header.Add("accept", config.Accept)
	req.Header.Add("content-type", config.ContentType)
	req.Header.Add("sec-ch-ua-mobile", config.SecChUaMobile)
	req.Header.Add("sec-fetch-dest", config.SecFetchDest)
	req.Header.Add("sec-fetch-mode", config.SecFetchMode)
	req.Header.Add("sec-fetch-site", config.SecFetchSite)
	req.Header.Add("user-agent", config.UserAgent)
	req.Header.Add("sec-ch-ua-platform", config.SecChUaPlatform)
	req.Header.Add("sec-ch-ua", config.SecChUa)
	req.Header.Add("cookie", config.Cookie)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return data, err
	}
	err = json.NewDecoder(resp.Body).Decode(&data)

	return data, err
}

func DownloadPatreonImage(urlStr string, fullFile string) error {
	out, err := os.Create(fullFile)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(urlStr)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	_, err = io.Copy(out, resp.Body)

	return err
}
