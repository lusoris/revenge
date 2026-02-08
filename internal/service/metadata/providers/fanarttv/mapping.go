package fanarttv

import (
	"strconv"

	"github.com/lusoris/revenge/internal/service/metadata"
)

// mapMovieImages converts Fanart.tv movie images to the internal format.
func mapMovieImages(resp *MovieResponse) *metadata.Images {
	if resp == nil {
		return nil
	}

	images := &metadata.Images{}

	// HD logos + standard logos → Logos
	images.Logos = append(images.Logos, mapFanartImages(resp.HDMovieLogos, "logo")...)
	images.Logos = append(images.Logos, mapFanartImages(resp.MovieLogos, "logo")...)

	// Posters
	images.Posters = mapFanartImages(resp.MoviePosters, "poster")

	// Backdrops (backgrounds)
	images.Backdrops = mapFanartImages(resp.MovieBackdrops, "backdrop")

	// Clearart, disc art, banners, thumbs → stored as additional images
	// Map banners as backdrops (wide format)
	images.Backdrops = append(images.Backdrops, mapFanartImages(resp.MovieBanners, "banner")...)

	// Thumbs as stills (landscape thumbnails)
	images.Stills = mapFanartImages(resp.MovieThumbs, "thumb")

	return images
}

// mapTVShowImages converts Fanart.tv TV show images to the internal format.
func mapTVShowImages(resp *TVShowResponse) *metadata.Images {
	if resp == nil {
		return nil
	}

	images := &metadata.Images{}

	// HD logos + clear logos → Logos
	images.Logos = append(images.Logos, mapFanartImages(resp.HDTVLogos, "logo")...)
	images.Logos = append(images.Logos, mapFanartImages(resp.ClearLogos, "logo")...)

	// Posters
	images.Posters = mapFanartImages(resp.TVPosters, "poster")

	// Backdrops (show backgrounds)
	images.Backdrops = mapFanartImages(resp.ShowBackdrops, "backdrop")
	images.Backdrops = append(images.Backdrops, mapFanartImages(resp.TVBanners, "banner")...)

	// Thumbs
	images.Stills = mapFanartImages(resp.TVThumbs, "thumb")

	// Clearart and character art → Profiles (artwork of characters/actors)
	images.Profiles = append(images.Profiles, mapFanartImages(resp.HDClearArt, "clearart")...)
	images.Profiles = append(images.Profiles, mapFanartImages(resp.ClearArt, "clearart")...)
	images.Profiles = append(images.Profiles, mapFanartImages(resp.CharacterArt, "characterart")...)

	return images
}

// mapSeasonImages extracts season-specific images from a TV show response.
func mapSeasonImages(resp *TVShowResponse, seasonNum int) *metadata.Images {
	if resp == nil {
		return nil
	}

	seasonStr := strconv.Itoa(seasonNum)
	images := &metadata.Images{}

	for _, img := range resp.SeasonPosters {
		if img.Season == seasonStr {
			images.Posters = append(images.Posters, mapSeasonImage(img, "poster"))
		}
	}
	for _, img := range resp.SeasonThumbs {
		if img.Season == seasonStr {
			images.Stills = append(images.Stills, mapSeasonImage(img, "thumb"))
		}
	}
	for _, img := range resp.SeasonBanners {
		if img.Season == seasonStr {
			images.Backdrops = append(images.Backdrops, mapSeasonImage(img, "banner"))
		}
	}

	if len(images.Posters) == 0 && len(images.Stills) == 0 && len(images.Backdrops) == 0 {
		return nil
	}
	return images
}

// mapFanartImages converts a slice of FanartImage to metadata.Image.
func mapFanartImages(imgs []FanartImage, _ string) []metadata.Image {
	if len(imgs) == 0 {
		return nil
	}

	result := make([]metadata.Image, 0, len(imgs))
	for _, img := range imgs {
		mi := metadata.Image{
			FilePath: img.URL,
		}
		if img.Lang != "" && img.Lang != "00" {
			mi.Language = &img.Lang
		}
		if likes, err := strconv.Atoi(img.Likes); err == nil {
			mi.VoteCount = likes
		}
		result = append(result, mi)
	}
	return result
}

// mapSeasonImage converts a SeasonImage to metadata.Image.
func mapSeasonImage(img SeasonImage, _ string) metadata.Image {
	mi := metadata.Image{
		FilePath: img.URL,
	}
	if img.Lang != "" && img.Lang != "00" {
		mi.Language = &img.Lang
	}
	if likes, err := strconv.Atoi(img.Likes); err == nil {
		mi.VoteCount = likes
	}
	return mi
}
