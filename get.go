package pocket

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// GetOptions specifies constraints for returning items.
type GetOptions struct {
	consumerKey string
	accessToken string
	ContentType GetContentType
	Count       int64
	DetailType  GetDetailType
	Domain      string
	Favorite    bool
	Offset      int64
	Search      string
	Since       time.Time
	Sort        GetSort
	State       GetState
	Tag         string
}
type shadowGetOptions struct {
	ConsumerKey string         `json:"consumer_key"`
	AccessToken string         `json:"access_token"`
	ContentType GetContentType `json:"contentType,omitempty"`
	Count       int64          `json:"count,omitempty"`
	DetailType  GetDetailType  `json:"detailType,omitempty"`
	Domain      string         `json:"domain,omitempty"`
	Favorite    int            `json:"favorite,omitempty"`
	Offset      int64          `json:"offset,omitempty"`
	Search      string         `json:"search,omitempty"`
	Since       int64          `json:"since,omitempty"`
	Sort        GetSort        `json:"sort,omitempty"`
	State       GetState       `json:"state,omitempty"`
	Tag         string         `json:"tag,omitempty"`
}

// GetContentType is the type for GetOptions.ContentType
type GetContentType string

// GetDetailType is the type for GetOptions.DetailType
type GetDetailType string

// GetSort is the type for GetOptions.Sort
type GetSort string

// GetState is the type for GetOptions.state
type GetState string

const (
	// GetContentTypeArticle only return articles
	GetContentTypeArticle GetContentType = "article"
	// GetContentTypeImage only return images
	GetContentTypeImage GetContentType = "image"
	// GetContentTypeVideo only return videos or articles with embedded videos
	GetContentTypeVideo GetContentType = "video"

	// GetDetailTypeComplete return all data about each item, including tags,
	// images, authors, videos, and more
	GetDetailTypeComplete GetDetailType = "complete"
	// GetDetailTypeSimple return basic information about each item, including
	// title, url, status, and more
	GetDetailTypeSimple GetDetailType = "simple"

	// GetSortNewest return items in order of newest to oldest
	GetSortNewest GetSort = "newest"
	// GetSortOldest return items in order of oldest to newest
	GetSortOldest GetSort = "oldest"
	// GetSortSite return items in order of url alphabetically
	GetSortSite GetSort = "site"
	// GetSortTitle return items in order of title alphabetically
	GetSortTitle GetSort = "title"

	// GetStateUnread only return unread items (default)
	GetStateUnread GetState = "unread"
	// GetStateArchive only return archived items
	GetStateArchive GetState = "archive"
	// GetStateAll return both unread and archived items
	GetStateAll GetState = "all"
)

// MarshalJSON is a custom marshaller
func (o GetOptions) MarshalJSON() ([]byte, error) {
	s := &shadowGetOptions{
		AccessToken: o.accessToken,
		ConsumerKey: o.consumerKey,
		Count:       o.Count,
		DetailType:  o.DetailType,
		Domain:      o.Domain,
		Favorite:    0,
		Offset:      o.Offset,
		Search:      o.Search,
		Since:       0,
		State:       o.State,
		Tag:         o.Tag,
	}

	if o.Favorite {
		s.Favorite = 1
	}
	var yearOne time.Time
	if o.Since != yearOne {
		s.Since = o.Since.Unix()
	}

	return json.Marshal(s)
}

// GetResponse is the type returned by a valid call to Get()
type GetResponse struct {
	Status int `json:"status"`
	List   map[string]struct {
		ItemID        string `json:"item_id"`
		ResolvedID    string `json:"resolved_id"`
		GivenURL      string `json:"given_url"`
		GivenTitle    string `json:"given_title"`
		Favorite      string `json:"favorite"`
		Status        string `json:"status"`
		ResolvedTitle string `json:"resolved_title"`
		ResolvedURL   string `json:"resolved_url"`
		Excerpt       string `json:"excerpt"`
		IsArticle     string `json:"is_article"`
		HasVideo      string `json:"has_video"`
		HasImage      string `json:"has_image"`
		WordCount     string `json:"word_count"`
		Images        struct {
			Num1 struct {
				ItemID  string `json:"item_id"`
				ImageID string `json:"image_id"`
				Src     string `json:"src"`
				Width   string `json:"width"`
				Height  string `json:"height"`
				Credit  string `json:"credit"`
				Caption string `json:"caption"`
			} `json:"1"`
		} `json:"images"`
		Tags map[string]struct {
			ItemID string `json:"item_id"`
			Tag    string `json:"tag"`
		} `json:"tags"`
		Videos struct {
			Num1 struct {
				ItemID  string `json:"item_id"`
				VideoID string `json:"video_id"`
				Src     string `json:"src"`
				Width   string `json:"width"`
				Height  string `json:"height"`
				Type    string `json:"type"`
				Vid     string `json:"vid"`
			} `json:"1"`
		} `json:"videos"`
	} `json:"list"`
}
type getResponseEmpty struct {
	Status     int           `json:"status"`
	Complete   int           `json:"complete"`
	List       []interface{} `json:"list"`
	Error      interface{}   `json:"error"`
	SearchMeta struct {
		SearchType string `json:"search_type"`
	} `json:"search_meta"`
	Since int `json:"since"`
}

// Get returns items from the user's pocket
func (p *PocketClient) Get(o *GetOptions) (*GetResponse, error) {
	o.accessToken = p.AccessToken
	o.consumerKey = p.ConsumerKey
	d, err := json.MarshalIndent(o, "", "  ")

	body := bytes.NewReader(d)
	url := fmt.Sprintf("%s/v3/get?consumer_key=%s&access_token=%s", p.URL, p.ConsumerKey, p.AccessToken)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, errors.New("Get: New Request: " + err.Error())
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("X-Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.New("Get: Client Do: " + err.Error())
	}

	if resp.StatusCode == 400 {
		err := resp.Header.Get("X-Error")
		return nil, errors.New("Get: Client Error: " + err)
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("Get: Client Code: " + resp.Status)
	}

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("Get: ReadAll: " + err.Error())
	}
	resp.Body.Close()

	var gr GetResponse
	err = json.Unmarshal(response, &gr)
	if err != nil {
		// If there was an error unmarshalling, maybe it's an empty result,
		// so try a different structure
		var er getResponseEmpty
		err2 := json.Unmarshal(response, &er)
		if err2 != nil {
			// Error a second time, return the first one as it's probably more useful
			return nil, errors.New("Get: Unmarshal: " + err.Error())
		}
		gr.Status = er.Status
		return &gr, nil
	}
	return &gr, nil
}
