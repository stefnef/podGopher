package rss

import (
	"database/sql"
	"podGopher/core/domain/model"
)

type PostgresRssOutAdapter struct {
	db *sql.DB
}

func (adapter *PostgresRssOutAdapter) GetRSSOrNil(showSlug string, distributionSlug string) (*model.RSS, error) {
	query := `
		SELECT 
			s.id, s.title, s.slug,
			d.id, d.show_id, d.title, d.slug,
			e.id, e.show_id, e.title
		FROM show s
		JOIN distribution d ON d.show_id = s.id
		LEFT JOIN episodes_distributions ed ON ed.distribution_id = d.id
		LEFT JOIN episode e ON e.id = ed.episode_id
		WHERE s.slug = $1 AND d.slug = $2
	`
	rows, err := adapter.db.Query(query, showSlug, distributionSlug)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	return parseRss(rows), nil
}

func parseRss(rows *sql.Rows) *model.RSS {
	var rss *model.RSS
	episodeSet := make(map[string]*model.Episode)

	for rows.Next() {
		var (
			sId, sTitle, sSlug          string
			dId, dShowId, dTitle, dSlug string
			eId, eShowId, eTitle        sql.NullString
		)

		_ = rows.Scan(
			&sId, &sTitle, &sSlug,
			&dId, &dShowId, &dTitle, &dSlug,
			&eId, &eShowId, &eTitle,
		)

		if rss == nil {
			rss = &model.RSS{
				Show: &model.Show{
					Id:    sId,
					Title: sTitle,
					Slug:  sSlug,
				},
				Distribution: &model.Distribution{
					Id:     dId,
					ShowId: dShowId,
					Title:  dTitle,
					Slug:   dSlug,
				},
			}
		}

		if eId.Valid {
			episodeSet[eId.String] = &model.Episode{
				Id:     eId.String,
				ShowId: eShowId.String,
				Title:  eTitle.String,
			}
		}
	}

	if rss == nil {
		return nil
	}

	for _, episode := range episodeSet {
		rss.Distribution.Episodes = append(rss.Distribution.Episodes, episode.Id)
		rss.Episodes = append(rss.Episodes, episode)
	}
	return rss
}

func NewPostgresRSSRepository(db *sql.DB) *PostgresRssOutAdapter {
	return &PostgresRssOutAdapter{db: db}
}
