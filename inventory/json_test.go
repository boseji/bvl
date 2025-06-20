// json_test.go - Part of Tests for the `inventory` Package
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
// Unit tests for JSON import/export functions
// Uses in-memory SQLite DB and temp JSON files
//

package inventory_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/boseji/bvl/inventory"
)

func setupJSONTestDB(t *testing.T) *inventory.InventoryDB {
	inv := inventory.NewInventoryDB(":memory:")
	if inv == nil {
		t.Fatal("failed to create InventoryDB")
	}
	return inv
}

func TestExportImportJSON(t *testing.T) {
	inv := setupJSONTestDB(t)
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

	// Export to JSON
	tmpfile := filepath.Join(os.TempDir(), "test_inventory_export.json")
	defer os.Remove(tmpfile)

	err = inv.ExportJSON(tmpfile)
	if err != nil {
		t.Fatalf("ExportJSON failed: %v", err)
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

	// Import from JSON
	err = inv.ImportJSON(tmpfile)
	if err != nil {
		t.Fatalf("ImportJSON failed: %v", err)
	}

	items, _ = inv.ListAll()
	if len(items) != 1 {
		t.Fatalf("expected 1 item after import, got %d", len(items))
	}

	if items[0].Description != "UPS" {
		t.Errorf("unexpected Description after import: %s", items[0].Description)
	}
}

func TestViewJSON(t *testing.T) {
	inv := setupJSONTestDB(t)
	defer inv.Close()

	item := inventory.Item{
		Description: "Router",
		Location:    "Rack 2",
		Status:      "Active",
		Remarks:     "setup",
	}
	_ = inv.AddItem(item)

	tmpfile := filepath.Join(os.TempDir(), "test_inventory_view.json")
	defer os.Remove(tmpfile)

	err := inv.ExportJSON(tmpfile)
	if err != nil {
		t.Fatalf("ExportJSON failed: %v", err)
	}

	// Just check ViewJSON runs without error
	err = inventory.ViewJSON(tmpfile)
	if err != nil {
		t.Fatalf("ViewJSON failed: %v", err)
	}
}

func TestExportJSONToString(t *testing.T) {
	inv := setupJSONTestDB(t)
	defer inv.Close()

	item := inventory.Item{
		Description: "Switch",
		Location:    "Rack 3",
		Status:      "Ready",
		Remarks:     "config done",
	}
	_ = inv.AddItem(item)

	jsonStr, err := inv.ExportJSONToString()
	if err != nil {
		t.Fatalf("ExportJSONToString failed: %v", err)
	}

	if len(jsonStr) == 0 {
		t.Fatalf("ExportJSONToString returned empty string")
	}
}

func TestImportJSONFromString(t *testing.T) {
	inv := setupJSONTestDB(t)
	defer inv.Close()

	jsonData := `[
      {
        "id": 1005,
        "description": "PDU",
        "location": "Rack 5",
        "status": "Installed",
        "remarks": "added"
      }
    ]`

	err := inv.ImportJSONFromString(jsonData)
	if err != nil {
		t.Fatalf("ImportJSONFromString failed: %v", err)
	}

	items, _ := inv.ListAll()
	if len(items) != 1 {
		t.Fatalf("expected 1 item after import, got %d", len(items))
	}

	if items[0].Description != "PDU" {
		t.Errorf("unexpected Description after import: %s", items[0].Description)
	}
}

func TestItemToJSONAndFromJSON(t *testing.T) {
	item := inventory.Item{
		ID:          1010,
		Description: "Firewall",
		Location:    "DC1",
		Status:      "Active",
		Remarks:     "setup complete",
	}

	jsonStr, err := item.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON failed: %v", err)
	}

	if len(jsonStr) == 0 {
		t.Fatalf("ToJSON returned empty string")
	}

	var newItem inventory.Item
	err = newItem.FromJSON(jsonStr)
	if err != nil {
		t.Fatalf("FromJSON failed: %v", err)
	}

	if newItem.Description != "Firewall" {
		t.Errorf("unexpected Description after FromJSON: %s", newItem.Description)
	}
}

func TestImportJSON_BadFile(t *testing.T) {
	inv := setupJSONTestDB(t)
	defer inv.Close()

	err := inv.ImportJSON("no_such_file.json")
	if err == nil {
		t.Fatalf("expected error for missing file")
	}
}

func TestExportJSON_BadPath(t *testing.T) {
	inv := setupJSONTestDB(t)
	defer inv.Close()

	err := inv.ExportJSON("/bad/path/test.json")
	if err == nil {
		t.Fatalf("expected error for bad path")
	}
}
