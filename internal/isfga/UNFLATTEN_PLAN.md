# Unflatten Hierarchy Implementation Plan

## Overview
Convert flat classification hierarchies (Kingdom, Phylum, Family, Genus, Species fields) to parent/child relationships (ParentID).

## Design Decisions

### Two-Phase Algorithm
1. **Phase 1**: Scan all rows, build parent map with unique parent taxa
2. **Phase 2**: Generate NameUsage records row-by-row

### ID Strategy
- Use real IDs from dataset when available (GenusID, FamilyID, etc.)
- Generate "sf-" prefixed IDs when missing
- "sf-" prefix distinguishes synthetic from original IDs

### Normalization
- Normalize rank-name for map keys: lowercase, space→underscore
- Preserve original capitalization in Name field
- Example: "Homo sapiens" → key: "species-homo_sapiens", Name: "Homo sapiens"

### Parent Determination
- Use Rank.Order() method to determine hierarchy
- Find highest populated rank < current rank
- Example: Species (Order=800) → Genus (Order=700) if Genus is populated

### Root Handling
- Create root entry with key "root": ID="", ParentKey=""
- Orphans (no higher classification) get ParentKey="root"
- Unresolvable parents fall back to root

### Inconsistencies
- First occurrence wins for conflicting IDs
- Typos in names create separate taxa (garbage in, garbage out)
- Log warnings but continue processing

## Data Structures

### ParentTaxonInfo

we will not need this type outside of the task, so we keep it private.
```go
type taxonInfo struct {
    ID        string  // Database ID: real ID or "sf-{counter}" or "sf-{rank}-{name}"
    ParentKey string  // Normalized parent map key: "family-hominidae", "root" for the root record.
    Name      string  // Original name with capitalization: "Homo"
}
```

### Parent Map
```
map[string]taxonInfo
Key: normalized "rank-name" (e.g., "genus-homo")
Value: taxonInfo struct
```

## Implementation Tasks

### 1. Setup
- [ ] Remove WithUnflatten field from config/config.go
- [ ] Add UnflattenHierarchy(ctx context.Context) error to Enricher interface

### 2. Foundation
- [ ] Add Order() method to Rank in pkg/coldp/enum_rank.go
  - Returns a slice sorted by rank position.
- [ ] Create normalization function
  - Lowercase, space→underscore, trim
- [ ] Create taxonInfo struct in internal/isfga/unflatten.go

### 3. Phase 1: Build Parent Map
- [ ] Scan all NameUsage/Taxon rows
- [ ] For each classification rank (Kingdom, Phylum, ..., Genus, Species):
  - Extract rank name and ID
  - Normalize to create map key
  - If key not in map:
    - Use real ID if available (KingdomID, GenusID, etc.)
    - Generate "sf-" ID if missing
    - Determine ParentKey (next higher populated rank)
    - Store ParentTaxonInfo
  - If key exists: skip (first occurrence wins)
- [ ] Add root entry: map["root"] = {ID: "", ParentKey: "", Name: ""}

### 4. Phase 2: Generate Output
- [ ] Create writtenParents set: map[string]bool
- [ ] For each original row:
  - Create leaf NameUsage with updated ParentID
  - Walk parent chain from leaf to root:
    ```go
    parentKey := determineParentKey(row)
    for parentKey != "" {
        if !writtenParents[parentKey] {
            parent := parentMap[parentKey]
            createParentNameUsage(parent)
            writtenParents[parentKey] = true
        }
        // Handle missing parent
        if parentInfo, exists := parentMap[parentKey]; exists {
            parentKey = parentInfo.ParentKey
        } else {
            log.Warn("Parent not found: %s", parentKey)
            break  // Unresolvable - leaf will point to root
        }
    }
    ```

### 5. Helper Functions Needed

#### determineParentKey(taxon)
1. Get taxon's rank (from Rank field or infer)
2. Get taxon's rank order value
3. Scan flat classification fields (Kingdom, Phylum, ..., Genus)
4. Find highest populated rank with Order() < taxon.Order()
5. Return normalized key for that rank

#### createParentNameUsage(parentInfo)
1. Create coldp.NameUsage struct
2. Set ID = parentInfo.ID
3. Set ParentID = parentMap[parentInfo.ParentKey].ID (or "" if root)
4. Set Rank = extract from map key
5. Set appropriate classification field (Genus, Family, etc.) = parentInfo.Name
6. Populate higher classification by walking ParentKey chain
7. Write to output SFGA

#### normalize(name)
```go
strings.TrimSpace(strings.Replace(strings.ToLower(name), " ", "_", -1))
```

### 6. Edge Cases
- [ ] Taxon with no Rank field: infer from populated classification fields
- [ ] Inconsistent hierarchies: Species1 has Subgenus, Species2 doesn't → handled by map
- [ ] Conflicting IDs: Genus="Homo" with GenusID="123" and GenusID="456" → first wins
- [ ] Missing parent in chain → set to root, log warning
- [ ] Typos in names → create separate taxa
- [ ] Empty/NULL classification fields → skip

### 7. Testing
- [ ] Test with complete flat hierarchy (Kingdom→Species)
- [ ] Test with partial hierarchy (only Genus, Family, Order)
- [ ] Test with inconsistent subgenus usage
- [ ] Test with real IDs vs generated IDs
- [ ] Test with typos/duplicates
- [ ] Test orphan taxa
- [ ] Test large dataset performance

## Example Flow

**Input:**
```
Row 1: ID="sp1", Rank=Species, Species="sapiens", Genus="Homo", GenusID="", Family="Hominidae", FamilyID="fam1"
Row 2: ID="sp2", Rank=Species, Species="erectus", Genus="Homo", GenusID="", Family="Hominidae", FamilyID="fam1"
```

**Phase 1 Map:**
```
"species-sapiens"    → {ID: "sp1" (from row), ParentKey: "genus-homo", Name: "sapiens"}
"species-erectus"    → {ID: "sp2" (from row), ParentKey: "genus-homo", Name: "erectus"}
"genus-homo"         → {ID: "sf-1" (generated), ParentKey: "family-hominidae", Name: "Homo"}
"family-hominidae"   → {ID: "fam1" (real ID), ParentKey: "", Name: "Hominidae"}
""                   → {ID: "", ParentKey: "", Name: ""}
```

**Phase 2 Output:**
```
NameUsage: ID="sp1", ParentID="sf-1", Species="sapiens", Genus="Homo", Family="Hominidae"
NameUsage: ID="sf-1", ParentID="fam1", Rank=Genus, Genus="Homo", Family="Hominidae"
NameUsage: ID="fam1", ParentID="", Rank=Family, Family="Hominidae"

NameUsage: ID="sp2", ParentID="sf-1", Species="erectus", Genus="Homo", Family="Hominidae"
(sf-1 and fam1 already written, skipped)
```
