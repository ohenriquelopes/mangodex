package mangodex

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	CreateMangaPath              = MangaListPath
	ViewMangaPath                = "manga/%s"
	UpdateMangaPath              = ViewMangaPath
	DeleteMangaPath              = ViewMangaPath
	UnfollowMangaPath            = "manga/%s/follow"
	FollowMangaPath              = UnfollowMangaPath
	MangaFeedPath                = "manga/%s/feed"
	MangaReadMarkersPath         = "manga/%s/read"
	GetRandomMangaPath           = "manga/random"
	TagListPath                  = "manga/tag"
	UpdateMangaReadingStatusPath = "manga/%s/status"
)

type MangaList struct {
	Results []MangaResponse `json:"results"`
	Limit   int             `json:"limit"`
	Offset  int             `json:"offset"`
	Total   int             `json:"total"`
}

type MangaResponse struct {
	Result        string         `json:"result"`
	Data          Manga          `json:"data"`
	Relationships []Relationship `json:"relationships"`
}

func (mr *MangaResponse) GetResult() string {
	return mr.Result
}

type Manga struct {
	ID         string          `json:"id"`
	Type       string          `json:"type"`
	Attributes MangaAttributes `json:"attributes"`
}

type Relationship struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type MangaAttributes struct {
	Title                  LocalisedString    `json:"title"`
	AltTitles              []*LocalisedString `json:"altTitles"`
	Description            LocalisedString    `json:"description"`
	IsLocked               bool               `json:"isLocked"`
	Links                  []*string          `json:"links"`
	OriginalLanguage       string             `json:"originalLanguage"`
	LastVolume             *string            `json:"lastVolume"`
	LastChapter            *string            `json:"lastChapter"`
	PublicationDemographic *string            `json:"publicationDemographic"`
	Status                 *string            `json:"status"`
	Year                   int                `json:"year"`
	ContentRating          *string            `json:"contentRating"`
	Tags                   []*LocalisedString `json:"tags"`
	Version                int                `json:"version"`
	CreatedAt              string             `json:"createdAt"`
	UpdatedAt              string             `json:"updatedAt"`
}

type LocalisedString struct {
	Property1 string `json:"property1"`
	Property2 string `json:"property2"`
}

type ChapterReadMarkersResponse struct {
	Result string   `json:"result"`
	Data   []string `json:"data"`
}

func (rmr *ChapterReadMarkersResponse) GetResult() string {
	return rmr.Result
}

type TagResponse struct {
	Result        string         `json:"result"`
	Data          Tag            `json:"data"`
	Relationships []Relationship `json:"relationships"`
}

func (tg *TagResponse) GetResult() string {
	return tg.Result
}

type Tag struct {
	ID         string        `json:"id"`
	Type       string        `json:"type"`
	Attributes TagAttributes `json:"attributes"`
}

type TagAttributes struct {
	Name    LocalisedString `json:"name"`
	Version int             `json:"version"`
}

type MangaReadingStatusResponse struct {
	Result   string            `json:"result"`
	Statuses map[string]string `json:"statuses"`
}

func (s *MangaReadingStatusResponse) GetResult() string {
	return s.Result
}

// CreateManga : Create a new manga.
// https://api.mangadex.org/docs.html#operation/post-manga
func (dc *DexClient) CreateManga(newManga io.Reader) (*MangaResponse, error) {
	return dc.CreateMangaContext(context.Background(), newManga)
}

// CreateMangaContext : CreateManga with custom context.
func (dc *DexClient) CreateMangaContext(ctx context.Context, newManga io.Reader) (*MangaResponse, error) {
	var mr MangaResponse
	err := dc.responseOp(ctx, http.MethodPost, CreateMangaPath, newManga, &mr)
	return &mr, err
}

// ViewManga : View a manga by ID.
// https://api.mangadex.org/docs.html#operation/get-manga-id
func (dc *DexClient) ViewManga(id string) (*MangaResponse, error) {
	return dc.ViewMangaContext(context.Background(), id)
}

// ViewMangaContext : ViewManga with custom context.
func (dc *DexClient) ViewMangaContext(ctx context.Context, id string) (*MangaResponse, error) {
	var mr MangaResponse
	err := dc.responseOp(ctx, http.MethodGet, fmt.Sprintf(ViewMangaPath, id), nil, &mr)
	return &mr, err
}

// UpdateManga : Update a Manga.
// https://api.mangadex.org/docs.html#operation/put-manga-id
func (dc *DexClient) UpdateManga(id string, upManga io.Reader) (*MangaResponse, error) {
	return dc.UpdateMangaContext(context.Background(), id, upManga)
}

// UpdateMangaContext : UpdateManga with custom context.
func (dc *DexClient) UpdateMangaContext(ctx context.Context, id string, upManga io.Reader) (*MangaResponse, error) {
	var mr MangaResponse
	err := dc.responseOp(ctx, http.MethodPut, fmt.Sprintf(UpdateMangaPath, id), upManga, &mr)
	return &mr, err
}

// DeleteManga : Delete a Manga through ID.
// https://api.mangadex.org/docs.html#operation/delete-manga-id
func (dc *DexClient) DeleteManga(id string) error {
	return dc.DeleteMangaContext(context.Background(), id)
}

// DeleteMangaContext : DeleteManga with custom context.
func (dc *DexClient) DeleteMangaContext(ctx context.Context, id string) error {
	return dc.responseOp(ctx, http.MethodDelete, fmt.Sprintf(DeleteMangaPath, id), nil, nil)
}

// UnfollowManga : Unfollow a Manga by ID.
// https://api.mangadex.org/docs.html#operation/delete-manga-id-follow
func (dc *DexClient) UnfollowManga(id string) error {
	return dc.UnfollowMangaContext(context.Background(), id)
}

// UnfollowMangaContext : UnfollowManga with custom context.
func (dc *DexClient) UnfollowMangaContext(ctx context.Context, id string) error {
	return dc.responseOp(ctx, http.MethodDelete, fmt.Sprintf(UnfollowMangaPath, id), nil, nil)
}

// FollowManga : Follow a Manga by ID.
// https://api.mangadex.org/docs.html#operation/post-manga-id-follow
func (dc *DexClient) FollowManga(id string) error {
	return dc.FollowMangaContext(context.Background(), id)
}

// FollowMangaContext : FollowManga with custom context.
func (dc *DexClient) FollowMangaContext(ctx context.Context, id string) error {
	return dc.responseOp(ctx, http.MethodPost, fmt.Sprintf(FollowMangaPath, id), nil, nil)
}

// MangaFeed : Get Manga feed by ID.
// https://api.mangadex.org/docs.html#operation/get-manga-id-feed
func (dc *DexClient) MangaFeed(id string, params url.Values) (*ChapterList, error) {
	return dc.MangaFeedContext(context.Background(), id, params)
}

// MangaFeedContext : MangaFeed with custom context.
func (dc *DexClient) MangaFeedContext(ctx context.Context, id string, params url.Values) (*ChapterList, error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = fmt.Sprintf(MangaFeedPath, id)

	// Set request parameters
	u.RawQuery = params.Encode()

	var l ChapterList
	_, err := dc.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &l)
	return &l, err
}

// MangaReadMarkers : Get list of Chapter IDs that are marked as read for a specified manga ID.
// https://api.mangadex.org/docs.html#operation/get-manga-chapter-readmarkers
func (dc *DexClient) MangaReadMarkers(id string) (*ChapterReadMarkersResponse, error) {
	return dc.MangaReadMarkersContext(context.Background(), id)
}

// MangaReadMarkersContext : MangaReadMarkers with custom context.
func (dc *DexClient) MangaReadMarkersContext(ctx context.Context, id string) (*ChapterReadMarkersResponse, error) {
	var rmr ChapterReadMarkersResponse
	err := dc.responseOp(ctx, http.MethodGet, fmt.Sprintf(MangaReadMarkersPath, id), nil, &rmr)
	return &rmr, err
}

// GetRandomManga : Return a random Manga.
// https://api.mangadex.org/docs.html#operation/get-manga-random
func (dc *DexClient) GetRandomManga() (*MangaResponse, error) {
	return dc.GetRandomMangaContext(context.Background())
}

// GetRandomMangaContext : GetRandomManga with custom context.
func (dc *DexClient) GetRandomMangaContext(ctx context.Context) (*MangaResponse, error) {
	var mr MangaResponse
	err := dc.responseOp(ctx, http.MethodGet, GetRandomMangaPath, nil, &mr)
	return &mr, err
}

// TagList : Get tag list.
// https://api.mangadex.org/docs.html#operation/get-manga-tag
func (dc *DexClient) TagList() (*TagResponse, error) {
	return dc.TagListContext(context.Background())
}

// TagListContext : TagList with custom context.
func (dc *DexClient) TagListContext(ctx context.Context) (*TagResponse, error) {
	var tg TagResponse
	err := dc.responseOp(ctx, http.MethodGet, TagListPath, nil, &tg)
	return &tg, err
}

// UpdateMangaReadingStatus : Update reading status for a manga.
func (dc *DexClient) UpdateMangaReadingStatus(id string, status ReadStatus) error {
	return dc.UpdateMangaReadingStatusContext(context.Background(), id, status)
}

// UpdateMangaReadingStatusContext : UpdateMangaReadingStatus with custom context.
func (dc *DexClient) UpdateMangaReadingStatusContext(ctx context.Context, id string, status ReadStatus) error {
	// Create required request body.
	req := map[string]ReadStatus{
		"status": status,
	}
	rbytes, err := json.Marshal(&req)
	if err != nil {
		return err
	}

	return dc.responseOp(ctx, http.MethodPost,
		fmt.Sprintf(UpdateMangaReadingStatusPath, id), bytes.NewBuffer(rbytes),
		nil)
}