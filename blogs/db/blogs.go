package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/ellemouton/thunder/blogs"
)

func Create(ctx context.Context, dbc *sql.DB, name, description string, text string) (int64, error) {
	tx, err := dbc.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	res, err := tx.ExecContext(ctx, "insert into articles_content set text=?", text)
	if err != nil {
		return 0, err
	}

	contentID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	res, err = tx.ExecContext(ctx, "insert into articles_info set name=?, description=?, "+
		"created_at=?, content_id=?", name, description, time.Now(), contentID)
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func LookupInfo(ctx context.Context, dbc *sql.DB, id int64) (*blogs.Info, error) {
	row := dbc.QueryRowContext(ctx, "select * from articles_info where id=?", id)

	var info blogs.Info
	err := row.Scan(&info.ID, &info.Name, &info.Description, &info.CreatedAt, &info.ContentID)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

func LookupContent(ctx context.Context, dbc *sql.DB, id int64) (*blogs.Content, error) {
	row := dbc.QueryRowContext(ctx, "select * from articles_content where id=?", id)

	var content blogs.Content
	err := row.Scan(&content.ID, &content.Text)
	if err != nil {
		return nil, err
	}

	return &content, nil
}

func ListAllInfoRev(ctx context.Context, dbc *sql.DB) (infos []*blogs.Info, err error) {
	rows, err := dbc.Query("select * from articles_info order by id desc")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		info := blogs.Info{}
		err = rows.Scan(&info.ID, &info.Name, &info.Description, &info.CreatedAt, &info.ContentID)
		if err != nil {
			return nil, err
		}
		infos = append(infos, &info)
	}
	return infos, rows.Err()
}

func UpdateBlog(ctx context.Context, dbc *sql.DB, id int64, name, abstract, content string) error {
	info, err := LookupInfo(ctx, dbc, id)
	if err != nil {
		return err
	}

	tx, err := dbc.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "update articles_content set text=? where id=?", content, info.ContentID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "update articles_info set name=?, description=? where id=?", name, abstract, id)
	if err != nil {
		return err
	}

	return tx.Commit()
}
