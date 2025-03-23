package response

import "bbs-go/internal/models"

func BuildTag(tag *models.Tag) *TagResponse {
	if tag == nil {
		return nil
	}
	return &TagResponse{Id: tag.Id, Name: tag.Name}
}

func BuildTags(tags []models.Tag) *[]TagResponse {
	if len(tags) == 0 {
		return nil
	}
	var responses []TagResponse
	for _, tag := range tags {
		responses = append(responses, *BuildTag(&tag))
	}
	return &responses
}
