// json.go - Part of the `inventory` Package
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
	"encoding/json"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

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

// ExportJSONstr allows to export the complete database as a JSON string.
func ExportJSONstr(db *sql.DB) (string, error) {
	rows, err := db.Query("SELECT id, description, location, status, remarks" +
		" FROM inventory ORDER BY id")
	if err != nil {
		return "", fmt.Errorf("query failed: %v", err)
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.ID, &item.Description, &item.Location,
			&item.Status, &item.Remarks)
		if err != nil {
			return "", fmt.Errorf("scan failed: %v", err)
		}
		items = append(items, item)
	}

	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return "", fmt.Errorf("JSON marshal failed: %v", err)
	}

	return string(data), nil
}

// ImportJSONstr helps to import a bulk of Records from a JSON array
// supplied as a string.
func ImportJSONstr(db *sql.DB, jsonStr string) error {
	var items []Item
	err := json.Unmarshal([]byte(jsonStr), &items)
	if err != nil {
		return fmt.Errorf("JSON unmarshal failed: %v", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("transaction failed: %v", err)
	}

	for _, item := range items {
		_, err := tx.Exec("INSERT OR REPLACE INTO inventory ("+
			"id, description, location, status, remarks)"+
			"VALUES (?, ?, ?, ?, ?)",
			item.ID, item.Description, item.Location, item.Status, item.Remarks)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("import failed for item %d: %v", item.ID, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit failed: %v", err)
	}

	return nil
}

// ExportJSONByID finds outputs a JSON sting for a given ID value of the
// Database provided it exists.
func ExportJSONByID(db *sql.DB, id int) (string, error) {
	var item Item
	row := db.QueryRow("SELECT id, description, location, status, remarks"+
		" FROM inventory WHERE id = ?", id)
	err := row.Scan(&item.ID, &item.Description, &item.Location, &item.Status,
		&item.Remarks)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("Item with ID %d not found", id)
		}
		return "", fmt.Errorf("query failed: %v", err)
	}

	data, err := json.MarshalIndent(item, "", "  ")
	if err != nil {
		return "", fmt.Errorf("JSON marshal failed: %v", err)
	}

	return string(data), nil
}

// ImportJSONRecord helps to add or modify a row in database. This is
// done on the basis of the supplied ID value.
func ImportJSONRecord(db *sql.DB, jsonStr string) error {
	var item Item
	err := json.Unmarshal([]byte(jsonStr), &item)
	if err != nil {
		return fmt.Errorf("JSON unmarshal failed: %v", err)
	}

	_, err = db.Exec("INSERT OR REPLACE INTO inventory ("+
		"id, description, location, status, remarks)"+
		"VALUES (?, ?, ?, ?, ?)",
		item.ID, item.Description, item.Location, item.Status, item.Remarks)
	if err != nil {
		return fmt.Errorf("insert/update failed for item %d: %v", item.ID, err)
	}

	return nil
}
