package lik

func IsEqualItems(itema Itemer, itemb Itemer) bool {
	if itema == nil && itemb == nil {
		return true
	} else if itema == nil || itemb == nil {
		return false
	} else if itema.IsSet() && itemb.IsSet() {
		return IsEqualSets(itema.ToSet(), itemb.ToSet())
	} else if itema.IsSet() || itemb.IsSet() {
		return false
	} else if itema.IsList() && itemb.IsList() {
		return IsEqualLists(itema.ToList(), itemb.ToList())
	} else if itema.IsList() || itemb.IsList() {
		return false
	} else {
		return itema.ToString() == itemb.ToString()
	}
}

func IsEqualSets(seta Seter, setb Seter) bool {
	if seta == nil && setb == nil {
		return true
	} else if seta == nil || setb == nil {
		return false
	} else if keysa, keysb := seta.SortKeys(), setb.SortKeys(); len(keysa) != len(keysb) {
		return false
	} else {
		for k, key := range keysa {
			if keyb := keysb[k]; keyb != key {
				return false
			} else if vala, valb := seta.GetItem(key), setb.GetItem(keyb); !IsEqualItems(vala, valb) {
				return false
			}
		}
		return true
	}
}

func IsEqualLists(lista Lister, listb Lister) bool {
	if lista == nil && listb == nil {
		return true
	} else if lista == nil || listb == nil {
		return false
	} else if lena, lenb := lista.Count(), listb.Count(); lena != lenb {
		return false
	} else {
		for k := 0; k < lena; k++ {
			if vala, valb := lista.GetItem(k), listb.GetItem(k); !IsEqualItems(vala, valb) {
				return false
			}
		}
		return true
	}
}
