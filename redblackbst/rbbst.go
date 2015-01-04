package redblackbst

// RedBlack is a sorted map built on a left leaning red black balanced
// search sorted map. It stores VType values, keyed by KType.
type RedBlack struct {
	root *node
}

// New creates a sorted map.
func New() *RedBlack { return &RedBlack{} }

// IsEmpty tells if the sorted map contains no key/value.
func (r RedBlack) IsEmpty() bool {
	return r.root == nil
}

// Size of the sorted map.
func (r RedBlack) Size() int { return size(r.root) }

// Clear all the values in the sorted map.
func (r *RedBlack) Clear() { r.root = nil }

// Put a value in the sorted map at key `k`. The old value at `k` is returned
// if the key was already present.
func (r *RedBlack) Put(k KType, v VType) (old VType, overwrite bool) {
	r.root, old, overwrite = put(r.root, k, v)
	return
}

func put(h *node, k KType, v VType) (_ *node, old VType, overwrite bool) {
	if h == nil {
		return newNode(k, v, 1, red), old, overwrite
	}

	cmp := k.Compare(h.key)
	if cmp < 0 {
		h.left, old, overwrite = put(h.left, k, v)
	} else if cmp > 0 {
		h.right, old, overwrite = put(h.right, k, v)
	} else {
		overwrite = true
		old = h.val
		h.val = v
	}

	if isRed(h.right) && !isRed(h.left) {
		h = rotateLeft(h)
	}
	if isRed(h.left) && isRed(h.left.left) {
		h = rotateRight(h)
	}
	if isRed(h.left) && isRed(h.right) {
		flipColors(h)
	}
	h.n = size(h.left) + size(h.right) + 1
	return h, old, overwrite
}

// Get a value from the sorted map at key `k`. Returns false
// if the key doesn't exist.
func (r RedBlack) Get(k KType) (VType, bool) {
	return loopGet(r.root, k)
}

func loopGet(h *node, k KType) (v VType, ok bool) {
	for h != nil {
		cmp := k.Compare(h.key)
		if cmp == 0 {
			return h.val, true
		} else if cmp < 0 {
			h = h.left
		} else if cmp > 0 {
			h = h.right
		}
	}
	return
}

// Has tells if a value exists at key `k`. This is short hand for `Get.
func (r RedBlack) Has(k KType) bool {
	_, ok := loopGet(r.root, k)
	return ok
}

// Min returns the smallest key/value in the sorted map, if it exists.
func (r RedBlack) Min() (k KType, v VType, ok bool) {
	if r.root == nil {
		return
	}
	h := min(r.root)
	return h.key, h.val, true
}

func min(x *node) *node {
	if x.left == nil {
		return x
	}
	return min(x.left)
}

// Max returns the largest key/value in the sorted map, if it exists.
func (r RedBlack) Max() (k KType, v VType, ok bool) {
	if r.root == nil {
		return
	}
	h := max(r.root)
	return h.key, h.val, true
}

func max(x *node) *node {
	if x.right == nil {
		return x
	}
	return max(x.right)
}

// Floor returns the largest key/value in the sorted map that is smaller than
// `k`.
func (r RedBlack) Floor(key KType) (k KType, v VType, ok bool) {
	x := floor(r.root, key)
	if x == nil {
		return
	}
	return x.key, x.val, true
}

func floor(h *node, k KType) *node {
	if h == nil {
		return nil
	}
	cmp := k.Compare(h.key)
	if cmp == 0 {
		return h
	}
	if cmp < 0 {
		return floor(h.left, k)
	}
	t := floor(h.right, k)
	if t != nil {
		return t
	}
	return h
}

// Ceiling returns the smallest key/value in the sorted map that is larger than
// `k`.
func (r RedBlack) Ceiling(key KType) (k KType, v VType, ok bool) {
	x := ceiling(r.root, key)
	if x == nil {
		return
	}
	return x.key, x.val, true
}

func ceiling(h *node, k KType) *node {
	if h == nil {
		return nil
	}
	cmp := k.Compare(h.key)
	if cmp == 0 {
		return h
	}
	if cmp > 0 {
		return ceiling(h.right, k)
	}
	t := ceiling(h.left, k)
	if t != nil {
		return t
	}
	return h
}

// Select key of rank k, meaning the k-th biggest KType in the sorted map.
func (r RedBlack) Select(key int) (k KType, v VType, ok bool) {
	x := nodeselect(r.root, key)
	if x == nil {
		return
	}
	return x.key, x.val, true
}

func nodeselect(x *node, k int) *node {
	if x == nil {
		return nil
	}
	t := size(x.left)
	if t > k {
		return nodeselect(x.left, k)
	} else if t < k {
		return nodeselect(x.right, k-t-1)
	} else {
		return x
	}
}

// Rank is the number of keys less than `k`.
func (r RedBlack) Rank(k KType) int {
	return keyrank(k, r.root)
}

func keyrank(k KType, h *node) int {
	if h == nil {
		return 0
	}
	cmp := k.Compare(h.key)
	if cmp < 0 {
		return keyrank(k, h.left)
	} else if cmp > 0 {
		return 1 + size(h.left) + keyrank(k, h.right)
	} else {
		return size(h.left)
	}
}

// Keys visit each keys in the sorted map, in order.
// It stops when visit returns false.
func (r RedBlack) Keys(visit func(KType, VType) bool) {
	min, _, ok := r.Min()
	if !ok {
		return
	}
	// if the min exists, then the max must exist
	max, _, _ := r.Max()
	r.RangedKeys(min, max, visit)
}

// RangedKeys visit each keys between lo and hi in the sorted map, in order.
// It stops when visit returns false.
func (r RedBlack) RangedKeys(lo, hi KType, visit func(KType, VType) bool) {
	keys(r.root, visit, lo, hi)
}

func keys(h *node, visit func(KType, VType) bool, lo, hi KType) bool {
	if h == nil {
		return true
	}
	cmplo := lo.Compare(h.key)
	cmphi := hi.Compare(h.key)
	if cmplo < 0 {
		if !keys(h.left, visit, lo, hi) {
			return false
		}
	}
	if cmplo <= 0 && cmphi >= 0 {
		if !visit(h.key, h.val) {
			return false
		}
	}
	if cmphi > 0 {
		if !keys(h.right, visit, lo, hi) {
			return false
		}
	}
	return true
}

// DeleteMin removes the smallest key and its value from the sorted map.
func (r *RedBlack) DeleteMin() (oldk KType, oldv VType, ok bool) {
	r.root, oldk, oldv, ok = deleteMin(r.root)
	if !r.IsEmpty() {
		r.root.color = black
	}
	return
}

func deleteMin(h *node) (_ *node, oldk KType, oldv VType, ok bool) {
	if h == nil {
		return nil, oldk, oldv, false
	}

	if h.left == nil {
		return nil, h.key, h.val, true
	}
	if !isRed(h.left) && !isRed(h.left.left) {
		h = moveRedLeft(h)
	}
	h.left, oldk, oldv, ok = deleteMin(h.left)
	return balance(h), oldk, oldv, ok
}

// DeleteMax removes the largest key and its value from the sorted map.
func (r *RedBlack) DeleteMax() (oldk KType, oldv VType, ok bool) {
	r.root, oldk, oldv, ok = deleteMax(r.root)
	if !r.IsEmpty() {
		r.root.color = black
	}
	return
}

func deleteMax(h *node) (_ *node, oldk KType, oldv VType, ok bool) {
	if h == nil {
		return nil, oldk, oldv, ok
	}
	if isRed(h.left) {
		h = rotateRight(h)
	}
	if h.right == nil {
		return nil, h.key, h.val, true
	}
	if !isRed(h.right) && !isRed(h.right.left) {
		h = moveRedRight(h)
	}
	h.right, oldk, oldv, ok = deleteMax(h.right)
	return balance(h), oldk, oldv, ok
}

// Delete key `k` from sorted map, if it exists.
func (r *RedBlack) Delete(k KType) (old VType, ok bool) {
	if r.root == nil {
		return
	}
	r.root, old, ok = delete(r.root, k)
	if !r.IsEmpty() {
		r.root.color = black
	}
	return
}

func delete(h *node, k KType) (_ *node, old VType, ok bool) {

	if h == nil {
		return h, old, false
	}

	if k.Compare(h.key) < 0 {
		if h.left == nil {
			return h, old, false
		}

		if !isRed(h.left) && !isRed(h.left.left) {
			h = moveRedLeft(h)
		}

		h.left, old, ok = delete(h.left, k)
		h = balance(h)
		return h, old, ok
	}

	if isRed(h.left) {
		h = rotateRight(h)
	}

	if k.Compare(h.key) == 0 && h.right == nil {
		return nil, h.val, true
	}

	if h.right != nil && !isRed(h.right) && !isRed(h.right.left) {
		h = moveRedRight(h)
	}

	if k.Compare(h.key) == 0 {

		var subk KType
		var subv VType
		h.right, subk, subv, ok = deleteMin(h.right)

		old, h.key, h.val = h.val, subk, subv
		ok = true
	} else {
		h.right, old, ok = delete(h.right, k)
	}

	h = balance(h)
	return h, old, ok
}

// deletions

func moveRedLeft(h *node) *node {
	flipColors(h)
	if isRed(h.right.left) {
		h.right = rotateRight(h.right)
		h = rotateLeft(h)
		flipColors(h)
	}
	return h
}

func moveRedRight(h *node) *node {
	flipColors(h)
	if isRed(h.left.left) {
		h = rotateRight(h)
		flipColors(h)
	}
	return h
}

func balance(h *node) *node {
	if isRed(h.right) {
		h = rotateLeft(h)
	}
	if isRed(h.left) && isRed(h.left.left) {
		h = rotateRight(h)
	}
	if isRed(h.left) && isRed(h.right) {
		flipColors(h)
	}
	h.n = size(h.left) + size(h.right) + 1
	return h
}

// nodes

const (
	red   = true
	black = false
)

type node struct {
	key         KType
	val         VType
	left, right *node
	n           int
	color       bool
}

func newNode(k KType, v VType, n int, color bool) *node {
	return &node{key: k, val: v, n: n, color: color}
}

func isRed(x *node) bool { return (x != nil) && (x.color == red) }

func rotateLeft(h *node) *node {
	x := h.right
	h.right = x.left
	x.left = h
	x.color = h.color
	h.color = red
	x.n = h.n
	h.n = 1 + size(h.left) + size(h.right)
	return x
}

func rotateRight(h *node) *node {
	x := h.left
	h.left = x.right
	x.right = h
	x.color = h.color
	h.color = red
	x.n = h.n
	h.n = 1 + size(h.left) + size(h.right)
	return x
}

func flipColors(h *node) {
	h.color = !h.color
	h.left.color = !h.left.color
	h.right.color = !h.right.color
}

func size(x *node) int {
	if x == nil {
		return 0
	}
	return x.n
}
