>
> ॐᳬ᳞ भूर्भुवः स्वः
>
> तत्स॑वि॒तुर्वरे॑ण्यं॒
>
> भर्गो॑ दे॒वस्य॑ धीमहि।
>
> धियो॒ यो नः॑ प्रचो॒दया॑त्॥
>

# बी.वी.एल - बोसजी के द्वारा रचित भंडार लेखांकन हेतु तन्त्राक्ष्।

> एक सुगम एवं उपयोगी भंडार संचालन हेतु तन्त्राक्ष्।

***एक रचनात्मक भारतीय उत्पाद ।***

<p align="center">
  <img src="../docs/assets/icon-full-color-mini.png" alt="Inventory CLI Logo">
</p>

## `bvl` - Boseji's Inventory Management Program

[![GitHub license](https://img.shields.io/github/license/boseji/bvl)](../LICENSE.txt)

Easy to use and useful stock, goods and materials handling software designed in `Golang` and a prototype in _Python_.

Simple GPLv2 CLI tool to manage inventory with `SQLite`, `CSV` import/export and logging.

---

## `inventory` Package for Inventory DB Implementation in Golang

### Core

* InventoryDB wrapper — safe, transactional DB access
* In-memory or file-based SQLite support
* Configurable sequence start (`IndexStart`)

### Data Model

* `Item` struct — ID, Description, Location, Status, Remarks (with `FormatRemarks()`)
* `FormatRemarks()` — consistent timestamped remarks
* JSON tags — for web/app/API compatibility

### Database Operations

* `AddItem()`
* `EditItem()`
* `DeleteItem()`
* `AppendItem()`
* `AppendRemarksEntry()`
* `GetItemByID()`
* `ListAll()`
* `ListItemsPaged()` — with pagination
* `NewItemIterator()` — with streaming Next()
* `ResetSequence()`

### CSV Support

* `ExportCSV()`
* `ImportCSV()`
* `ViewCSV()`
* InventoryDB wrappers
* CLI-friendly and Excel-friendly CSV format

### JSON Support

* `ExportJSON()`
* `ImportJSON()`
* `ViewJSON()`
* `ExportJSONToString()`
* `ImportJSONFromString()`
* InventoryDB wrappers
* Web/Electron/API/jq friendly JSON format

### Item JSON Helpers

* `Item.ToJSON()`
* `Item.FromJSON()`

### Test Suite

* `db_test.go` — core DB functions
* `inventorydb_test.go` — InventoryDB methods
* `csv_test.go` — CSV
* `json_test.go` — JSON
* Full error path coverage
* Rollback scenarios covered

### Ready for:

* CLI tools
* Web frontends
* Electron apps
* REST APIs
* jq pipelines
* Automation scripts
* CI/CD integration

---

## License

This project is released under the GNU General Public License v2. See the [LICENSE](../LICENSE.txt) file for details.

Sources: <https://github.com/boseji/bvl>

`bvl` - Boseji's Inventory Management Program

Copyright (C) 2025 by Abhijit Bose (aka. Boseji)

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License version 2 only
as published by the Free Software Foundation.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.

You should have received a copy of the GNU General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

SPDX-License-Identifier: `GPL-2.0-only`

Full Name: `GNU General Public License v2.0 only`

Please visit <https://spdx.org/licenses/GPL-2.0-only.html> for details.

---
