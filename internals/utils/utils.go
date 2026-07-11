package utils

import (
	"GoNewsScrapper/internals/database"
)

func RemoveNewsDuplicatesBySlug(news []database.CreateNewsParams) []database.CreateNewsParams {
	results := make([]database.CreateNewsParams, 0, len(news))
	seen := map[string]struct{}{}

	for _, newsItem := range news {
		if _, ok := seen[newsItem.Slug]; !ok {
			results = append(results, newsItem)
			seen[newsItem.Slug] = struct{}{}
		}
	}

	return results
}
