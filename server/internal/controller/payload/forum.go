package payload

import "bbs-go/internal/model"

type ForumResponse struct {
	Id          int64  `json:"id,omitempty"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Logo        string `json:"logo,omitempty"`
	Description string `json:"description,omitempty"`
}

func BuildForum(forum *model.Forum) *ForumResponse {
	if forum == nil {
		return nil
	}
	return &ForumResponse{
		Id:          forum.Id,
		Name:        forum.Name,
		Slug:        forum.Slug,
		Logo:        forum.Logo,
		Description: forum.Description,
	}
}

func BuildForumList(forums []model.Forum) []ForumResponse {
	if len(forums) == 0 {
		return nil
	}
	var ret []ForumResponse
	for _, forum := range forums {
		ret = append(ret, *BuildForum(&forum))
	}
	return ret
}
