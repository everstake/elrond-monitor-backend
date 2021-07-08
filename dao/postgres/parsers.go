package postgres

import (
	"github.com/Masterminds/squirrel"
	"github.com/everstake/elrond-monitor-backend/dao/dmodels"
)

func (db Postgres) GetParsers() (parsers []dmodels.Parser, err error) {
	q := squirrel.Select("*").From(dmodels.ParsersTable)
	err = db.find(&parsers, q)
	if err != nil {
		return nil, err
	}
	return parsers, nil
}

func (db Postgres) GetParser(title string) (parser dmodels.Parser, err error) {
	q := squirrel.Select("*").From(dmodels.ParsersTable).
		Where(squirrel.Eq{"par_title": title})
	err = db.first(&parser, q)
	return parser, err
}

func (db Postgres) UpdateParserHeight(parser dmodels.Parser) error {
	q := squirrel.Update(dmodels.ParsersTable).
		Where(squirrel.Eq{"par_id": parser.ID}).
		SetMap(map[string]interface{}{
			"par_height": parser.Height,
		})
	return db.update(q)
}
