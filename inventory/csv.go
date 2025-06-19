// csv.go - Part of the `inventory` Package
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
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

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
