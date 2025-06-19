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

// ExportCSV writes all inventory items to CSV file.
//
// Fields:
// id, description, location, status, remarks
//
// Returns error if writing fails.
func ExportCSV(db *sql.DB, filename string) error {
	items, err := ListAll(db)
	if err != nil {
		return fmt.Errorf("export csv failed: %v", err)
	}

	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("could not create file: %v", err)
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	header := []string{
		"id", "description", "location", "status", "remarks",
	}
	writer.Write(header)

	for _, item := range items {
		record := []string{
			fmt.Sprintf("%d", item.ID),
			item.Description,
			item.Location,
			item.Status,
			item.Remarks,
		}
		writer.Write(record)
	}

	fmt.Printf("exported %d items to %s\n", len(items), filename)
	return nil
}

// ImportCSV reads inventory items from CSV file.
//
// Existing records with same ID will be replaced.
// CSV must have same field order as exported.
//
// Returns error if reading fails.
func ImportCSV(db *sql.DB, filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("could not open file: %v", err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("could not read csv: %v", err)
	}

	if len(records) < 1 {
		return fmt.Errorf("csv is empty")
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("transaction failed: %v", err)
	}

	for i, row := range records {
		if i == 0 {
			continue // skip header
		}

		if len(row) < 5 {
			tx.Rollback()
			return fmt.Errorf(
				"invalid csv row %d: expected 5 fields", i)
		}

		_, err := tx.Exec(`
            INSERT OR REPLACE INTO inventory
            (id, description, location, status, remarks)
            VALUES (?, ?, ?, ?, ?)`,
			row[0], row[1], row[2], row[3], row[4])
		if err != nil {
			tx.Rollback()
			return fmt.Errorf(
				"import failed at row %d: %v", i, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit failed: %v", err)
	}

	fmt.Printf(
		"imported %d items from %s\n", len(records)-1, filename)
	return nil
}

// ViewCSV displays CSV file to stdout.
//
// No database required.
// Useful for debugging or piping to jq.
func ViewCSV(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("could not open file: %v", err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("could not read csv: %v", err)
	}

	for i, row := range records {
		if i == 0 {
			fmt.Println("--- CSV Header ---")
		}
		fmt.Printf("%d: %s\n", i, row)
	}

	fmt.Printf(
		"displayed %d rows from %s\n", len(records), filename)
	return nil
}
