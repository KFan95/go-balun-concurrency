package compute

import (
	"database/internal/database/query"
	"errors"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"go.uber.org/zap"
)

type ParserInterface interface {
	Parse(string) (query.QueryInterface, error)
}

type Parser struct {
	logger           *zap.Logger
	participleParser *participle.Parser[pQuery]
}

type pQuery struct {
	SetQuery *pSetQuery `  "SET" @@`
	GetQuery *pGetQuery `| "GET" @@`
	DelQuery *pDelQuery `| "DEL" @@`
}

type pSetQuery struct {
	Key   string `@Key`
	Value string `@Key`
}

type pGetQuery struct {
	Key string `@Key`
}

type pDelQuery struct {
	Key string `@Key`
}

func NewParser(logger *zap.Logger) (ParserInterface, error) {
	if logger == nil {
		return nil, errors.New("logger is required")
	}

	simplelexer := lexer.MustSimple([]lexer.SimpleRule{
		{Name: "Key", Pattern: `[a-zA-Z\d\*\/_]+`},
		{Name: "Value", Pattern: `[a-zA-Z\d\*\/_]+`},
		{Name: "Whitespace", Pattern: `\s+`},
	})

	participleParser := participle.MustBuild[pQuery](
		participle.Lexer(simplelexer),
		participle.Elide("Whitespace"),
	)

	return &Parser{
		logger:           logger,
		participleParser: participleParser,
	}, nil
}

func (p *Parser) Parse(queryStr string) (query.QueryInterface, error) {
	queryStr = strings.Trim(queryStr, " \r\n\t")

	p.logger.Info("parsing query", zap.String("query", queryStr))

	q, err := p.participleParser.ParseString("", queryStr)

	if err != nil {
		return nil, err
	}

	if q.GetQuery != nil {
		p.logger.Info("parsed get query", zap.Reflect("query", q.GetQuery))

		return &query.Get{
			Key: q.GetQuery.Key,
		}, nil
	} else if q.SetQuery != nil {
		p.logger.Info("parsed set query", zap.Reflect("query", q.SetQuery))

		return &query.Set{
			Key:   q.SetQuery.Key,
			Value: q.SetQuery.Value,
		}, nil
	} else if q.DelQuery != nil {
		p.logger.Info("parsed del query", zap.Reflect("query", q.DelQuery))

		return &query.Del{
			Key: q.DelQuery.Key,
		}, nil
	}

	return nil, errors.New("invalid query")
}
