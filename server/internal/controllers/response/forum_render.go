package response

import "bbs-go/internal/models"

func BuildForum(forum *models.Forum) *ForumResponse {
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

func BuildForumList(forums []models.Forum) []ForumResponse {
	if len(forums) == 0 {
		return nil
	}
	var ret []ForumResponse
	for _, forum := range forums {
		ret = append(ret, *BuildForum(&forum))
	}
	return ret
}
