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

// DB Access Layer
//
// Provides database access functions for Inventory CLI.
// All "write" functions accept an Execer interface to support transactions.
//
// Conventions:
// - Line width <= 80 characters
// - All errors lowercase, no punctuation
// - Documentation is verbose
//

package inventory

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/boseji/bsg/gen"
	_ "github.com/mattn/go-sqlite3"
)

// IndexStart defines the starting value for auto-incremented IDs.
const (
	// Starting value of the Index
	IndexStart = 1000
)

// Execer defines something that can Exec SQL.
// Both *sql.DB and *sql.Tx implement this.
type Execer interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

// InventoryDB wraps *sql.DB and provides safe transaction helpers.
//
// Users do not need to work with *sql.DB directly.
type InventoryDB struct {
	db *sql.DB
}

// NewInventoryDB opens or creates the database and returns InventoryDB.
//
// Ensures the table exists, sequence is initialized.
// Returns a ready-to-use InventoryDB wrapper.
//
// Usage:
//
//	inv := NewInventoryDB("inventory.db")
//
// Notes:
// - Underlying connection is stored in inv.db
// - Close() must be called when finished
// - Table creation is idempotent
func NewInventoryDB(dbFile string) *InventoryDB {
	db := OpenDB(dbFile)
	return &InventoryDB{db: db}
}

// WithTransaction executes the given function inside a transaction.
//
// Usage:
//
//	err := inv.WithTransaction(func(tx Execer) error {
//	    err := AddItem(tx, item)
//	    if err != nil {
//	        return err
//	    }
//	    return AppendRemarksEntry(tx, item.ID, "added new")
//	})
//
// If fn() returns error:
// - Transaction is rolled back
//
// If fn() returns nil:
// - Transaction is committed
//
// Notes:
// - Use for any group of changes that must be atomic
// - If the DB fails, returns error
func (inv *InventoryDB) WithTransaction(
	fn func(tx Execer) error) error {

	tx, err := inv.db.Begin()
	if err != nil {
		return fmt.Errorf("begin tx failed: %v", err)
	}

	err = fn(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit tx failed: %v", err)
	}

	return nil
}

// DB returns the underlying *sql.DB (for read-only queries).
// Use only when needed, e.g. for GetItemByID.
func (inv *InventoryDB) DB() *sql.DB {
	return inv.db
}

// Close closes the underlying database connection.
func (inv *InventoryDB) Close() error {
	return inv.db.Close()
}

// OpenDB opens or creates the SQLite database file at dbFile path.
//
// It ensures that the 'inventory' table exists with the required fields:
// - id          INTEGER PRIMARY KEY AUTOINCREMENT
// - description TEXT
// - location    TEXT
// - status      TEXT
// - remarks     TEXT
//
// It also ensures that the autoincrement sequence is initialized:
// - If the sequence is missing, sets it to IndexStart.
//
// Usage:
//
//	db := OpenDB("inventory.db")
//
// Notes:
// - Returns a *sql.DB connection (ready to use)
// - Fails fatally if the database cannot be opened or schema is invalid
// - Table creation is idempotent (safe to call multiple times)
// - Auto-increment starts from IndexStart (default 1000)
func OpenDB(dbFile string) *sql.DB {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	// Create the inventory table if not present
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
		log.Fatalf("failed to create table: %v", err)
	}

	// Initialize sequence only if not already set
	_, err = db.Exec(`
    INSERT INTO sqlite_sequence (name, seq)
    SELECT 'inventory', ?
    WHERE NOT EXISTS (
        SELECT 1 FROM sqlite_sequence WHERE name = 'inventory'
    );`, IndexStart)
	if err != nil {
		log.Printf("could not init sequence: %v", err)
	}

	return db
}

// AppendItem inserts a new item if it does not exist,
// or replaces an existing item with the same ID.
//
// Works with both *sql.DB and *sql.Tx.
func AppendItem(exec Execer, item Item) error {
	_, err := exec.Exec(`
        INSERT OR REPLACE INTO inventory
        (id, description, location, status, remarks)
        VALUES (?, ?, ?, ?, ?)`,
		item.ID, item.Description, item.Location,
		item.Status, item.FormatRemarks())
	if err != nil {
		return fmt.Errorf("insert or replace failed: %v", err)
	}
	return nil
}

// AppendRemarksEntry appends a new log entry to the item's
// remarks field, using the standard timestamp format.
//
// The entry is formatted as:
//
//	[YYYY-MM-DD HH:MM] message
//
// This function uses SQL to append without loading
// the entire remarks field:
//
//	SET remarks = COALESCE(remarks, '') || char(10) || ?
//
// Usage:
//
//	err := AppendRemarksEntry(tx, 1002, "replaced battery")
//
// Resulting remarks field:
//
//	(previous remarks)
//	[2025-06-20 16:55] replaced battery
//
// Notes:
// - Does not modify other fields (description, location, status)
// - If item ID does not exist, no rows are updated
// - Use when you only want to add an audit/log entry
func AppendRemarksEntry(exec Execer, id int, message string) error {
	t := gen.BST().Format("2006-01-02 15:04")
	formatted := fmt.Sprintf("[%s] %s", t, message)

	// Appends new entry to remarks field:
	// COALESCE(remarks, '') → ensure string (not null)
	// char(10) → newline character
	// '||' → SQLite concat
	// Result: remarks = old + '\n' + new entry
	_, err := exec.Exec(`
        UPDATE inventory
        SET remarks = 
            COALESCE(remarks, '') || char(10) || ?
        WHERE id = ?`,
		formatted, id)
	if err != nil {
		return fmt.Errorf("append to remarks failed: %v", err)
	}
	return nil
}

// AddItem inserts a new item into the inventory table.
//
// The ID is assigned automatically (auto-increment).
// The remarks field is always stored in timestamped format
// by calling item.FormatRemarks().
//
// Usage:
//
//	item := Item{
//	    Description: "New inverter",
//	    Location:    "Warehouse 1",
//	    Status:      "Operational",
//	    Remarks:     "installed and tested",
//	}
//	err := AddItem(tx, item)
//
// Resulting remarks field:
//
//	[2025-06-20 16:45] installed and tested
//
// Notes:
// - If used inside a transaction, pass tx as exec
// - If using plain db connection, pass db as exec
// - Remarks will always follow consistent format
func AddItem(exec Execer, item Item) error {
	_, err := exec.Exec(`
        INSERT INTO inventory
        (description, location, status, remarks)
        VALUES (?, ?, ?, ?)`,
		item.Description, item.Location,
		item.Status, item.FormatRemarks())
	if err != nil {
		return fmt.Errorf("insert failed: %v", err)
	}
	return nil
}

// EditItem updates the item's fields (description, location, status)
// and appends the new remarks text to the existing remarks field.
//
// Remarks field acts as an append-only log:
//   - Previous remarks are preserved
//   - New entry is appended with timestamp format:
//     [YYYY-MM-DD HH:MM] message
//
// This function does not load existing remarks in Go;
// it performs the append using SQL:
//
//	SET remarks = COALESCE(remarks, '') || char(10) || ?
//
// Usage:
//
//	item := Item{
//	    ID: 1002,
//	    Description: "Updated inverter",
//	    Location:    "Warehouse 3",
//	    Status:      "Operational",
//	    Remarks:     "maintenance check completed",
//	}
//	err := EditItem(tx, item)
//
// Resulting remarks field:
//
//	(previous remarks)
//	[2025-06-20 16:22] maintenance check completed
//
// Notes:
// - If item ID does not exist, no rows are updated
// - If used inside transaction (tx), pass tx as exec
// - To append a single new log entry, use AppendRemarksEntry()
// - To display remarks nicely, use item.FormatRemarks()
func EditItem(exec Execer, item Item) error {
	_, err := exec.Exec(`
        UPDATE inventory
        SET description = ?, location = ?,
            status = ?,
            remarks = COALESCE(remarks, '') || char(10) || ?
        WHERE id = ?`,
		item.Description, item.Location,
		item.Status,
		item.FormatRemarks(),
		item.ID)
	if err != nil {
		return fmt.Errorf("update failed: %v", err)
	}
	return nil
}

// DeleteItem deletes a row from inventory by ID.
func DeleteItem(exec Execer, id int) error {
	_, err := exec.Exec(`
        DELETE FROM inventory
        WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete failed: %v", err)
	}
	return nil
}

// ResetSequence resets the auto-increment index to IndexStart.
func ResetSequence(exec Execer) error {
	_, err := exec.Exec(`
        UPDATE sqlite_sequence
        SET seq = ?
        WHERE name = 'inventory'`, IndexStart)
	if err != nil {
		return fmt.Errorf("reset sequence failed: %v", err)
	}
	return nil
}

// ListAll returns all items from the inventory table,
// ordered by ID. This is a simple bulk fetch function.
//
// Use cautiously for very large databases. For pagination,
// use ListItemsPaged() or ItemIterator().
func ListAll(db *sql.DB) ([]Item, error) {
	rows, err := db.Query(`
        SELECT id, description, location, status, remarks
        FROM inventory ORDER BY id`)
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
			return nil, fmt.Errorf("scan failed: %v", err)
		}
		items = append(items, item)
	}
	return items, nil
}

// GetItemByID returns a single item with the given ID.
func GetItemByID(db *sql.DB, id int) (Item, error) {
	var item Item
	row := db.QueryRow(`
        SELECT id, description, location, status, remarks
        FROM inventory WHERE id = ?`, id)
	err := row.Scan(
		&item.ID, &item.Description,
		&item.Location, &item.Status, &item.Remarks)
	if err != nil {
		if err == sql.ErrNoRows {
			return item, fmt.Errorf("item %d not found", id)
		}
		return item, fmt.Errorf("query failed: %v", err)
	}
	return item, nil
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
