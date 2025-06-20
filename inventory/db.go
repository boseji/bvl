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

// AppendItem inserts or replaces an item in the inventory table,
// using the provided ID. If an item with the same ID already exists,
// it will be replaced with the new values.
//
// The remarks field is processed through item.FormatRemarks()
// to ensure consistent timestamped formatting.
//
// Typical usage:
//
//	item := Item{
//	    ID:          1234,
//	    Description: "UPS 3KVA",
//	    Location:    "Rack 5",
//	    Status:      "Operational",
//	    Remarks:     "installed new unit",
//	}
//	err := AppendItem(tx, item)
//
// Resulting record:
//
//	id          = 1234
//	description = "UPS 3KVA"
//	location    = "Rack 5"
//	status      = "Operational"
//	remarks     = "[2025-06-20 12:30] installed new unit"
//
// Use cases:
//
// - To update an existing record fully (replace)
// - To insert a new record with known ID
// - To bulk insert/update items
//
// Notes:
//
// - Safe to call repeatedly with the same item
// - Will replace existing record (INSERT OR REPLACE)
// - Does not check for ID conflicts beyond replacement
// - Remarks field will always be formatted via FormatRemarks()
// - If ID is not set, use AddItem() instead
// - Works with both *sql.DB and *sql.Tx.
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
// - Works with both *sql.DB and *sql.Tx.
func AppendRemarksEntry(exec Execer, id int, message string) error {
	t := gen.BST().Format("2006-01-02 15:04")
	formatted := fmt.Sprintf("[%s] %s", t, message)

	// Appends new entry to remarks field:
	// COALESCE(remarks, '') → ensure string (not null)
	// char(10) → newline character
	// '||' → SQLite concat
	// Result: remarks = old + '\n' + new entry
	res, err := exec.Exec(`
        UPDATE inventory
        SET remarks = 
            COALESCE(remarks, '') || char(10) || ?
        WHERE id = ?`,
		formatted, id)
	if err != nil {
		return fmt.Errorf("append to remarks failed: %v", err)
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("append failed: no such ID %d", id)
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
// - Works with both *sql.DB and *sql.Tx.
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
// - Works with both *sql.DB and *sql.Tx.
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

// DeleteItem removes a record from the inventory table by ID.
//
// If the specified ID does not exist, the operation is a no-op
// (no error is returned).
//
// Typical usage:
//
//	err := inv.DeleteItem(1234)
//
// Result:
//
// - If item with id = 1234 exists → record is deleted
// - If no such item → nothing is done, no error
//
// Use cases:
//
// - To permanently remove an inventory record
// - To clean up old or duplicate items
// - To reset part of the inventory manually
//
// Notes:
//
// - This is a destructive operation (cannot be undone)
// - Should typically be logged via remarks before use
// - Use AppendRemarksEntry() if you want an audit trail before delete
// - Works with both *sql.DB and *sql.Tx.
func DeleteItem(exec Execer, id int) error {
	_, err := exec.Exec(`
        DELETE FROM inventory
        WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete failed: %v", err)
	}
	return nil
}

// ResetSequence wraps ResetSequence with automatic transaction.
//
// Resets the auto-increment sequence for the inventory table
// back to IndexStart (default: 1000).
//
// Typically used after manually clearing records, or for test setups.
//
// Usage:
//
//	err := inv.ResetSequence()
//
// Result:
//
// - Sets the internal sqlite_sequence counter for 'inventory' table
// - Next inserted record will use ID = IndexStart + 1
//
// Use cases:
//
// - After deleting all items (clear inventory)
// - For test environments to reset IDs
// - To reinitialize an empty database
//
// Notes:
//
// - Does not delete records (use DeleteItem or manual purge first)
// - Safe to call multiple times
// - Has no effect if records still exist with higher IDs
// - Works with both *sql.DB and *sql.Tx.
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

// ListAll returns all items in the inventory table, sorted by ID.
//
// This is a read-only operation. It does not require a transaction.
// It can be used for reporting, exporting, or displaying all items.
//
// Usage:
//
//	items, err := inv.ListAll()
//	if err != nil {
//	    // handle error
//	}
//	for _, item := range items {
//	    fmt.Println(item.ID, item.Description, item.Status)
//	}
//
// Result:
//
// - Returns []Item containing all inventory records
// - Sorted by id ASC (oldest first)
//
// Use cases:
//
// - To display the full inventory
// - To export data to CSV, JSON
// - For reports or dashboards
//
// Notes:
//
//   - If the table is empty, returns an empty slice (no error)
//   - The remarks field will be returned in raw form
//     (use item.FormatRemarks() for display)
//   - This method does not paginate large inventories
//     (use ListItemsPaged for that)
//   - Use cautiously for very large databases. For pagination,
//     use ListItemsPaged() or ItemIterator().
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

// GetItemByID returns a single item from the inventory table
// that matches the given ID.
//
// If no item is found with the given ID, returns an error:
//
//	"item <id> not found"
//
// Typical usage:
//
//	item, err := inv.GetItemByID(1234)
//	if err != nil {
//	    // handle error (not found, or query error)
//	} else {
//	    fmt.Println(item.Description, item.Status)
//	}
//
// Result:
//
// - If item exists → returns populated Item struct
// - If not found → returns zero-value Item + error
//
// Use cases:
//
// - To display or edit a specific inventory item
// - To retrieve details for audit or reporting
// - To check existence of an item by ID
//
// Notes:
//
//   - This is a read-only query (no transaction needed)
//   - The remarks field is returned as raw string
//     (use item.FormatRemarks() for formatted display)
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

// ListItemsPaged returns a slice of items after a given starting ID,
// up to a specified limit.
//
// This is a read-only operation. It does not require a transaction.
//
// Usage:
//
//	items, err := inv.ListItemsPaged(lastID, 10)
//	if err != nil {
//	    // handle error
//	}
//	for _, item := range items {
//	    fmt.Printf("%d: %s\n", item.ID, item.Description)
//	}
//
// Result:
//
// - Returns up to 'limit' number of items with id > afterID
// - Results are sorted by id ASC
//
// Use cases:
//
// - For paging through large inventories
// - For implementing UI pagination
// - For batch export or processing
//
// Notes:
//
// - If no items match the query, returns an empty slice
// - Use afterID = 0 to start from beginning
// - If fewer than 'limit' items remain, returns as many as available
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
