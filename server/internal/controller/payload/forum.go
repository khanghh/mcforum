package payload

import "bbs-go/internal/model"

type ForumResponse struct {
	Id          int64  `json:"id,omitempty"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Logo        string `json:"logo,omitempty"`
	Description string `json:"description,omitempty"`
	CanWrite    bool   `json:"canWrite,omitempty"`
}

func BuildForum(forum *model.Forum, rank int) *ForumResponse {
	if forum == nil {
		return nil
	}
	return &ForumResponse{
		Id:          forum.ID,
		Name:        forum.Name,
		Slug:        forum.Slug,
		Logo:        forum.Logo,
		Description: forum.Description,
		CanWrite:    rank >= forum.WriteRank,
	}
}

func BuildForumList(forums []model.Forum, rank int) []ForumResponse {
	if len(forums) == 0 {
		return nil
	}
	var ret []ForumResponse
	for _, forum := range forums {
		ret = append(ret, *BuildForum(&forum, rank))
	}
	return ret
}
