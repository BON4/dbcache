package repo

import (
	"context"
	"database/sql"
	"dbcache/models"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type PsqlRepository struct {
	DB *sql.DB
}

func NewPsqlRepository(database *sql.DB) Repository {
	return &PsqlRepository{DB: database}
}

func (p *PsqlRepository) GetItem(ctx context.Context, itemId string) (*models.Item, error) {
	var item models.Item
	var sqlStr string = "SELECT * from Items where id = $1"
	err := p.DB.QueryRowContext(ctx, sqlStr, itemId).Scan(&item.Id, &item.TransportId, &item.Number)

	switch {
	case err == sql.ErrNoRows:
		return nil, errors.New("Not found")
	case err != nil:
		return nil, err
	default:
		return &item, nil
	}
}

func (p *PsqlRepository) GetTransport(ctx context.Context, transportId string) (*models.Transport, error) {
	var transport models.Transport
	var sqlStr string = "SELECT * from Transport where transport_id = $1"
	err := p.DB.QueryRowContext(ctx, sqlStr, transportId).Scan(&transport.Id, &transport.Number)

	switch {
	case err == sql.ErrNoRows:
		return nil, errors.New("Not found")
	case err != nil:
		return nil, err
	default:
		return &transport, nil
	}
}

func (p *PsqlRepository) GetTransportItemView(ctx context.Context, transportId string) (*models.TransportItemView, error) {
	var sqlStr string = "select * from transport_view where transport_id = $1"

	rows, err := p.DB.QueryContext(ctx, sqlStr, transportId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var transport models.Transport
	var items []models.Item = make([]models.Item, 0)

	for rows.Next() {
		var item models.Item

		if err := rows.Scan(&transport.Id, &transport.Number, &item.Id, &item.Number); err != nil {
			return nil, err
		}

		item.TransportId = transport.Id

		items = append(items, item)
	}

	return &models.TransportItemView{Transport: transport, Items: items}, nil
}

func (p *PsqlRepository) CreateAlotItems(ctx context.Context, itemsConf models.CreateAlotItems) error {
	tx, err := p.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	items := make([]string, 0, itemsConf.Length*2)

	for i := 0; i < itemsConf.Length; i++ {
		var itemString string = fmt.Sprintf(`(%s, '%s')`, itemsConf.TransportId, uuid.New().String())
		items = append(items, itemString)
	}

	smt := fmt.Sprintf("INSERT INTO Items(transport_id, number) values %s", strings.Join(items, ","))

	_, txerr := tx.ExecContext(ctx, smt)

	if txerr != nil {
		tx.Rollback()
		return txerr
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
