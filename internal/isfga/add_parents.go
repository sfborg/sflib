package isfga

import (
	"context"
	"errors"
	"log/slog"
	"strconv"
	"strings"

	"github.com/gnames/gnlib"
	"github.com/gnames/gnlib/ent/nomcode"
	"github.com/gnames/gnparser"
	"github.com/sfborg/sflib/pkg/coldp"
	"github.com/sfborg/sflib/pkg/parser"
	"github.com/sfborg/sflib/pkg/sfga"
	"golang.org/x/sync/errgroup"
)

func (a *isfga) addParents(emptySfga sfga.Archive) error {
	var err error
	var nodes map[string]*node
	var maxSfID int
	nodes, maxSfID, err = a.buildNodes()
	if err != nil {
		return err
	}
	err = a.addNameUsages(emptySfga, nodes, maxSfID)
	if err != nil {
		return err
	}

	return nil
}

type node struct {
	ID         string
	tmpID      int
	sciName    string
	rank       coldp.Rank
	parentName string
	path       []coldp.NameUsage
}

func (a *isfga) buildNodes() (map[string]*node, int, error) {
	// string is scientific name
	nameMap := make(map[string]*node)
	var tmpID int

	var maxSfID int

	g, ctx := errgroup.WithContext(context.Background())
	chIn := make(chan coldp.NameUsage)

	g.Go(func() error {
		for nu := range chIn {
			if strings.HasPrefix(nu.ID, "sf-") {
				i, err := strconv.Atoi(nu.ID[3:])
				if err == nil && i > maxSfID {
					maxSfID = i
				}
			}

			path := getPathFromNameUsage(nu)
			if existing := nameMap[nu.ScientificName]; existing != nil {
				// If synthetic node was created before real record,
				// update it with the real ID
				if existing.tmpID > 0 && nu.ID != "" {
					existing.ID = nu.ID
					existing.tmpID = 0
				} else {
					slog.Warn(
						"Duplicated scientific name, skipping", "name", nu.ScientificName,
					)
				}
				continue
			}
			n := node{
				ID: nu.ID, sciName: nu.ScientificName,
				rank: nu.Rank, path: path,
			}
			tmpID = processNode(nameMap, n, tmpID)
		}
		return nil
	})

	err := a.LoadNameUsages(ctx, chIn)
	if err != nil {
		return nil, 0, err
	}
	close(chIn)

	err = g.Wait()
	if err != nil && !errors.Is(err, context.Canceled) {
		return nil, 0, err
	}

	return nameMap, maxSfID, nil
}

func processNode(nameMap map[string]*node, n node, tmpID int) int {
	nameMap[n.sciName] = &n
	if len(n.path) == 0 {
		return tmpID
	}
	path := n.path
	var parentName string
	for i := range path {
		nu := path[i]
		pathNode := node{
			sciName:    nu.ScientificName,
			rank:       nu.Rank,
			parentName: parentName,
			path:       path[:i],
		}
		parentName = nu.ScientificName
		if nameMap[nu.ScientificName] == nil {
			tmpID++
			pathNode.tmpID = tmpID
			nameMap[nu.ScientificName] = &pathNode
		}
	}
	n.parentName = parentName
	return tmpID
}

// formatSubgenus constructs the full subgenus scientific name in the format
// "Genus (Subgenus)". Returns empty string if subgenus is empty.
func formatSubgenus(genus, subgenus string) string {
	if subgenus == "" {
		return ""
	}
	if genus == "" {
		return subgenus
	}
	return genus + " (" + subgenus + ")"
}

// getPathFromNameUsage creates a slice of all parent taxa included into flat
// classification from a NameUsage. Empty records and self-references (where
// the path entry rank matches the record's rank) are ignored.
func getPathFromNameUsage(nu coldp.NameUsage) []coldp.NameUsage {
	taxa := []coldp.NameUsage{
		{Rank: coldp.Realm, ScientificName: nu.Realm},
		{Rank: coldp.Kingdom, ScientificName: nu.Kingdom},
		{Rank: coldp.Phylum, ScientificName: nu.Phylum},
		{Rank: coldp.Subphylum, ScientificName: nu.Subphylum},
		{Rank: coldp.Class, ScientificName: nu.Class},
		{Rank: coldp.Subclass, ScientificName: nu.Subclass},
		{Rank: coldp.Order, ScientificName: nu.Order},
		{Rank: coldp.Suborder, ScientificName: nu.Suborder},
		{Rank: coldp.Superfamily, ScientificName: nu.Superfamily},
		{Rank: coldp.Family, ScientificName: nu.Family},
		{Rank: coldp.Subfamily, ScientificName: nu.Subfamily},
		{Rank: coldp.Tribe, ScientificName: nu.Tribe},
		{Rank: coldp.Subtribe, ScientificName: nu.Subtribe},
		{Rank: coldp.Genus, ScientificName: nu.Genus},
		{Rank: coldp.Subgenus, ScientificName: formatSubgenus(nu.Genus, nu.Subgenus)},
		{Rank: coldp.Section, ScientificName: nu.Section},
		{Rank: coldp.Species, ScientificName: nu.Species},
	}
	var idx int
	for i := range taxa {
		// Skip empty entries
		if taxa[i].ScientificName == "" {
			continue
		}
		// Skip self-references: when flat hierarchy rank matches record's rank
		if taxa[i].Rank == nu.Rank {
			continue
		}
		taxa[idx] = taxa[i]
		idx++
	}
	return taxa[:idx]
}

func (a *isfga) addNameUsages(
	emptySfga sfga.Archive,
	ids map[string]*node,
	maxSfID int,
) error {
	var err error
	g, ctx := errgroup.WithContext(context.Background())
	chIn := make(chan coldp.NameUsage)
	chOut := make(chan []coldp.NameUsage)
	// for now we do not parallelize
	pp := parser.Pool(1)

	g.Go(func() error {
		defer close(chOut)
		return a.addParentData(ctx, pp, chIn, chOut, ids, maxSfID)
	})

	g.Go(func() error {
		return saveNameUsages(ctx, emptySfga, chOut)
	})

	err = a.LoadNameUsages(ctx, chIn)
	if err != nil {
		return err
	}
	close(chIn)

	err = g.Wait()
	if err != nil && !errors.Is(err, context.Canceled) {
		return err
	}

	return nil
}

func (a *isfga) addParentData(
	ctx context.Context,
	pp map[nomcode.Code]chan gnparser.GNparser,
	chIn <-chan coldp.NameUsage,
	chOut chan<- []coldp.NameUsage,
	nameMap map[string]*node,
	maxSfID int,
) error {
	var err error
	var nus []coldp.NameUsage
	batchSize := a.cfg.BatchSize
	processed := make(gnlib.Set[string])
	batch := make([]coldp.NameUsage, 0, batchSize)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case nu, ok := <-chIn:
			if !ok {
				// send remaining batch
				if len(batch) > 0 {
					select {
					case <-ctx.Done():
						return ctx.Err()
					case chOut <- batch:
					}
				}
				return nil
			}
			nus, err = a.addMissingData(nu, pp, nameMap, processed, maxSfID)
			if err != nil {
				return err
			}
			batch = append(batch, nus...)
			if len(batch) >= batchSize {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case chOut <- batch:
				}
				batch = make([]coldp.NameUsage, 0, batchSize)
			}
		}
	}
}

func (a *isfga) addMissingData(
	nu coldp.NameUsage,
	pp map[nomcode.Code]chan gnparser.GNparser,
	nameMap map[string]*node,
	processed gnlib.Set[string],
	maxSfID int,
) ([]coldp.NameUsage, error) {
	var res []coldp.NameUsage

	var parser gnparser.GNparser
	if nu.ParseQuality.Valid {
		parser = <-pp[nu.Code]
		defer returnParser(pp[nu.Code], parser)
	}

	path := getPathFromNameUsage(nu)
	// get node with missing data
	n := nameMap[nu.ScientificName]

	// convert IDs to int-based IDs
	if n.parentName != "" {
		nu.ParentID = nodeWithID(n.parentName, maxSfID, nameMap).ID
	}

	// populate flat classification IDs for existing records
	addFlatClassificationIDs(&nu, nameMap, maxSfID)

	// add parent NameUsages that haven't been processed yet
	// only for synthetic nodes (tmpID > 0), not for existing NameUsages
	for _, v := range path {
		nUp := nameMap[v.ScientificName]
		if nUp == nil {
			continue
		}
		// ensure ID is set (generates if needed)
		_ = nodeWithID(v.ScientificName, maxSfID, nameMap)
		if processed.Has(nUp.ID) {
			continue
		}
		// skip nodes that already exist as real NameUsages (tmpID == 0)
		if nUp.tmpID == 0 {
			processed.Add(nUp.ID)
			continue
		}
		v.ID = nUp.ID
		v.ScientificNameString = v.ScientificName
		v.Code = nu.Code
		v.TaxonomicStatus = coldp.AcceptedTS
		nParent := nodeWithID(nUp.parentName, maxSfID, nameMap)
		v.ParentID = nParent.ID
		// populate flat classification from path
		populateFlatClassification(&v, nUp.path)
		addFlatClassificationIDs(&v, nameMap, maxSfID)

		if nu.ParseQuality.Valid {
			v.Amend(parser)
		}

		processed.Add(v.ID)
		res = append(res, v)
	}

	processed.Add(nu.ID)
	res = append(res, nu)

	return res, nil
}

// populateFlatClassification sets flat classification fields on a NameUsage
// based on the path of parent taxa.
func populateFlatClassification(nu *coldp.NameUsage, path []coldp.NameUsage) {
	for _, p := range path {
		switch p.Rank {
		case coldp.Realm:
			nu.Realm = p.ScientificName
		case coldp.Kingdom:
			nu.Kingdom = p.ScientificName
		case coldp.Phylum:
			nu.Phylum = p.ScientificName
		case coldp.Subphylum:
			nu.Subphylum = p.ScientificName
		case coldp.Class:
			nu.Class = p.ScientificName
		case coldp.Subclass:
			nu.Subclass = p.ScientificName
		case coldp.Order:
			nu.Order = p.ScientificName
		case coldp.Suborder:
			nu.Suborder = p.ScientificName
		case coldp.Superfamily:
			nu.Superfamily = p.ScientificName
		case coldp.Family:
			nu.Family = p.ScientificName
		case coldp.Subfamily:
			nu.Subfamily = p.ScientificName
		case coldp.Tribe:
			nu.Tribe = p.ScientificName
		case coldp.Subtribe:
			nu.Subtribe = p.ScientificName
		case coldp.Genus:
			nu.Genus = p.ScientificName
		case coldp.Subgenus:
			nu.Subgenus = p.ScientificName
		case coldp.Section:
			nu.Section = p.ScientificName
		case coldp.Species:
			nu.Species = p.ScientificName
		}
	}
}

// addFlatClassificationIDs populates the ID fields for flat classification
// ranks for existing NameUsage records (not synthetic ones).
func addFlatClassificationIDs(
	nu *coldp.NameUsage,
	nameMap map[string]*node,
	maxSfID int,
) {
	// helper to get ID from any node in nameMap
	getID := func(name string) string {
		if name == "" {
			return ""
		}
		n := nameMap[name]
		if n == nil {
			return ""
		}
		if n.ID == "" {
			n.ID = "sf-" + strconv.Itoa(maxSfID+n.tmpID)
		}
		return n.ID
	}

	nu.RealmID = getID(nu.Realm)
	nu.KingdomID = getID(nu.Kingdom)
	nu.PhylumID = getID(nu.Phylum)
	nu.SubphylumID = getID(nu.Subphylum)
	nu.ClassID = getID(nu.Class)
	nu.SubclassID = getID(nu.Subclass)
	nu.OrderID = getID(nu.Order)
	nu.SuborderID = getID(nu.Suborder)
	nu.SuperfamilyID = getID(nu.Superfamily)
	nu.FamilyID = getID(nu.Family)
	nu.SubfamilyID = getID(nu.Subfamily)
	nu.TribeID = getID(nu.Tribe)
	nu.SubtribeID = getID(nu.Subtribe)
	nu.GenusID = getID(nu.Genus)
	nu.SubgenusID = getID(formatSubgenus(nu.Genus, nu.Subgenus))
	nu.SectionID = getID(nu.Section)
	nu.SpeciesID = getID(nu.Species)
}

func nodeWithID(name string, maxSfID int, nameMap map[string]*node) *node {
	if name == "" {
		return &node{}
	}
	n := nameMap[name]
	if n.ID != "" {
		return n
	}
	intID := maxSfID + n.tmpID
	n.ID = "sf-" + strconv.Itoa(intID)
	return n
}

func saveNameUsages(
	ctx context.Context,
	a sfga.Archive,
	chOut <-chan []coldp.NameUsage,
) error {
	// ensure connection
	_, err := a.Connect()
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case batch, ok := <-chOut:
			if !ok {
				return nil
			}
			if err := a.InsertNameUsages(batch); err != nil {
				return err
			}
		}
	}
}

func returnParser(ch chan<- gnparser.GNparser, p gnparser.GNparser) {
	ch <- p
}
