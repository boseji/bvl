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
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/boseji/bsg/gen"
	_ "github.com/mattn/go-sqlite3"
)

const (
	// Starting value of the Index
	IndexStart = 1000
)

// OpenDB creates the Database connection. It would also create the
// database if one does not exists.
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

// List all item just list all the items in the Database.
func ListAllItems(db *sql.DB) {
	rows, err := db.Query("SELECT id, description, location, status, remarks" +
		" FROM inventory ORDER BY id")
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	defer rows.Close()

	fmt.Printf("%-5s %-40s %-25s %-15s %s\n", "ID", "Description",
		"Location", "Status", "Remarks")
	fmt.Println(strings.Repeat("-", 110))
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.ID, &item.Description, &item.Location,
			&item.Status, &item.Remarks)
		if err != nil {
			log.Fatalf("Scan failed: %v", err)
		}
		fmt.Printf("%-5d %-40s %-25s %-15s %s\n", item.ID, item.Description,
			item.Location, item.Status, item.Remarks)
	}
}

// AddItem is used to add a new row in the database
func AddItem(db *sql.DB, desc, loc, status, remarks string) {
	stmt, err := db.Prepare("INSERT INTO inventory (" +
		"description, location, status, remarks)" +
		" VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Fatalf("Prepare failed: %v", err)
	}
	defer stmt.Close()

	// Log like look
	ts := gen.BST().Format("[2006-01-02 15:04:05]")
	newRemarks := fmt.Sprintf("%s %s", ts, remarks)

	res, err := stmt.Exec(desc, loc, status, newRemarks)
	if err != nil {
		log.Fatalf("Insert failed: %v", err)
	}
	id, _ := res.LastInsertId()
	fmt.Printf("Item added with ID %d\n", id)
}

// EditItem helps to alter a particular part of the row in Database.
// If a particular field need not be changed then it can be left blank
// and the same will be skipped in the update command.
func EditItem(db *sql.DB, id int, desc, loc, status, remarks string) {
	parts := []string{}
	args := []interface{}{}
	if desc != "" {
		parts = append(parts, "description = ?")
		args = append(args, desc)
	}
	if loc != "" {
		parts = append(parts, "location = ?")
		args = append(args, loc)
	}
	if status != "" {
		parts = append(parts, "status = ?")
		args = append(args, status)
	}
	if remarks != "" {
		parts = append(parts, "remarks = ?")
		args = append(args, remarks)
	}
	if len(parts) == 0 {
		fmt.Println("Nothing to update.")
		return
	}
	args = append(args, id)
	stmt := fmt.Sprintf("UPDATE inventory SET %s WHERE id = ?",
		strings.Join(parts, ", "))
	_, err := db.Exec(stmt, args...)
	if err != nil {
		log.Fatalf("Update failed: %v", err)
	}
	fmt.Printf("Item %d updated\n", id)
}

// DeleteItem helps to delete an item based on its ID
func DeleteItem(db *sql.DB, id int) {
	_, err := db.Exec("DELETE FROM inventory WHERE id = ?", id)
	if err != nil {
		log.Fatalf("Delete failed: %v", err)
	}
	fmt.Printf("Item %d deleted\n", id)
}

// LogRemark helps to add a message like a Log into the Remarks field
func LogRemark(db *sql.DB, id int, message string) {
	var current string
	row := db.QueryRow("SELECT remarks FROM inventory WHERE id = ?", id)
	err := row.Scan(&current)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	ts := gen.BST().Format("[2006-01-02 15:04:05]")
	newRemarks := current + "\n" + fmt.Sprintf("%s %s", ts, message)
	_, err = db.Exec("UPDATE inventory SET remarks = ? WHERE id = ?",
		newRemarks, id)
	if err != nil {
		log.Fatalf("Update failed: %v", err)
	}
	fmt.Printf("Log appended to item %d\n", id)
}

// ExportCSV allows to transform all the Data into a CSV file.
func ExportCSV(db *sql.DB, csvFile string) {
	rows, err := db.Query("SELECT id, description, location, status, remarks" +
		" FROM inventory ORDER BY id")
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	defer rows.Close()

	f, err := os.Create(csvFile)
	if err != nil {
		log.Fatalf("Failed to create CSV: %v", err)
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	writer.Write([]string{"id", "description", "location", "status", "remarks"})

	count := 0
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.ID, &item.Description, &item.Location,
			&item.Status, &item.Remarks)
		if err != nil {
			log.Fatalf("Scan failed: %v", err)
		}
		record := []string{
			fmt.Sprintf("%d", item.ID),
			item.Description,
			item.Location,
			item.Status,
			item.Remarks,
		}
		writer.Write(record)
		count++
	}
	fmt.Printf("Exported %d rows to %s\n", count, csvFile)
}

// ImportCSV allows the data to be read back from a CSV file
func ImportCSV(db *sql.DB, csvFile string) {
	f, err := os.Open(csvFile)
	if err != nil {
		log.Fatalf("Failed to open CSV: %v", err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Failed to read CSV: %v", err)
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Transaction failed: %v", err)
	}

	for i, row := range records {
		if i == 0 {
			continue // header
		}
		_, err := tx.Exec("INSERT OR REPLACE INTO inventory ("+
			"id, description, location, status, remarks)"+
			"VALUES (?, ?, ?, ?, ?)",
			row[0], row[1], row[2], row[3], row[4])
		if err != nil {
			log.Fatalf("Import failed at row %d: %v", i, err)
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Fatalf("Commit failed: %v", err)
	}
	fmt.Printf("Imported %d records from %s\n", len(records)-1, csvFile)
}

// ExportJSONFile would allow the full database to be exported as a JSON file
func ExportJSONFile(db *sql.DB, jsonFile string) {
	rows, err := db.Query("SELECT id, description, location, status, remarks FROM inventory ORDER BY id")
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.ID, &item.Description, &item.Location, &item.Status, &item.Remarks)
		if err != nil {
			log.Fatalf("Scan failed: %v", err)
		}
		items = append(items, item)
	}

	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		log.Fatalf("JSON marshal failed: %v", err)
	}

	err = os.WriteFile(jsonFile, data, 0644)
	if err != nil {
		log.Fatalf("Failed to write JSON file: %v", err)
	}

	fmt.Printf("Exported %d items to %s\n", len(items), jsonFile)
}

// ImportJSONFile can be used to Import a set of records into the database
// using a .JSON file.
func ImportJSONFile(db *sql.DB, jsonFile string) {
	data, err := os.ReadFile(jsonFile)
	if err != nil {
		log.Fatalf("Failed to read JSON file: %v", err)
	}

	var items []Item
	err = json.Unmarshal(data, &items)
	if err != nil {
		log.Fatalf("JSON unmarshal failed: %v", err)
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Transaction failed: %v", err)
	}

	for _, item := range items {
		_, err := tx.Exec(`INSERT OR REPLACE INTO inventory (id, description, location, status, remarks)
            VALUES (?, ?, ?, ?, ?)`, item.ID, item.Description, item.Location, item.Status, item.Remarks)
		if err != nil {
			log.Fatalf("Import failed for item %d: %v", item.ID, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatalf("Commit failed: %v", err)
	}

	fmt.Printf("Imported %d items from %s\n", len(items), jsonFile)
}
