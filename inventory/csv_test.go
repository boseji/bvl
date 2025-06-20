// csv_test.go - Part of Tests for the `inventory` Package
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

//
// Unit tests for CSV import/export functions
// Uses in-memory SQLite DB and temp CSV files
//

package inventory_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/boseji/bvl/inventory"
)

func setupCSVTestDB(t *testing.T) *inventory.InventoryDB {
	inv := inventory.NewInventoryDB(":memory:")
	if inv == nil {
		t.Fatal("failed to create InventoryDB")
	}
	return inv
}

func TestExportImportCSV(t *testing.T) {
	inv := setupCSVTestDB(t)
	defer inv.Close()

	// Add test data
	item := inventory.Item{
		Description: "UPS",
		Location:    "Rack 1",
		Status:      "Operational",
		Remarks:     "installed",
	}
	err := inv.AddItem(item)
	if err != nil {
		t.Fatalf("AddItem failed: %v", err)
	}

	// Export to CSV
	tmpfile := filepath.Join(os.TempDir(), "test_inventory_export.csv")
	defer os.Remove(tmpfile)

	err = inv.ExportCSV(tmpfile)
	if err != nil {
		t.Fatalf("ExportCSV failed: %v", err)
	}

	// Clear table before import
	err = inv.WithTransaction(func(tx inventory.Execer) error {
		_, err := tx.Exec(`DELETE FROM inventory`)
		return err
	})
	if err != nil {
		t.Fatalf("clear table failed: %v", err)
	}

	items, _ := inv.ListAll()
	if len(items) != 0 {
		t.Fatalf("expected empty DB after clear, got %d", len(items))
	}

	// Import from CSV
	err = inv.ImportCSV(tmpfile)
	if err != nil {
		t.Fatalf("ImportCSV failed: %v", err)
	}

	items, _ = inv.ListAll()
	if len(items) != 1 {
		t.Fatalf("expected 1 item after import, got %d", len(items))
	}

	if items[0].Description != "UPS" {
		t.Errorf("unexpected Description after import: %s", items[0].Description)
	}
}

func TestViewCSV(t *testing.T) {
	inv := setupCSVTestDB(t)
	defer inv.Close()

	item := inventory.Item{
		Description: "Router",
		Location:    "Rack 2",
		Status:      "Active",
		Remarks:     "setup",
	}
	_ = inv.AddItem(item)

	tmpfile := filepath.Join(os.TempDir(), "test_inventory_view.csv")
	defer os.Remove(tmpfile)

	err := inv.ExportCSV(tmpfile)
	if err != nil {
		t.Fatalf("ExportCSV failed: %v", err)
	}

	// Just check ViewCSV runs without error
	err = inventory.ViewCSV(tmpfile)
	if err != nil {
		t.Fatalf("ViewCSV failed: %v", err)
	}
}

func TestImportCSV_BadFile(t *testing.T) {
	inv := setupCSVTestDB(t)
	defer inv.Close()

	err := inv.ImportCSV("no_such_file.csv")
	if err == nil {
		t.Fatalf("expected error for missing file")
	}
}

func TestExportCSV_BadPath(t *testing.T) {
	inv := setupCSVTestDB(t)
	defer inv.Close()

	err := inv.ExportCSV("/bad/path/test.csv")
	if err == nil {
		t.Fatalf("expected error for bad path")
	}
}
