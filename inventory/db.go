// db.go - Part of the `inventory` Package
//
//     ॐ भूर्भुवः स्वः
//     तत्स॑वि॒तुर्वरे॑ण्यं॒
//    भर्गो॑ दे॒वस्य॑ धीमहि।
//   धियो॒ यो नः॑ प्रचो॒दया॑त्॥
//
//
//  बी.वी.एल - बोसजी के द्वारा रचित भंडार लेखांकन हेतु तन्त्राक्ष्।
// =============================================
//
// एक सुगम एवं उपयोगी भंडार संचालन हेतु तन्त्राक्ष्।
//
// एक रचनात्मक भारतीय उत्पाद ।
//
// bvl - Boseji's Inventory Management Program
//
// Easy to use and useful stock, goods and materials handling software.
//
// Sources
// -------
// https://github.com/boseji/bvl
//
// License
// -------
//
//   bvl - Boseji's Inventory Management Program.
//   Copyright (C) 2025 by Abhijit Bose (aka. Boseji)
//
//   This program is free software: you can redistribute it and/or modify
//   it under the terms of the GNU General Public License version 2 only
//   as published by the Free Software Foundation.
//
//   This program is distributed in the hope that it will be useful,
//   but WITHOUT ANY WARRANTY; without even the implied warranty of
//   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
//
//   You should have received a copy of the GNU General Public License
//   along with this program. If not, see <https://www.gnu.org/licenses/>.
//
//  SPDX-License-Identifier: GPL-2.0-only
//  Full Name: GNU General Public License v2.0 only
//  Please visit <https://spdx.org/licenses/GPL-2.0-only.html> for details.
//

package inventory

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const (
	// Starting value of the Index
	IndexStart = 1000
)

// OpenDB opens (or creates) the SQLite database at dbFile.
//
// It ensures the inventory table is created, and initializes the
// sequence to 1000 if missing.
//
// Returns an open *sql.DB handle.
func OpenDB(dbFile string) *sql.DB {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	// Create the DB if It does not exists
	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS inventory (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        description TEXT,
        location TEXT,
        status TEXT,
        remarks TEXT
    );
    `)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	// Generate Query:
	// INSERT INTO sqlite_sequence (name, seq)
	// SELECT 'inventory', 1000
	// WHERE NOT EXISTS (SELECT 1 FROM sqlite_sequence WHERE name = 'inventory')
	queryStr := "INSERT INTO sqlite_sequence (name, seq)" +
		fmt.Sprintf("SELECT 'inventory', %d", IndexStart) +
		"WHERE NOT EXISTS (" +
		"SELECT 1 FROM sqlite_sequence WHERE name = 'inventory')"

	// Create the Index start if one does not exists
	_, err = db.Exec(queryStr)
	if err != nil {
		log.Printf("Note: could not init sequence: %v", err)
	}

	return db
}

// GetItemByID returns a single Item by its ID.
//
// If not found, returns an error.
func GetItemByID(db *sql.DB, id int) (Item, error) {
	var item Item
	row := db.QueryRow(`
        SELECT id, description, location, status, remarks
        FROM inventory WHERE id = ?`, id)
	err := row.Scan(
		&item.ID, &item.Description, &item.Location,
		&item.Status, &item.Remarks)
	if err != nil {
		if err == sql.ErrNoRows {
			return item, fmt.Errorf("Item %d not found", id)
		}
		return item, fmt.Errorf("query failed: %v", err)
	}
	return item, nil
}

// SearchItems performs full-text search across:
// Description, Location, Status.
//
// Returns all matching Items.
func SearchItems(db *sql.DB, search string) ([]Item, error) {
	like := "%" + search + "%"
	rows, err := db.Query(`
        SELECT id, description, location, status, remarks
        FROM inventory
        WHERE description LIKE ?
           OR location LIKE ?
           OR status LIKE ?
        ORDER BY id`, like, like, like)
	if err != nil {
		return nil, fmt.Errorf("search query failed: %v", err)
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		err := rows.Scan(
			&item.ID, &item.Description, &item.Location,
			&item.Status, &item.Remarks)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %v", err)
		}
		items = append(items, item)
	}
	return items, nil
}

// SearchByDescription returns Items where description LIKE search.
func SearchByDescription(db *sql.DB, search string) ([]Item, error) {
	like := "%" + search + "%"
	rows, err := db.Query(`
        SELECT id, description, location, status, remarks
        FROM inventory
        WHERE description LIKE ?
        ORDER BY id`, like)
	if err != nil {
		return nil, fmt.Errorf(
			"search by description failed: %v", err)
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.ID, &item.Description,
			&item.Location, &item.Status, &item.Remarks)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %v", err)
		}
		items = append(items, item)
	}
	return items, nil
}

// SearchByLocation returns Items where location LIKE search.
func SearchByLocation(db *sql.DB, search string) ([]Item, error) {
	like := "%" + search + "%"
	rows, err := db.Query(`
        SELECT id, description, location, status, remarks
        FROM inventory
        WHERE location LIKE ?
        ORDER BY id`, like)
	if err != nil {
		return nil, fmt.Errorf(
			"search by location failed: %v", err)
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.ID, &item.Description,
			&item.Location, &item.Status, &item.Remarks)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %v", err)
		}
		items = append(items, item)
	}
	return items, nil
}

// SearchByStatus returns Items where status LIKE search.
func SearchByStatus(db *sql.DB, search string) ([]Item, error) {
	like := "%" + search + "%"
	rows, err := db.Query(`
        SELECT id, description, location, status, remarks
        FROM inventory
        WHERE status LIKE ?
        ORDER BY id`, like)
	if err != nil {
		return nil, fmt.Errorf(
			"search by status failed: %v", err)
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.ID, &item.Description,
			&item.Location, &item.Status, &item.Remarks)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %v", err)
		}
		items = append(items, item)
	}
	return items, nil
}

// SearchByRemarks returns Items where remarks LIKE search.
func SearchByRemarks(db *sql.DB, search string) ([]Item, error) {
	like := "%" + search + "%"
	rows, err := db.Query(`
        SELECT id, description, location, status, remarks
        FROM inventory
        WHERE remarks LIKE ?
        ORDER BY id`, like)
	if err != nil {
		return nil, fmt.Errorf(
			"search by remarks failed: %v", err)
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.ID, &item.Description,
			&item.Location, &item.Status, &item.Remarks)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %v", err)
		}
		items = append(items, item)
	}
	return items, nil
}

// ListAll returns all items from the inventory table,
// ordered by ID. This is a simple bulk fetch function.
//
// Use cautiously for very large databases. For pagination,
// use ListItemsPaged() or ItemIterator().
func ListAll(db *sql.DB) ([]Item, error) {
	rows, err := db.Query(`
        SELECT id, description, location, status, remarks
        FROM inventory
        ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("query failed: %v", err)
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.ID, &item.Description,
			&item.Location, &item.Status, &item.Remarks)
		if err != nil {
			return nil, fmt.Errorf("qcan failed: %v", err)
		}
		items = append(items, item)
	}

	return items, nil
}

// ListItemsPaged returns up to limit Items,
// starting after given ID (afterID).
//
// Useful for paging through large inventories.
func ListItemsPaged(
	db *sql.DB, afterID int, limit int) ([]Item, error) {

	rows, err := db.Query(`
        SELECT id, description, location, status, remarks
        FROM inventory
        WHERE id > ?
        ORDER BY id
        LIMIT ?`, afterID, limit)
	if err != nil {
		return nil, fmt.Errorf("paged query failed: %v", err)
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		err := rows.Scan(
			&item.ID, &item.Description, &item.Location,
			&item.Status, &item.Remarks)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %v", err)
		}
		items = append(items, item)
	}
	return items, nil
}

// ItemIterator streams matching records one-by-one.
//
// Provides: Next(), Close().
type ItemIterator struct {
	rows *sql.Rows
}

// NewItemIterator creates an iterator matching WHERE clause.
//
// Example:
//
//	it, err := NewItemIterator(db, "WHERE status LIKE ?", "%Available%")
func NewItemIterator(
	db *sql.DB, whereClause string, args ...interface{},
) (*ItemIterator, error) {

	query := `
        SELECT id, description, location, status, remarks
        FROM inventory `
	if whereClause != "" {
		query += whereClause
	}
	query += " ORDER BY id"

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("iterator query failed: %v", err)
	}

	return &ItemIterator{rows: rows}, nil
}

// Next returns next Item from iterator.
//
// ok=false when no more rows.
// Always call Close() after use.
func (it *ItemIterator) Next() (Item, bool, error) {
	var item Item
	if it.rows.Next() {
		err := it.rows.Scan(
			&item.ID, &item.Description, &item.Location,
			&item.Status, &item.Remarks)
		if err != nil {
			return item, false, fmt.Errorf("iterator scan failed: %v", err)
		}
		return item, true, nil
	}
	return item, false, nil
}

// Close closes the iterator.
func (it *ItemIterator) Close() error {
	return it.rows.Close()
}

// AddItem inserts a new Item.
func AddItem(db *sql.DB, item Item) error {
	_, err := db.Exec(`
        INSERT INTO inventory
        (description, location, status, remarks)
        VALUES (?, ?, ?, ?)`,
		item.Description, item.Location,
		item.Status, item.Remarks)
	if err != nil {
		return fmt.Errorf("insert failed: %v", err)
	}
	return nil
}

// EditItem updates an existing Item.
func EditItem(db *sql.DB, item Item) error {
	_, err := db.Exec(`
        UPDATE inventory
        SET description = ?, location = ?,
            status = ?, remarks = ?
        WHERE id = ?`,
		item.Description, item.Location,
		item.Status, item.Remarks, item.ID)
	if err != nil {
		return fmt.Errorf("update failed: %v", err)
	}
	return nil
}

// DeleteItem deletes an Item by ID.
func DeleteItem(db *sql.DB, id int) error {
	_, err := db.Exec(`
        DELETE FROM inventory
        WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete failed: %v", err)
	}
	return nil
}

// ResetSequence sets the ID auto-increment
// to restart from 1001.
func ResetSequence(db *sql.DB) error {
	_, err := db.Exec(
		fmt.Sprintf("UPDATE sqlite_sequence"+
			"SET seq = %d"+
			"WHERE name = 'inventory'", IndexStart))
	if err != nil {
		return fmt.Errorf("reset sequence failed: %v", err)
	}
	return nil
}
