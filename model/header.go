package model

import "fmt"

type Header struct {
	ID    uint   `json:"id,omitempty" db:"id"`
	Name  string `json:"name" db:"name"`
	Value string `json:"value" value:"value"`
}

func (c *Conn) SaveHeader(header *Header, userID, routeID uint) (err error) {
	stmt := fmt.Sprintf(`
		INSERT INTO	header(user_id, route_id, name, value)
		VALUES ($1, $2, pgp_sym_encrypt($3, '%v'), pgp_sym_encrypt($4, '%v'))
		RETURNING id
	`, PG_CRYPT_KEY, PG_CRYPT_KEY)

	row := c.db.QueryRow(stmt, userID, routeID, header.Name, header.Value)

	err = row.Scan(&header.ID)
	if err != nil {
		return
	}

	return nil
}

func (c *Conn) UpdateHeader(header Header, userID, routeID uint) (err error) {
	stmt := fmt.Sprintf(`
		UPDATE header
		SET name = pgp_sym_encrypt($1, '%v'), value = pgp_sym_encrypt($2, '%v')
		WHERE id= $3 AND user_id = $4 AND route_id = $5
	`, PG_CRYPT_KEY, PG_CRYPT_KEY)

	_, err = c.db.Exec(stmt, header.Name, header.Value, header.ID, userID, routeID)
	if err != nil {
		return
	}

	return nil
}

func (c *Conn) DeleteHeader(headerName string, userID, routeID uint) (err error) {
	stmt, err := c.db.Prepare(fmt.Sprintf(`
		DELETE FROM header WHERE "name" = pgp_sym_encrypt($1, '%v') AND user_id = $2 AND route_id = $3
		`, PG_CRYPT_KEY),
	)

	if err != nil {
		return
	}

	_, err = stmt.Exec(headerName, userID, routeID)
	if err != nil {
		return
	}

	return nil
}

func (c *Conn) GetHeader(userID, routeID uint) (headers []Header, err error) {
	stmt := fmt.Sprintf(`
		SELECT id, pgp_sym_decrypt(name::bytea, '%v') as name, pgp_sym_decrypt(value::bytea, '%v') as value
		FROM header
		WHERE user_id = $1 AND route_id = $2
	`, PG_CRYPT_KEY, PG_CRYPT_KEY)

	rows, err := c.db.Query(stmt, userID, routeID)
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var header Header

		err := rows.Scan(&header.ID, &header.Name, &header.Value)
		if err != nil {
			return nil, err
		}

		headers = append(headers, header)
	}

	if err = rows.Err(); err != nil {
		return
	}

	return
}
