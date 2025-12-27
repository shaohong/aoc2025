package main

import (
	"math/rand"
	"time"
)

func rotate90Square(p Polyomino) Polyomino {
	n := len(p.cells)
	out := Polyomino{id: p.id, cells: make([][]int, n)}
	for r := 0; r < n; r++ {
		out.cells[r] = make([]int, n)
		for c := 0; c < n; c++ {
			out.cells[r][c] = p.cells[n-1-c][r]
		}
	}
	return out
}

func uniqueRotationsSquare(p Polyomino) []Polyomino {
	seen := map[string]bool{}
	out := make([]Polyomino, 0)
	cur := p
	for i := 0; i < 4; i++ {
		k := polyKey(cur)
		if !seen[k] {
			seen[k] = true
			out = append(out, cur)
		}
		cur = rotate90Square(cur)
	}
	return out
}

func polyKey(p Polyomino) string {
	b := make([]byte, 0, len(p.cells)*len(p.cells[0])+len(p.cells))
	for r := 0; r < len(p.cells); r++ {
		for c := 0; c < len(p.cells[0]); c++ {
			if p.cells[r][c] == 1 {
				b = append(b, '1')
			} else {
				b = append(b, '0')
			}
		}
		b = append(b, '/')
	}
	return string(b)
}

type placement struct {
	typeID     int // which polyomino type
	rotIdx     int // which rotation of that type
	topR, topC int
	shape      Polyomino
	score      int
}

// CanPackWithCounts tries to place EXACTLY counts[i] copies of polyTypes[i] on the grid
// without overlap, within empty cells (grid==0). Full coverage is NOT required.
//
// grid: 0 empty, 1 blocked/occupied. Placements write 1s.
// polyTypes: each type is NxN square, all same N.
// counts: how many of each type to place.
//
// Returns (ok, outGrid, placements). If ok==false, outGrid/placements represent the best attempt found.
func CanPackWithCounts(grid [][]int, polyTypes []Polyomino, counts []int) (bool, [][]int, []placement) {
	if len(polyTypes) != len(counts) {
		return false, cloneGrid(grid), nil
	}

	N := len(polyTypes[0].cells) // assume NxN
	// Precompute unique rotations per type
	typeRots := make([][]Polyomino, len(polyTypes))
	for i := range polyTypes {
		typeRots[i] = uniqueRotationsSquare(polyTypes[i])
	}

	// total number of polyominos to place
	totalNeeded := 0
	for _, c := range counts {
		totalNeeded += c
	}

	// First: deterministic greedy+repair
	ok, gBest, plBest := runAttempt(grid, typeRots, counts, N, attemptConfig{
		tries:        1,
		repairSteps:  max(2000, totalNeeded*80),
		removeMin:    1,
		removeMax:    4,
		emptyScanCap: 0,
		seed:         1,
	})
	if ok {
		return true, gBest, plBest
	}

	// Then: randomized attempts (common in packing feasibility)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	bestGrid := gBest
	bestPls := plBest
	bestPlacedCells := countOnesDelta(grid, gBest)

	randomTries := 120
	for t := 0; t < randomTries; t++ {
		seed := r.Int63()
		ok2, g2, pl2 := runAttempt(grid, typeRots, counts, N, attemptConfig{
			tries:        1,
			repairSteps:  max(600, 60*totalNeeded),
			removeMin:    2,
			removeMax:    6,
			emptyScanCap: 0, // scan all empties
			seed:         seed,
		})
		if ok2 {
			return true, g2, pl2
		}
		score := countOnesDelta(grid, g2)
		if score > bestPlacedCells {
			bestPlacedCells = score
			bestGrid = g2
			bestPls = pl2
		}
	}

	return false, bestGrid, bestPls
}

func runAttempt(grid [][]int, typeRots [][]Polyomino, counts []int, N int, cfg attemptConfig) (bool, [][]int, []placement) {
	// 1. Greedy step: pick the next placement that looks best right now
	// 	* implemented by findBestNextPlacement(...)
	//  * it scans empty “anchor” cells and tries all remaining piece types + rotations
	//  * chooses the placement with the highest heuristic score (scoreContact)

	// 2. Repair step (when stuck): if it can’t place anything, it “undoes” the last k placements and tries again
	//  * remove k last placements (k randomized in a range)
	//  * restore counts
	//  * continue

	r := rand.New(rand.NewSource(cfg.seed))

	g := cloneGrid(grid)
	rem := append([]int(nil), counts...)
	pls := make([]placement, 0)

	remainingPieces := 0
	for _, c := range rem {
		remainingPieces += c
	}

	steps := 0
	for remainingPieces > 0 && steps < cfg.repairSteps {
		steps++

		pl, ok := findBestNextPlacement(g, typeRots, rem, N, r, cfg.emptyScanCap)
		if ok {
			applyPlacement(g, pl, N)
			pls = append(pls, pl)
			rem[pl.typeID]--
			remainingPieces--
			continue
		}

		// No feasible placement anywhere for remaining pieces: repair by removing last k pieces
		if len(pls) == 0 {
			return false, g, pls
		}
		k := cfg.removeMin + r.Intn(max(1, cfg.removeMax-cfg.removeMin+1))
		if k > len(pls) {
			k = len(pls)
		}
		for i := 0; i < k; i++ {
			last := pls[len(pls)-1]
			pls = pls[:len(pls)-1]
			unapplyPlacement(g, last, N)
			rem[last.typeID]++
			remainingPieces++
		}
	}

	return remainingPieces == 0, g, pls
}

// findBestNextPlacement scans empty cells and tries all remaining piece types (and rotations).
// Picks the best-scoring placement found. If none exist, returns ok=false.
func findBestNextPlacement(g [][]int, typeRots [][]Polyomino, rem []int, N int, r *rand.Rand, emptyScanCap int) (placement, bool) {
	H, W := len(g), len(g[0])
	bestFound := false
	var best placement

	scanned := 0

	// You can change scan order (e.g. bottom-left) for different behavior.
	for rr := 0; rr < H; rr++ {
		for cc := 0; cc < W; cc++ {
			if g[rr][cc] != 0 {
				continue
			}
			scanned++
			if emptyScanCap > 0 && scanned > emptyScanCap {
				// keep runtime bounded; we tried a bunch of anchors already
				break
			}
			anchorR, anchorC := rr, cc

			// Try each type with remaining copies
			// Randomize type order slightly to escape repeated failures
			typeOrder := make([]int, 0, len(typeRots))
			for t := range typeRots {
				if rem[t] > 0 {
					typeOrder = append(typeOrder, t)
				}
			}
			r.Shuffle(len(typeOrder), func(i, j int) { typeOrder[i], typeOrder[j] = typeOrder[j], typeOrder[i] })

			for _, t := range typeOrder {
				rotsForType := typeRots[t]
				for ri, sh := range rotsForType {
					// Enumerate placements that cover the anchor: map each filled cell (r,c) to anchor
					for pr := 0; pr < N; pr++ {
						for pc := 0; pc < N; pc++ {
							if sh.cells[pr][pc] == 0 {
								continue
							}
							topR := anchorR - pr
							topC := anchorC - pc
							if canPlace(g, sh, topR, topC, N) {
								pl := placement{
									typeID: t,
									rotIdx: ri,
									topR:   topR,
									topC:   topC,
									shape:  sh,
									score:  scoreContact(g, sh, topR, topC, N),
								}
								if !bestFound || pl.score > best.score {
									bestFound = true
									best = pl
								}
							}
						}
					}
				}
			}

			// Small speed optimization: if we found *something* at an early anchor,
			// often good enough; comment out if you want more global searching.
			if bestFound {
				return best, true
			}
		}
	}

	return placement{}, false
}

/* ------------------------------ grid ops ------------------------------ */

func canPlace(grid [][]int, sh Polyomino, topR, topC, N int) bool {
	H, W := len(grid), len(grid[0])
	if topR < 0 || topC < 0 || topR+N > H || topC+N > W {
		return false
	}
	for r := 0; r < N; r++ {
		for c := 0; c < N; c++ {
			if sh.cells[r][c] == 1 && grid[topR+r][topC+c] != 0 {
				return false
			}
		}
	}
	return true
}

func applyPlacement(grid [][]int, pl placement, N int) {
	sh := pl.shape
	for r := 0; r < N; r++ {
		for c := 0; c < N; c++ {
			if sh.cells[r][c] == 1 {
				grid[pl.topR+r][pl.topC+c] = 1
			}
		}
	}
}

func unapplyPlacement(grid [][]int, pl placement, N int) {
	sh := pl.shape
	for r := 0; r < N; r++ {
		for c := 0; c < N; c++ {
			if sh.cells[r][c] == 1 {
				grid[pl.topR+r][pl.topC+c] = 0
			}
		}
	}
}

func cloneGrid(g [][]int) [][]int {
	out := make([][]int, len(g))
	for i := range g {
		out[i] = make([]int, len(g[0]))
		copy(out[i], g[i])
	}
	return out
}

// counts how many cells became occupied (1) compared to original grid
func countOnesDelta(orig, now [][]int) int {
	H, W := len(orig), len(orig[0])
	d := 0
	for r := 0; r < H; r++ {
		for c := 0; c < W; c++ {
			if orig[r][c] == 0 && now[r][c] == 1 {
				d++
			}
		}
	}
	return d
}

/* ------------------------------ attempt engine ------------------------------ */

type attemptConfig struct {
	tries        int
	repairSteps  int
	removeMin    int
	removeMax    int
	emptyScanCap int
	seed         int64
}

/* ------------------------------ scoring ------------------------------ */

// Simple, effective: maximize contact with boundary or already-occupied cells.
// (You can add “avoid tiny holes” penalties later if needed.)
func scoreContact(g [][]int, sh Polyomino, topR, topC, N int) int {
	H, W := len(g), len(g[0])
	contact := 0
	dirs := [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

	for r := 0; r < N; r++ {
		for c := 0; c < N; c++ {
			if sh.cells[r][c] == 0 {
				continue
			}
			gr, gc := topR+r, topC+c
			for _, d := range dirs {
				nr, nc := gr+d[0], gc+d[1]
				if nr < 0 || nr >= H || nc < 0 || nc >= W {
					contact++
				} else if g[nr][nc] == 1 {
					contact++
				}
			}
		}
	}
	return contact
}
