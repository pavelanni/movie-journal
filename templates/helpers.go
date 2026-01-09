// Package templates provides template helpers and rendering utilities for the Movie Journal application.
package templates

import "github.com/pavelanni/movie-journal/internal/models"

func getWatchedDate(entry *models.DiaryEntry) string {
	if entry != nil {
		return entry.WatchedDate.Format("2006-01-02")
	}
	return ""
}

func getMovieTitle(entry *models.DiaryEntry) string {
	if entry != nil {
		return entry.Movie.Title
	}
	return ""
}

func getWatchedLocation(entry *models.DiaryEntry) string {
	if entry != nil {
		return entry.WatchedLocation
	}
	return ""
}

func getWatchedWith(entry *models.DiaryEntry) string {
	if entry != nil {
		return entry.WatchedWith
	}
	return ""
}

func getNotes(entry *models.DiaryEntry) string {
	if entry != nil {
		return entry.Notes
	}
	return ""
}

func getStarClass(rating int) string {
	switch {
	case rating >= 4:
		return "w-4 h-4 text-green-400"
	case rating == 3:
		return "w-4 h-4 text-yellow-400"
	default:
		return "w-4 h-4 text-red-400"
	}
}
