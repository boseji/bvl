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
	"os"
)

// ExportCSV writes all inventory records to a CSV file.
//
// Usage:
//
//	err := ExportCSV(db, "inventory.csv")
//
// The CSV will have the following columns:
//
//	id, description, location, status, remarks
//
// Existing file will be overwritten.
//
// Returns error if file cannot be written or query fails.
func ExportCSV(db *sql.DB, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("create csv failed: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"id", "description", "location", "status", "remarks"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("write csv header failed: %v", err)
	}

	rows, err := db.Query(`SELECT id, description, location, status, remarks FROM inventory ORDER BY id`)
	if err != nil {
		return fmt.Errorf("query inventory failed: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ID, &item.Description, &item.Location, &item.Status, &item.Remarks); err != nil {
			return fmt.Errorf("scan failed: %v", err)
		}
		record := []string{
			fmt.Sprintf("%d", item.ID),
			item.Description,
			item.Location,
			item.Status,
			item.Remarks,
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("write csv row failed: %v", err)
		}
	}

	return nil
}

// ImportCSV reads inventory records from a CSV file and imports them.
//
// Existing records with matching IDs will be replaced.
//
// Usage:
//
//	err := ImportCSV(db, "inventory.csv")
//
// CSV format must have columns:
//
//	id, description, location, status, remarks
//
// Each row is imported using AppendItem().
//
// Returns error on file error, parse error, or DB error.
func ImportCSV(exec Execer, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("open csv failed: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("read csv failed: %v", err)
	}

	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) != 5 {
			return fmt.Errorf("csv row %d has wrong column count", i)
		}

		var item Item
		fmt.Sscanf(row[0], "%d", &item.ID)
		item.Description = row[1]
		item.Location = row[2]
		item.Status = row[3]
		item.Remarks = row[4]

		if err := AppendItem(exec, item); err != nil {
			return fmt.Errorf("import row %d failed: %v", i, err)
		}
	}

	return nil
}

// ViewCSV prints the content of a CSV file to stdout.
//
// Usage:
//
//	err := ViewCSV("inventory.csv")
//
// The output is formatted as columns:
//
//	id  description  location  status  remarks
//
// Errors are returned if the file cannot be read.
func ViewCSV(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("open csv failed: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("read csv failed: %v", err)
	}

	for _, row := range rows {
		fmt.Printf("%-5s %-20s %-15s %-15s %-s\n", row[0], row[1], row[2], row[3], row[4])
	}

	return nil
}

// ExportCSV writes all inventory records to CSV using InventoryDB.
//
// Usage:
//
//	err := inv.ExportCSV("inventory.csv")
//
// Same as ExportCSV() raw.
func (inv *InventoryDB) ExportCSV(filename string) error {
	return ExportCSV(inv.db, filename)
}

// ImportCSV imports inventory records from CSV using InventoryDB.
//
// Usage:
//
//	err := inv.ImportCSV("inventory.csv")
//
// The import runs inside a transaction.
func (inv *InventoryDB) ImportCSV(filename string) error {
	return inv.WithTransaction(func(tx Execer) error {
		return ImportCSV(tx, filename)
	})
}
